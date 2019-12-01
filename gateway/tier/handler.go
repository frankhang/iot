package main

import (
	"context"
	"github.com/frankhang/util/tcp"
)

var DefaultTierHanlder TierHandler
//tierHandler implements Hanlder
type TierHandler struct {
	driver *TireDriver
}

func (th *TierHandler) Handle(ctx context.Context, cc *tcp.ClientConn, data []byte) error {
	return nil
}