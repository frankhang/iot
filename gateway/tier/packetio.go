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
		driver: driver,
	}
}





func (p *TierPacketIO) ReadPacket() ([]byte, error) {
	var header [8]byte

	waitTimeout := time.Duration(p.driver.cfg.ReadTimeout)
	if waitTimeout > 0 {
		if err := p.BufReadConn.SetReadDeadline(time.Now().Add(waitTimeout)); err != nil {
			return nil, err
		}
	}
	if _, err := io.ReadFull(p.BufReadConn, header[:]); err != nil {
		return nil, errors.Trace(err)
	}

	length := int(header[5])

	data := make([]byte, length - 8)
	if waitTimeout > 0 {
		if err := p.BufReadConn.SetReadDeadline(time.Now().Add(waitTimeout)); err != nil {
			return nil, err
		}
	}
	if _, err := io.ReadFull(p.BufReadConn, data); err != nil {
		return nil, errors.Trace(err)
	}
	return data, nil

}

func (p *TierPacketIO) WritePacket(data []byte) error {
	return nil
}

