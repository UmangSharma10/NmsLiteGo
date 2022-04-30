package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
)

func System(client *ssh.Client) {

	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	systemDataAll, err := session.Output("ps axo state | grep \"R\" | wc -l && ps axo state | grep \"B\" | wc -l")
	if err != nil {
		panic(err)
	}
	systemDataString := string(systemDataAll)
	splitsystemData := strings.Split(systemDataString, "\n")

	session, err = client.NewSession()
	if err != nil {
		panic(err)
	}
	uptimeData, err := session.Output("uptime -p | awk '{print $4*3600 + $4*60 }'")
	uptimeDataString := string(uptimeData)
	uptimeDataSplit := strings.Split(uptimeDataString, "\n")

	session, err = client.NewSession()
	if err != nil {
		panic(err)
	}
	OsData, err := session.Output("hostnamectl | grep \"hostname\"")
	OsDataString := string(OsData)
	OsNameDataSplit := strings.Split(OsDataString, ":")
	OsNameSplit := strings.Split(OsNameDataSplit[1], "\n")
	//OsNameSplit2 := strings.Split(OsNameSplit[1], "\n")

	session, err = client.NewSession()
	if err != nil {
		panic(err)
	}
	OsVersionData, err := session.Output("hostnamectl | grep \"Operating System\"")
	OsVersionDataString := string(OsVersionData)
	OsVersionDataSplit := strings.Split(OsVersionDataString, ":")
	OsVersionSplit := strings.Split(OsVersionDataSplit[1], "\n")
	//ThreadCount
	session, err = client.NewSession()
	if err != nil {
		panic(err)
	}
	systemThreadData, err := session.Output("ps -eLf | wc -l")
	systemThreadDataString := string(systemThreadData)
	systemThreadDataSplit := strings.Split(systemThreadDataString, ":")
	sysThread := strings.Split(systemThreadDataSplit[0], "\n")
	result := map[string]interface{}{
		"system.running.process":  splitsystemData[0],
		"system.blocking.process": splitsystemData[1],
		"system.uptime":           uptimeDataSplit[0],
		"system.thread":           sysThread[0],
		"system.os.name":          OsNameSplit[0],
		"system.os.version":       OsVersionSplit[0],
	}

	bytes, _ := json.Marshal(result)
	fmt.Println(string(bytes))

}
