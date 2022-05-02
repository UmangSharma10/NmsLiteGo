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

		if credMap["category"] == "discovery" {
			var bval = SSH.Discovery(credMap)
			fmt.Println(bval)
		}
		/*else if(credMap["category"] == "polling"){
		}*/

		//eyJkZXZpY2UiOiJsaW51eCIsImhvc3QiOiJsb2NhbGhvc3QiLCJwb3J0IjoiMjIiLCJ1c2VyIjoidW1hbmciLCJwYXNzd29yZCI6Ik1pbmRAMTIzIiwiZGlzY292ZXJ5IjoidHJ1ZSJ9

	} else if string(credMap["device"]) == "windows" {

		Winrm.Discovery(credMap)

	} else if string(credMap["device"]) == "network" {

		if credMap["category"] == "discovery" {
			var value = Snmp.Discovery(credMap)
			fmt.Println(value)
		} else if credMap["category"] == "polling" {
			var result = Snmp.Polling(credMap)
			fmt.Println(result)
		}

	}

}
