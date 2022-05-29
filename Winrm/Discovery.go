package Winrm

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"log"
	"strings"
)

func Discovery(credMaps map[string]interface{}) {
	port := int(credMaps["port"].(float64))

	endpoint := winrm.NewEndpoint(credMaps["ip.address"].(string), port, false, false, nil, nil, nil, 0)

	client, err := winrm.NewClient(endpoint, credMaps["user"].(string), credMaps["password"].(string))

	defer func() {

		if r := recover(); r != nil {
			res := make(map[string]interface{})
			res["status"] = "failed"
			res["status.code"] = "200"
			res["error"] = r

			bytes, _ := json.Marshal(res)

			stringEncode := b64.StdEncoding.EncodeToString(bytes)
			log.SetFlags(0)
			log.Print(stringEncode)

		}

	}()

	commandfordisk := "Get-WmiObject win32_logicaldisk | Foreach-Object {$_.DeviceId,$_.Freespace,$_.Size -join \" \"}"

	_, _, _, errClient := client.RunPSWithString(commandfordisk, "")

	result := make(map[string]interface{})

	if err != nil {

		log.Fatal(err)

	}
	if errClient != nil {

		log.SetFlags(0)

		err := errClient.Error()

		subStringUnknownErrorPort := "unknown error"

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

		} else if strings.Contains(err, subStringUnknownErrorPort) {
			result["status"] = "failed"

			result["error"] = errClient.Error()

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

	result["ip.address"] = credMaps["ip.address"]

	result["metric.type"] = credMaps["metric.type"]

	result["port"] = credMaps["port"]

	result["user"] = credMaps["user"]

	result["password"] = credMaps["password"]

	data, _ := json.Marshal(result)

	stringEncode := b64.StdEncoding.EncodeToString(data)

	fmt.Println(stringEncode)

}
