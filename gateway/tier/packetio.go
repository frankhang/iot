package main

import (
	"github.com/frankhang/util/errors"
	"github.com/frankhang/util/tcp"
	"io"
)

//tierPacketIO implements PacketReader and PacketWriter
type TierPacketIO struct {
	*tcp.PacketIO

	driver *TireDriver
}

func NewTierPacketIO(packetIO *tcp.PacketIO, driver *TireDriver) *TierPacketIO {
	return &TierPacketIO{
		PacketIO: packetIO,
		driver:   driver,
	}
}

func (p *TierPacketIO) ReadPacket() ([]byte, error) {
	var header [9]byte

	p.SetReadTimeout()
	if _, err := io.ReadFull(p.BufReadConn, header[:]); err != nil {
		//fmt.Printf("ReadFull: %s\n" , errors.Trace(err))
		return nil, errors.Trace(err)
		//return nil, err
	}

	//fmt.Printf("header9: %s\n", hex.EncodeToString(header[:]))
	length := int(header[6])

	//println("length: ", length)


	//buf := make([]byte, length-8)
	buf := p.Alloc.AllocWithLen(length-9, length-9)

	p.SetReadTimeout()
	if _, err := io.ReadFull(p.BufReadConn, buf); err != nil {
		return nil, errors.Trace(err)
	}

	//fmt.Printf("buf: %s\n", hex.EncodeToString(buf[:]))

	//generate whole packet
	data := p.Alloc.AllocWithLen(length, length)
	copy(data, header[:])
	copy(data[9:], buf)

	//fmt.Printf("ReadPacket: [%s]\n", hex.EncodeToString(data))
	return data, nil

}

func (p *TierPacketIO) WritePacket(data []byte) error {

	if n, err := p.Write(data); err != nil {
		errors.Log(errors.Trace(err))
		return errors.Trace(errors.ErrBadConn.GenWithStackByArgs(p.ConnectionID))
	} else if n != len(data) {
		return errors.Trace(errors.ErrBadConn.GenWithStackByArgs(p.ConnectionID))
	} else {
		return nil
	}

}
