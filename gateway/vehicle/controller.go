package main

import (
	"context"
	"encoding/binary"
	"github.com/frankhang/util/errors"
	"github.com/frankhang/util/tcp"
)

type Controller struct {
	*PacketIO

	ctx context.Context
	cc  *tcp.ClientConn
}

func generateResponse(header []byte, data []byte) []byte {

	rd := make([]byte, len(data))
	copy(rd, data)
	if rd[0] != '(' {
		return nil
	}
	switch rd[3] {
	case 'h':
		rd[3] = 'H'
		return fillCrc(rd)
	case 'c':
		rd[3] = 'C'
		return fillCrc(rd[:len(rd)-2])
	case 'S':
		rd[3] = rd[4]
		return fillCrc(rd[:len(rd)-1])
	case 'r':
		rd[3] = '1'
		return fillCrc(rd)
	default:
		return nil
	}

}

func fillCrc(data []byte) []byte {

	var crcLen int
	if data[0] == '(' {
		crcLen = len(data) - 3
		data[len(data)-1] = ')'
	} else {
		crcLen = len(data) - 2
	}

	//crc16 := util.Crc16(data[:crcLen])
	crc16 := uint16(0)
	binary.BigEndian.PutUint16(data[crcLen:], crc16)

	return data

}

func (c *Controller) Protocol1(header []byte, data []byte) (err error) {

	packet := generateResponse(header, data)

	if len(packet) > 0 {
		if err = c.WritePacket(c.ctx, packet); err != nil {
			err = errors.Trace(err)
			return
		}
	}

	return
}

func (c *Controller) Protocol2(header []byte, data []byte) (err error) {
	packet := generateResponse(header, data)

	if len(packet) > 0 {
		if err = c.WritePacket(c.ctx, packet); err != nil {
			err = errors.Trace(err)
			return
		}
	}

	return

}

func (c *Controller) Protocol3(header []byte, data []byte) (err error) {
	packet := generateResponse(header, data)

	if len(packet) > 0 {
		if err = c.WritePacket(c.ctx, packet); err != nil {
			err = errors.Trace(err)
			return
		}
	}

	return

}

func (c *Controller) Protocol4(header []byte, data []byte) (err error) {
	packet := generateResponse(header, data)

	if len(packet) > 0 {
		if err = c.WritePacket(c.ctx, packet); err != nil {
			err = errors.Trace(err)
			return
		}
	}

	return

}

func (c *Controller) Data1(header []byte, data []byte) (err error) {
	packet := generateResponse(header, data)

	if len(packet) > 0 {
		if err = c.WritePacket(c.ctx, packet); err != nil {
			err = errors.Trace(err)
			return
		}
	}

	return

}
