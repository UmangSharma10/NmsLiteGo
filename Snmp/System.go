package Snmp

import (
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"os"
)

var list []string

func System(client *g.GoSNMP) {
	var listofInterface []map[string]interface{}
	interfaceOids := []string{
		".1.3.6.1.2.1.2.2.1.1", ".1.3.6.1.2.1.2.2.1.2", ".1.3.6.1.2.1.2.2.1.3", ".1.3.6.1.2.1.2.2.1.5", ".1.3.6.1.2.1.2.2.1.6", ".1.3.6.1.2.1.2.2.1.7", ".1.3.6.1.2.1.2.2.1.8", ".1.3.6.1.2.1.2.2.1.8", ".1.3.6.1.2.1.2.2.1.14", ".1.3.6.1.2.1.2.2.1.20", ".1.3.6.1.2.1.2.2.1.16", ".1.3.6.1.2.1.2.2.1.10"}
	err := client.Walk(".1.3.6.1.2.1.2.2.1.1", Walk)
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
			case g.IPAddress:
				listData[variable.Name] = w.Variables[key].Value

			default:

			}
		}
		listofInterface = append(listofInterface, listData)
	}
}

func Walk(pdu g.SnmpPDU) error {
	str := fmt.Sprintf("%v", pdu.Value)
	list = append(list, str)
	return nil
}
