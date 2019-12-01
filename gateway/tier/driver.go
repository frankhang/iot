package main

import (
	"crypto/tls"
	"github.com/frankhang/util/tcp"
)

// TireDriver implements tcp.IDriver.
type TireDriver struct {
	cfg *tcp.Config
	tierIO *TierPacketIO
	tierHandler *TierHandler

}

// NewTireDriver creates a new TireDriver.
func NewTireDriver(cfg *tcp.Config, tierIO *TierPacketIO, tierHandler *TierHandler) *TireDriver {
	driver := &TireDriver{
		cfg: cfg,
		tierIO: tierIO,
		tierHandler:tierHandler,
	}
	tierIO.driver = driver
	tierHandler.driver = driver

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

func (td *TireDriver) GetPacketReader() tcp.PacketReader {
	return td.tierIO
}

func (td *TireDriver) GetPacketWriter() tcp.PacketWriter {
	return td.tierIO
}


func (td *TireDriver) SetPacketIO(packetIO *tcp.PacketIO) {
	td.tierIO.PacketIO = packetIO

}

func (td *TireDriver) GetHandler() tcp.Handler {
	return td.tierHandler
}


