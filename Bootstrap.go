package main

import (
	"NmsLite/SSH"
	"NmsLite/Snmp"
	"NmsLite/Winrm"
	"encoding/base64"
	json "encoding/json"
	"os"
)

func main() {

	recevedARG1 := os.Args[1]

	jsonDecodedString, err := base64.StdEncoding.DecodeString(recevedARG1)

	if err != nil {
		panic(err)
	}

	var credMap map[string]interface{}

	err = json.Unmarshal(jsonDecodedString, &credMap)
	if err != nil {

	}

	if string(credMap["metric.type"].(string)) == "linux" {

		if credMap["category"] == "discovery" {
			SSH.GetDiscovery(credMap)
		} else if credMap["category"] == "polling" {
			SSH.Polling(credMap)
		}

	} else if string(credMap["metric.type"].(string)) == "windows" {
		if credMap["category"] == "discovery" {
			Winrm.Discovery(credMap)
		} else if credMap["category"] == "polling" {
			Winrm.Polling(credMap)
		}

	} else if string(credMap["metric.type"].(string)) == "network" {

		if credMap["category"] == "discovery" {
			Snmp.GetDiscovery(credMap)

		} else if credMap["category"] == "polling" {
			Snmp.Polling(credMap)
		}

	}

}
