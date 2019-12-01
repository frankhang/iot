package main

import (
	//"github.com/frankhang/util/tcp"
)

//tierPacket implements PacketReader and PacketWriter
type TierPacketIO struct {
	tierDriver *TireDriver

}

var DefaultTierPacketIO TierPacketIO

func (p *TierPacketIO) ReadOnePacket() ([]byte, error) {


	return nil, nil
}



func (p *TierPacketIO) WritePacket(data []byte) error {
	return nil
}