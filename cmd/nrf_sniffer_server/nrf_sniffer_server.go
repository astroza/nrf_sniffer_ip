/* (c) 2018 felipe@astroza.cl - See LICENSE
 */
package main

import (
	"log"

	"github.com/astroza/nrf_sniffer_ip/serial"
	"github.com/astroza/nrf_sniffer_ip/server"
)

func main() {
	listener := server.CreateTCPServer()
	for {
		clientConn, err := server.WaitForAClient(listener)
		if err != nil {
			log.Printf("Can't accept a new client: %v", err)
		}
		serialPort, err := serial.OpenPort()
		if err != nil {
			log.Fatalf("Can't open serial port: %v", err)
		}
	clientLoop:
		for {
			select {
			case buffer := <-serialPort.In:
				_, err := clientConn.Handle.Write(buffer)
				if err != nil {
					log.Println(err)
					break clientLoop
				}
			case buffer := <-clientConn.In:
				_, err := serialPort.Handle.Write(buffer)
				if err != nil {
					log.Fatal(err)
				}
			case err := <-clientConn.Error:
				log.Println(err)
				break clientLoop
			case err := <-serialPort.Error:
				log.Fatal(err)
			}
		}
	}
}
