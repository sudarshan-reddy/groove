package main

import (
	"fmt"
	"log"

	"github.com/sudarshan-reddy/groove"
	"github.com/sudarshan-reddy/groove/dht"
)

func main() {
	aq, err := groove.InitGroove(0x04)
	if err != nil {
		log.Fatal(err)
	}
	defer aq.Close()

	reading, err := aq.AnalogRead(14)
	if err != nil {
		log.Fatal(err)
	}

	smbData, err := dht.ReadDHT()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(smbData.CelsiusTemp, reading)
}
