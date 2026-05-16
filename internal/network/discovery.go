package network

import (
	"log"

	"github.com/grandcat/zeroconf"
)

func StartDiscovery() {
	server, err := zeroconf.Register(
		"MediaServer",
		"_mediaserver._tcp",
		"local.",
		9000,
		[]string{"https=true"},
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("mDNS ativo")

	_ = server
}
