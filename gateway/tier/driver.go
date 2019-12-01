package main

import (
	"crypto/tls"
	"github.com/frankhang/util/tcp"
)

// TireDriver implements tcp.IDriver.
type TireDriver struct {
	tierIO *TierPacketIO
}

// NewTireDriver creates a new TireDriver.
func NewTireDriver(tierIO *TierPacketIO) *TireDriver {
	driver := &TireDriver{
		tierIO: tierIO,
	}
	tierIO.tierDriver = driver
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

func (td *TireDriver) GetPacketReader() tcp.PacketReader {
	return td.tierIO
}

func (td *TireDriver) GetPacketWriter() tcp.PacketWriter {
	return td.tierIO
}

func (td *TireDriver) OpenCtx(connID uint64, capability uint32, collation uint8, dbname string, tlsState *tls.ConnectionState) (tcp.QueryCtx, error) {
	return nil, nil
}
