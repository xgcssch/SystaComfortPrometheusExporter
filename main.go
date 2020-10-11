package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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

			_, addr, err := pc.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}

			//fmt.Printf("packet-received: bytes=%d from=%s\n", n, addr.String())

			type ReceivePacket struct {
				MacAddress [6]byte    // 0-5
				Counter    uint16     // 6-7
				_          [8]byte    // 8-15
				PacketType byte       // 16
				_          [7]byte    // 17-23
				Values     [256]int32 // 24-1048
			}

			type ResponsePacket struct {
				MacAddress   [6]byte // 0-5
				InputCounter uint16  // 6-7
				_            [4]byte // 8-11
				PacketType   byte    // 12
				_            [3]byte // 13-15
				PacketID     uint16  // 16-17
				Counter      uint16  // 18-19
			}

			var dp ReceivePacket

			br := bytes.NewReader(buffer)
			binary.Read(br, binary.LittleEndian, &dp)

			switch dp.PacketType {
			case 0:
				//fmt.Printf("Initial Packet\n")
			case 1:
				var Aussentemperator = float64(dp.Values[0]) / 10
				var VorlaufIst = float64(dp.Values[1]) / 10
				var RuecklaufIst = float64(dp.Values[2]) / 10
				fmt.Printf("Außentemperatur:%f\n", Aussentemperator)
				fmt.Printf("Vorlauf:%f\n", VorlaufIst)
				fmt.Printf("Rücklauf:%f\n", RuecklaufIst)

				//fmt.Printf("PacketType:%d\n", dp.PacketType)
				//for i := 0; i < 256; i++ {
				//	fmt.Printf("%d -> %f\n", i, float64(dp.Values[i])/10)
				//}
			case 2:
				var Aussentemperator = float64(dp.Values[60]) / 10
				fmt.Printf("ST Außentemperatur:%f\n", Aussentemperator)

				//fmt.Printf("PacketType:%d\n", dp.PacketType)
				//for i := 0; i < 256; i++ {
				//	fmt.Printf("%d -> %f\n", i, float64(dp.Values[i])/10)
				//}
			default:
				//fmt.Printf("Unknown PacketType:%d\n", dp.PacketType)
			}

			var ReturnID uint16 = (uint16(dp.MacAddress[5]) << 8) + uint16(dp.MacAddress[4]) + macOffset
			var Counter uint16 = dp.Counter + counterOffset
			rp := ResponsePacket{dp.MacAddress, dp.Counter, [4]byte{0, 0, 0, 0}, 0x01, [3]byte{0, 0, 0}, ReturnID, Counter}
			bw := bytes.NewBuffer(make([]byte, 0))
			binary.Write(bw, binary.LittleEndian, &rp)
			//
			pc.WriteTo(bw.Bytes(), addr)
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

	http.Handle("/metrics", promhttp.Handler())
	s := &http.Server{
		Addr:           ":2112",
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		BaseContext:    func(net.Listener) context.Context { return ctx },
	}
	log.Fatal(s.ListenAndServe())
	//http.ListenAndServe(":2112", nil)

	err := server(ctx, ":22460")
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Ended ...")
	//!-
}
