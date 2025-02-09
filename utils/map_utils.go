package utils

import (
	"MP2/types"
	"encoding/json"
)

/*
*

	MachineMapToBytes() converts a map of Hostname(string) to MachineInfo to a JSON byte array

	Parameters:
		m: a map of Hostname(string) to MachineInfo

	Returns:
		[]byte: a JSON byte array

*
*/
func MachineMapToBytes(m map[string]types.MachineInfo) ([]byte, error) {
	return json.Marshal(m)
}

/*
*

	BytesToMachineMap() converts a JSON byte array to a map of Hostname(string) to MachineInfo

	Parameters:
		b: a JSON byte array

	Returns:
		map[string]types.MachineInfo: a map of Hostname(string) to MachineInfo

*
*/
func BytesToMachineMap(b []byte) (map[string]types.MachineInfo, error) {
	var m map[string]types.MachineInfo
	err := json.Unmarshal(b, &m)
	return m, err
}
