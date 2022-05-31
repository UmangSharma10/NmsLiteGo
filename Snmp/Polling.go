package Snmp

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func Polling(credMaps map[string]interface{}) {

	var version = g.Version2c

	switch credMaps["version"].(string) {

	case "v1":

		version = g.Version1

		break
	case "v2c":

		version = g.Version2c

		break
	}

	var community = "public"
	switch credMaps["community"].(string) {

	case "public":

		community = "public"

		break
	case "private":

		community = "private"

		break
	}
	port := uint16(credMaps["port"].(float64))
	// Build our own GoSNMP struct, rather than using g.Default.
	// Do verbose logging of packets.
	params := &g.GoSNMP{
		Target:    credMaps["ip.address"].(string),
		Port:      uint16(port),
		Community: community,
		Version:   version,
		Timeout:   time.Duration(2) * time.Second,
		//Logger:    g.NewLogger(log.New(os.Stdout, "", 0)),
	}
	err := params.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)

	}

	defer func() {

		if r := recover(); r != nil {
			res := make(map[string]interface{})

			res["error"] = r

			bytes, _ := json.Marshal(res)

			stringEncode := b64.StdEncoding.EncodeToString(bytes)
			log.SetFlags(0)
			log.Print(stringEncode)

		}

	}()

	data := make(map[string]interface{})
	var result = ""
	if credMaps["metricGroup"] == "Interface" {
		result = fetchInterface(params)

	} else if credMaps["metricGroup"] == "System" {
		result = fetchSystem(params)
	}
	data["monitorId"] = credMaps["monitorId"]
	data["metricGroup"] = credMaps["metricGroup"]
	data["metric.type"] = credMaps["metric.type"]
	data["value"] = result

	dataMarshal, errMarshal := json.Marshal(data)
	if errMarshal != nil {
		res := make(map[string]interface{})
		res["error"] = errMarshal.Error()
		bytes, _ := json.Marshal(res)

		stringEncode := b64.StdEncoding.EncodeToString(bytes)
		log.SetFlags(0)
		log.Print(stringEncode)

	}

	stringEncode := b64.StdEncoding.EncodeToString(dataMarshal)

	fmt.Println(stringEncode)
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

		for _, variable := range w.Variables {
			vname := fmt.Sprintf("%v", variable.Name)
			ch := strings.SplitAfter(vname, ".1.3.6.1.2.1.2.2.1.")
			arr := strings.Split(ch[1], ".")
			choice, _ := strconv.Atoi(arr[0])

			switch choice {
			case 2:
				listData["interface.Description"] = string(variable.Value.([]byte))

			case 3:
				if variable.Value == 1 {
					listData["interface.type"] = "other"
				}
				if variable.Value == 6 {
					listData["interface.type"] = "ethernetCsmacd"
				}
				if variable.Value == 135 {
					listData["interface.type"] = "l2vlan"
				}
				if variable.Value == 53 {
					listData["interface.type"] = "propVirtual"
				}

			case 5:
				listData["interface.speed"] = variable.Value
			case 7:
				if variable.Value == 2 {
					listData["interface.admin.status"] = "down"
				}
				if variable.Value == 1 {
					listData["interface.admin.status"] = "up"
				}

			case 8:
				if variable.Value == 2 {
					listData["interface.operating.status"] = "down"
				}
				if variable.Value == 1 {
					listData["interface.operating.status"] = "up"
				}
			case 14:
				if variable.Value == nil {
					listData["ifInError"] = ""
				}

				if variable.Value == 0 {
					listData["ifInError"] = variable.Value
				}

			case 16:
				if variable.Value == 0 {
					listData["interface.out.octetes"] = variable.Value
				} else {
					listData["interface.out.octetes"] = variable.Value
				}

			case 20:
				if variable.Value == 0 {
					listData["ifOutError"] = variable.Value
				}

				if variable.Value == nil {
					listData["ifOutError"] = ""
				}

			case 10:
				if variable.Value == 0 {
					listData["interface.in.octetes"] = variable.Value
				} else {
					listData["interface.in.octetes"] = variable.Value
				}
			default:

			}

			listofInterface = append(listofInterface, listData)
		}
	}
	result := map[string]interface{}{
		"fetchInterface": listofInterface,
	}

	bytes, errMarshal := json.Marshal(result)

	if errMarshal != nil {
		res := make(map[string]interface{})
		res["error"] = errMarshal.Error()
		bytes, _ := json.Marshal(res)

		stringEncode := b64.StdEncoding.EncodeToString(bytes)
		log.SetFlags(0)
		log.Print(stringEncode)

	}
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

	bytes, errMarshal := json.Marshal(result)

	if errMarshal != nil {
		res := make(map[string]interface{})
		res["error"] = errMarshal.Error()
		bytes, _ := json.Marshal(res)

		stringEncode := b64.StdEncoding.EncodeToString(bytes)
		log.SetFlags(0)
		log.Print(stringEncode)

	}
	return string(bytes)

}
