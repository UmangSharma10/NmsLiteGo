package Winrm

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"log"
	"strconv"
	"strings"
)

func Discovery(credMaps map[string]string) {
	port, errPort := strconv.Atoi(credMaps["port"])

	endpoint := winrm.NewEndpoint(credMaps["ip.address"], port, false, false, nil, nil, nil, 0)

	client, err := winrm.NewClient(endpoint, credMaps["user"], credMaps["password"])

	commandfordisk := "Get-WmiObject win32_logicaldisk | Foreach-Object {$_.DeviceId,$_.Freespace,$_.Size -join \" \"}"

	_, _, _, errClient := client.RunPSWithString(commandfordisk, "")

	result := make(map[string]interface{})

	if err != nil {

		log.Fatal(err)

	}
	if errClient != nil {

		log.SetFlags(0)

		err := errClient.Error()

		subStringPortError := "connection refused"

		subStringDialError := "invalid content type"

		if strings.Contains(err, subStringPortError) {

			result["status"] = "failed"

			result["error"] = "Port Invalid, Connection refused"

			result["status.code"] = "400"

			data, _ := json.Marshal(result)

			stringEncode := b64.StdEncoding.EncodeToString(data)

			log.SetFlags(0)

			log.Fatal(stringEncode)

		} else if strings.Contains(err, subStringDialError) {

			result["status"] = "failed"

			result["error"] = "invalid content type, user,password or ip.address does not match each other"

			result["status.code"] = "401"

			data, _ := json.Marshal(result)

			stringEncode := b64.StdEncoding.EncodeToString(data)

			log.SetFlags(0)

			log.Fatal(stringEncode)

		}

	} else {

		result["status"] = "success"

		result["status.code"] = "200"
	}

	if errPort != nil {

		log.Fatal(errPort)

	}

	result["ip.address"] = credMaps["ip.address"]

	result["metric.type"] = credMaps["metric.type"]

	result["port"] = credMaps["port"]

	result["user"] = credMaps["user"]

	result["password"] = credMaps["password"]

	data, _ := json.Marshal(result)

	stringEncode := b64.StdEncoding.EncodeToString(data)

	fmt.Println(stringEncode)

}
