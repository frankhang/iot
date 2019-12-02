package main

import (
	"context"
	"github.com/frankhang/util/tcp"
)

type Controller struct {
	ctx context.Context
	cc *tcp.ClientConn
}

func (ctl *Controller) TirePressureReport(data []byte) {
	
}

func (ctl *Controller) TireReplaceAck(data []byte) {
	
}



