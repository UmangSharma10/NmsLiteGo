package Winrm

import (
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"strconv"
	"strings"
)

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func Cpu(credMaps map[string]string) {

	port, _ := strconv.Atoi(credMaps["port"])
	endpoint := winrm.NewEndpoint(credMaps["host"], port, false, false, nil, nil, nil, 0)

	client, err := winrm.NewClient(endpoint, credMaps["user"], credMaps["password"])

	if err != nil {
		panic(err)
	}
	output := ""
	ac := "Get-WmiObject win32_OperatingSystem |%{\"{0} {1} {2} {3}\" -f $_.totalvisiblememorysize, $_.freephysicalmemory, $_.totalvirtualmemorysize, $_.freevirtualmemory}"
	output, _, _, err = client.RunPSWithString(ac, "")

	myStringArray := strings.Split(standardizeSpaces(output), " ")
	//var memorylist map[string]string

	totalVisibleMemory := myStringArray[0]
	physicalMemory := myStringArray[1]
	totalVirtualMemory := myStringArray[2]
	freeVirtualMemory := myStringArray[3]

	commandfordisk := "Get-WmiObject win32_logicaldisk | Foreach-Object {$_.DeviceId,$_.Freespace,$_.Size -join \" \"}"
	command := " "
	command, _, _, err = client.RunPSWithString(commandfordisk, "")
	//fmt.Println(command)
	var disklist []map[string]string
	array := strings.Split(command, "\n")
	for _, v := range array {
		splits := strings.Split(v, " ")

		if len(splits) == 3 {
			temp := map[string]string{
				"disk":       splits[0],
				"Free-space": splits[1],
				"Size":       splits[2],
			}
			disklist = append(disklist, temp)
		}

		if len(splits) != 3 {
			temp := map[string]string{
				"disk":       splits[0],
				"Free-space": "0",
				"Size":       "0",
			}
			disklist = append(disklist, temp)
		}

	}

	command1 := "Get-WmiObject win32_processor | select SystemName, LoadPercentage"
	command, _, _, err = client.RunPSWithString(command1, "")
	//fmt.Println(command)
	var getprocess []map[string]string
	getarray := strings.Split(command1, " ")

	tempprocessor := map[string]string{
		"System_Name":    getarray[0],
		"LoadPercentage": getarray[1],
	}

	getprocess = append(getprocess, tempprocessor)

	result := map[string]interface{}{
		"Disk":               disklist,
		"totalVisualMemory":  totalVisibleMemory,
		"physicalMemory":     physicalMemory,
		"totalVirtualMemory": totalVirtualMemory,
		"freeVirtualMemory":  freeVirtualMemory,
		"process":            getprocess,
	}

	bytes, err := json.Marshal(result)
	fmt.Println(bytes)
}
