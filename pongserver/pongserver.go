package pongserver

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"github.com/panjf2000/gnet/v2"
	"reacpf.org/udspong/codec"
)

type PongServer struct {
	gnet.BuiltinEventEngine
	eng          gnet.Engine
	addr         string
	multicore    bool
	connected    int32
	disconnected int32
	// channelMap   map[string]gnet.Conn
	channelMap sync.Map
	db         sync.Map
}

const (
	GET_NUM         = 1
	COMMAND_NUM     = 2
	BYTE_VALUE_NUM  = 3
	ROUTE_VALUE_NUM = 4
	REG_NUM         = 5
)

func (s *PongServer) OnBoot(e gnet.Engine) (action gnet.Action) {
	s.eng = e
	// s.channelMap = make(map[string]gnet.Conn)
	s.channelMap = sync.Map{}
	s.db = sync.Map{}
	log.Printf("server start on %s", s.addr)
	return
}

func (s *PongServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	c.SetContext(new(codec.PongCodec))
	atomic.AddInt32(&s.connected, 1)
	return
}

func (s *PongServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		log.Printf("error when close conn %s,%v \n", c.RemoteAddr().String(), err)
	}
	atomic.AddInt32(&s.disconnected, 1)
	atomic.AddInt32(&s.connected, -1)
	return
}

var state [2][]byte
var inx = 0

func (s *PongServer) OnTraffic(c gnet.Conn) gnet.Action {

	codec := c.Context().(*codec.PongCodec)
	for c.InboundBuffered() > 1 {
		data, num := codec.Decode(c)
		switch num {
		case GET_NUM:
			v, _ := s.db.Load(string(data))
			fmt.Printf("set to store %s and return %s \n", string(data), v)
		case COMMAND_NUM:
			fmt.Printf("do command %s \n", string(data))
		case BYTE_VALUE_NUM:
			//
			state[inx] = data
			inx = inx + 1
			if inx == 2 {
				inx = 0
				s.db.Store(string(state[0]), string(state[1]))
			}
			fmt.Println(string(data))
		case ROUTE_VALUE_NUM:
			key, value := ParseRouteValues(data)
			if key != "" {
				// kc := s.channelMap[key]
				kc, _ := s.channelMap.Load(key)
				log.Printf("route msg to %s\n", key)
				kc.(gnet.Conn).AsyncWrite(value, func(c gnet.Conn, err error) error {
					return nil
				})
			}
		case REG_NUM:
			key := string(data)
			// s.channelMap[key] = c
			s.channelMap.Store(key, c)
			log.Printf("reg %s\n", key)
		default:
			fmt.Println(num)
		}

	}

	c.Write([]byte("done"))
	return gnet.None
}

func ParseRouteValues(buf []byte) (string, []byte) {
	i := bytes.Index(buf, []byte(" "))
	if i == -1 {
		return "", nil
	}

	return string(buf[0:i]), buf[i+1:]
}
