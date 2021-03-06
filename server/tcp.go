/* (c) 2018 felipe@astroza.cl - See LICENSE
 */
package server

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/astroza/nrf_sniffer_ip/common"
)

const defaultServerPort = "5311"

func getServerPort() string {
	var serverPort string
	if serverPort = os.Getenv("SERVER_PORT"); serverPort == "" {
		serverPort = defaultServerPort
	}
	return serverPort
}

func CreateTCPServer() net.Listener {

	server, err := net.Listen("tcp", fmt.Sprintf(":%v", getServerPort()))
	if err != nil {
		log.Fatalf("Can't listen TCP port: %v", err.Error())
	}
	return server
}

func WaitForAClient(listener net.Listener) (*common.ReadWriterChannel, error) {
	conn, err := listener.Accept()
	if err != nil {
		return nil, err
	}
	r := common.ReadWriterChannel{}
	r.Init(conn)
	return &r, nil
}
