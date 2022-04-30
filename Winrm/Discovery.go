package Winrm

import (
	"github.com/masterzen/winrm"
	"strconv"
)

func Discovery(credMaps map[string]string) bool {

	if credMaps["discovery"] == "true" {
		port, _ := strconv.Atoi(credMaps["port"])
		endpoint := winrm.NewEndpoint(credMaps["host"], port, false, false, nil, nil, nil, 0)

		client, err := winrm.NewClient(endpoint, credMaps["user"], credMaps["password"])

		if err != nil {
			return false
		}

		Disk(client)

	}
	return true

}
