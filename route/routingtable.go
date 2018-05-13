package route

import (
	"github.com/chenminjian/spider/container"
	"github.com/chenminjian/spider/config"
	"net"
)

type Routingtable struct {
	bucketSize  int
	cachedNodes *container.SyncedMap
}

func NewRoutingTable() *Routingtable {
	t := &Routingtable{
		bucketSize:  8,
		cachedNodes: container.NewSyncedMap(),
	}
	return t
}

func (rt *Routingtable) Insert(node *config.Node) {
	rt.cachedNodes.Set(node.Addr.String(), node)
}

func (rt *Routingtable) Delete(addr *net.UDPAddr) {
	rt.cachedNodes.Delete(addr.String())
}

func (rt *Routingtable) GetNeighbors(target string) []*config.Node {
	// fmt.Printf("func:GetNeighbors,node size:%d\r\n", rt.cachedNodes.Len())

	length := rt.cachedNodes.Len()
	if length > 8 {
		length = 8
	}

	var arr1 = make([]*config.Node, length)

	num := 0
	for item := range rt.cachedNodes.Iter() {
		if num >= length {
			continue
		}
		arr1[num] = item.Val.(*config.Node)
		num++
	}

	return arr1
}

func (rt *Routingtable) Iter() <-chan container.MapItem {
	return rt.cachedNodes.Iter()
}

func (rt *Routingtable) Len() int {
	return rt.cachedNodes.Len()
}
