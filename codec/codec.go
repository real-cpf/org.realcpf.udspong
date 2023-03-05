package codec

import (
	"encoding/binary"
	"fmt"

	"github.com/panjf2000/gnet/v2"
)

type PongCodec struct{}

func (codec *PongCodec) Decode(c gnet.Conn) ([]byte, int) {
	var start uint16 = binary.BigEndian.Uint16(Warp(c, 2))
	if start == 0 {
		return nil, 0
	}
	var len uint32 = binary.BigEndian.Uint32(Warp(c, 4))
	data := Warp(c, int(len))

	c.Discard(2)

	return data, int(start)
}

func Warp(c gnet.Conn, len int) []byte {
	data, _ := c.Next(len)
	return data
}

func (codec *PongCodec) Encode(buf []byte) ([]byte, error) {

	i := 0
	total := len(buf)
	for i < total {
		var start uint16 = binary.BigEndian.Uint16(buf[i : i+2])
		i = i + 2
		if 0 == start {
			continue
		}

		bodyLen := binary.BigEndian.Uint32(buf[i : i+4])
		i = i + 4
		sub := buf[i : i+int(bodyLen)]
		fmt.Println(sub)

	}
	return nil, nil
}
