package helpers

import (
	"MP2/disseminator"
	"MP2/types"
	"time"
)

/*
*

	GetMachines() allows for thread-safe reading of the Machines map.

*
*/
func GetMachines() map[string]types.MachineInfo {
	disseminator.Mu.RLock() // Acquire a read lock
	defer disseminator.Mu.RUnlock()
	// Return a copy of the map to avoid concurrent modification issues
	copyMachines := make(map[string]types.MachineInfo)
	for k, v := range disseminator.Machines {
		copyMachines[k] = v
	}
	return copyMachines
}

/*
*

	GetChangesMapMachines() allows for thread-safe reading of the Machines map.

*
*/
func GetChangesMapMachines() map[string]types.MachineInfo {
	disseminator.CMLock.RLock() // Acquire a read lock
	defer disseminator.CMLock.RUnlock()
	// Return a copy of the map to avoid concurrent modification issues
	copyMachines := make(map[string]types.MachineInfo)
	for k, v := range disseminator.ChangesMap {
		copyMachines[k] = v
	}
	return copyMachines
}

/*
*

	ReduceTTLForChangesMap() reduces the TTL of the machines in the ChangesMap by 1

*
*/
func ReduceTTLForChangesMap() {
	disseminator.CMLock.Lock()
	var deleteKeys []string
	for k, v := range disseminator.ChangesMap {
		if v.TTL > 0 {
			v.TTL = v.TTL - 1
			disseminator.ChangesMap[k] = v
		}
		if v.TTL == 0 {
			deleteKeys = append(deleteKeys, k)
		}
	}
	for _, k := range deleteKeys {
		delete(disseminator.ChangesMap, k)
	}
	disseminator.CMLock.Unlock()
}

/*
*

	CompareAndUpdateMap() compares the new map with the current map and updates the current map with the new map
	based on few conditions
	This is a thread-safe operation on the Machines map.

*
*/
func CompareAndUpdateMap(kMachinesMap map[string]types.MachineInfo) {
	disseminator.Mu.Lock()
	disseminator.CMLock.Lock()

	for k, v := range kMachinesMap {
		if currentV, ok := disseminator.Machines[k]; !ok {
			if _, ok := disseminator.DeadMachines[k]; !ok {
				if v.Status != types.Dead {
					v.Timestamp = time.Now().UnixMilli()
					disseminator.Machines[k] = v
					v.TTL = 10
					disseminator.ChangesMap[k] = v
				}
			}
		} else {
			if v.Status != currentV.Status {
				if v.Incarnation > currentV.Incarnation || v.Status != types.Alive {
					// Check if k is machineid
					// increment inc++, status = alive
					// else
					// TO-DO: what happens if a node was suspected during the Ping-Ack+S and now the mode is suspect off.
					if k == disseminator.MachineId && v.Status == types.Suspect {
						v.Incarnation = v.Incarnation + 1
						v.Status = types.Alive
					}
					if currentV.Status != types.Dead {
						v.Timestamp = time.Now().UnixMilli()
						disseminator.Machines[k] = v
						v.TTL = 10
						disseminator.ChangesMap[k] = v
					}
				}
			}
		}
	}
	disseminator.Mu.Unlock()
	disseminator.CMLock.Unlock()
}

/*
*

	UpdateMachineStatus() updates the status of a machine in the Machines map
	This is a thread-safe operation on the Machines map.

*
*/
func UpdateMachineStatus(machineId string, status string) {
	disseminator.Mu.Lock()
	disseminator.CMLock.Lock()
	if disseminator.Machines[machineId].Status == types.Alive {
		disseminator.Machines[machineId] = types.MachineInfo{
			Status:      status,
			Timestamp:   time.Now().UnixMilli(),
			Incarnation: disseminator.Machines[machineId].Incarnation,
			TTL:         disseminator.Machines[machineId].TTL,
		}

		disseminator.ChangesMap[machineId] =
			types.MachineInfo{
				Status:      disseminator.Machines[machineId].Status,
				Timestamp:   disseminator.Machines[machineId].Timestamp,
				Incarnation: disseminator.Machines[machineId].Incarnation,
				TTL:         10,
			}
	}
	disseminator.CMLock.Unlock()
	disseminator.Mu.Unlock()
}

func GetXAliveMachines(x int) map[string]types.MachineInfo {
	disseminator.Mu.RLock()
	defer disseminator.Mu.RUnlock()
	aliveMachines := make(map[string]types.MachineInfo)
	for k, v := range disseminator.Machines {
		if v.Status == types.Alive {
			aliveMachines[k] = v
		}
		if len(aliveMachines) == x {
			break
		}
	}
	return aliveMachines
}

func GetXDeadMachines(x int) map[string]types.MachineInfo {
	disseminator.Mu.RLock()
	defer disseminator.Mu.RUnlock()
	deadMachines := make(map[string]types.MachineInfo)
	for k, v := range disseminator.Machines {
		if v.Status == types.Dead {
			deadMachines[k] = v
		}
		if len(deadMachines) == x {
			break
		}
	}
	return deadMachines
}

func GetXSuspectMachines(x int) map[string]types.MachineInfo {
	disseminator.Mu.RLock()
	defer disseminator.Mu.RUnlock()
	suspectMachines := make(map[string]types.MachineInfo)
	for k, v := range disseminator.Machines {
		if v.Status == types.Suspect {
			suspectMachines[k] = v
		}
		if len(suspectMachines) == x {
			break
		}
	}
	return suspectMachines
}
