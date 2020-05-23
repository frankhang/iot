package main

import (
	"github.com/frankhang/util/config"
	"github.com/frankhang/util/tcp"
)

// Driver implements tcp.IDriver.
type Driver struct {
	cfg *config.Config
}

// NewDriver creates a new Driver.
func NewDriver(cfg *config.Config) *Driver {
	driver := &Driver{
		cfg: cfg,
	}

	return driver
}

// Context implements QueryCtx.
type Context struct {
	currentDB string
}



func (td *Driver) OpenCtx(connID uint64, capability uint32, collation uint8, dbname string, tlsState *tls.ConnectionState) (tcp.QueryCtx, error) {
	return nil, nil
}

func (td *Driver) GeneratePacketIO(cc *tcp.ClientConn) *tcp.PacketIO {
	packetIO := tcp.NewPacketIO(cc)

	tPacketIO := NewPacketIO(packetIO, td)
	tHandler := NewHandler(tPacketIO, td)

	packetIO.PacketReader = tPacketIO
	packetIO.PacketWriter = tPacketIO

	cc.Handler = tHandler

	return packetIO
}
