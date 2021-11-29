package handle

import (
	"baby-fried-rice/internal/pkg/module/live/config"
	"baby-fried-rice/internal/pkg/module/live/log"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
	"github.com/pkg/errors"
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

func CreateSession(sdp string, bizId, accountId string, video bool) (swapSdp string, err error) {
	log.Logger.Debug(fmt.Sprintf("user %v start create session %v start", accountId, bizId))
	var sessionID = fmt.Sprintf("%v:%v", bizId, accountId)
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

	log.Logger.Debug(fmt.Sprintf("user %v create session %v success", accountId, bizId))
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

func JoinSession(sdp string, bizId, accountId string, video bool) (swapSdp string, err error) {
	log.Logger.Debug(fmt.Sprintf("user %v start join session %v", accountId, bizId))
	var sessionID = fmt.Sprintf("%v:%v", bizId, accountId)

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
	log.Logger.Debug(fmt.Sprintf("user %v join session %v success", accountId, bizId))
	return
}
