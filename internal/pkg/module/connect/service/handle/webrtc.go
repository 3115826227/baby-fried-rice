package handle

import (
	"baby-fried-rice/internal/pkg/module/connect/log"
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
	localTrackMap        = new(sync.Map)
	peerConnectionConfig = webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}
)

func CreateSession(sdp string, id int64, accountId string) (swapSdp string, err error) {
	log.Logger.Debug(fmt.Sprintf("user %v start create session %v start", accountId, id))
	var sessionID = fmt.Sprintf("%v:%v", id, accountId)
	offer := webrtc.SessionDescription{}
	Decode(sdp, &offer)

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(peerConnectionConfig)
	if err != nil {
		return
	}

	// Allow us to receive 1 video track
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		return
	}

	localTrackChan := make(chan *webrtc.TrackLocalStaticRTP)
	// Set a handler for when a new remote track starts, this just distributes all our packets
	// to connected peers
	peerConnection.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
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

		// Create a local track, all our SFU clients will be fed via this track
		localTrack, newTrackErr := webrtc.NewTrackLocalStaticRTP(remoteTrack.Codec().RTPCodecCapability, "video", sessionID)
		if newTrackErr != nil {
			panic(newTrackErr)
		}
		localTrackChan <- localTrack

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
	go func() {
		localTrack := <-localTrackChan
		//data, err1 := json.Marshal(localTrack)
		//if err1 != nil {
		//	log.Logger.Error(err1.Error())
		//	return
		//}
		//if err1 = cache.GetCache().Add(sessionID, string(data)); err1 != nil {
		//	log.Logger.Error(err1.Error())
		//	return
		//}
		log.Logger.Debug(fmt.Sprintf("sessionID: %v receive local track ", sessionID))
		localTrackMap.Store(sessionID, localTrack)
	}()
	log.Logger.Debug(fmt.Sprintf("user %v create session %v success", accountId, id))
	return
}

func JoinSession(sdp string, id int64, accountId string) (swapSdp string, err error) {
	log.Logger.Debug(fmt.Sprintf("user %v start join session %v", accountId, id))
	var sessionID = fmt.Sprintf("%v:%v", id, accountId)
	//val, err := cache.GetCache().Get(sessionID)
	//if err != nil {
	//	log.Logger.Error(err.Error())
	//	return
	//}
	//var localTrack *webrtc.TrackLocalStaticRTP
	//if err = json.Unmarshal([]byte(val), &localTrack); err != nil {
	//	log.Logger.Error(err.Error())
	//	return
	//}
	value, exist := localTrackMap.Load(sessionID)
	if !exist {
		err = errors.New(fmt.Sprintf("sessionID %v isn't exist", sessionID))
		return
	}
	localTrack := value.(*webrtc.TrackLocalStaticRTP)

	recvOnlyOffer := webrtc.SessionDescription{}
	Decode(sdp, &recvOnlyOffer)
	var peerConnection *webrtc.PeerConnection
	if peerConnection, err = webrtc.NewPeerConnection(peerConnectionConfig); err != nil {
		return
	}
	var rtpSender *webrtc.RTPSender
	if rtpSender, err = peerConnection.AddTrack(localTrack); err != nil {
		return
	}
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
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
