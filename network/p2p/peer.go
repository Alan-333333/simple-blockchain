package p2p

import (
	"io"
	"net"
)

type Peer struct {
	ID   string
	Conn net.Conn

	// 发送队列
	sendQueue chan []byte

	// 关闭标志
	closed chan bool
}

// 创建一个新的Peer
func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		ID:        GeneratePeerID(),
		Conn:      conn,
		sendQueue: make(chan []byte),
		closed:    make(chan bool),
	}
}

// 启动peer,处理消息读取和发送
func (p *Peer) start() {
	go p.ReadLoop()
	go p.WriteLoop()
}

// 读取循环
func (p *Peer) ReadLoop() {
	for {
		// 1. 从conn读取数据
		data := make([]byte, 1024)
		n, err := p.Conn.Read(data)
		if err != nil {
			if err == io.EOF {
				// 对端关闭
				p.Close()
				return
			}
			// 其他错误
			continue
		}

		// 2. 解码消息
		msg := DecodeMessage(data[:n])
		// 3. 处理消息
		HandleMessage(&msg)

		// 检查关闭状态
		if isClosed(p.closed) {
			break
		}
	}
}

// 写入循环
func (p *Peer) WriteLoop() {
	for {
		select {
		// 1. 从sendQueue取数据
		case data := <-p.sendQueue:
			// 2. 发送数据
			p.Conn.Write(data)
		// 检查关闭状态
		case <-p.closed:
			break
		}
	}
}

// 发送数据
func (p *Peer) Send(data []byte) {
	p.sendQueue <- data
}

// 关闭连接
func (p *Peer) Close() {
	p.closed <- true
}

// 检查关闭状态
func isClosed(ch chan bool) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}

func (p *Peer) SendPing() {
	p.sendQueue <- EncodeMessage(MsgTypePing, "ping")
}
