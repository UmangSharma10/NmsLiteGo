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

	var credMap map[string]string

	err = json.Unmarshal(jsonDecodedString, &credMap)
	if err != nil {

	}

	if string(credMap["metric.type"]) == "linux" {

		if credMap["category"] == "discovery" {
			SSH.Discovery(credMap)
		} else if credMap["category"] == "polling" {
			SSH.Polling(credMap)
		}

	} else if string(credMap["metric.type"]) == "windows" {
		if credMap["category"] == "discovery" {
			Winrm.Discovery(credMap)
		}

	} else if string(credMap["metric.type"]) == "network" {

		if credMap["category"] == "discovery" {
			Snmp.Discovery(credMap)

		} else if credMap["category"] == "polling" {
			Snmp.Polling(credMap)
		}

	}

}
