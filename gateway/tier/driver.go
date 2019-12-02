package main

import (
	"crypto/tls"
	"github.com/frankhang/util/tcp"
)

// TireDriver implements tcp.IDriver.
type TireDriver struct {
	cfg *tcp.Config
}

// NewTireDriver creates a new TireDriver.
func NewTireDriver(cfg *tcp.Config) *TireDriver {
	driver := &TireDriver{
		cfg: cfg,
	}

	return driver
}

// TireContext implements QueryCtx.
type TireContext struct {
	currentDB string
}

// TireStatement implements PreparedStatement.
type TireStatement struct {
	id  uint32
	ctx *TireContext
}

func (td *TireDriver) OpenCtx(connID uint64, capability uint32, collation uint8, dbname string, tlsState *tls.ConnectionState) (tcp.QueryCtx, error) {
	return nil, nil
}

func (td *TireDriver) GeneratePacketIO(cc *tcp.ClientConn) *tcp.PacketIO {
	packetIO := tcp.NewPacketIO(cc)

	tierPacketIO := NewTierPacketIO(packetIO, td)
	tierHandler := NewTierHandler(tierPacketIO, td)

	packetIO.PacketReader = tierPacketIO
	packetIO.PacketWriter = tierPacketIO

	cc.Handler = tierHandler

	return packetIO
}
