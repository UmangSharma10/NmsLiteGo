package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
)

func Disk(client *ssh.Client) {
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
	fmt.Println(string(bytes))

}
