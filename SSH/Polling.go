package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"strconv"
	"strings"
	"time"
)

func Polling(credMaps map[string]string) string {
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
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer func(sshClient *ssh.Client) {
		err := sshClient.Close()
		if err != nil {

		}
	}(sshClient)
	var result = ""

	if credMaps["metricGroup"] == "fetchCpu" {
		result = fetchCpu(sshClient)
	} else if credMaps["metricGroup"] == "fetchMemory" {
		result = fetchMemory(sshClient)
	} else if credMaps["metricGroup"] == "fetchProcess" {
		result = fetchProcess(sshClient)
	} else if credMaps["metricGroup"] == "fetchSystem" {
		result = fetchSystem(sshClient)
	} else if credMaps["metricGroup"] == "fetchDisk" {
		result = fetchDisk(sshClient)
	}
	return result
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
			"system.cpu.core.name":         split1[4],
			"system.cpu.core.user.percent": split1[5],
			"system.cpu.core.idle.percent": split1[13],
			"system.cpu.core.sys.percent":  split1[7],
		}

		getCpuMap = append(getCpuMap, tempCpu)

	}

	result := map[string]interface{}{
		"system.cpu.percent":      splitSystemUserPercent[0],
		"system.cpu.user.percent": splitCpuUserPercent[0],
		"system.cpu.idle.percent": splitCpuIdlePercentage[0],
		"system.cpu.sys.percent":  splitCpuSysPercentage[0],
		"system.cpu.core":         getCpuMap,
	}

	bytes, _ := json.Marshal(result)
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

	bytes, _ := json.Marshal(result)
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
	OsVersionData, err := session.Output("hostnamectl | grep \"Operating fetchSystem\"")
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
	diskFreePercent, err := session.Output("df --total | grep \"total\" |  awk '{print (($4/$2)*100)}'")
	if err != nil {
		panic(err)
	}
	diskFreeBytestemp := string(diskFreePercent)
	splitdiskFreeBytePercent := strings.Split(standardizeSpaces(diskFreeBytestemp), " ")

	//DisktotalBytes
	session, err = client.NewSession()
	if err != nil {
		panic(err)
	}
	diskTotal, err := session.Output("df --total | grep \"total\" |  awk '{print $4}'")
	if err != nil {
		panic(err)
	}
	diskFreeBytes := string(diskTotal)
	splitdiskTotal := strings.Split(standardizeSpaces(diskFreeBytes), " ")

	//DiskUserBytes
	session, err = client.NewSession()
	if err != nil {
		panic(err)
	}
	diskUserByte, err := session.Output("df --total | grep \"total\" |  awk '{print $3}'")
	if err != nil {
		panic(err)
	}
	diskUserBytestring := string(diskUserByte)
	splitdiskUserByte := strings.Split(standardizeSpaces(diskUserBytestring), " ")

	//TotalFreeBytes
	session, err = client.NewSession()
	if err != nil {
		panic(err)
	}
	diskTotalFreeByte, err := session.Output("df --total | grep \"total\" |  awk '{print $4}'")
	if err != nil {
		panic(err)
	}
	diskFreeBytestring := string(diskTotalFreeByte)
	splitdiskFreeTotalByte := strings.Split(standardizeSpaces(diskFreeBytestring), " ")

	//Volume
	diskVolumeAll := string(disk)
	diskVolumeAllTemp := strings.Split(diskVolumeAll, "\n")
	flag := 0

	for _, v := range diskVolumeAllTemp {
		if flag <= 0 {
			flag++
			continue
		}

		splitDiskVolume := strings.Split(standardizeSpaces(v), " ")
		if len(splitDiskVolume) < 6 {
			break
		}

		tempDisk := map[string]string{
			"disk.volume.used.percent": splitDiskVolume[4],
			"disk.volume.free.percent": splitdiskFreeBytePercent[0],
			"disk.volume.total.bytes":  splitDiskVolume[1],
			"disk.volume.free.bytes":   splitDiskVolume[3],
			"disk.volume.used.bytes":   splitDiskVolume[2],
		}

		getCpuMap = append(getCpuMap, tempDisk)
	}

	result := map[string]interface{}{
		//"disk.utilization" : ,
		"disk.total.bytes": splitdiskTotal[0],
		"disk.used.bytes":  splitdiskUserByte[0],
		"disk.free.bytes":  splitdiskFreeTotalByte[0],
		"disk.volume":      getCpuMap,
	}
	bytes, _ := json.Marshal(result)
	return string(bytes)

}
