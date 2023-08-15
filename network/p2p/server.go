package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"time"

	"github.com/Alan-333333/simple-blockchain/transaction"
)

const PingInterval = 1 * time.Second

type Server struct {
	port int

	peers map[string]*Peer

	messages chan *Message

	onTx func(transaction.Transaction)
}

func NewServer(port int) *Server {
	return &Server{
		port:     port,
		peers:    make(map[string]*Peer),
		messages: make(chan *Message),
	}
}

// 启动服务器
func (s *Server) Start() {
	// 1. 监听端口
	listener, _ := net.Listen("tcp", fmt.Sprintf(":%d", s.port))

	// 2. 接收连接请求
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		peer := NewPeer(conn)
		fmt.Println(s.peers)
		go s.onConnect(peer)
	}
}

// 处理连接
func (s *Server) onConnect(peer *Peer) {

	// 1. 建立连接
	// 已经由net.Conn建立

	// 2. 交换版本
	// myVersion := Version{ /*...*/ }
	// peer.Send(myVersion)
	// peerVersion := peer.ReceiveVersion()

	// 3. 添加到peers
	s.peers[peer.ID] = peer

	// 4. 启动peer
	peer.start()

	// 5. 启动ping检查
	go func() {
		lastPing := time.Now()
		for {
			if time.Now().Sub(lastPing) > PingInterval {
				peer.SendPing()
				lastPing = time.Now()
			}
		}
	}()

}

func (s *Server) SetOnTx(callback func(transaction.Transaction)) {
	s.onTx = callback
}

// 广播消息到所有peers
func (s *Server) Broadcast(msgType int, msg *Message) {
	for _, peer := range s.peers {
		peer.Send(EncodeMessage(msgType, msg))
	}
}

// 处理接收到的消息
func HandleMessage(msg *Message) {
	switch msg.MsgType {
	case MsgTypeVersion:
		// 处理版本
		var payload Version
		decoder := gob.NewDecoder(bytes.NewBuffer(msg.Data))
		decoder.Decode(&payload)
	case MsgTypeTx:
		// 处理Tx
	case MsgTypeBlock:
	// 处理block
	case MsgTypePing:
		fmt.Println(msg.Data)
	}
}

// 广播时调用
func (s *Server) BroadcastTx(tx transaction.Transaction) {
	if s.onTx != nil {
		s.onTx(tx)
	}
}
