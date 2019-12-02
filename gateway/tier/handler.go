package main

import (
	"context"
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
	handler.ctl = &Controller{}
	return handler
}

func (th *TierHandler) Handle(ctx context.Context, cc *tcp.ClientConn, data []byte) error {

	ctl := th.ctl
	ctl.ctx = ctx
	ctl.cc = cc

	cmd := data[0]
	//dispach cmd process logic to controller
	switch cmd {
	case 0x55:
		ctl.TirePressureReport(data)
	case 0x57:
		ctl.TireReplaceAck(data)
	}

	return nil
}
