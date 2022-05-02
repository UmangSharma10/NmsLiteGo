package Snmp

import (
	g "github.com/gosnmp/gosnmp"
	"log"
	"net"
	"strconv"
	"time"
)

func Discovery(credMaps map[string]string) string {

	if credMaps["category"] == "discovery" {
		port, _ := strconv.Atoi(credMaps["port"])
		// Build our own GoSNMP struct, rather than using g.Default.
		// Do verbose logging of packets.
		params := &g.GoSNMP{
			Target:    credMaps["host"],
			Port:      uint16(port),
			Community: "public",
			Version:   g.Version2c,
			Timeout:   time.Duration(2) * time.Second,
			//Logger:    g.NewLogger(log.New(os.Stdout, "", 0)),
		}
		err := params.Connect()

		_, err2 := params.Get([]string{"1.3.6.1.2.1.1.1.0"})
		if err2 != nil {
			return "failed"
		}
		if err != nil {
			log.Fatalf("Connect() err: %v", err)
			return "failed"
		}
		defer func(Conn net.Conn) {
			err := Conn.Close()
			if err != nil {

			}
		}(params.Conn)

		//TODO: ASk what to do
	}
	return "success"

}
