package main

import (
	"context"
	"github.com/frankhang/util/hack"
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

func (th *TierHandler) Handle(ctx context.Context, cc *tcp.ClientConn, header []byte, data []byte) error {

	ctl := th.ctl
	ctl.cc = cc

	sum := util.Sum(header)
	sum += util.Sum(data[:len(data)-3])

	ctx = logutil.WithKeyValue(ctx, "sum", strconv.Itoa(sum))

	cmd := hack.String(header[:2])
	//dispach cmd process logic to controller
	switch cmd {
	case "55":
		ctl.ctx = logutil.WithKeyValue(ctx, "method", "TirePressureReport")
		return ctl.TirePressureReport(header, data)
	case "57":
		ctl.ctx = logutil.WithKeyValue(ctx, "method", "TireReplaceAck")
		return ctl.TireReplaceAck(header, data)
	}

	logutil.Logger(ctx).Warn("no controller method found")
	return nil
}
