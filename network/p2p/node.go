package p2p

import (
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/Alan-333333/simple-blockchain/transaction"
)

type Node struct {
	ID     string
	IP     string
	Port   int
	Server *Server
}

// 生成随机节点ID
func GenerateNodeID() string {
	// 使用crypto/rand生成随机数作为节点ID
	buf := make([]byte, 32)
	rand.Read(buf)
	return fmt.Sprintf("%x", buf)
}

// 生成随机PeerID
func GeneratePeerID() string {
	// 使用crypto/rand生成随机数作为节点ID
	buf := make([]byte, 8)
	rand.Read(buf)
	return fmt.Sprintf("%x", buf)
}

// 创建新节点
func NewNode(ip string, port int) *Node {
	id := GenerateNodeID()
	server := NewServer(port)
	return &Node{
		ID:     id,
		IP:     ip,
		Port:   port,
		Server: server,
	}
}

// 节点连接到网络
func (node *Node) Listen() {
	// 调用p2p/server的Listen启动监听
	server := NewServer(node.Port)

	server.SetOnTx(node.BroadcastTx)

	server.Start()
}

// 节点连接到peer
func (node *Node) Connect(ip string, port int) {
	// 创建连接
	conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(port))
	if err != nil {
		log.Println("Connecting to peer failed:", err)
		return
	}
	// 使用p2p/peer封装连接
	peer := NewPeer(conn)
	// 处理连接
	node.handleConn(peer)
}

// 将peer添加到节点的连接列表
func (node *Node) handleConn(p *Peer) {
	node.Server.Peers[p.ID] = p
	p.start()
	go p.TimerPing()
}

// 广播交易到网络
func (node *Node) BroadcastTx(tx transaction.Transaction) {

	// 构造消息体
	msg := &Message{
		MsgType: MsgTypeTx,
		Data:    EncodeMessage(MsgTypeTx, tx),
	}
	// 广播消息
	node.Server.Broadcast(MsgTypeTx, msg)
}
