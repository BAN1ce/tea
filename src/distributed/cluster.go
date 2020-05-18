package distributed

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/memberlist"
	"os"
	"strings"
)

var members = flag.String("members", "", "comma seperated list of members")

var localNodeName string

var localNode *memberlist.Node

var cluster *memberlist.Memberlist

// 集群

// 客户端与节点的关系

var broadcasts *memberlist.TransmitLimitedQueue

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
		RetransmitMult: 1,
	}

	node := m.LocalNode()
	localNode = node
	fmt.Printf("Local Node %s %s:%d\n", localNodeName, node.Addr, node.Port)
}

/**
获取本地节点名
*/
func GetLocalNodeName() string {
	return localNodeName
}

func GetCluster() *memberlist.Memberlist {

	return cluster
}
