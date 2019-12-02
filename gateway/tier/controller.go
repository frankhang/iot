package main

import (
	"context"
	"fmt"
	"github.com/frankhang/util/logutil"
	"github.com/frankhang/util/tcp"
	"go.uber.org/zap"
)

type Controller struct {
	*TierPacketIO

	ctx context.Context
	cc  *tcp.ClientConn
}

func (c *Controller) TirePressureReport(data []byte) error {
	logutil.Logger(c.ctx).Info("TirePressureReport: [%x].", zap.ByteString("packetData", data))
	fmt.Printf("TirePressureReport: [%x].", data)

	c.WritePacket([]byte{0x56, 0xAA, 0x00, 0xff, 0xEE, 0xEE})

	dd := c.Alloc.Alloc(3)
	dd = append(dd, byte(9)) //length
	dd = append(dd, byte(0)) //tire number
	dd = append(dd, byte(1)) //user id

	c.WritePacket(dd)

	return nil

}

func (c *Controller) TireReplaceAck(data []byte) error {
	logutil.Logger(c.ctx).Info("TireReplaceAck: [%x].", zap.ByteString("packetData", data))
	fmt.Printf("TireReplaceAck: [%x].", data)

	return nil

}
