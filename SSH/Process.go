package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
)

func Process(client *ssh.Client) {

	var getProcessMap []map[string]string

	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	psaux, err := session.Output("ps aux")
	if err != nil {
		panic(err)
	}
	psauxString := string(psaux)

	splitprocess := strings.Split(psauxString, "\n")

	flag := 0
	for _, v := range splitprocess {

		if flag < 1 {
			flag++
			continue
		}
		splitN := strings.SplitN(standardizeSpaces(v), " ", 11)
		if len(splitN) <= 10 {
			break
		}

		temp1 := map[string]string{
			"process.user":           splitN[0],
			"process.pid":            splitN[1],
			"process.memory.percent": splitN[3],
			"process.command":        splitN[10],
			"process.cpu.percent":    splitN[2],
		}

		getProcessMap = append(getProcessMap, temp1)

	}

	result := map[string]interface{}{
		"process": getProcessMap,
	}

	bytes, _ := json.Marshal(result)
	fmt.Println(string(bytes))

}
