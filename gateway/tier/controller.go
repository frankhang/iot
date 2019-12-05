package main

import (
	"context"
	"fmt"
	"github.com/frankhang/util/logutil"
	"github.com/frankhang/util/tcp"
	"github.com/frankhang/util/util"
	"go.uber.org/zap"
)

type Controller struct {
	*TierPacketIO

	ctx context.Context
	cc  *tcp.ClientConn
}

func (c *Controller) TirePressureReport(data []byte) error {
	//logutil.Logger(c.ctx).Info("TirePressureReport: [%x].", zap.ByteString("packetData", data))
	fmt.Printf("TirePressureReport: [%x].\n", data)
	fmt.Printf("TirePressureReport: sum=[%x]\n", util.Sum(data[:len(data)-1]))

	var s int

	h := []byte{0x56, 0xAA, 0x00, 0xff, 0xEE, 0xEE}
	s += util.Sum(h)
	if err := c.WritePacket(h); err != nil {
		return err
	}

	dd := c.Alloc.Alloc(3 + 1)
	dd = append(dd, byte(9+1)) //length
	dd = append(dd, byte(0))   //tire number
	dd = append(dd, byte(1))   //user id
	s += util.Sum(dd[:3])
	//check sum
	dd = append(dd, byte(s))

	if err := c.WritePacket(dd); err != nil {
		return err
	}

	return nil

}

func (c *Controller) TireReplaceAck(data []byte) error {
	logutil.Logger(c.ctx).Info("TireReplaceAck: [%x].", zap.ByteString("packetData", data))
	fmt.Printf("TireReplaceAck: [%x].\n", data)

	return nil

}
