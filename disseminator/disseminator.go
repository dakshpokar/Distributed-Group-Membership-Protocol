package disseminator

import (
	"MP2/logger"
	"MP2/types"
	"MP2/utils"
	"strconv"
	"time"

	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

/*
*

	Machines is a map that stores the information of all the machines in the group
	This is a global variable that is accessed by multiple goroutines
	This is specifically used by the Disseminator to keep track of all the machines in the group

	Do not directly access this variable, use the functions in the helpers package to access this variable
	Functions in the helpers package are thread-safe and will handle the locking of the Machines map

*
*/
var (
	Machines = make(map[string]types.MachineInfo)
	Mu       sync.RWMutex
)

var (
	ChangesMap = make(map[string]types.MachineInfo)
	CMLock     sync.RWMutex
)

var (
	DeadMachines = make(map[string]types.MachineInfo)
	DMLock       sync.RWMutex
)

var MachineId = ""

var FalsePositives = 0

/*
*

	PrintLog() prints the log message with the Disseminator tag

*
*/
func PrintLog(message string) {
	// fmt.Println("Disseminator | " + message)
	logger.LogToFile("Disseminator | " + message + "\n")
}

/*
*

	StartCleanupThread() starts a cleanup thread that removes dead machines from the group

*
*/
func StartCleanupThread() {
	PrintLog("StartCleanupThread | Starting cleanup thread...")
	for {
		time.Sleep(10 * time.Second)
		Mu.Lock()
		for key, value := range Machines {
			if value.Status == types.Dead {
				PrintLog("StartCleanupThread | Removing dead machine: " + key)
				delete(Machines, key)
				DeadMachines[key] = value
			}
		}
		Mu.Unlock()
	}
}

/*
*

	StartSuspectCleanupThread() starts a cleanup thread that converts suspected machines to dead from the group

*
*/
func StartSuspectCleanupThread() {
	PrintLog("StartSuspectCleanupThread | Starting cleanup thread...")
	for {
		time.Sleep(3000 * time.Millisecond)
		types.Su.RLock()
		if types.Suspicion {
			Mu.RLock()
			var toUpdateKeys []string
			for key, value := range Machines {
				if value.Status == types.Suspect && value.Timestamp < time.Now().UnixMilli()-3000 {
					toUpdateKeys = append(toUpdateKeys, key)
				}
			}
			Mu.RUnlock()
			if len(toUpdateKeys) > 0 {
				Mu.Lock()
				for _, key := range toUpdateKeys {
					value := Machines[key]
					value.Status = types.Dead
					PrintLog("StartSuspectCleanupThread | Converting suspect to dead machine: " + key)
					Machines[key] = value
				}
				Mu.Unlock()
			}
		}
		types.Su.RUnlock()
	}
}

/*
*

	JoinGroup() allows Disseminator makes a call to Introducer to join the group

*
*/
func JoinGroup() {
	introducerAddr := os.Getenv("INTRODUCER_ADDR") + ":6911"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", introducerAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	PrintLog("JoinGroup | Connected to introducer at " + introducerAddr)
	conn.Write(
		utils.RequestToJSONBytes(
			types.Request{
				Req_type: "join",
				Data:     conn.LocalAddr().String(),
			},
		),
	)

	HandleJoinRequest(conn)
}

/*
*

	HandleJoinRequest() handles the response from the Introducer for the join request made by Disseminator

*
*/
func HandleJoinRequest(conn net.Conn) {
	for {
		PrintLog("HandleJoinRequest | Waiting for response")
		reader := bufio.NewReader(conn)
		req, err := reader.ReadBytes('\r')
		request := utils.JSONBytesToRequest(req)
		if err != nil {
			PrintLog("HandleJoinRequest | Error reading from connection:" + err.Error())
			return
		}
		PrintLog("HandleJoinRequest | Received response: " + request.Req_type + " from " + request.Data)
		if request.Req_type == "initiate_join" {
			Mu.Lock()
			Machines, err = utils.BytesToMachineMap(request.Byte_Data)
			MachineId = request.Data
			Mu.Unlock()
			types.Su.Lock()
			types.Suspicion, _ = strconv.ParseBool(request.Meta_Data)
			types.Su.Unlock()

			go StartCleanupThread()
			go StartSuspectCleanupThread()
			if err != nil {
				PrintLog("HandleJoinRequest | Error converting bytes to machines map: " + err.Error())
				return
			}
			PrintLog("HandleJoinRequest | Added machine to group, machines: " + fmt.Sprintf("%v", Machines))
			break
		}
	}
	defer conn.Close()
}

/*
*

	RequestToJoinGroup() sends a request to the Introducer to join the group

*
*/
func RequestToJoinGroup() {
	if os.Getenv("INTRODUCER_ADDR") == "" {
		PrintLog("RequestToJoinGroup | Please set the INTRODUCER_ADDR environment variable")
		os.Exit(1)
	}
	JoinGroup()
}
