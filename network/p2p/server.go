package p2p

import (
	"fmt"
	"net"
	"sync"
	"time"

	blockchain "github.com/Alan-333333/simple-blockchain/block/chain"
	"github.com/Alan-333333/simple-blockchain/transaction"
	"github.com/Alan-333333/simple-blockchain/wallet"
)

const PingInterval = 1 * time.Second

type Server struct {
	port int

	Peers    map[string]*Peer
	peerLock sync.Mutex
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
	s.AddPeer(peer)

	// 4. 启动peer
	peer.start()

	// 5. 启动ping检查
	go peer.TimerPing()

	// 6.监听消息
	go s.readPeerMsg(peer)

}

func (s *Server) AddPeer(peer *Peer) {

	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.Peers[peer.ID] = peer
}

func (s *Server) GetPeers() map[string]*Peer {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	peers := make(map[string]*Peer, len(s.Peers))
	for k, v := range s.Peers {
		peers[k] = v
	}
	return peers
}

// 广播消息到所有peers
func (s *Server) Broadcast(msgType int, data interface{}, readPeer *Peer) {
	peers := s.GetPeers()
	for _, peer := range peers {
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
		// 增加版本校验逻辑
		version := DecodeVersion(msg.Data)
		if version.Version != VERSION {
			delete(s.Peers, readPeer.ID)
		}
		return
	case MsgTypeTx:
		// 处理Tx
		tx, err := DecodeTransaction(msg.Data)
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

		if err != nil {
			fmt.Println(err)
			return
		}
		// 校验
		if !bc.IsValidBlock(block) {
			return
		}

		// 添加到区块链
		bc.AddBlock(block)

		bc.Save()

		// 广播给其他节点
		s.Broadcast(MsgTypeBlock, msg.Data, readPeer)
	case MsgTypeWallet:
		wallet, err := wallet.DecodeWallet(msg.Data)

		if err != nil {
			fmt.Println(err)
			return
		}

		wallet.Save()

		s.Broadcast(MsgTypeWallet, msg.Data, readPeer)

	case MsgTypePing:
		return
	}
}
