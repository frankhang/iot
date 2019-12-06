package main

import (
	"bufio"
	"fmt"
	"github.com/frankhang/util/arena"
	"io"
	"net"
	"time"
)

const (
	defaultReaderSize = 4096
	defaultWriterSize = 4096
)



type bufferedReadConn struct {
	net.Conn
	rb *bufio.Reader
}

func newBufferedReadConn(conn net.Conn) *bufferedReadConn {
	return &bufferedReadConn{
		Conn: conn,
		rb:   bufio.NewReaderSize(conn, defaultReaderSize),
	}
}
func (conn *bufferedReadConn) Read(b []byte) (n int, err error) {
	return conn.rb.Read(b)
}


var (
//bufio.NewWriterSize(p.BufReadConn, defaultWriterSize)
)

func main() {

	//conn, err := net.Dial("tcp", "iot.cectiy.com:10001")
	conn, err := net.Dial("tcp", "localhost:10001")
	if err != nil {
		fmt.Println("check server")
		return
	}
	defer conn.Close()

	alloc := arena.NewAllocator(32 * 1024)

	bufReadConn := newBufferedReadConn(conn)
	bufWriter := bufio.NewWriterSize(bufReadConn, defaultWriterSize)

	//s := []byte("abc")
	//conn.Write(s)

	h := []byte{0x55, 0xAA, 0xAA, 0xBB, 0xEE, 0xEE}

	//fmt.Printf("h： %s\n", hex.EncodeToString(h))
	if _, err = bufWriter.Write(h); err != nil {
		println("write h error:", err)
		return
	}

	dd := alloc.Alloc(3 + 3 + 1)
	dd = append(dd, byte(9 + 3 + 1)) //length
	dd = append(dd, byte(1)) //tire number
	dd = append(dd, byte(1)) //user id

	//tier information
	dd  = append(dd, 0x32)
	dd  = append(dd, byte(1)) //temperature
	dd  = append(dd, byte(1)) //pressure

	//check sum
	dd = append(dd, 0xff)

	//fmt.Printf("dd： %s\n", hex.EncodeToString(dd))
	if _, err = bufWriter.Write(dd); err != nil {
		fmt.Printf("write dd error: %s", err)
		return
	}

	bufWriter.Flush()

	//Read packet from server

	var header [9]byte

	waitTimeout := time.Duration(3 * time.Second)

	if err := bufReadConn.SetReadDeadline(time.Now().Add(waitTimeout)); err != nil {
		//fmt.Printf("SetReadDeadline error: %s", err)
		return
	}

	if _, err := io.ReadFull(bufReadConn, header[:]); err != nil {
		fmt.Printf("read head error: %v", err)
		return
	}

	fmt.Printf("read header: [%x].\n", header)

	//length := int(header[6])


	tireNum := int(header[7])

	if tireNum > 0 {
		//todo for tierNum > 0
	}

	if err := bufReadConn.SetReadDeadline(time.Now().Add(waitTimeout)); err != nil {
		fmt.Printf("SetReadDeadline error: %s\n", err)
		return
	}

	var checksum [1]byte

	if _, err := io.ReadFull(bufReadConn, checksum[:]); err != nil {
		fmt.Printf("read head error: %s", err)
		return
	}

	fmt.Printf("read sum: [%x].\n", checksum[:])


	bufReadConn.Close()

}
