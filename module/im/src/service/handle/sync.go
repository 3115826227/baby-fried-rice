package handle

import "github.com/3115826227/baby-fried-rice/module/im/src/service/model"

/*
	消息同步监听器
*/

type MessageSyncListener struct {
	MessageSyncChan chan model.ChatMessageReceive
	stopNotify      chan bool
}

var (
	listener *MessageSyncListener
)

func init() {
	listener = newMessageSyncListener()
	go func() {
		listener.start()
	}()
}

func newMessageSyncListener() *MessageSyncListener {
	return &MessageSyncListener{
		MessageSyncChan: make(chan model.ChatMessageReceive, 20000),
		stopNotify:      make(chan bool, 1),
	}
}

func SyncMessage(message model.ChatMessageReceive) {
	listener.MessageSyncChan <- message
}

func StopListener() {
	listener.stop()
}

func (listener *MessageSyncListener) start() {
	for {
		select {
		case <-listener.stopNotify:
			//todo chan中的数据需要处理
			return
		case messageReceive := <-listener.MessageSyncChan:
			if _, exist := ConnectionMap[messageReceive.Receive]; exist {
				ConnectionMap[messageReceive.Receive].WriteJSON(&messageReceive)
			}
		}
	}
}

func (listener *MessageSyncListener) stop() {
	listener.stopNotify <- true
}
