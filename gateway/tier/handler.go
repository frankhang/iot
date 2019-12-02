package main

import (
	"context"
	"github.com/frankhang/util/tcp"
)


//tierHandler implements Hanlder
type TierHandler struct {
	*TierPacketIO

	driver *TireDriver
}

func NewTierHandler(tierPacketIO *TierPacketIO, driver *TireDriver) *TierHandler {
	return &TierHandler{TierPacketIO: tierPacketIO, driver: driver}
}


func (th *TierHandler) Handle(ctx context.Context, cc *tcp.ClientConn, data []byte) error {
	return nil
}