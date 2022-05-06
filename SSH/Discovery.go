package SSH

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"strconv"
	"strings"
	"time"
)

func Discovery(credMaps map[string]string) {

	port, errPort := strconv.Atoi(credMaps["port"])

	sshHost := credMaps["ip.address"]

	sshUser := credMaps["user"]

	sshPassword := credMaps["password"]

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

	defer func(sshClient *ssh.Client) {

		err := sshClient.Close()

		if err != nil {

			result["status"] = "failed"

			result["Error"] = err.Error()
		} else {

			result["status"] = "success"

			result["status.code"] = "200"
		}
	}(sshClient)

	if errPort != nil {

		result["status"] = "failed"

		result["error"] = "Port invalid"

		result["status.code"] = "400"

		data, _ := json.Marshal(result)

		stringEncode := b64.StdEncoding.EncodeToString(data)

		log.SetFlags(0)

		log.Fatal(stringEncode)

	}

	if errorDial != nil {

		log.SetFlags(0)

		err := errorDial.Error()
		//log.Fatal(err)
		subStringPortError := "connection refused"

		subStringDialError := "handshake failed"

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

		result["error"] = "Discovery Failed, command uname did not work."

		result["status.code"] = "400"

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
