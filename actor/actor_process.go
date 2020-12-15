package actor

import (
	"sync/atomic"

	"github.com/AsynkronIT/protoactor-go/mailbox"
)

// 一个Actor Process可以看作Erlang里面一个独立的逻辑进程
type ActorProcess struct {
	mailbox mailbox.Mailbox
	dead    int32
}

func NewActorProcess(mailbox mailbox.Mailbox) *ActorProcess {
	return &ActorProcess{mailbox: mailbox}
}

func (ref *ActorProcess) SendUserMessage(pid *PID, message interface{}) {
	ref.mailbox.PostUserMessage(message)
}
func (ref *ActorProcess) SendSystemMessage(pid *PID, message interface{}) {
	ref.mailbox.PostSystemMessage(message)
}

func (ref *ActorProcess) Stop(pid *PID) {
	atomic.StoreInt32(&ref.dead, 1)
	ref.SendSystemMessage(pid, stopMessage)
}
