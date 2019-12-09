package main

import (
	"context"
	"github.com/frankhang/util/logutil"
	"github.com/frankhang/util/tcp"
	"github.com/frankhang/util/util"
)

type Controller struct {
	*TierPacketIO

	ctx context.Context
	cc  *tcp.ClientConn
}

func (c *Controller) TirePressureReport(data []byte) error {
	logutil.Logger(c.ctx).Info("controller")

	var s int

	h := []byte{0x56, 0xAA, 0x00, 0xff, 0xEE, 0xEE}
	s += util.Sum(h)
	if err := c.WritePacket(c.ctx, h); err != nil {
		return err
	}

	dd := c.Alloc.Alloc(3 + 1)
	dd = append(dd, byte(9+1)) //length
	dd = append(dd, byte(0))   //tire number
	dd = append(dd, byte(1))   //user id
	s += util.Sum(dd[:3])
	//check sum
	dd = append(dd, byte(s))

	if err := c.WritePacket(c.ctx, dd); err != nil {
		return err
	}

	return nil

}

func (c *Controller) TireReplaceAck(data []byte) error {
	logutil.Logger(c.ctx).Info("controller")
	return nil

}
