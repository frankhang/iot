package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/frankhang/util/hack"
	"github.com/frankhang/util/util"
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
	url = flag.String("url", "localhost:10002", "host:port")
)

func main() {

	crc := util.CrcCcittFfff([]byte("123456789"))
	fmt.Printf("crcffff of 123456789:%x\n", crc)

	flag.Parse()

	//conn, err := net.Dial("tcp", "iot.cectiy.com:10001")
	println("connecting : " + *url)
	conn, err := net.Dial("tcp", *url)
	if err != nil {
		fmt.Println("check server")
		return
	}
	defer conn.Close()

	//alloc := arena.NewAllocator(32 * 1024)

	bufReadConn := newBufferedReadConn(conn)
	bufWriter := bufio.NewWriterSize(bufReadConn, defaultWriterSize)

	//s := []byte("abc")
	//conn.Write(s)

	dd := []string{
		"Aev00AAAA00",
		"(evh00)",
		"(evc1100)",
		"Bev00BBBBBBBBBBBBBBBB00",
		"(evS100)",
		"(evr00)",

	}

	for _, d := range dd {
		data := createPacket(hack.Slice(d))
		println("sending...")
		fmt.Printf("packet = [%x]\n", data)
		fmt.Printf("packetStr = [%s]\n", data)

		//fmt.Printf("dd： %s\n", hex.EncodeToString(dd))
		if _, err = bufWriter.Write(data); err != nil {
			fmt.Printf("write data error: %s", err)
			return
		}

		bufWriter.Flush()
		if d[0] == '(' {
			if err = readResponse(bufReadConn, data); err != nil {
				fmt.Printf("read response error: %s", err)
				return
			}
		}

	}

	//waitTimeout := time.Duration(10 * time.Second)
	//
	//println("Reading...")
	//var d [1]byte
	//for {
	//	if err := bufReadConn.SetReadDeadline(time.Now().Add(waitTimeout)); err != nil {
	//		fmt.Printf("SetReadDeadline error: %v\n", err)
	//		return
	//	}
	//	if _, err := io.ReadFull(bufReadConn, d[:]); err != nil {
	//		fmt.Printf("ReadFull error: %v\n", err)
	//		return
	//	}
	//	//print(hack.String(d[:]))
	//	fmt.Printf("%c", d[0])
	//
	//}

}

func readResponse(bufReader *bufferedReadConn, request []byte) (err error) {
	waitTimeout := time.Duration(3 * time.Second)

	var data [7]byte

	if err = bufReader.SetReadDeadline(time.Now().Add(waitTimeout)); err != nil {
		return
	}
	if _, err = io.ReadFull(bufReader, data[:]); err != nil {
		return
	}
	fmt.Printf("readResponse: [%x]\n", data)
	fmt.Printf("readResponse Str: [%s]\n", data)
	if data[len(data)-1] != ')' {
		err = fmt.Errorf("end code should be )")
		return
	}

	ss := len(data) - 3
	crc := util.Crc16(data[:ss])
	//expectedCrc := binary.BigEndian.Uint16(data[ss:])

	//if crc != expectedCrc {
	//	err = fmt.Errorf("crc check error, %d != %d", crc, expectedCrc)
	//	return
	//}
	fmt.Printf("readResponse: crc =%d\n", crc)

	switch request[3] {
	case 'h':
		if data[3] == 'H' {
			fmt.Printf("***** Got response for protocol 1 *****\n")
		} else {
			err = fmt.Errorf("Got err response for protocol 1\n")
		}
	case 'c':
		if data[3] == 'C' {
			fmt.Printf("***** Got response for protocol 2 *****\n")
		} else {
			err = fmt.Errorf("Got err response for protocol 2\n")
		}
	case 'S':
		fmt.Printf("***** Got response for protocol 3 *****\n")
	case 'r':
		fmt.Printf("***** Got response for protocol 4 *****\n")

	}


	return
}

func createPacket(d []byte) []byte {
	var crcLen int

	data := make([]byte, len(d))
	copy(data, d)

	if data[0] == '(' {
		crcLen = len(data) - 3
	} else {
		crcLen = len(data) - 2
		len := len(data)
		binary.BigEndian.PutUint16(data[3:], uint16(len))
	}

	crc16 := util.CrcCcittFfff(data[:crcLen])
	//crc16 := uint16(0)
	binary.BigEndian.PutUint16(data[crcLen:], crc16)

	return data
}
