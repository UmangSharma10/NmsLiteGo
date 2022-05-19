package Snmp

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"log"
	"time"
)

func Discovery(credMaps map[string]interface{}) {
	var version = g.Version2c

	switch credMaps["version"] {

	case "v1":

		version = g.Version1

		break
	case "v2c":

		version = g.Version2c

		break
	}

	port, _ := credMaps["port"].(int64)

	params := &g.GoSNMP{

		Target: credMaps["ip.address"].(string),

		Port: uint16(port),

		Community: credMaps["community"].(string),

		Version: version,

		Timeout: time.Duration(2) * time.Second,
		//Logger:    g.NewLogger(log.New(os.Stdout, "", 0)),
	}

	result := make(map[string]interface{})

	_, errGet := params.Get([]string{".1.3.6.1.2.1.1.1.0"})

	if errGet != nil {

		result["status"] = "failed"

		result["error"] = "Port Invalid, Connection refused"

		result["status.code"] = "400"

		data, _ := json.Marshal(result)

		stringEncode := b64.StdEncoding.EncodeToString(data)

		log.SetFlags(0)

		log.Fatal(stringEncode)

	} else {

		result["status"] = "success"

		result["status.code"] = "200"
	}

	result["ip.address"] = credMaps["ip.address"]
	result["port"] = credMaps["port"]
	result["community"] = credMaps["community"]
	result["version"] = credMaps["version"]

	data, _ := json.Marshal(result)
	stringEncode := b64.StdEncoding.EncodeToString(data)
	fmt.Println(stringEncode)
}
