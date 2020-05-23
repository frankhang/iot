package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/frankhang/util/errors"
	"github.com/frankhang/util/tcp"
	"io"
)

const (
	sizeHead = 27
	locSize  = 18
)

//tierPacketIO implements PacketReader and PacketWriter
type PacketIO struct {
	*tcp.PacketIO

	driver *Driver
}

func NewPacketIO(packetIO *tcp.PacketIO, driver *Driver) *PacketIO {
	return &PacketIO{
		PacketIO: packetIO,
		driver:   driver,
	}
}

func (p *PacketIO) ReadPacket(ctx context.Context) (header []byte, data []byte, err error) {

	var b [1]byte
	var h [4]byte
	var dd []byte

	p.SetReadTimeout()

	bufReader := p.BufReadConn.BufReader
	if _, err = io.ReadFull(bufReader, b[:]); err != nil {
		err = errors.Trace(err)
		return
	}

	if b[0] == '(' {

		p.SetReadTimeout()
		if dd, err = bufReader.ReadBytes(')'); err != nil {
			err = errors.Trace(err)
			return
		}

		if len(dd) < 6 {
			err = fmt.Errorf("ReadPacket: Unexcepted len %d, [%x]", len(dd), dd)
			err = errors.Trace(err)
			return
		}

		ll := len(dd) + 1
		data = p.Alloc.AllocWithLen(ll+1, ll+1)
		data[0] = '('
		copy(data[1:], dd)


		return
	}

	p.SetReadTimeout()
	if _, err = io.ReadFull(bufReader, h[:]); err != nil {
		err = errors.Trace(err)
		return
	}

	size := int(binary.BigEndian.Uint16(h[2:]))

	data = p.Alloc.AllocWithLen(5+size+2, 5+size+2)

	data[0] = '('
	copy(data[1:], h[:])
	p.SetReadTimeout()
	if _, err = io.ReadFull(bufReader, data[5:]); err != nil {
		err = errors.Trace(err)
		return
	}

	return
}
