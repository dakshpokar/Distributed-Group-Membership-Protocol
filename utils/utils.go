package utils

import (
	"MP2/types"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
*

	GetHostName() returns the hostname from a key

	Parameters:
		key: a string

	Returns:
		string: the hostname

*
*/
func GetHostName(key string) string {
	return strings.SplitN(key, ":", 2)[0]
}

/*
*

	FilteredHostName() returns the hostname from a key

	Parameters:
		hostName: a string

	Returns:
		string: the hostname

*
*/
func FilteredHostName(hostName string) string {
	if hostName[len(hostName)-1] == '.' {
		return hostName[:len(hostName)-1]
	}
	if hostName == "localhost" {
		hostName = os.Getenv("INTRODUCER_ADDR")
	}
	return hostName
}

func PrintSuspectMachinesMap(machines map[string]types.MachineInfo) {
	var machineIds []string
	for machineId := range machines {
		machineIds = append(machineIds, machineId)
	}
	sort.Strings(machineIds)

	for _, machineId := range machineIds {
		machineInfo := machines[machineId]
		if machineInfo.Status == types.Suspect {
			fmt.Print("Machine ID: " + machineId + " ---- ")
			fmt.Print("Status: " + machineInfo.Status + " ---- ")
			fmt.Print("Timestamp: " + strconv.FormatInt(machineInfo.Timestamp, 10) + " ---- ")
			fmt.Print("Incarnation Number: " + strconv.Itoa(machineInfo.Incarnation) + " ---- ")
		}
	}
}

func PrintMachinesMap(machines map[string]types.MachineInfo) {
	var machineIds []string
	for machineId := range machines {
		machineIds = append(machineIds, machineId)
	}
	sort.Strings(machineIds)

	for _, machineId := range machineIds {
		machineInfo := machines[machineId]
		fmt.Print("Machine ID: " + machineId + " ---- ")
		fmt.Print("Status: " + machineInfo.Status + " ---- ")
		fmt.Print("Timestamp: " + strconv.FormatInt(machineInfo.Timestamp, 10) + " ---- ")
		fmt.Print("Incarnation Number: " + strconv.Itoa(machineInfo.Incarnation) + " ---- ")
		fmt.Println()
	}
}
