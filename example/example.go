package main

import (
	"fmt"
	"log"
	"time"

	"github.com/MasandeM/sps30"

	"go.bug.st/serial"
)

func main() {

	mode := &serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	log.Println("Connecting to UART")

	uart, err := serial.Open("/dev/ttyUSB0", mode) //should be read from a config file or something
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully Connected")

	device := sps30.New(uart)

	//Create a struct that is passed to a function then is populated. th
	version_info := sps30.VersionInfo{}
	err = device.ReadVersion(&version_info)

	if err != nil {
		log.Fatal("Error reading version information: ", err)
	}

	fmt.Printf("FW: %d.%d, HW: %d, SHDLC: %d.%d\n",
		version_info.FirmwarMajor,
		version_info.FirmwarMinor,
		version_info.HardwarRevision,
		version_info.SHDLCMajor,
		version_info.SHDLCMinor)

	err = device.StartMeasurement()
	if err != nil {
		log.Fatal("error starting measurement")
	}

	measurement := sps30.Measurement{}
	for {

		err = device.ReadMeasurement(&measurement)
		if err != nil {
			fmt.Printf("[-] error reading measurement: %v\n", err)
		} else {
			fmt.Printf(`
measured values:
				%0.2f pm1.0
				%0.2f pm2.5
				%0.2f pm4.0
				%0.2f pm10.0
				%0.2f nc0.5
				%0.2f nc1.0
				%0.2f nc2.5
				%0.2f nc4.5
				%0.2f nc10.0
				%0.2f typical particle size
`,
				measurement.Mc1p0, measurement.Mc2p5, measurement.Mc4p0, measurement.Mc10p0, measurement.Nc0p5,
				measurement.Nc1p0, measurement.Nc2p5, measurement.Nc4p0, measurement.Nc10p0,
				measurement.TypicalParticleSize)
		}

		time.Sleep(time.Second)
	}
}
