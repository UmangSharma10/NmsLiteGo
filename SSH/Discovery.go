package SSH

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"strconv"
	"time"
)

func Discovery(credMaps map[string]string) bool {

	if credMaps["discovery"] == "true" {
		port, _ := strconv.Atoi(credMaps["port"])
		sshHost := credMaps["host"]
		sshUser := credMaps["user"]
		sshPassword := credMaps["password"]
		sshPort := port
		// Create SSHP login configuration
		config := &ssh.ClientConfig{
			Timeout:         10 * time.Second, //ssh connection time out time is one second, if SSH validation error returns in one second
			User:            sshUser,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Config: ssh.Config{Ciphers: []string{
				"aes128-ctr", "aes192-ctr", "aes256-ctr",
			}},
			//HostKeyCallback: hostKeyCallBackFunc(h.Host),
		}
		config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}

		// dial gets SSH client
		addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
		sshClient, err := ssh.Dial("tcp", addr, config)

		session, err := sshClient.NewSession()
		if err != nil {
			panic(err)
		}
		_, err = session.Output("uname")
		if err != nil {
			return false
		}

		if err != nil {
			return false
		}
		defer func(sshClient *ssh.Client) {
			err := sshClient.Close()
			if err != nil {

			}
		}(sshClient)

		//Call Polling
		//Cpu(sshClient)
		//Disk(sshClient)
		//Memory(sshClient)
		//Process(sshClient)
		System(sshClient)
	} else if credMaps["discovery"] == "false" {
		//call polling
	}

	return true
}
