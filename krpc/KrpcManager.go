package krpc

import (
	"github.com/chenminjian/spider/config"
	"fmt"
	"github.com/chenminjian/spider/utils"
	"github.com/chenminjian/spider/dht"
	"net"
	"time"
	"github.com/chenminjian/spider/bencoded"
	"github.com/chenminjian/spider/container"
	"errors"
	"encoding/json"
)

var Mgr *KrpcManager

const (
	pingType         = "ping"
	findNodeType     = "find_node"
	getPeersType     = "get_peers"
	announcePeerType = "announce_peer"
)

type KrpcManager struct {
	transIdIdx uint16               // transaction id index
	maxTransId uint16               // max transaction id
	queryChan  chan *config.Query   // query request chan
	respChan   chan *config.Resp    // query resp chan
	dht        *dht.Dht             // dht
	transMap   *container.SyncedMap // transaction map
	Metric     *config.Metric       // metric
}

func Init() (*KrpcManager) {
	d := dht.NewDht()

	m := &KrpcManager{
		transIdIdx: 1,
		maxTransId: 40000,
		queryChan:  make(chan *config.Query, 32),
		respChan:   make(chan *config.Resp, 32),
		dht:        d,
		transMap:   container.NewSyncedMap(),
		Metric:     config.NewMetric(),
	}
	d.SetHandler(m)

	Mgr = m
	return m
}

func (k *KrpcManager) Run() {
	fmt.Println("Run")

	go k.dht.Run()

	ticker := time.Tick(time.Second * 30)

	go func() {
		for {
			select {
			case q := <-k.queryChan:
				go k.query(q)
			case r := <-k.respChan:
				go k.reply(r)
			case <-ticker:
				// go k.fullTableFindNode()
			}
		}
	}()
}

func (k *KrpcManager) fullTableFindNode() {
	fmt.Printf("func:fullTableFindNode,table size:%d\r\n", k.dht.RoutingTable.Len())
	for item := range k.dht.RoutingTable.Iter() {
		node := item.Val.(*config.Node)
		k.FindNode(node, node.Id[:16]+k.dht.Node.Id[16:], utils.RandomString(20))
	}
}

func (k *KrpcManager) query(query *config.Query) {
	// fmt.Println("func:query")

	transId := query.Data["t"].(string)
	trans := config.NewTransaction(query)
	k.transMap.Set(transId, trans)

	err := k.send(query.Node.Addr, query.Data)
	if err != nil {
		fmt.Println("send query error!")
		return
	}

	success := false
	select {
	case <-trans.Resp:
		// fmt.Println("func:query,receive resp")
		success = true
	case <-time.After(time.Second * 30):
		// fmt.Println("over time!")
	}

	if !success {
		// k.dht.BlackList.Insert(query.Node.Addr.IP.String(), query.Node.Addr.Port)
		// k.dht.RoutingTable.Delete(query.Node.Addr)
	}

}

func (k *KrpcManager) reply(resp *config.Resp) {
	fmt.Println("func:reply---")

	k.send(resp.Node.Addr, resp.Data)
}

func (k *KrpcManager) send(addr *net.UDPAddr, data map[string]interface{}) error {
	k.dht.Conn.SetWriteDeadline(time.Now().Add(time.Second * 15))

	_, err := k.dht.Conn.WriteToUDP([]byte(bencoded.Encode(data)), addr)
	if err != nil {
		fmt.Println("WriteToUDP error!")
	}
	// fmt.Println("after WriteToUDP!")
	return err
}

func (k *KrpcManager) getNewTransactionId() string {
	k.transIdIdx = (k.transIdIdx + 1) % k.maxTransId
	return string(utils.Uint16ToBytes(k.transIdIdx))
}

func (k *KrpcManager) FindNode(node *config.Node, id, target string) {
	fmt.Printf("func:FindNode,addr:%s\r\n", node.Addr.String())

	k.sendQuery(
		node,
		findNodeType,
		map[string]interface{}{
			"id":     id,
			"target": target,
		},
	)
}

func (k *KrpcManager) sendQuery(node *config.Node, queryType string, content map[string]interface{}) {
	transId := k.getNewTransactionId()
	data := k.makeQuery(
		transId,
		queryType,
		content,
	)

	k.queryChan <- &config.Query{
		Node: node,
		Data: data,
	}
}

func (k *KrpcManager) makeQuery(transId string, queryType string, content map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"t": transId,
		"y": "q",
		"q": queryType,
		"a": content,
	}
}

func (k *KrpcManager) sendResp(node *config.Node, transId string, content map[string]interface{}) {
	data := k.makeResp(
		transId,
		content,
	)

	k.respChan <- &config.Resp{
		Node: node,
		Data: data,
	}
}

func (k *KrpcManager) makeResp(transId string, content map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"t": transId,
		"y": "r",
		"r": content,
	}
}

func (k *KrpcManager) Handle(packet *config.Packet) {
	fmt.Println("func:Handle")
	obj, err := bencoded.Decode(packet.Data)
	if err != nil {
		fmt.Println("func:Handle error! ------------")
		// fmt.Printf("content:%s\r\n", string(packet.Data))
		return
	}

	// str, _ := json.Marshal(obj)
	// fmt.Printf("func:Handle,addr:%s\r\n", packet.Addr.String())
	// fmt.Println("receive:" + string(str))

	data := obj.(map[string]interface{})

	if data["y"] == nil {
		fmt.Println("func:Handle error! ------------")
		return
	}
	y := data["y"].(string)

	if y == "r" {
		k.handleResp(packet, data)
	} else if y == "q" {
		k.handleQuery(packet, data)
	} else {
		fmt.Println("handleError")
	}
}

func (k *KrpcManager) Join() {
	for _, addr := range k.dht.EntryNodes {
		raddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			continue
		}

		k.FindNode(
			&config.Node{Addr: raddr},
			k.dht.Node.Id,
			k.dht.Node.Id,
		)
	}
}

func (k *KrpcManager) handleResp(packet *config.Packet, data map[string]interface{}) {
	// obj, _ := bencoded.Decode(packet.Data)
	// data := obj.(map[string]interface{})

	if data["t"] == nil {
		fmt.Println("resp t lacks!")
		return
	}
	transId := data["t"].(string)
	val, ok := k.transMap.Get(transId)
	if !ok {
		return
	}

	trans := val.(*config.Transaction)

	trans.Resp <- 0

	q := trans.Data["q"].(string)

	if data["r"] == nil {
		fmt.Println("handle resp error!")
		return
	}
	r := data["r"].(map[string]interface{})

	id := r["id"].(string)
	node, err := config.NewNode(id, utils.GenAddress(packet.Addr.IP.String(), packet.Addr.Port))
	if err != nil {
		fmt.Println("node id error!")
		return
	}
	k.dht.RoutingTable.Insert(node)

	switch q {
	case findNodeType:
		if r["nodes"] != nil {
			k.keepFinding(r["nodes"].(string))
		}
	case announcePeerType:
	case getPeersType:
	case pingType:
		fmt.Println("impossible")
	}
}

func (k *KrpcManager) handleQuery(packet *config.Packet, data map[string]interface{}) {
	// fmt.Println("func:handleQuery--------")

	str, _ := json.Marshal(data)
	fmt.Println("receive:" + string(str))

	if data["q"] == nil {
		fmt.Println("handle query q lacks!")
		return
	}

	q := data["q"].(string)
	switch q {
	case pingType:
		k.replyPing(data, packet.Addr)
	case findNodeType:
		k.replyFindNode(data, packet.Addr)
	case getPeersType:
		k.replyGetPeers(data, packet.Addr)
	case announcePeerType:
		fmt.Println("receive announce_peer----")
		k.Metric.Receive.AnnouncePeer++
	default:
		fmt.Println("receive other type query!-------")
	}
}

func (k *KrpcManager) replyPing(data map[string]interface{}, addr *net.UDPAddr) error {
	fmt.Println("func:replyPing--------")

	if data["a"] == nil {
		fmt.Println("replyPing error -------")
		return errors.New("replyPing error")
	}
	if data["t"] == nil {
		fmt.Println("replyPing error -------")
		return errors.New("replyPing error")
	}
	a := data["a"].(map[string]interface{})

	if a["id"] == nil {
		fmt.Println("replyPing error -------")
		return errors.New("replyPing error")
	}

	id := a["id"].(string)
	t := data["t"].(string)

	node, err := config.NewNode(id, utils.GenAddress(addr.IP.String(), addr.Port))
	if err != nil {
		fmt.Println("replyPing error -------")
		return errors.New("replyPing error")
	}

	k.Metric.Receive.Ping++

	k.sendResp(node, t, map[string]interface{}{
		"id": id[:16] + k.dht.Node.Id[16:],
	})

	return nil
}

func (k *KrpcManager) replyGetPeers(data map[string]interface{}, addr *net.UDPAddr) {
	fmt.Println("func:replyGetPeers--------")

	if data["a"] == nil {
		fmt.Println("reply getPeers a lacks!")
		return
	}
	a := data["a"].(map[string]interface{})

	if a["id"] == nil {
		fmt.Println("reply getPeers a lacks!")
		return
	}
	if data["t"] == nil {
		fmt.Println("reply getPeers a lacks!")
		return
	}
	if a["info_hash"] == nil {
		fmt.Println("reply getPeers a lacks!")
		return
	}

	k.Metric.Receive.GetPeers++

	id := a["id"].(string)
	t := data["t"].(string)
	infoHash := a["info_hash"].(string)

	node, err := config.NewNode(id, utils.GenAddress(addr.IP.String(), addr.Port))
	if err != nil {
		fmt.Println("replyGetPeers error -------")
		return
	}
	k.sendResp(node, t, map[string]interface{}{
		"id":    infoHash[:16] + k.dht.Node.Id[16:],
		"token": "spider",
		"nodes": utils.RandomString(26 * 3),
	})
}

func (k *KrpcManager) replyFindNode(data map[string]interface{}, addr *net.UDPAddr) {
	fmt.Println("func:replyFindNode--------")

	if data["a"] == nil {
		fmt.Println("replyFindNode data error -------")
		return
	}
	if data["t"] == nil {
		fmt.Println("replyFindNode data error -------")
		return
	}
	a := data["a"].(map[string]interface{})

	if a["id"] == nil {
		fmt.Println("replyFindNode data error -------")
		return
	}
	id := a["id"].(string)
	t := data["t"].(string)
	// target := a["target"].(string)

	node, err := config.NewNode(id, utils.GenAddress(addr.IP.String(), addr.Port))
	if err != nil {
		fmt.Println("replyFindNode error -------")
		return
	}

	k.Metric.Receive.FindNode++

	k.sendResp(node, t, map[string]interface{}{
		"id":    id[:16] + k.dht.Node.Id[16:],
		"nodes": utils.RandomString(26 * 3),
	})
}

func (k *KrpcManager) keepFinding(nodes string) error {
	// fmt.Println("func:keepFinding")

	if len(nodes)%26 != 0 {
		fmt.Println("nodes' length should be divisible by 26!")
		return errors.New("nodes' length should be divisible by 26!")
	}

	num := len(nodes) / 26
	for i := 0; i < num; i++ {
		nodeStr := string(nodes[i*26 : (i+1)*26])
		node, _ := config.NewNodeFromCompactNodeInfo(nodeStr)

		if k.dht.BlackList.Has(node.Addr.IP.String(), node.Addr.Port) {
			fmt.Printf("ip:%s,port:%d is in black list!\r\n", node.Addr.IP.String(), node.Addr.Port)
			continue
		}
		// k.dht.RoutingTable.Insert(node)
		// fmt.Printf("func:keepFinding,table len:%d\r\n", k.dht.RoutingTable.Len())
		k.FindNode(node, node.Id[:16]+k.dht.Node.Id[16:], node.Id[:12]+utils.RandomString(8))
	}

	/*
	target := utils.RandomString(20)
	neighbors := k.dht.RoutingTable.GetNeighbors(target)
	for _, item := range neighbors {
		k.FindNode(item, k.dht.Node.Id[:16]+utils.RandomString(4), target)
	}
	*/

	return nil
}
