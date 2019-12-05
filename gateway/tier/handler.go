package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/frankhang/util/tcp"
)

//tierHandler implements Hanlder
type TierHandler struct {
	*TierPacketIO

	driver *TireDriver

	ctl *Controller
}

func NewTierHandler(tierPacketIO *TierPacketIO, driver *TireDriver) *TierHandler {
	handler := &TierHandler{TierPacketIO: tierPacketIO, driver: driver}
	handler.ctl = &Controller{TierPacketIO: tierPacketIO}
	return handler
}

func (th *TierHandler) Handle(ctx context.Context, cc *tcp.ClientConn, data []byte) error {

	ctl := th.ctl
	ctl.ctx = ctx
	ctl.cc = cc

	fmt.Printf("Handle: 【%s】\n", hex.EncodeToString(data))
	cmd := data[0]
	//dispach cmd process logic to controller
	switch cmd {
	case 0x55:
		return ctl.TirePressureReport(data)
	case 0x57:
		return ctl.TireReplaceAck(data)
	}

	return nil
}
