package handle

import (
	"baby-fried-rice/internal/pkg/module/im/config"
	"baby-fried-rice/internal/pkg/module/im/log"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"io"
	"io/ioutil"
	"sync"
	"time"
)

const (
	RtcpPLIInterval = time.Second * 3
	Compress        = false
	rtcpPLIInterval = time.Second * 3
)

func zip(in []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(in)
	if err != nil {
		panic(err)
	}
	err = gz.Flush()
	if err != nil {
		panic(err)
	}
	err = gz.Close()
	if err != nil {
		panic(err)
	}
	return b.Bytes()
}

func unzip(in []byte) []byte {
	var b bytes.Buffer
	_, err := b.Write(in)
	if err != nil {
		panic(err)
	}
	r, err := gzip.NewReader(&b)
	if err != nil {
		panic(err)
	}
	res, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return res
}

func Decode(in string, obj interface{}) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}

	if Compress {
		b = unzip(b)
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
}

func Encode(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	if Compress {
		b = zip(b)
	}

	return base64.StdEncoding.EncodeToString(b)
}

var (
	videoLocalTrackMap   = new(sync.Map)
	audioLocalTrackMap   = new(sync.Map)
	peerConnectionConfig = webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: config.GetConfig().Stuns,
			},
			{
				URLs:       config.GetConfig().Turn.URLs,
				Username:   config.GetConfig().Turn.Username,
				Credential: config.GetConfig().Turn.Credential,
			},
		},
	}
)

type SessionWebRTC interface {
}

type SessionWebRTCProvider struct {
	userPeerConnectionLock sync.RWMutex
	userPeerConnections    map[string]*userPeerConnection
}

type userPeerConnection struct {
	sessionId   int64
	video       bool
	joinLock    sync.RWMutex
	joins       map[string]*joinUserPeerConnection
	conn        *webrtc.PeerConnection
	videoStatus atomic.Bool
	videoTrack  *webrtc.TrackLocalStaticRTP
	audioStatus atomic.Bool
	audioTrack  *webrtc.TrackLocalStaticRTP
}

type joinUserPeerConnection struct {
	sessionId               int64
	video                   bool
	conn                    *webrtc.PeerConnection
	receiveVideoTrackCancel context.CancelFunc
	receiveAudioTrackCancel context.CancelFunc
}

func (upc *userPeerConnection) getVideoStatus() bool {
	return upc.videoStatus.Load()
}

func (upc *userPeerConnection) getAudioStatus() bool {
	return upc.audioStatus.Load()
}

func NewSessionWebRTCProvider() SessionWebRTC {
	return &SessionWebRTCProvider{}
}

func (provider *SessionWebRTCProvider) getAccountIdUpc(accountId string, sessionId int64) (*userPeerConnection, error) {
	provider.userPeerConnectionLock.RLock()
	upc, exist := provider.userPeerConnections[accountId]
	if !exist {
		return nil, fmt.Errorf("user web rtc not exist")
	}
	upc.joinLock.RLock()
	return upc, nil
}

func (provider *SessionWebRTCProvider) setAccountIdUpc(accountId string, conn *userPeerConnection, sessionId int64) {
	provider.userPeerConnectionLock.Lock()
	provider.userPeerConnections[accountId] = conn
	provider.userPeerConnectionLock.Unlock()
}

func (provider *SessionWebRTCProvider) getJoinAccountIdUpc(accountId, joinAccountId string, sessionId int64) (*userPeerConnection, *joinUserPeerConnection, error) {
	provider.userPeerConnectionLock.RLock()
	upc, exist := provider.userPeerConnections[accountId]
	if !exist {
		return nil, nil, fmt.Errorf("user web rtc not exist")
	}
	upc.joinLock.RLock()
	var joinUpc *joinUserPeerConnection
	joinUpc, exist = upc.joins[joinAccountId]
	if !exist {
		return nil, nil, fmt.Errorf("you didn't join user' webrtc")
	}
	upc.joinLock.RUnlock()
	provider.userPeerConnectionLock.RUnlock()
	return upc, joinUpc, nil
}

func (provider *SessionWebRTCProvider) setJoinAccountIdUpc(accountId, joinAccountId string, joinUpc *joinUserPeerConnection) error {
	provider.userPeerConnectionLock.Lock()
	upc, exist := provider.userPeerConnections[accountId]
	if !exist {
		return fmt.Errorf("user web rtc not exist")
	}
	upc.joinLock.Lock()
	provider.userPeerConnections[accountId].joins[joinAccountId] = joinUpc
	upc.joinLock.Unlock()
	provider.userPeerConnectionLock.Unlock()
	return nil
}

// 创建会话中某个用户的webrtc
func (provider *SessionWebRTCProvider) CreateSessionUserWebRTC(sdp, accountId string, sessionId int64, video bool) (swapSdp string, err error) {
	offer := webrtc.SessionDescription{}
	Decode(sdp, &offer)
	peerConnection, err := webrtc.NewPeerConnection(peerConnectionConfig)
	if err != nil {
		return
	}
	var upc = &userPeerConnection{
		sessionId: sessionId,
		video:     video,
		joinLock:  sync.RWMutex{},
		joins:     make(map[string]*joinUserPeerConnection),
		conn:      peerConnection,
	}
	if video {
		if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
			return
		}
		upc.videoStatus.Store(true)
	}
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		return
	}
	upc.audioStatus.Store(true)
	provider.setAccountIdUpc(accountId, upc, sessionId)
	videoLocalTrackChan := make(chan *webrtc.TrackLocalStaticRTP)
	audioLocalTrackChan := make(chan *webrtc.TrackLocalStaticRTP)
	var key = fmt.Sprintf("%v-%v", accountId, sessionId)
	peerConnection.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		var kind = remoteTrack.Kind().String()
		go func() {
			ticker := time.NewTicker(rtcpPLIInterval)
			for range ticker.C {
				if rtcpSendErr := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(remoteTrack.SSRC())}}); rtcpSendErr != nil {
					log.Logger.Error(rtcpSendErr.Error())
				}
			}
		}()
		localTrack, newTrackErr := webrtc.NewTrackLocalStaticRTP(remoteTrack.Codec().RTPCodecCapability, kind, key)
		if newTrackErr != nil {
			panic(newTrackErr)
		}
		switch kind {
		case webrtc.RTPCodecTypeVideo.String():
			videoLocalTrackChan <- localTrack
		case webrtc.RTPCodecTypeAudio.String():
			audioLocalTrackChan <- localTrack
		default:
			err = webrtc.ErrUnknownType
			panic(err)
		}
		rtpBuf := make([]byte, 1400)
		for {
			i, _, readErr := remoteTrack.Read(rtpBuf)
			if readErr != nil {
				panic(readErr)
			}

			var write = false
			if kind == webrtc.RTPCodecTypeVideo.String() {
				write = upc.getVideoStatus()
			} else if kind == webrtc.RTPCodecTypeAudio.String() {
				write = upc.getAudioStatus()
			}
			if write {
				if _, err = localTrack.Write(rtpBuf[:i]); err != nil && !errors.Is(err, io.ErrClosedPipe) {
					panic(err)
				}
			}

		}
	})
	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		return
	}
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		return
	}
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		return
	}
	<-gatherComplete
	swapSdp = Encode(*peerConnection.LocalDescription())
	if video {
		go func(accountId string) {
			localTrack := <-videoLocalTrackChan
			upc, err = provider.getAccountIdUpc(accountId, sessionId)
			if err != nil {
				return
			}
			upc.videoTrack = localTrack
			provider.setAccountIdUpc(accountId, upc, sessionId)
		}(accountId)
	}
	go func(accountId string) {
		localTrack := <-audioLocalTrackChan
		upc, err = provider.getAccountIdUpc(accountId, sessionId)
		if err != nil {
			return
		}
		upc.audioTrack = localTrack
		provider.setAccountIdUpc(accountId, upc, sessionId)
	}(accountId)
	return
}

// 关闭会话中自己的webrtc
func (provider *SessionWebRTCProvider) CloseSessionWebRTC(accountId string, sessionId int64) error {
	return nil
}

// 关闭会话中自己的视频轨
func (provider *SessionWebRTCProvider) CloseWebRTCVideoTrack(accountId string, sessionId int64) error {
	upc, err := provider.getAccountIdUpc(accountId, sessionId)
	if err != nil {
		return err
	}
	if upc.getVideoStatus() {
		upc.videoStatus.Store(false)
	}
	return nil
}

// 打开会话中自己的视频轨
func (provider *SessionWebRTCProvider) OpenWebRTCVideoTrack(accountId string, sessionId int64) error {
	upc, err := provider.getAccountIdUpc(accountId, sessionId)
	if err != nil {
		return err
	}
	if !upc.getVideoStatus() {
		upc.videoStatus.Store(true)
	}
	return nil
}

// 关闭会话中自己的音频轨
func (provider *SessionWebRTCProvider) CloseWebRTCAudioTrack(accountId string, sessionId int64) error {
	upc, err := provider.getAccountIdUpc(accountId, sessionId)
	if err != nil {
		return err
	}
	if upc.getAudioStatus() {
		upc.audioStatus.Store(false)
	}
	return nil
}

// 打开会话中自己的音频轨
func (provider *SessionWebRTCProvider) OpenWebRTCAudioTrack(accountId string, sessionId int64) error {
	upc, err := provider.getAccountIdUpc(accountId, sessionId)
	if err != nil {
		return err
	}
	if !upc.getAudioStatus() {
		upc.audioStatus.Store(true)
	}
	return nil
}

// 加入会话中某个用户的webrtc（不包括自己）
func (provider *SessionWebRTCProvider) JoinSessionUserWebRTC(sdp, accountId, joinAccountId string, sessionId int64, video bool) (swapSdp string, err error) {
	var peerConnection *webrtc.PeerConnection
	if peerConnection, err = webrtc.NewPeerConnection(peerConnectionConfig); err != nil {
		return
	}
	var joinUpc = &joinUserPeerConnection{
		sessionId: sessionId,
		video:     video,
		conn:      peerConnection,
	}
	if err = provider.setJoinAccountIdUpc(accountId, joinAccountId, joinUpc); err != nil {
		return
	}
	if video {
		if err = provider.OpenUserWebRTCVideoTrack(accountId, joinAccountId, sessionId); err != nil {
			return
		}
	}
	if err = provider.OpenUserWebRTCAudioTrack(accountId, joinAccountId, sessionId); err != nil {
		return
	}
	recvOnlyOffer := webrtc.SessionDescription{}
	Decode(sdp, &recvOnlyOffer)
	if err = peerConnection.SetRemoteDescription(recvOnlyOffer); err != nil {
		return
	}
	var answer webrtc.SessionDescription
	if answer, err = peerConnection.CreateAnswer(nil); err != nil {
		return
	}
	var gatherComplete <-chan struct{}
	gatherComplete = webrtc.GatheringCompletePromise(peerConnection)
	if err = peerConnection.SetLocalDescription(answer); err != nil {
		return
	}
	<-gatherComplete
	swapSdp = Encode(*peerConnection.LocalDescription())
	return
}

// 关闭会话中某个用户的视频轨
func (provider *SessionWebRTCProvider) CloseUserWebRTCVideoTrack(accountId, joinAccountId string, sessionId int64) (err error) {
	_, joinUpc, err := provider.getJoinAccountIdUpc(accountId, joinAccountId, sessionId)
	if err != nil {
		return err
	}
	if joinUpc.receiveVideoTrackCancel == nil {
		return nil
	}
	joinUpc.receiveVideoTrackCancel()
	joinUpc.receiveVideoTrackCancel = nil
	return provider.setJoinAccountIdUpc(accountId, joinAccountId, joinUpc)
}

// 打开会话中某个用户的视频轨
func (provider *SessionWebRTCProvider) OpenUserWebRTCVideoTrack(accountId, joinAccountId string, sessionId int64) error {
	upc, joinUpc, err := provider.getJoinAccountIdUpc(accountId, joinAccountId, sessionId)
	if err != nil {
		return err
	}
	// 视频轨已经被打开
	if joinUpc.receiveVideoTrackCancel != nil {
		return nil
	}
	if !upc.video {
		return fmt.Errorf("user didn't open video")
	}
	// 添加用户的视频轨
	var videoRtpSender *webrtc.RTPSender
	if videoRtpSender, err = joinUpc.conn.AddTrack(upc.videoTrack); err != nil {
		return err
	}
	// 创建上下文，cancel方法以及视频轨接收函数
	ctx, cancel := context.WithCancel(context.Background())
	joinUpc.receiveVideoTrackCancel = cancel
	go func(ctx context.Context) {
		rtcpBuf := make([]byte, 1500)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if _, _, rtcpErr := videoRtpSender.Read(rtcpBuf); rtcpErr != nil {
					return
				}
			}
		}
	}(ctx)
	return provider.setJoinAccountIdUpc(accountId, joinAccountId, joinUpc)
}

// 关闭会话中某个用户的音频轨
func (provider *SessionWebRTCProvider) CloseUserWebRTCAudioTrack(accountId, joinAccountId string, sessionId int64) (err error) {
	_, joinUpc, err := provider.getJoinAccountIdUpc(accountId, joinAccountId, sessionId)
	if err != nil {
		return err
	}
	// 音频轨已经被关闭
	if joinUpc.receiveAudioTrackCancel == nil {
		return nil
	}
	joinUpc.receiveAudioTrackCancel()
	joinUpc.receiveAudioTrackCancel = nil
	return provider.setJoinAccountIdUpc(accountId, joinAccountId, joinUpc)
}

// 打开会话中某个用户的音频轨
func (provider *SessionWebRTCProvider) OpenUserWebRTCAudioTrack(accountId, joinAccountId string, sessionId int64) error {
	upc, joinUpc, err := provider.getJoinAccountIdUpc(accountId, joinAccountId, sessionId)
	if err != nil {
		return err
	}
	// 音频轨已经被打开
	if joinUpc.receiveAudioTrackCancel != nil {
		return nil
	}
	var audioRtpSender *webrtc.RTPSender
	if audioRtpSender, err = joinUpc.conn.AddTrack(upc.audioTrack); err != nil {
		return err
	}
	// 创建上下文，cancel方法以及视频轨接收函数
	ctx, cancel := context.WithCancel(context.Background())
	joinUpc.receiveVideoTrackCancel = cancel
	go func(ctx context.Context) {
		rtcpBuf := make([]byte, 1500)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if _, _, rtcpErr := audioRtpSender.Read(rtcpBuf); rtcpErr != nil {
					return
				}
			}
		}
	}(ctx)
	return provider.setJoinAccountIdUpc(accountId, joinAccountId, joinUpc)
}

// 关闭会话中某个用户的webrtc（不包括自己）
func (provider *SessionWebRTCProvider) CloseSessionUserWebRTC(accountId, joinAccountId string, sessionId int64) error {
	if accountId == joinAccountId {
		return nil
	}
	upc, err := provider.getAccountIdUpc(accountId, sessionId)
	if err != nil {
		return err
	}
	if upc.video {
		if err = provider.CloseUserWebRTCVideoTrack(accountId, joinAccountId, sessionId); err != nil {
			return err
		}
	}
	if err = provider.CloseUserWebRTCAudioTrack(accountId, joinAccountId, sessionId); err != nil {
		return err
	}
	upc.joinLock.Lock()
	delete(upc.joins, joinAccountId)
	upc.joinLock.Unlock()
	provider.setAccountIdUpc(accountId, upc, sessionId)
	return nil
}

// 离开会话的webrtc
func (provider *SessionWebRTCProvider) LeaveSessionWebRTC(joinAccountId string, sessionId int64) (err error) {
	return
}

func CreateSession(sdp string, id int64, accountId string, video bool) (swapSdp string, err error) {
	log.Logger.Debug(fmt.Sprintf("user %v start create session %v start", accountId, id))
	var sessionID = fmt.Sprintf("%v:%v", id, accountId)
	offer := webrtc.SessionDescription{}
	Decode(sdp, &offer)

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(peerConnectionConfig)
	if err != nil {
		return
	}

	// Allow us to receive 1 video track and 1 audio track
	if video {
		if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
			return
		}
	}
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		return
	}

	videoLocalTrackChan := make(chan *webrtc.TrackLocalStaticRTP)
	audioLocalTrackChan := make(chan *webrtc.TrackLocalStaticRTP)
	// Set a handler for when a new remote track starts, this just distributes all our packets
	// to connected peers
	peerConnection.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		var kind = remoteTrack.Kind().String()
		// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
		// This can be less wasteful by processing incoming RTCP events, then we would emit a NACK/PLI when a viewer requests it
		go func() {
			ticker := time.NewTicker(rtcpPLIInterval)
			for range ticker.C {
				if rtcpSendErr := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(remoteTrack.SSRC())}}); rtcpSendErr != nil {
					log.Logger.Error(rtcpSendErr.Error())
				}
			}
		}()
		log.Logger.Debug("peer connection on track remote track, the kind is " + kind)

		// Create a local track, all our SFU clients will be fed via this track
		localTrack, newTrackErr := webrtc.NewTrackLocalStaticRTP(remoteTrack.Codec().RTPCodecCapability, kind, sessionID)
		if newTrackErr != nil {
			panic(newTrackErr)
		}
		switch kind {
		case webrtc.RTPCodecTypeVideo.String():
			videoLocalTrackChan <- localTrack
		case webrtc.RTPCodecTypeAudio.String():
			audioLocalTrackChan <- localTrack
		default:
			err = webrtc.ErrUnknownType
			panic(err)
		}

		rtpBuf := make([]byte, 1400)
		for {
			i, _, readErr := remoteTrack.Read(rtpBuf)
			if readErr != nil {
				panic(readErr)
			}

			// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
			if _, err = localTrack.Write(rtpBuf[:i]); err != nil && !errors.Is(err, io.ErrClosedPipe) {
				panic(err)
			}
		}
	})

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		return
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		return
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		return
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete
	swapSdp = Encode(*peerConnection.LocalDescription())
	if video {
		go func() {
			localTrack := <-videoLocalTrackChan
			log.Logger.Debug(fmt.Sprintf("sessionID: %v receive video local track ", sessionID))
			videoLocalTrackMap.Store(sessionID, localTrack)
		}()
	}
	go func() {
		localTrack := <-audioLocalTrackChan
		log.Logger.Debug(fmt.Sprintf("sessionID: %v receive audio local track ", sessionID))
		audioLocalTrackMap.Store(sessionID, localTrack)
	}()

	log.Logger.Debug(fmt.Sprintf("user %v create session %v success", accountId, id))
	return
}

func getVideoLocalTrack(sessionID string) (*webrtc.TrackLocalStaticRTP, error) {
	value, exist := videoLocalTrackMap.Load(sessionID)
	if !exist {
		err := errors.New(fmt.Sprintf("video sessionID %v isn't exist", sessionID))
		return nil, err
	}
	localTrack := value.(*webrtc.TrackLocalStaticRTP)
	return localTrack, nil
}

func getAudioLocalTrack(sessionID string) (*webrtc.TrackLocalStaticRTP, error) {
	value, exist := audioLocalTrackMap.Load(sessionID)
	if !exist {
		err := errors.New(fmt.Sprintf("audio sessionID %v isn't exist", sessionID))
		return nil, err
	}
	localTrack := value.(*webrtc.TrackLocalStaticRTP)
	return localTrack, nil
}

func JoinSession(sdp string, id int64, accountId string, video bool) (swapSdp string, err error) {
	log.Logger.Debug(fmt.Sprintf("user %v start join session %v", accountId, id))
	var sessionID = fmt.Sprintf("%v:%v", id, accountId)

	var videoLocalTrack *webrtc.TrackLocalStaticRTP
	if video {
		if videoLocalTrack, err = getVideoLocalTrack(sessionID); err != nil {
			return
		}
	}
	var audioLocalTrack *webrtc.TrackLocalStaticRTP
	if audioLocalTrack, err = getAudioLocalTrack(sessionID); err != nil {
		return
	}
	recvOnlyOffer := webrtc.SessionDescription{}
	Decode(sdp, &recvOnlyOffer)
	var peerConnection *webrtc.PeerConnection
	if peerConnection, err = webrtc.NewPeerConnection(peerConnectionConfig); err != nil {
		return
	}
	if video {
		var videoRtpSender *webrtc.RTPSender
		if videoRtpSender, err = peerConnection.AddTrack(videoLocalTrack); err != nil {
			return
		}
		go func() {
			rtcpBuf := make([]byte, 1500)
			for {
				if _, _, rtcpErr := videoRtpSender.Read(rtcpBuf); rtcpErr != nil {
					return
				}
			}
		}()
	}
	var audioRtpSender *webrtc.RTPSender
	if audioRtpSender, err = peerConnection.AddTrack(audioLocalTrack); err != nil {
		return
	}
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := audioRtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()
	if err = peerConnection.SetRemoteDescription(recvOnlyOffer); err != nil {
		return
	}
	var answer webrtc.SessionDescription
	if answer, err = peerConnection.CreateAnswer(nil); err != nil {
		return
	}
	var gatherComplete <-chan struct{}
	gatherComplete = webrtc.GatheringCompletePromise(peerConnection)
	if err = peerConnection.SetLocalDescription(answer); err != nil {
		return
	}
	<-gatherComplete
	swapSdp = Encode(*peerConnection.LocalDescription())
	log.Logger.Debug(fmt.Sprintf("user %v join session %v success", accountId, id))
	return
}
