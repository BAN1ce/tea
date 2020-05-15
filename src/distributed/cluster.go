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

var localNodeName string

var localNode *Node

// 集群
var Cluster Nodes

// 客户端与节点的关系
var ClientUidNode sync.Map

var broadcasts *memberlist.TransmitLimitedQueue

func init() {

	Cluster = Nodes{
		M:     make(map[string]*Node),
		Mutex: sync.RWMutex{},
	}
}

/**
M localNodeName => *Node
*/
type Nodes struct {
	M     map[string]*Node
	Mutex sync.RWMutex
}

type Node struct {
	*memberlist.Node
	clientUid sync.Map
	mutex     sync.Mutex
}

func NewNode(n *memberlist.Node) *Node {
	m := new(Node)
	m.Node = n
	return m
}

func MemberJoin(n *memberlist.Node) {
	Cluster.Mutex.Lock()
	node := NewNode(n)
	Cluster.M[n.Name] = node
	Cluster.Mutex.Unlock()

}

func nodeLeave(n *memberlist.Node) {
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
	localNodeName = c.Name
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
	localNode = NewNode(node)
	fmt.Printf("Local member %s:%d\n", node.Addr, node.Port)
}

func BroadcastSub() {

	fmt.Println("client connect broadcast")

	broadcasts.QueueBroadcast(&broadcast{
		msg:    append([]byte("sub sub")),
		notify: nil,
	})

}

func GetLocalNodeName() string {

	return localNodeName
}
