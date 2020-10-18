//
//
// Derived from a perl script developed by Klaus.Schmidinger@tvdr.de

package udpserve

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

func transformToIndicator(indicator int32) float64 {
	if indicator == 0 {
		return float64(0)
	}

	return float64(1)
}
func transformBoolToIndicator(indicator bool) float64 {
	if indicator {
		return float64(0)
	}

	return float64(1)
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

	http.Handle("/metrics", promhttp.Handler())
	s := &http.Server{
		Addr:           ":2112",
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		BaseContext:    func(net.Listener) context.Context { return ctx },
	}
	go func() {
		log.Fatal(s.ListenAndServe())
	}()
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
				// 0: Fühler Außentemperatur / OK
				systacomfortOutsideTemperatureCelsius.Set(float64(dp.Values[0]) / 10)
				// 1: Vorlauf Heizung / OK
				systacomfortHeatercircuitSupplyTemperatureCelsius.Set(float64(dp.Values[1]) / 10)
				// 2: Rücklauf Heizung / Looks good, but value does not match S-touch
				systacomfortHeatercircuitReturnTemperatureCelsius.Set(float64(dp.Values[2]) / 10)
				// 3=Brauchwasser TWO / OK
				systacomfortTapwaterTemperatureCelsius.Set(float64(dp.Values[3]) / 10)
				// 6=Zirkulation / OK
				systacomfortCirculationTemperatureCelsius.Set(float64(dp.Values[6]) / 10)
				// 9=Raumtemperatur / OK
				systacomfortInsideTemperatureCelsius.Set(float64(dp.Values[9]) / 10)
				// 11=Kollektor
				// 12=KesselVorlauf / OK
				var BoilerSupplyTemperatureCelsius = float64(dp.Values[12]) / 10
				systacomfortBoilerSupplyTemperatureCelsius.Set(BoilerSupplyTemperatureCelsius)
				// 13=KesselRuecklauf / OK
				var BoilerReturnTemperatureCelsius = float64(dp.Values[13]) / 10
				systacomfortBoilerReturnTemperatureCelsius.Set(BoilerReturnTemperatureCelsius)
				// 22=BrauchwasserSoll / OK
				systacomfortTapwaterTargetTemperatureCelsius.Set(float64(dp.Values[22]) / 10)
				// 23=InnenSoll / OK
				systacomfortInsideTargetTemperatureCelsius.Set(float64(dp.Values[23]) / 10)
				// 34=KesselSoll
				systacomfortBoilerTargetTemperatureCelsius.Set(float64(dp.Values[34]) / 10)
				// 36=Betriebsart 0=Heiprogramm 1, 1=Heiprogramm 2, 2=Heiprogramm 3, 3=Dauernd Normal, 4=Dauernd Komfort, 5=Dauernd Absenken
				systacomfortModeInfo.Set(float64(dp.Values[36]))
				// 39=Raumtemperatur normal (soll)
				// 40=Raumtemperatur komfort (soll)
				// 41=Raumtemperatur abgesenkt (soll)
				// 48=Fusspunkt / OK
				systacomfortHeatCurveRootPointCelsius.Set(float64(dp.Values[48]) / 10)
				// 50=Steilheit / OK
				systacomfortHeatCurveGradientCelsius.Set(float64(dp.Values[50]) / 10)
				// 52=Max. Vorlauftemperatur
				systacomfortBoilerSupplyUpperLimitCelsius.Set(float64(dp.Values[52]) / 10)
				// 53=Heizgrenze Heizbetrieb
				// 54=Heizgrenze Absenken
				// 55=Frostschutz Aussentemperatur
				// 56=Vorhaltezeit Aufheizen
				// 57=Raumeinfluss
				// 58=Ueberhoehung Kessel
				// 59=Spreizung Heizkreis
				// 60=Minimale Drehzahl Pumpe PHK
				// 62=Mischer Laufzeit
				// 149=Brauchwassertemperatur normal
				// 150=Brauchwassertemperatur komfort
				// 155=Brauchwassertemperatur Schaltdifferenz
				// 158=Nachlauf Pumpe PK/LP
				// 162=Min. Laufzeit Kessel
				// 169=Nachlaufzeit Pumpe PZ
				// 171=Zirkulation Schaltdifferenz
				// 179=Betriebszeit Kessel (Stunden) / OK
				// 180=Betriebszeit Kessel (Minuten) / OK
				systacomfortBoilerRuntimeSeconds.Set((float64(dp.Values[179])*6 + float64(dp.Values[180])/10) * 60)
				// 181=Anzahl Brennerstarts / OK
				systacomfortBoilerStartsTotal.Set(float64(dp.Values[181]))
				// 220=Aktive Relais (Bitpattern)
				//   RelaisHeizkreispumpe    = 0x0001
				//   RelaisLadepumpe         = 0x0080
				//   RelaisZirkulationspumpe = 0x0100
				//   RelaisKessel            = 0x0200
				// Brenner aktiv wenn RelaisKessel && (KesselVorlauf - KesselRuecklauf > 2);
				relayState := dp.Values[220]
				fmt.Printf("relais -> %d\n", dp.Values[220])
				systacomfortBoilerHeatercircuitPumpInfo.Set(transformToIndicator(relayState & 0x0001))
				systacomfortBoilerLoadpumpInfo.Set(transformToIndicator(relayState & 0x0080))
				systacomfortBoilerCirculationPumpInfo.Set(transformToIndicator(relayState & 0x0100))
				systacomfortBoilerActiveInfo.Set(transformToIndicator(relayState & 0x0200))

				systacomfortBoilerTorchInfo.Set(transformBoolToIndicator((relayState&0x0200) != 0 && (BoilerSupplyTemperatureCelsius-BoilerReturnTemperatureCelsius > 2)))

				// 228=Fehlerstatus (255 = OK)
				systacomfortBoilerErrorInfo.Set(float64(dp.Values[228]))
				fmt.Printf("228 -> %d\n", dp.Values[228])
				// 248=Status
				systacomfortBoilerStatusInfo.Set(float64(dp.Values[248]))
				// All State candidates
				// fmt.Printf("42 -> %d\n", dp.Values[42])
				// fmt.Printf("232 -> %d\n", dp.Values[232])
				// fmt.Printf("231 -> %d\n", dp.Values[231])
				// fmt.Printf("248 -> %d\n", dp.Values[248])

				//fmt.Printf("PacketType:%d\n", dp.PacketType)
				//for i := 0; i < 256; i++ {
				//	fmt.Printf("%d -> %f\n", i, float64(dp.Values[i])/10)
				//}
			case 2:
				// 56: TSA / OK
				systacomfortSolarpanelTemperatureCelsius.Set(float64(dp.Values[56]) / 10)
				// 57: Solar Rücklauf / OK
				systacomfortSolarpanelReturnTemperatureCelsius.Set(float64(dp.Values[57]) / 10)
				// 59: Solar Vorlauf / OK
				systacomfortSolarpanelSupplyTemperatureCelsius.Set(float64(dp.Values[59]) / 10)
				// 60: Temparatur Kollektor TAM / OK
				systacomfortSolarpanelOutsideTemperatureCelsius.Set(float64(dp.Values[60]) / 10)
				// 63: Maximale Kollektortemperatur / OK
				systacomfortSolarpanelMaximumTemperatureCelsius.Set(float64(dp.Values[63]) / 10)
				// 86: Solar Speicher oben TWO (=Warmwasser ist) / OK

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

// StartupServer asdfasdfasdf
func StartupServer(port int) {
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
