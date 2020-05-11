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
var Cluster Members

func init() {

	Cluster = Members{
		M:     make(map[string]*Member),
		Mutex: sync.RWMutex{},
	}
}

type Members struct {
	M     map[string]*Member
	Mutex sync.RWMutex
}

type Member struct {
	*memberlist.Node
	mutex sync.Mutex
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

	node := m.LocalNode()
	fmt.Printf("Local member %s:%d\n", node.Addr, node.Port)
}
