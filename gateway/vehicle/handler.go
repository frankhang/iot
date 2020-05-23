package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/frankhang/util/errors"
	"github.com/frankhang/util/logutil"
	"github.com/frankhang/util/tcp"
	"github.com/frankhang/util/util"
	l "github.com/sirupsen/logrus"
	"go.uber.org/zap"

)

//tierHandler implements Hanlder
type Handler struct {
	*PacketIO

	driver *Driver

	ctl *Controller
}

func NewHandler(tPacketIO *PacketIO, driver *Driver) *Handler {
	handler := &Handler{PacketIO: tPacketIO, driver: driver}
	handler.ctl = &Controller{PacketIO: tPacketIO}
	return handler
}

func (th *Handler) Handle(ctx context.Context, cc *tcp.ClientConn, header []byte, data []byte) (err error) {

	ctl := th.ctl
	ctl.cc = cc

	//check crc
	var ss int
	if data[0] == '(' {
		ss = len(data) - 3
	} else {
		ss = len(data) - 2
	}

	crc := util.Crc16(data[:ss])
	expectedCrc := binary.BigEndian.Uint16(data[ss:])


	//ctx = logutil.WithInt(ctx, "length", len(header)+len(data))
	//ctx = logutil.WithString(ctx, "packet", fmt.Sprintf("%x%x", header, data))
	//ctx = logutil.WithInt(ctx, "sum", sum)

	if l.GetLevel() >= l.DebugLevel {
		logutil.Logger(ctx).Debug("Packet received",
			zap.Int("size", len(data)),
			zap.String("packet", fmt.Sprintf("%x%x", header, data)),
			zap.Uint16("crc", crc),
			zap.Uint16("expected", expectedCrc),
			zap.String("packetStr", fmt.Sprintf("%s%s", header, data)),
		)
	}
	if crc != expectedCrc {
		err = fmt.Errorf("Handle: crc check error, %d != %d, data=[%x]", crc, expectedCrc, data)
		err = errors.Trace(err)
		return
	}

	//cmd := hack.String(header[:2])
	////dispach cmd process logic to controller
	//switch cmd {
	//case "55":
	//	ctl.ctx = logutil.WithString(ctx, "method", "TirePressureReport")
	//	err := ctl.TirePressureReport(header, data)
	//	return errors.Trace(err)
	//case "57":
	//	ctl.ctx = logutil.WithString(ctx, "method", "TireReplaceAck")
	//	err := ctl.TireReplaceAck(header, data)
	//	return errors.Trace(err)
	//}
	//
	//logutil.Logger(ctx).Warn("no controller method found")
	return nil
}
