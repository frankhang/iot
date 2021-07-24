package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/frankhang/iot/gateway/vehicle/models/dto"
	"github.com/frankhang/util/errors"
	"github.com/frankhang/util/logutil"
	"github.com/frankhang/util/tcp"
	"github.com/frankhang/util/util"
	l "github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
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

	//check crc & get cmd
	var ss int
	var cmd byte
	var cmdStr string
	if data[0] == '(' {
		ss = len(data) - 3
		cmd = data[3]
	} else {
		ss = len(data) - 2
		cmd = data[0]
	}

	cmdStr = string([]byte{cmd})
	crc := binary.BigEndian.Uint16(data[ss:])
	expectedCrc := util.CrcCcittFfff(data[:ss])

	//ctx = logutil.WithInt(ctx, "length", len(header)+len(data))
	//ctx = logutil.WithString(ctx, "packet", fmt.Sprintf("%x%x", header, data))
	//ctx = logutil.WithInt(ctx, "sum", sum)

	if l.GetLevel() >= l.DebugLevel {
		logutil.Logger(ctx).Debug("Packet received",
			zap.String("cmd", cmdStr),
			zap.Int("size", len(data)),
			zap.String("packet", fmt.Sprintf("%x%x", header, data)),
			zap.Uint16("crc", crc),
			zap.Uint16("expectedCrc", expectedCrc),
			zap.String("packetStr", fmt.Sprintf("%s%s", header, data)),
		)
	}
	if crc != expectedCrc {
		err = fmt.Errorf("Handle: crc check error, %d != %d, data=[%x]", crc, expectedCrc, data)
		err = errors.Trace(err)
		return
	}

	logutil.BgLogger().Warn("listening ...")
	if data[0] == '(' {
		switch cmd {
		case 'h':
			ctl.ctx = logutil.WithString(ctx, "method", "Protocol1")
			err := ctl.Protocol1(header, data)
			return errors.Trace(err)
		case 'c':
			ctl.ctx = logutil.WithString(ctx, "method", "Protocol2")
			err := ctl.Protocol2(header, data)
			return errors.Trace(err)
		case 'S':
			ctl.ctx = logutil.WithString(ctx, "method", "Protocol3")
			err := ctl.Protocol3(header, data)
			return errors.Trace(err)
		case 'r':
			ctl.ctx = logutil.WithString(ctx, "method", "Protocol4")
			err := ctl.Protocol4(header, data)
			return errors.Trace(err)
		default:
			logutil.Logger(ctx).Warn("no controller method found")
		}

	} else {
		ctl.ctx = logutil.WithString(ctx, "method", "Data1")
		err := ctl.Data1(header, data)
		return errors.Trace(err)

	}

	tt := time.Now().UnixNano() / 1e6
	transferUrl := th.driver.cfg.TransferUrl
	pp := string(data)

	transfer := &dto.Transfer{Timestamp: tt, Cmd: cmdStr, Package: pp}

	var js []byte
	if js, err = json.Marshal(transfer); err != nil {
		return errors.Trace(err)
	}
	jsonStr := string(js)

	if l.GetLevel() >= l.DebugLevel {
		logutil.Logger(ctx).Debug("Transfer  to url",
			zap.String("url", transferUrl),
			//zap.Int64("time", tt),
			//zap.String("cmd", cmdStr),
			//zap.String("package", pp),
			zap.String("msjon", jsonStr),
		)
	}
	token := "2017#halouhuandian.com#2021"

	c := &http.Client{
		Timeout: time.Second * 3,
	}
	var resp *http.Response
	if resp, err = c.PostForm(transferUrl, url.Values{"token": {token}, "mjson": {jsonStr}}); err != nil {
		return errors.Trace(err)
	}
	defer resp.Body.Close()

	status := resp.StatusCode

	if status != 200 {
		logutil.Logger(ctx).Error("transfer encounter critical problem",
			zap.String("url", transferUrl),
			zap.String("status", resp.Status))
		return
	}

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)

	if l.GetLevel() >= l.DebugLevel {
		logutil.Logger(ctx).Debug("Response from url",
			zap.String("url", transferUrl),
			zap.String("body", string(body)),

		)
	}


	return nil
}
