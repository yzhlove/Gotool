package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func main() {

	var port = 5678
	go server(port)
	time.Sleep(time.Second * 1)
	client(port)
}

type PacketType byte

const (
	DataPacket PacketType = iota // 数据包
	AckPacket                    // ACK 回复包
)

type PacketHeader struct {
	SeqNum uint32 // 序列号
	Type   PacketType
}

type Packet struct {
	Header PacketHeader
	Data   []byte
}

func (p *Packet) Marshal() []byte {
	buf := make([]byte, 5+len(p.Data))
	binary.BigEndian.PutUint32(buf[:4], p.Header.SeqNum)
	buf[4] = byte(p.Header.Type)
	copy(buf[5:], p.Data)
	return buf
}

func (p *Packet) Unmarshal(buf []byte) error {
	if len(buf) < 5 {
		return fmt.Errorf("invalid packet")
	}
	p.Header.SeqNum = binary.BigEndian.Uint32(buf[:4])
	p.Header.Type = PacketType(buf[4])
	p.Data = buf[5:]
	return nil
}

const (
	maxSeqNum = 1024                   // 最大序列号
	rto       = 200 * time.Millisecond // 重传时间
)

var (
	errSessClosed = errors.New("session closed! ")
)

type SendSession struct {
	sync.Mutex
	conn        *net.UDPConn
	remoteAddr  *net.UDPAddr
	nextSeqNum  uint32
	pendingAcks map[uint32]*pendingPacket
	closed      bool
}

type pendingPacket struct {
	sendTime time.Time
	data     []byte
}

func (s *SendSession) Send(data []byte) error {
	s.Lock()
	defer s.Unlock()
	if s.closed {
		return errSessClosed
	}

	seq := s.nextSeqNum
	s.nextSeqNum = (s.nextSeqNum + 1) % maxSeqNum
	p := &Packet{
		Header: PacketHeader{
			SeqNum: seq,
			Type:   DataPacket,
		},
		Data: data,
	}

	msg := p.Marshal()
	if _, err := s.conn.WriteToUDP(msg, s.remoteAddr); err != nil {
		return err
	}

	s.pendingAcks[seq] = &pendingPacket{
		sendTime: time.Now(),
		data:     data,
	}
	log.Println("send udp data: ", string(data))
	return nil
}

const MtuSize = 1500

func (s *SendSession) ackReceiver() {

	buf := make([]byte, MtuSize)

	for {

		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			if s.closed {
				return
			}
			log.Println("receive udp data failed! ", addr.String(), err)
			continue
		}

		if addr.String() != s.remoteAddr.String() {
			log.Println("receive udp data from unknown addr: ", addr.String())
			continue
		}

		p := &Packet{}
		if err = p.Unmarshal(buf[:n]); err != nil {
			log.Println("unmarshal udp data failed! ", err)
			continue
		}

		if p.Header.Type == AckPacket {
			s.Lock()
			seq := p.Header.SeqNum
			if _, ok := s.pendingAcks[seq]; ok {
				log.Println("receive ack: ", seq)
				delete(s.pendingAcks, seq)
			} else {
				log.Println("receive ack: ", seq, "but not found")
			}
			s.Unlock()
		}
	}
}

func (s *SendSession) retransmitter() {
	tk := time.NewTicker(rto >> 1)
	defer tk.Stop()
	for range tk.C {
		if s.closed {
			return
		}

		s.Lock()

		now := time.Now()
		for seq, pp := range s.pendingAcks {
			if now.Sub(pp.sendTime) > rto {
				log.Println("retransmit: ", seq)
				if _, err := s.conn.WriteToUDP(pp.data, s.remoteAddr); err != nil {
					log.Println("retransmit failed! ", err)
				}
				pp.sendTime = now
			}
		}
		s.Unlock()
	}
}

type RecvSession struct {
	sync.Mutex
	conn        *net.UDPConn
	expectedSeq uint32
}

func (r *RecvSession) Listen(handle func([]byte)) error {

	buf := make([]byte, MtuSize)
	for {
		n, addr, err := r.conn.ReadFromUDP(buf)
		if err != nil {
			return err
		}

		p := &Packet{}
		if err = p.Unmarshal(buf[:n]); err != nil {
			log.Println("unmarshal udp data failed! ", err)
			continue
		}

		if p.Header.Type == DataPacket {
			go r.handleDataPacket(p, addr, handle)
		}
	}
}

func (r *RecvSession) handleDataPacket(p *Packet, addr *net.UDPAddr, handle func([]byte)) {

	r.Lock()
	defer r.Unlock()

	seq := p.Header.SeqNum
	log.Println("receive udp data: ", seq)

	pkt := &Packet{Header: PacketHeader{SeqNum: seq, Type: AckPacket}}
	bytes := pkt.Marshal()
	if _, err := r.conn.WriteToUDP(bytes, addr); err != nil {
		log.Println("send udp data failed! ", err)
	}

	if seq == r.expectedSeq {
		handle(p.Data)
		r.expectedSeq = (r.expectedSeq + 1) % maxSeqNum
	} else {
		log.Println("receive udp data out of order: ", seq, "expected: ", r.expectedSeq)
	}
}

func server(port int) {

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Println("listen udp on ", port, " ...")

	recv := &RecvSession{conn: l}
	if err = recv.Listen(func(data []byte) {
		log.Println("receive udp data: ", string(data))
	}); err != nil {
		log.Fatal(err)
	}
}

func client(port int) {

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	send := &SendSession{
		conn:        conn,
		remoteAddr:  addr,
		nextSeqNum:  0,
		pendingAcks: make(map[uint32]*pendingPacket),
	}

	for i := 0; i < 10; i++ {
		if err := send.Send([]byte(fmt.Sprintf("hello udp#%d", i+1))); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Millisecond * 500)
	}

	time.Sleep(time.Second * 6)
}
