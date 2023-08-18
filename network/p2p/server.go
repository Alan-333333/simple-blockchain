package p2p

import (
	"fmt"
	"net"
	"time"

	blockchain "github.com/Alan-333333/simple-blockchain/block/chain"
	"github.com/Alan-333333/simple-blockchain/transaction"
)

const PingInterval = 1 * time.Second

type Server struct {
	port int

	Peers map[string]*Peer

	onTx func(*transaction.Transaction)

	onBlock func(*blockchain.Block)
}

func NewServer(port int) *Server {
	return &Server{
		port:  port,
		Peers: make(map[string]*Peer),
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
		go s.onConnect(peer)
	}
}

// 处理连接
func (s *Server) onConnect(peer *Peer) {

	// 1. 建立连接
	// 已经由net.Conn建立

	// 2. 交换版本
	peer.SendVersion()

	// 3. 添加到peers
	s.Peers[peer.ID] = peer

	// 4. 启动peer
	peer.start()

	// 5. 启动ping检查
	go peer.TimerPing()

	// 6.监听消息
	go s.readPeerMsg(peer)

	fmt.Println("peers", s.Peers)
}

func (s *Server) SetOnTx(callback func(*transaction.Transaction)) {
	s.onTx = callback
}

func (s *Server) SetOnBlock(callback func(block *blockchain.Block)) {
	s.onBlock = callback
}

// 广播消息到所有peers
func (s *Server) Broadcast(msgType int, data interface{}, readPeer *Peer) {
	for _, peer := range s.Peers {
		// 过滤接收方peer，防止循环广播
		if readPeer != nil && readPeer.ID == peer.ID {
			continue
		}
		peer.Send(EncodeMessage(msgType, data))
	}
}

func (s *Server) readPeerMsg(peer *Peer) {
	for {
		msg := <-peer.msgChan
		s.handleMessage(msg, peer)
	}
}

// 处理接收到的消息
func (s *Server) handleMessage(msg *Message, readPeer *Peer) {

	switch msg.MsgType {
	case MsgTypeVersion:
		// 处理版本
		version := DecodeVersion(msg.Data)
		fmt.Println("version:", version)

	case MsgTypeTx:
		// 处理Tx
		tx, err := DecodeTransaction(msg.Data)
		fmt.Println("tx", tx)
		if err != nil {
			return
		}
		// 校验
		if !tx.IsValid() {
			return
		}
		// 添加到交易池
		txPool := transaction.GetTxPool()
		txPool.AddTx(tx)

		// 广播给其他节点
		s.Broadcast(MsgTypeTx, msg.Data, readPeer)
	case MsgTypeBlock:

		bc := blockchain.GetBlockchain()
		// 解码
		block, err := DecodeBlock(msg.Data)

		fmt.Println("block", block)
		if err != nil {
			return
		}
		// 校验
		if !bc.IsValidBlock(block) {
			return
		}

		// 添加到区块链
		bc.AddBlock(block)

		// 广播给其他节点
		s.Broadcast(MsgTypeBlock, msg.Data, readPeer)

	case MsgTypePing:
		fmt.Println("MsgTypePing:", string(msg.Data))
	}
}

// 广播时调用
func (s *Server) BroadcastTx(tx *transaction.Transaction) {
	if s.onTx != nil {
		s.onTx(tx)
	}
}

func (s *Server) BroadcastBlock(block *blockchain.Block) {
	if s.onBlock != nil {
		s.onBlock(block)
	}
}
