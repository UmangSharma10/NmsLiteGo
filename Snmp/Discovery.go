package Snmp

import (
	g "github.com/gosnmp/gosnmp"
	"net"
	"strconv"
	"time"
)

func Discovery(credMaps map[string]string) bool {

	if credMaps["discovery"] == "true" {
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
		if err != nil {
			return false
			//log.Fatalf("Connect() err: %v", err)
		}
		defer func(Conn net.Conn) {
			err := Conn.Close()
			if err != nil {

			}
		}(params.Conn)

		System(params)
		// polling
	} else if credMaps["discovery"] == "false" {
		//polling
	}
	return true
}
