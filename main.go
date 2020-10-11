package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const maxBufferSize = 1024 * 64
const counterOffset uint16 = 0x3FBF
const macOffset uint16 = 0x8E83
const replyMsgLength = 20

func memset(a []byte, v byte) {
	for i := range a {
		a[i] = v
	}
}

func server(ctx context.Context, address string) (err error) {
	// ListenPacket provides us a wrapper around ListenUDP so that
	// we don't need to call `net.ResolveUDPAddr` and then subsequentially
	// perform a `ListenUDP` with the UDP address.
	//
	// The returned value (PacketConn) is pretty much the same as the one
	// from ListenUDP (UDPConn) - the only difference is that `Packet*`
	// methods and interfaces are more broad, also covering `ip`.
	pc, err := net.ListenPacket("udp", address)
	if err != nil {
		return
	}

	// `Close`ing the packet "connection" means cleaning the data structures
	// allocated for holding information about the listening socket.
	defer pc.Close()

	doneChan := make(chan error, 1)

	// Given that waiting for packets to arrive is blocking by nature and we want
	// to be able of canceling such action if desired, we do that in a separate
	// go routine.
	go func() {
		for {
			buffer := make([]byte, maxBufferSize)

			n, addr, err := pc.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}

			fmt.Printf("packet-received: bytes=%d from=%s\n", n, addr.String())

			type DataPacket struct {
				//MacAddress [6]byte // 0-5
				//Counter    uint16  // 6-7
				//_          [8]byte // 8-15
				//PacketType uint16  // 16
				//_          [7]byte // 17-23
				V1         [15]byte
				PacketType byte
				V2         [8]byte
				Values     [232]int32
			}

			var dp DataPacket

			br := bytes.NewReader(buffer)
			binary.Read(br, binary.LittleEndian, &dp)

			if dp.PacketType != 0 {
				fmt.Printf("PacketType:%d\n", dp.PacketType)
			} else {
				fmt.Printf("Initial Packet\n")
			}

			replyBuffer := bytes.Repeat([]byte{0}, replyMsgLength)
			replyBuffer[12] = 0x01

			var i int = 88
			fmt.Printf("%d", i)
		}
	}()

	select {
	case <-ctx.Done():
		//fmt.Println("cancelled")
		//err = ctx.Err()
		err = nil
	case err = <-doneChan:
	}

	return
}

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func(cancel context.CancelFunc) {
		//sig := <-sigs
		<-sigs
		//fmt.Println()
		//fmt.Println(sig)

		cancel()

		done <- true
	}(cancel)
	log.Print("Starting ...")

	err := server(ctx, ":22460")
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Ended ...")
	//!-
}
