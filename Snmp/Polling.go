package Snmp

import (
	"encoding/json"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func Polling(credMaps map[string]string) string {
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
		log.Fatalf("Connect() err: %v", err)

	}
	defer func(Conn net.Conn) {
		err := Conn.Close()
		if err != nil {

		}
	}(params.Conn)
	var result = ""
	if credMaps["metricGroup"] == "fetchInterface" {
		result = fetchInterface(params)
		//fmt.Println(result)
	} else if credMaps["metricGroup"] == "fetchSystem" {
		result = fetchSystem(params)
		// fmt.Println(result)
	}

	return result
}

var list []string

func fetchInterface(client *g.GoSNMP) string {
	var listofInterface []map[string]interface{}
	interfaceOids := []string{
		".1.3.6.1.2.1.2.2.1.1", ".1.3.6.1.2.1.2.2.1.2", ".1.3.6.1.2.1.2.2.1.3", ".1.3.6.1.2.1.2.2.1.5", ".1.3.6.1.2.1.2.2.1.6", ".1.3.6.1.2.1.2.2.1.7", ".1.3.6.1.2.1.2.2.1.8", ".1.3.6.1.2.1.2.2.1.14", ".1.3.6.1.2.1.2.2.1.20", ".1.3.6.1.2.1.2.2.1.16", ".1.3.6.1.2.1.2.2.1.10"}
	err := client.Walk(".1.3.6.1.2.1.2.2.1.1", walk)
	if err != nil {
		os.Exit(1)
	}

	for i := 0; i < len(list); i++ {
		var newArray []string
		var listData = make(map[string]interface{})
		for _, value := range interfaceOids {
			value = value + "." + list[i]
			newArray = append(newArray, value)
		}

		w, err := client.Get(newArray)
		if err != nil {
		}

		for key, variable := range w.Variables {
			switch variable.Type {
			case g.Integer:
				listData[variable.Name] = w.Variables[key].Value
			case g.OctetString:
				listData[variable.Name] = string(w.Variables[key].Value.([]byte))
			case g.Gauge32:
				listData[variable.Name] = w.Variables[key].Value

			case g.Counter32:
				listData[variable.Name] = w.Variables[key].Value

			default:

			}
		}
		listofInterface = append(listofInterface, listData)
	}
	result := map[string]interface{}{
		"fetchInterface": listofInterface,
	}

	bytes, _ := json.Marshal(result)
	return string(bytes)

}

func walk(pdu g.SnmpPDU) error {
	str := fmt.Sprintf("%v", pdu.Value)
	list = append(list, str)
	return nil
}

func fetchSystem(client *g.GoSNMP) string {

	sysnameOid := "1.3.6.1.2.1.1.5.0"
	sysdescOid := "1.3.6.1.2.1.1.1.0"
	sysuptimeOid := "1.3.6.1.2.1.1.3.0"
	syslocationOid := "1.3.6.1.2.1.1.6.0"
	sysOid := ".1.3.6.1.2.1.1.2.0"

	oids := []string{sysnameOid, sysdescOid, sysuptimeOid, syslocationOid, sysOid}
	sysGet, err2 := client.Get(oids)
	if err2 != nil {
	}
	sysNameTemp := ""
	sysDecsTemp := ""
	sysupTimeTemp := ""
	sysLocationtemp := ""
	sysOidtemp := ""
	for i, variable := range sysGet.Variables {

		switch i {
		case 0:
			sysNameTemp = string(variable.Value.([]byte))
		case 1:
			sysDecsTemp = string(variable.Value.([]byte))
		case 2:
			str := fmt.Sprintf("%v", variable.Value)
			sysupTimeTemp = str
		case 3:
			sysLocationtemp = string(variable.Value.([]byte))
		case 4:
			sysOidtemp = variable.Value.(string)

		}
	}

	result := map[string]interface{}{
		"system.name":        sysNameTemp,
		"system.description": sysDecsTemp,
		"system.uptime":      sysupTimeTemp,
		"system.location":    sysLocationtemp,
		"system.oid":         sysOidtemp,
	}

	bytes, _ := json.Marshal(result)
	return string(bytes)

}
