package distributed

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/memberlist"
	"os"
	"strings"
	"sync"
)

var members = flag.String("members", "", "comma seperated list of members")

var localMemberName string

// 集群
var Cluster Members

// 客户端与节点的关系
var ClientUidNode sync.Map

var broadcasts *memberlist.TransmitLimitedQueue

func init() {

	Cluster = Members{
		M:     make(map[string]*Member),
		Mutex: sync.RWMutex{},
	}
}

/**
M localMemberName => *Member
*/
type Members struct {
	M     map[string]*Member
	Mutex sync.RWMutex
}

type Member struct {
	*memberlist.Node
	clientUid sync.Map
	mutex     sync.Mutex
}

func NewMember(n *memberlist.Node) *Member {
	m := new(Member)
	m.Node = n
	return m
}

func MemberJoin(n *memberlist.Node) {
	Cluster.Mutex.Lock()
	member := NewMember(n)
	Cluster.M[n.Name] = member
	Cluster.Mutex.Unlock()

}

func MemberLeave(n *memberlist.Node) {
	Cluster.Mutex.Lock()
	delete(Cluster.M, n.Name)
	Cluster.Mutex.Unlock()
}

func BootCluster() {

	hostname, _ := os.Hostname()
	c := memberlist.DefaultLocalConfig()
	c.Events = &eventDelegate{}
	c.Delegate = &delegate{}
	c.BindPort = 0
	uid, _ := uuid.NewUUID()
	c.Name = hostname + "-" + uid.String()
	localMemberName = c.Name
	m, err := memberlist.Create(c)
	if err != nil {
		panic(err)
	}

	if len(*members) > 0 {
		parts := strings.Split(*members, ",")
		_, err := m.Join(parts)
		if err != nil {
			return
		}
	}
	broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return m.NumMembers()
		},
		RetransmitMult: 3,
	}

	node := m.LocalNode()
	fmt.Printf("Local member %s:%d\n", node.Addr, node.Port)
}

func BroadcastSub() {

	fmt.Println("client connect broadcast")

	broadcasts.QueueBroadcast(&broadcast{
		msg:    append([]byte("sub sub")),
		notify: nil,
	})

}

func GetLocalMemberName() string {

	return localMemberName
}
