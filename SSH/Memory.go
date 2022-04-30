package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
)

func Memory(client *ssh.Client) {

	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	totalMemoryALl, err := session.Output("free -b | grep \"Mem\" |  awk '{print $2}' && free -b | grep \"Mem\" |  awk '{print $3}' && free -b | grep \"Mem\" |  awk '{print $4}' && free -b | grep \"Mem\" |  awk '{print $7}' && free -b | grep \"Mem\" |  awk '{print $4}' && free -b | grep \"Mem\" |  awk '{print ($7)}' && free -b | grep \"Mem\" |  awk '{print ($3/$2)*100}' && free -b | grep \"Mem\" |  awk '{print ($4/$2)*100}' && free -b | grep \"Swap\" |  awk '{print ($2)}'")
	if err != nil {
		panic(err)
	}

	allmemory := string(totalMemoryALl)
	allmemorysplit := strings.Split(allmemory, "\n")

	result := map[string]interface{}{
		"memory.installed":    allmemorysplit[0],
		"memory.used":         allmemorysplit[1],
		"memory.free":         allmemorysplit[2],
		"memory.available":    allmemorysplit[3],
		"memory.used.percent": allmemorysplit[6],
		"memory.free.percent": allmemorysplit[7],
		"memory.swap.total":   allmemorysplit[8],
	}

	bytes, _ := json.Marshal(result)
	fmt.Println(string(bytes))

}
