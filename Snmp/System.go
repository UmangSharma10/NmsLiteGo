package Snmp

import (
	"encoding/json"
	"fmt"
	"github.com/gosnmp/gosnmp"
)

func System(client *gosnmp.GoSNMP) {

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
	fmt.Println(string(bytes))

}
