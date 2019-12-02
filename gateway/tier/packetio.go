package main

import (
	"github.com/frankhang/util/errors"
	"github.com/frankhang/util/tcp"
	"io"
	"time"
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

	waitTimeout := time.Duration(p.driver.cfg.ReadTimeout)
	if waitTimeout > 0 {
		if err := p.BufReadConn.SetReadDeadline(time.Now().Add(waitTimeout)); err != nil {
			return nil, err
		}
	}
	if _, err := io.ReadFull(p.BufReadConn, header[:]); err != nil {
		return nil, errors.Trace(err)
	}

	length := int(header[6])

	if length == 0 {
		return header[:], nil
	}

	//buf := make([]byte, length-8)
	buf := p.Alloc.AllocWithLen(length-9, length-9)
	if waitTimeout > 0 {
		if err := p.BufReadConn.SetReadDeadline(time.Now().Add(waitTimeout)); err != nil {
			return nil, err
		}
	}
	if _, err := io.ReadFull(p.BufReadConn, buf); err != nil {
		return nil, errors.Trace(err)
	}

	//generate whole packet
	data := p.Alloc.AllocWithLen(length, length)
	copy(data, header[:])
	copy(data[8:], buf)
	return buf, nil

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
