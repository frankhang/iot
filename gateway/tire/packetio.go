package main

import (
	"context"
	"fmt"
	"github.com/frankhang/util/errors"
	"github.com/frankhang/util/hack"
	"github.com/frankhang/util/logutil"
	"github.com/frankhang/util/tcp"
	"go.uber.org/zap"
	"io"
	"strconv"
)

const (
	sizeHead = 27
	locSize  = 18
)

//tierPacketIO implements PacketReader and PacketWriter
type TierPacketIO struct {
	*tcp.PacketIO

	driver *TireDriver
}

func NewTierPacketIO(packetIO *tcp.PacketIO, driver *TireDriver) *TierPacketIO {
	return &TierPacketIO{
		PacketIO: packetIO,
		driver:   driver,
	}
}

func (p *TierPacketIO) ReadPacket(ctx context.Context) (header []byte, data []byte, err error) {

	var size int

	header = p.Alloc.AllocWithLen(sizeHead, sizeHead)
	p.SetReadTimeout()
	if _, err = io.ReadFull(p.BufReadConn, header[:]); err != nil {
		return nil, nil, errors.Trace(err)
	}

	//locStr := strings.TrimSpace(hack.String(header[locSize : locSize+3]))
	locStr := hack.String(header[locSize : locSize+2])
	if size, err = strconv.Atoi(locStr); err != nil {
		return nil, nil, errors.Trace(err)
	}

	logutil.Logger(ctx).Debug("ReadPacket",
		zap.Int("sizeInHeader", size),
		zap.String("header", fmt.Sprintf("%x", header)),
	)

	data = p.Alloc.AllocWithLen(size-sizeHead, size-sizeHead)

	p.SetReadTimeout()
	if _, err = io.ReadFull(p.BufReadConn, data); err != nil {
		return nil, nil, errors.Trace(err)
	}

	return
}
