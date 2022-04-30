package main

import (
	"NmsLite/SSH"
	"NmsLite/Snmp"
	"NmsLite/Winrm"
	"encoding/base64"
	json "encoding/json"
	"fmt"
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

	//fmt.Println("it will collect addresss")

	if string(credMap["device"]) == "linux" {

		var bval = SSH.Discovery(credMap)
		fmt.Println(bval)

	} else if string(credMap["device"]) == "windows" {

		Winrm.Discovery(credMap)

	} else if string(credMap["device"]) == "network" {

		Snmp.Discovery(credMap)

	}

}
