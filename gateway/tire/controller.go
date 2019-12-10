package main

import (
	"context"
	"fmt"
	"github.com/frankhang/util/hack"
	"github.com/frankhang/util/logutil"
	"github.com/frankhang/util/tcp"
	"github.com/frankhang/util/util"
)

type Controller struct {
	*TierPacketIO

	ctx context.Context
	cc  *tcp.ClientConn
}

func (c *Controller) TirePressureReport(header []byte, data []byte) error {
	logutil.Logger(c.ctx).Info("controller")

	h := []byte{0x35, 0x36, 0x20, 0x41, 0x41, 0x20, 0x39, 0x30, 0x20, 0x38, 0x43, 0x20, 0x30, 0x33, 0x20, 0x30, 0x34, 0x20}

	size := len(h) + 3*4
	dd := c.Alloc.Alloc(size)
	dd = append(dd, h...)
	dd = append(dd, hack.Slice(fmt.Sprintf("%-3d", size))...) //size
	dd = append(dd, hack.Slice(fmt.Sprintf("%02d ", 0))...)   //tire number
	dd = append(dd, 0x30, 0x31, 0x20)                         //user id

	//check sum
	s := util.Sum(dd)
	dd = append(dd, hack.Slice(fmt.Sprintf("%3d", s))...)

	if err := c.WritePacket(c.ctx, dd); err != nil {
		return err
	}

	return nil

}

func (c *Controller) TireReplaceAck(header []byte, data []byte) error {
	logutil.Logger(c.ctx).Info("controller")
	return nil

}
