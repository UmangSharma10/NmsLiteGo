package Winrm

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/masterzen/winrm"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func Polling(credMaps map[string]interface{}) {
	port := int(credMaps["port"].(float64))

	endpoint := winrm.NewEndpoint(credMaps["ip.address"].(string), port, false, false, nil, nil, nil, 0)

	client, err := winrm.NewClient(endpoint, credMaps["user"].(string), credMaps["password"].(string))

	if err != nil {
		panic(err)
	}

	defer func() {

		if r := recover(); r != nil {
			res := make(map[string]interface{})

			res["error"] = r

			bytes, _ := json.Marshal(res)

			stringEncode := b64.StdEncoding.EncodeToString(bytes)
			log.SetFlags(0)
			log.Print(stringEncode)

		}

	}()

	data := make(map[string]interface{})

	var result = ""
	if credMaps["metricGroup"] == "Cpu" {
		result = fetchCpu(client)
	} else if credMaps["metricGroup"] == "Memory" {
		result = fetchMemory(client)
	} else if credMaps["metricGroup"] == "Process" {
		result = fetchProcess(client)
	} else if credMaps["metricGroup"] == "Disk" {
		result = fetchDisk(client)
	} else if credMaps["metricGroup"] == "System" {
		result = fetchSystem(client)
	}

	data["monitorId"] = credMaps["monitorId"]
	data["metricGroup"] = credMaps["metricGroup"]
	data["metric.type"] = credMaps["metric.type"]
	data["value"] = result

	dataMarshal, _ := json.Marshal(data)

	stringEncode := b64.StdEncoding.EncodeToString(dataMarshal)

	fmt.Println(stringEncode)

}

func fetchCpu(client *winrm.Client) string {

	result := make(map[string]interface{})

	var cores []map[string]interface{}

	ac := "(Get-Counter '\\Processor(*)\\% Idle Time','\\Processor(*)\\% Processor Time','\\Processor(*)\\% user time' -ErrorAction SilentlyContinue).countersamples | Format-List  -Property Path,Cookedvalue;"

	output := ""

	output, _, _, _ = client.RunPSWithString(ac, "")

	re := regexp.MustCompile("Path\\s*:\\s\\\\\\\\\\w*-\\w*\\\\\\w*\\((\\S*)\\)\\\\([\\w\\d\\s%]+)\\n\\w*\\s\\:\\s(\\d*.\\d*)")

	value := re.FindAllStringSubmatch(output, -1)

	var count = 3

	size := len(value) / count

	for coreCount := 0; coreCount < len(value)/count; coreCount++ {

		count := 0

		core := make(map[string]interface{})

		if value[coreCount][1] == "_total" {

			result["system.cpu.idle.percent"] = value[coreCount][3]

			result["system.cpu.process.percent"] = value[count+size][3]

			result["system.cpu.user.percent"] = strings.Split(value[count+size+size][3], "\r")[0]

		} else {

			core["core.name"] = value[coreCount][1]

			core["core.idle.percent"] = value[coreCount][3]

			core["core.process.percent"] = value[count+size][3]

			core["core.user.percent"] = strings.Split(value[count+size+size][3], "\r")[0]

			cores = append(cores, core)
		}
	}

	result["system.cpu.core"] = cores

	bytes, _ := json.Marshal(result)

	return string(bytes)
}

func fetchMemory(client *winrm.Client) string {
	result := make(map[string]interface{})

	output := ""
	ac := "Get-WmiObject win32_OperatingSystem |%{\"{0} {1} {2} {3}\" -f $_.totalvisiblememorysize, $_.freephysicalmemory, $_.totalvirtualmemorysize, $_.freevirtualmemory}"
	output, _, _, _ = client.RunPSWithString(ac, "")

	myStringArray := strings.Split(standardizeSpaces(output), " ")

	totalVisibleMemory := myStringArray[0]
	physicalMemory := myStringArray[1]
	totalVirtualMemory := myStringArray[2]
	freeVirtualMemory := myStringArray[3]

	result["total.visible.memory"] = totalVisibleMemory
	result["physical.memory"] = physicalMemory
	result["total.virtual.memory"] = totalVirtualMemory
	result["free.virtual.memory"] = freeVirtualMemory

	bytes, _ := json.Marshal(result)

	return string(bytes)
}

func fetchProcess(client *winrm.Client) string {
	result := make(map[string]interface{})

	var processData []map[string]interface{}

	output := ""

	ac := "(Get-Counter '\\Process(*)\\ID Process','\\Process(*)\\% Processor Time','\\Process(*)\\Thread Count' -ErrorAction SilentlyContinue).countersamples | Format-List  -Property Path,Cookedvalue;"

	output, _, _, _ = client.RunPSWithString(ac, "")

	re := regexp.MustCompile("Path\\s*:\\s\\\\\\\\(\\w*-\\w*)\\\\\\w*\\((\\S*)\\)\\\\([\\w\\d\\s%]+)\\n\\w*\\s\\:\\s(\\d*)")

	value := re.FindAllStringSubmatch(output, -1)

	processData = append(processData, result)

	var count int

	for processCount := 0; processCount < len(value); processCount++ {

		tempProcessDatai := make(map[string]interface{})

		tempProcessDataj := make(map[string]interface{})

		processName := value[processCount][2]

		for j := 0; j < len(processData); j++ {

			tempProcessDatai = processData[j]

			if tempProcessDatai[processName] != nil {

				count = 1

				break

			} else {

				count = 0

			}

		}

		if count == 0 {

			tempProcessDataj["process.name"] = processName

			if (value[processCount][3]) == "id process\r" {

				tempProcessDataj["process.id"] = value[processCount][4]

			} else if value[processCount][3] == "% processor time\r" {

				tempProcessDataj["process.processor.time.percent"] = value[processCount][4]

			} else if value[processCount][3] == "thread count\r" {

				tempProcessDataj["process.thread.count"] = value[processCount][4]
			}

			processData = append(processData, tempProcessDataj)

		} else {

			if (value[processCount][3]) == "id process\r" {

				tempProcessDatai["process.id"] = value[processCount][4]

			} else if value[processCount][3] == "% processor time\r" {

				tempProcessDatai["process.processor.time.percent"] = value[processCount][4]

			} else if value[processCount][3] == "thread count\r" {

				tempProcessDatai["process.thread.count"] = value[processCount][4]
			}

			count = 1

			processData = append(processData, tempProcessDatai)
		}

	}

	processData = processData[1:]

	size := (len(processData)) / 3

	var values []map[string]interface{}

	for kCount := 0; kCount < len(processData)/3; kCount = kCount + 1 {

		count := kCount

		tempMapData := make(map[string]interface{})

		tempMapData = processData[kCount]

		tempMapData["process.processor.time.percent"] = processData[count+size]["process.processor.time.percent"]

		tempMapData["process.thread.count"] = processData[count+size+size]["process.thread.count"]

		values = append(values, tempMapData)
	}
	result["process"] = values

	data, _ := json.Marshal(result)

	return string(data)
}

func fetchDisk(client *winrm.Client) string {

	result := make(map[string]interface{})

	output := ""

	ac := "Get-WmiObject win32_logicaldisk |Foreach-Object {$_.DeviceId,$_.Freespace,$_.Size}"

	output, _, _, _ = client.RunPSWithString(ac, "")

	res := strings.Split(output, "\r\n")

	var disks []map[string]interface{}

	var usedBytes int64

	var totalBytes int64

	for diskCount := 0; diskCount < len(res); diskCount = diskCount + 3 {

		disk := make(map[string]interface{})

		disk["disk.name"] = strings.Split(res[diskCount], ":")[0]

		if (diskCount+1) > len(res) || res[diskCount+1] == "" {

			disk["disk.free.bytes"] = 0

			disk["disk.total.bytes"] = 0

			disk["disk.available.bytes"] = 0

			disk["disk.used.percent"] = 0

			disk["disk.free.percent"] = 0

			disks = append(disks, disk)

			break
		}

		bytes, _ := strconv.ParseInt(res[diskCount+1], 10, 64)

		usedBytes = usedBytes + bytes

		disk["disk.available.bytes"], _ = strconv.ParseInt(res[diskCount+1], 10, 64)

		bytes, _ = strconv.ParseInt(res[diskCount+2], 10, 64)

		totalBytes = totalBytes + bytes

		disk["disk.total.bytes"] = bytes

		disk["disk.used.bytes"] = (disk["disk.total.bytes"]).(int64) - (disk["disk.available.bytes"]).(int64)

		disk["disk.used.percent"] = (float64((float64((disk["disk.total.bytes"]).(int64)) - float64((disk["disk.used.bytes"]).(int64))) / float64((disk["Disk.Total.Bytes"].(int64))))) * 100

		disk["disk.free.percent"] = 100 - disk["disk.used.percent"].(float64)

		disks = append(disks, disk)
	}

	result["disk.total.bytes"] = totalBytes

	result["disk.used.bytes"] = usedBytes

	result["disk.available.bytes"] = totalBytes - usedBytes

	result["disk.used.percent"] = ((float64(totalBytes) - float64(usedBytes)) / float64(totalBytes)) * 100

	result["disk.available.percent"] = 100.00 - (result["Disk.Used.Percent"]).(float64)

	result["disk"] = disks

	data, _ := json.Marshal(result)

	return string(data)

}

func fetchSystem(client *winrm.Client) string {
	result := make(map[string]interface{})
	a := "aa"
	output := ""
	ac := "(Get-WmiObject win32_operatingsystem).name;(Get-WMIObject win32_operatingsystem).version;whoami;(Get-WMIObject win32_operatingsystem).LastBootUpTime;" // Command jo humko run karna hain
	output, _, _, _ = client.RunPSWithString(ac, a)
	res1 := strings.Split(output, "\n")
	result["system.os.name"] = strings.Split(res1[0], "\r")[0]
	result["system.os.version"] = strings.Split(res1[1], "\r")[0]
	result["system.user.name"] = strings.Split(res1[2], "\r")[0]
	result["system.up.time"] = strings.Split(res1[3], "\r")[0]
	result["status"] = "success"
	data, _ := json.Marshal(result)
	return string(data)
}
