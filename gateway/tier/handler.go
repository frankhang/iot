package main

import (
	"context"
	"github.com/frankhang/util/logutil"
	"github.com/frankhang/util/tcp"
	"github.com/frankhang/util/util"
	"strconv"
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
	ctl.cc = cc

	//fmt.Printf("Handle: 【%s】\n", hex.EncodeToString(data))


	cmd := data[0]
	sum := util.Sum(data[:len(data)-1])

	//ctx = logutil.WithKeyValue(ctx, "param", fmt.Sprintf("%x", data))
	ctx = logutil.WithKeyValue(ctx, "sum", strconv.Itoa(sum))

	//dispach cmd process logic to controller
	switch cmd {
	case 0x55:
		ctl.ctx = logutil.WithKeyValue(ctx, "method", "TirePressureReport")
		return ctl.TirePressureReport(data)
	case 0x57:
		ctl.ctx = logutil.WithKeyValue(ctx, "method", "TireReplaceAck")
		return ctl.TireReplaceAck(data)
	}

	logutil.Logger(ctx).Warn("no controller method found")
	return nil
}
