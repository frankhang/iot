package main

import (
	"bufio"
	"fmt"
	"github.com/frankhang/util/arena"
	"github.com/frankhang/util/hack"
	"github.com/frankhang/util/util"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	defaultReaderSize = 4096
	defaultWriterSize = 4096

	sizeHeader = 27
	locSize    = 18
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

	conn, err := net.Dial("tcp", "iot.cectiy.com:10001")
	//conn, err := net.Dial("tcp", "localhost:10001")
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

	h := []byte{0x35, 0x35, 0x20, 0x41, 0x41, 0x20, 0x39, 0x30, 0x20, 0x38, 0x39, 0x20, 0x30, 0x33, 0x20, 0x30, 0x34, 0x20}

	size := len(h) + 3*4 + 9*2
	println("size = ", size)

	dd := alloc.Alloc(size)

	dd = append(dd, h...)
	dd = append(dd, hack.Slice(fmt.Sprintf("%-3d", size))...) //size
	dd = append(dd, hack.Slice(fmt.Sprintf("%-3d", 2))...)   //tire number
	dd = append(dd, 0x30, 0x31, 0x20) //user id
	dd = append(dd, 0x31, 0x31, 0x20, 0x30, 0x30, 0x20, 0x30, 0x30, 0x20)// data of first tire
	dd = append(dd, 0x31, 0x32, 0x20, 0x30, 0x31, 0x20, 0x30, 0x30, 0x20)// data of second tire


	sum := util.Sum(dd)
	dd = append(dd, hack.Slice(fmt.Sprintf("%3d", sum))...) //check sum

	fmt.Printf("packet = [%x]\n", dd)

	//fmt.Printf("dd： %s\n", hex.EncodeToString(dd))
	if _, err = bufWriter.Write(dd); err != nil {
		fmt.Printf("write dd error: %s", err)
		return
	}

	bufWriter.Flush()

	//Read packet from server

	var header [sizeHeader]byte

	waitTimeout := time.Duration(3 * time.Second)

	if err := bufReadConn.SetReadDeadline(time.Now().Add(waitTimeout)); err != nil {
		fmt.Printf("SetReadDeadline error: %v\n", err)
		return
	}

	if _, err := io.ReadFull(bufReadConn, header[:]); err != nil {
		fmt.Printf("read head error: %v\n", err)
		return
	}

	fmt.Printf("read header: [%x].\n", header)

	s := hack.String(header[locSize : locSize+3])
	if size, err = strconv.Atoi(strings.TrimSpace(s)); err != nil {
		fmt.Printf("get size error: %v\n", err)
		return
	}

	data := alloc.AllocWithLen(size -sizeHeader, size -sizeHeader)
	if _, err := io.ReadFull(bufReadConn, data); err != nil {
		fmt.Printf("read head error: %s", err)
		return
	}

	fmt.Printf("read data: [%x].\n", data)


	tireNum := int(header[7])

	if tireNum > 0 {
		//todo for tierNum > 0
	}

	if err := bufReadConn.SetReadDeadline(time.Now().Add(waitTimeout)); err != nil {
		fmt.Printf("SetReadDeadline error: %s\n", err)
		return
	}

	bufReadConn.Close()

}
