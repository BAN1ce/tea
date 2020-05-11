package distributed

import (
	"fmt"
	"github.com/hashicorp/memberlist"
)

type delegate struct {
}

type eventDelegate struct {
}

func (d delegate) NodeMeta(limit int) []byte {

	return nil
}

func (d delegate) NotifyMsg([]byte) {
	return
}

func (d delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return nil
}

func (d delegate) LocalState(join bool) []byte {

	return  nil
}

func (d delegate) MergeRemoteState(buf []byte, join bool) {
	return
}

func (e eventDelegate) NotifyJoin(n *memberlist.Node) {
	//fixme
	MemberJoin(n)
}

func (e eventDelegate) NotifyLeave(n *memberlist.Node) {
	MemberLeave(n)
}

func (e eventDelegate) NotifyUpdate(*memberlist.Node) {
	fmt.Println("Notify Update")
}
