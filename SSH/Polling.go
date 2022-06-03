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

func Polling(credMaps map[string]interface{}) {
	port := uint16(credMaps["port"].(float64))
	sshHost := credMaps["ip.address"].(string)
	sshUser := credMaps["user"].(string)
	sshPassword := credMaps["password"].(string)
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
	sshClient, errorDial := ssh.Dial("tcp", addr, config)

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

	resultStatus := make(map[string]interface{})
	if errorDial != nil {

		log.SetFlags(0)

		err := errorDial.Error()
		//log.Fatal(err)
		subStringPortError := "connection refused"

		subStringDialError := "handshake failed"

		subStringUnknownErrorPort := "unknown error"

		if strings.Contains(err, subStringPortError) {

			resultStatus["status"] = "failed"

			resultStatus["error"] = "Port Invalid, Connection refused"

			resultStatus["status.code"] = "400"

			data, _ := json.Marshal(resultStatus)

			stringEncode := b64.StdEncoding.EncodeToString(data)

			log.SetFlags(0)

			log.Fatal(stringEncode)

		} else if strings.Contains(err, subStringDialError) {

			resultStatus["status"] = "failed"

			resultStatus["error"] = "ssh Handshake Failed, user,password or ip.address does not match each other"

			resultStatus["status.code"] = "400"

			data, _ := json.Marshal(resultStatus)

			stringEncode := b64.StdEncoding.EncodeToString(data)

			log.SetFlags(0)

			log.Fatal(stringEncode)

		} else if strings.Contains(err, subStringUnknownErrorPort) {
			resultStatus["status"] = "failed"

			resultStatus["error"] = errorDial.Error()

			resultStatus["status.code"] = "401"

			data, _ := json.Marshal(resultStatus)

			stringEncode := b64.StdEncoding.EncodeToString(data)

			log.SetFlags(0)

			log.Fatal(stringEncode)
		}

	} else {

		resultStatus["status"] = "success"

		resultStatus["status.code"] = "200"
	}

	data := make(map[string]interface{})
	var result = ""

	if credMaps["metricGroup"] == "Cpu" {
		result = fetchCpu(sshClient)
	} else if credMaps["metricGroup"] == "Memory" {
		result = fetchMemory(sshClient)
	} else if credMaps["metricGroup"] == "Process" {
		result = fetchProcess(sshClient)
	} else if credMaps["metricGroup"] == "System" {
		result = fetchSystem(sshClient)
	} else if credMaps["metricGroup"] == "Disk" {
		result = fetchDisk(sshClient)
	}
	data["monitorId"] = credMaps["monitorId"]
	data["metricGroup"] = credMaps["metricGroup"]
	data["metric.type"] = credMaps["metric.type"]
	data["value"] = result

	dataMarshal, errMarshal := json.Marshal(data)
	if errMarshal != nil {
		res := make(map[string]interface{})
		res["error"] = errMarshal.Error()
		bytes, _ := json.Marshal(res)

		stringEncode := b64.StdEncoding.EncodeToString(bytes)
		log.SetFlags(0)
		log.Print(stringEncode)

	}

	stringEncode := b64.StdEncoding.EncodeToString(dataMarshal)

	fmt.Println(stringEncode)
}
func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func fetchCpu(client *ssh.Client) string {
	var getCpuMap []map[string]string
	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}

	session, err = client.NewSession()
	cpuUserPercent, err := session.Output("mpstat | grep \"all\"")
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
	cpuCore, err := session.Output("mpstat -P ALL")
	if err != nil {
		panic(err)
	}

	cpuCoreSplit := string(cpuCore)
	cpuCoreSplittemp := strings.Split(cpuCoreSplit, "\n")
	flag := 0
	for _, value := range cpuCoreSplittemp {
		if flag < 3 {
			flag++
			continue
		}

		split1 := strings.Split(standardizeSpaces(value), " ")
		if len(split1) < 14 {
			break
		}

		tempCpu := map[string]string{
			"cpu.core.name":         split1[3],
			"cpu.core.user.percent": split1[4],
			"cpu.core.idle.percent": split1[13],
			"cpu.core.sys.percent":  split1[6],
		}

		getCpuMap = append(getCpuMap, tempCpu)

	}

	result := map[string]interface{}{
		"cpu.percent":      splitSystemUserPercent[0],
		"cpu.user.percent": splitCpuUserPercent[4],
		"cpu.idle.percent": splitCpuUserPercent[13],
		"cpu.sys.percent":  splitCpuUserPercent[6],
		"cpu.core":         getCpuMap,
	}

	bytes, errMarshal := json.Marshal(result)

	if errMarshal != nil {
		res := make(map[string]interface{})
		res["error"] = errMarshal.Error()
		bytes, _ := json.Marshal(res)

		stringEncode := b64.StdEncoding.EncodeToString(bytes)
		log.SetFlags(0)
		log.Print(stringEncode)

	}
	return string(bytes)

}
func fetchMemory(client *ssh.Client) string {

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

	bytes, errMarshal := json.Marshal(result)

	if errMarshal != nil {
		res := make(map[string]interface{})
		res["error"] = errMarshal.Error()
		bytes, _ := json.Marshal(res)

		stringEncode := b64.StdEncoding.EncodeToString(bytes)
		log.SetFlags(0)
		log.Print(stringEncode)

	}
	return string(bytes)

}
func fetchProcess(client *ssh.Client) string {

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
	for _, value := range splitprocess {

		if flag < 1 {
			flag++
			continue
		}
		splitN := strings.SplitN(standardizeSpaces(value), " ", 11)
		if len(splitN) <= 10 {
			break
		}

		tempData := map[string]string{
			"process.user":           splitN[0],
			"process.pid":            splitN[1],
			"process.memory.percent": splitN[3],
			"process.command":        splitN[10],
			"process.cpu.percent":    splitN[2],
		}

		getProcessMap = append(getProcessMap, tempData)

	}

	result := map[string]interface{}{
		"process": getProcessMap,
	}

	bytes, errMarshal := json.Marshal(result)

	if errMarshal != nil {
		res := make(map[string]interface{})
		res["error"] = errMarshal.Error()
		bytes, _ := json.Marshal(res)

		stringEncode := b64.StdEncoding.EncodeToString(bytes)
		log.SetFlags(0)
		log.Print(stringEncode)

	}

	return string(bytes)

}
func fetchSystem(client *ssh.Client) string {

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

	bytes, errMarshal := json.Marshal(result)

	if errMarshal != nil {
		res := make(map[string]interface{})
		res["error"] = errMarshal.Error()
		bytes, _ := json.Marshal(res)

		stringEncode := b64.StdEncoding.EncodeToString(bytes)
		log.SetFlags(0)
		log.Print(stringEncode)

	}
	return string(bytes)

}
func fetchDisk(client *ssh.Client) string {
	var getCpuMap []map[string]string

	//DiskALL
	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	disk, err := session.Output("df --total")
	if err != nil {
		panic(err)
	}

	//DiskFreeBytePercent
	session, err = client.NewSession()
	if err != nil {
		panic(err)
	}
	diskFreePercent, err := session.Output("df --total | grep \"total\" |  awk '{print (($4/$2)*100)}';df --total | grep \"total\" |  awk '{print $3 \" \" $4}'")
	if err != nil {
		panic(err)
	}
	diskFreeBytestemp := string(diskFreePercent)
	splitdiskFreeBytePercent := strings.Split(standardizeSpaces(diskFreeBytestemp), "\n")
	splitdiskUserByte := strings.Split(standardizeSpaces(splitdiskFreeBytePercent[0]), " ")
	intsplitUno, _ := strconv.Atoi(splitdiskUserByte[1])
	intsplitDos, _ := strconv.Atoi(splitdiskUserByte[2])

	userFreeData := intsplitUno - intsplitDos

	//Volume
	diskVolumeAll := string(disk)
	diskVolumeAllTemp := strings.Split(diskVolumeAll, "\n")
	flag := 0

	for _, value := range diskVolumeAllTemp {
		if flag <= 0 {
			flag++
			continue
		}

		splitDiskVolume := strings.Split(standardizeSpaces(value), " ")
		if len(splitDiskVolume) < 6 {
			break
		}

		tempDisk := map[string]string{
			"disk.volume.used.percent": splitDiskVolume[4],
			"disk.volume.free.percent": splitdiskUserByte[0],
			"disk.volume.total.bytes":  splitDiskVolume[1],
			"disk.volume.free.bytes":   splitDiskVolume[3],
			"disk.volume.used.bytes":   splitDiskVolume[2],
		}

		getCpuMap = append(getCpuMap, tempDisk)
	}

	result := map[string]interface{}{
		//"disk.utilization" : ,
		"disk.total.bytes": splitdiskUserByte[0],
		"disk.used.bytes":  splitdiskUserByte[1],
		"disk.free.bytes":  userFreeData,
		"disk.volume":      getCpuMap,
	}
	bytes, errMarshal := json.Marshal(result)

	if errMarshal != nil {
		res := make(map[string]interface{})
		res["error"] = errMarshal.Error()
		bytes, _ := json.Marshal(res)

		stringEncode := b64.StdEncoding.EncodeToString(bytes)
		log.SetFlags(0)
		log.Print(stringEncode)

	}
	return string(bytes)

}
