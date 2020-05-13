package distributed

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/memberlist"
)

type delegate struct {
}

type eventDelegate struct {
}

type broadcast struct {
	msg    []byte
	notify chan<- struct{}
}

func (b *broadcast) Invalidates(other memberlist.Broadcast) bool {
	return false
}

func (b *broadcast) Message() []byte {
	return b.msg
}

func (b *broadcast) Finished() {
	if b.notify != nil {
		close(b.notify)
	}
}

func (d delegate) NodeMeta(limit int) []byte {

	return nil
}

func (d delegate) NotifyMsg(msg []byte) {
	fmt.Println("NotifyMsg", msg)

	update := newUpdate()
	if err := json.Unmarshal(msg, update); err != nil {

		fmt.Println(update)
		handle(update)
	}
	return
}

func (d delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return broadcasts.GetBroadcasts(overhead, limit)
}

func (d delegate) LocalState(join bool) []byte {

	return nil
}

func (d delegate) MergeRemoteState(buf []byte, join bool) {

	fmt.Println("MergeRemoteState", buf)
	return
}

func (e eventDelegate) NotifyJoin(n *memberlist.Node) {
	//fixme

	fmt.Println("somebody join", n.Name)
	MemberJoin(n)
}

func (e eventDelegate) NotifyLeave(n *memberlist.Node) {
	MemberLeave(n)
}

func (e eventDelegate) NotifyUpdate(*memberlist.Node) {
	fmt.Println("Notify Update")
}
