package introducer

import (
	"MP2/disseminator"
	"MP2/helpers"
	"MP2/logger"
	"MP2/types"
	"MP2/utils"
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
*

	PrintLog() prints the log message with the Introducer tag

*
*/
func PrintLog(message string) {
	// fmt.Println("Introducer | " + message)
	logger.LogToFile("Introducer | " + message + "\n")
}

/*
*

	HandleMachineJoinRequest() handles the join request from a machine and assigns a machine ID to the machine

*
*/
func HandleMachineJoinRequest(conn net.Conn) {
	PrintLog("HandleMachineJoinRequest | Accepted connection from " + conn.RemoteAddr().String())
	for {
		req, err := bufio.NewReader(conn).ReadString('\r')
		if err != nil {
			fmt.Println(err)
			return
		}

		request := utils.JSONBytesToRequest([]byte(req))

		PrintLog("HandleMachineJoinRequest | Received: " + request.Req_type + " from " + request.Data)

		hostName, err := net.LookupAddr(strings.Split(conn.RemoteAddr().String(), ":")[0])
		if err != nil {
			PrintLog("HandleMachineJoinRequest | Error looking up address: " + err.Error())
			return
		}
		PrintLog(strings.Join(hostName, ""))

		machineId := utils.FilteredHostName(hostName[0]) + ":" + strconv.FormatInt(time.Now().UnixMilli(), 10)

		PrintLog("Assigned Machine ID | Machine ID: " + machineId)

		disseminator.Mu.Lock()
		disseminator.CMLock.Lock()
		disseminator.Machines[machineId] = types.MachineInfo{
			Status:      types.Alive,
			Timestamp:   time.Now().UnixMilli(),
			Incarnation: 0,
			TTL:         1,
		}
		disseminator.ChangesMap[machineId] = types.MachineInfo{
			Status:      types.Alive,
			Timestamp:   time.Now().UnixMilli(),
			Incarnation: 0,
			TTL:         10,
		}
		disseminator.CMLock.Unlock()
		disseminator.Mu.Unlock()

		PrintLog("HandleMachineJoinRequest | Added machine to group, machines: " + fmt.Sprintf("%v", disseminator.Machines))

		machines, err := utils.MachineMapToBytes(helpers.GetMachines())

		if err != nil {
			PrintLog("HandleMachineJoinRequest | Error converting machines map to bytes: " + err.Error())
			return
		}

		if request.Req_type == "join" {
			types.Su.RLock()
			conn.Write(
				utils.RequestToJSONBytes(
					types.Request{
						Req_type:  "initiate_join",
						Data:      machineId,
						Meta_Data: strconv.FormatBool(types.Suspicion),
						Byte_Data: machines,
					},
				),
			)
			types.Su.RUnlock()
		}
	}
}

/*
*

	StartIntroducer() starts the introducer and listens for incoming connections on port 6911

*
*/
func StartIntroducer() {
	// Start Introduer at Port 6911
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":6911")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	PrintLog("StartIntroducer | Listening on: " + tcpAddr.String())
	defer listener.Close()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		// Accept new connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		// Handle new connections to Introducer in a Goroutine for concurrency
		go HandleMachineJoinRequest(conn)
	}
}
