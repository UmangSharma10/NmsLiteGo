package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
)

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
func Cpu(client *ssh.Client) {
	var getCpuMap []map[string]string
	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}

	session, err = client.NewSession()
	cpuUserPercent, err := session.Output("mpstat | awk 'NR==4{print $5}'")
	if err != nil {
		panic(err)
	}
	cpuUPercent := string(cpuUserPercent)
	splitCpuUserPercent := strings.Split(standardizeSpaces(cpuUPercent), " ")

	session, err = client.NewSession()
	systemUserPercent, err := session.Output("mpstat | awk 'NR == 4 {print ($5+$6+$7+$8+$9+$10+$11+$12+$13)}'")
	if err != nil {
		panic(err)
	}
	systemUpercent := string(systemUserPercent)
	splitSystemUserPercent := strings.Split(standardizeSpaces(systemUpercent), " ")

	session, err = client.NewSession()
	cpuIdlePercent, err := session.Output("mpstat | awk 'NR==4 {print $14}'")
	if err != nil {
		panic(err)
	}
	cpuIdlePercentage := string(cpuIdlePercent)
	splitCpuIdlePercentage := strings.Split(standardizeSpaces(cpuIdlePercentage), " ")

	session, err = client.NewSession()
	cpuSysPercent, err := session.Output("mpstat | awk 'NR==4 {print $7}'")
	if err != nil {
		panic(err)
	}
	cpuSysPercentage := string(cpuSysPercent)
	splitCpuSysPercentage := strings.Split(standardizeSpaces(cpuSysPercentage), " ")

	session, err = client.NewSession()
	cpuCore, err := session.Output("mpstat -P ALL")
	if err != nil {
		panic(err)
	}

	cpuCoreSplit := string(cpuCore)
	cpuCoreSplittemp := strings.Split(cpuCoreSplit, "\n")
	flag := 0
	for _, v := range cpuCoreSplittemp {
		if flag < 3 {
			flag++
			continue
		}

		split1 := strings.Split(standardizeSpaces(v), " ")
		if len(split1) < 14 {
			break
		}

		tempCpu := map[string]string{
			"cpu.core.name":         split1[4],
			"cpu.core.user.percent": split1[5],
			"cpu.core.idle.percent": split1[13],
			"cpu.core.sys.percent":  split1[7],
		}

		getCpuMap = append(getCpuMap, tempCpu)

	}

	result := map[string]interface{}{
		"cpu.percent":      splitSystemUserPercent[0],
		"cpu.user.percent": splitCpuUserPercent[0],
		"cpu.idle.percent": splitCpuIdlePercentage[0],
		"cpu.sys.percent":  splitCpuSysPercentage[0],
		"cpu.core":         getCpuMap,
	}

	bytes, _ := json.Marshal(result)
	fmt.Println(string(bytes))

}
