package main

import (
	"context"
	"fmt"
	"github.com/frankhang/util/errors"
	"github.com/frankhang/util/hack"
	"github.com/frankhang/util/logutil"
	"github.com/frankhang/util/tcp"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type Controller struct {
	*TierPacketIO

	ctx context.Context
	cc  *tcp.ClientConn
}

func (c *Controller) TirePressureReport(header []byte, data []byte) error {

	s := hack.String(header[6 : 6+3])
	pressureLimit, err := strconv.Atoi(strings.TrimSpace(s))
	errors.MustNil(errors.Trace(err))
	ctx := logutil.WithInt(c.ctx, "pressureLimit", pressureLimit)

	s = hack.String(header[9 : 9+3])
	temperatureLimit, err := strconv.Atoi(strings.TrimSpace(s))
	errors.MustNil(errors.Trace(err))
	ctx = logutil.WithInt(ctx, "temperatureLimit", temperatureLimit)

	deviceId := fmt.Sprintf("%x", header[12:12+6])
	ctx = logutil.WithString(ctx, "deviceId", deviceId)

	s = hack.String(header[locSize : locSize+3])
	packetSize, err := strconv.Atoi(strings.TrimSpace(s))
	errors.MustNil(errors.Trace(err))
	ctx = logutil.WithInt(ctx, "packetSize", packetSize)

	s = hack.String(header[locSize+3 : locSize+3+3])
	tireNum, err := strconv.Atoi(strings.TrimSpace(s))
	errors.MustNil(errors.Trace(err))
	ctx = logutil.WithInt(ctx, "tireNum", tireNum)

	s = hack.String(header[locSize+6 : locSize+6+3])
	userId, err := strconv.Atoi(strings.TrimSpace(s))
	errors.MustNil(errors.Trace(err))
	ctx = logutil.WithInt(ctx, "userId", userId)

	logutil.Logger(ctx).Info("controller")

	if len(data) >= 3 {
		loc := 0
		for i := 0; i < tireNum; i++ {
			if data[loc] == 0 && data[loc+1] == 0 {
				break
			}

			tireName := "T" + strconv.Itoa(i)
			s = hack.String(data[loc : loc+3])
			tireLoc, err := strconv.Atoi(strings.TrimSpace(s))
			errors.MustNil(errors.Trace(err))
			c := logutil.WithInt(c.ctx, "loc", tireLoc)

			s = hack.String(data[loc+3 : loc+3+3])
			tirePressure, err := strconv.Atoi(strings.TrimSpace(s))
			errors.MustNil(errors.Trace(err))
			c = logutil.WithInt(c, "pressure", tirePressure)

			s = hack.String(data[loc+6 : loc+6+3])
			tireTemperature, err := strconv.Atoi(strings.TrimSpace(s))
			errors.MustNil(errors.Trace(err))
			c = logutil.WithInt(c, "temperature", tireTemperature)

			logutil.Logger(c).Debug(tireName)

			loc += 9
		}

	}

	h := []byte{0x35, 0x36, 0x20, 0x41, 0x41, 0x20, 0x39, 0x30, 0x20, 0x38, 0x43, 0x20, 0x30, 0x33, 0x20, 0x30, 0x34, 0x20}

	size := len(h) + 3*4
	dd := c.Alloc.Alloc(size)
	dd = append(dd, h...)
	dd = append(dd, hack.Slice(fmt.Sprintf("%-3d", size))...) //size
	dd = append(dd, hack.Slice(fmt.Sprintf("%-3d", 0))...)    //tire number
	dd = append(dd, 0x30, 0x31, 0x20)                         //user id

	//sum := util.Sum(dd)
	dd = append(dd, hack.Slice(fmt.Sprintf("%3d", 0))...) //check sum
	//dd = append(dd, hack.Slice(fmt.Sprintf("%3d", sum))...) //check sum

	logutil.Logger(c.ctx).Info("controller ack",
		zap.Int("size", size),
		zap.Int("len", len(dd)),
		zap.String("packet", fmt.Sprintf("%x", dd)))

	if err := c.WritePacket(c.ctx, dd); err != nil {
		return errors.Trace(err)
	}

	return nil

}

func (c *Controller) TireReplaceAck(header []byte, data []byte) error {
	logutil.Logger(c.ctx).Info("controller")
	return nil

}
