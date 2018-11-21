/* (c) 2018 felipe@astroza.cl - See LICENSE
 */
package server

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Announcer struct {
	Run chan bool
}

func (announcer *Announcer) Init() {
	announcer.Run = make(chan bool)
	broadcastAddrs := make([]net.IP, 0, 4)
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatalln(err)
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				// IPv4 Only
				if v.IP.To4() != nil {
					broadcastAddr := make(net.IP, 4)
					for i, n := range v.IP[12:16] {
						broadcastAddr[i] = n | ^v.Mask[i]
					}
					broadcastAddrs = append(broadcastAddrs, broadcastAddr)
				}

			}
		}
	}
	go func(addrs []net.IP) {
		var isRunning bool
		disabledTick := make(<-chan time.Time)
		var ticker *time.Ticker
		tick := disabledTick
		for {
			select {
			case isRunning = <-announcer.Run:
				if isRunning {
					ticker = time.NewTicker(5 * time.Second)
					tick = ticker.C
				} else if ticker != nil {
					ticker.Stop()
					ticker = nil
					tick = disabledTick
				}
			case <-tick:
				announcer.doSend(addrs)
			}
		}
	}(broadcastAddrs)
}

func (announcer *Announcer) doSend(addrs []net.IP) {
	for _, addr := range addrs {
		handle, err := net.Dial("udp", fmt.Sprintf("%v:%v", addr, defaultServerPort))
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Fprintf(handle, "TCP:%v", getServerPort())
		handle.Close()
	}
}
