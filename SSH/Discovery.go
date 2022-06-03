package SSH

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"strings"
	"time"
)

func GetDiscovery(credMaps map[string]interface{}) {

	port := uint16(credMaps["port"].(float64))

	sshHost := credMaps["ip.address"].(string)

	sshUser := credMaps["user"].(string)

	sshPassword := credMaps["password"].(string)

	sshPort := port
	// Create SSHP login configuration
	config := &ssh.ClientConfig{

		Timeout: 10 * time.Second, //ssh connection time out time is one second, if SSH validation error returns in one second

		User: sshUser,

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		Config: ssh.Config{Ciphers: []string{
			"aes128-ctr", "aes192-ctr", "aes256-ctr",
		}},
		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}

	config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}

	// dial gets SSH client
	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)

	sshClient, errorDial := ssh.Dial("tcp", addr, config)

	result := make(map[string]interface{})

	defer func() {

		if deferError := recover(); deferError != nil {
			res := make(map[string]interface{})
			res["status"] = "failed"
			res["status.code"] = "200"
			res["error"] = deferError

			bytes, _ := json.Marshal(res)

			stringEncode := b64.StdEncoding.EncodeToString(bytes)
			log.SetFlags(0)
			log.Print(stringEncode)

		}

	}()

	if errorDial != nil {

		log.SetFlags(0)

		err := errorDial.Error()
		//log.Fatal(err)
		subStringPortError := "connection refused"

		subStringDialError := "handshake failed"

		subStringUnknownErrorPort := "unknown error"

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

			result["error"] = "ssh Handshake Failed, user,password or ip.address does not match each other"

			result["status.code"] = "400"

			data, _ := json.Marshal(result)

			stringEncode := b64.StdEncoding.EncodeToString(data)

			log.SetFlags(0)

			log.Fatal(stringEncode)

		} else if strings.Contains(err, subStringUnknownErrorPort) {
			result["status"] = "failed"

			result["error"] = errorDial.Error()

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

	session, err := sshClient.NewSession()

	if err != nil {

		result["status"] = "failed"

		result["error"] = err.Error()

		result["status.code"] = "400"

	} else {

		result["status"] = "success"

		result["status.code"] = "200"
	}
	_, err = session.Output("uname")

	if err != nil {

		result["status"] = "failed"

		result["error"] = "discovery Failed, command uname did not work."

		result["status.code"] = "400"

	} else {

		result["status"] = "success"

		result["status.code"] = "200"
	}

	result["dis.id"] = credMaps["dis.id"]
	result["ip.address"] = credMaps["ip.address"]

	result["metric.type"] = credMaps["metric.type"]

	result["port"] = credMaps["port"]

	result["user"] = credMaps["user"]

	result["password"] = credMaps["password"]

	data, _ := json.Marshal(result)

	stringEncode := b64.StdEncoding.EncodeToString(data)

	fmt.Println(stringEncode)
}
