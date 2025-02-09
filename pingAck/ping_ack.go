package pingAck

import (
	"MP2/disseminator"
	"MP2/helpers"
	"MP2/logger"
	"MP2/types"
	"MP2/utils"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	Bytes     int64
	StartTime time.Time
	BrLock    sync.RWMutex
)

var PingCount = 0
var AckCount = 0

var DropProbability = 0.0

func PrintLog(message string) {
	// fmt.Println("Ping_Ack | " + message)
	logger.LogToFile("Ping_Ack | " + message + "\n")
}

func PingAck() {
	var wg sync.WaitGroup
	wg.Add(2)
	// fmt.Println("in PingAck")
	go Server(&wg)
	time.Sleep(2 * time.Second)
	go Client(&wg)
	wg.Wait()
	// fmt.Println("after client")
	time.Sleep(2 * time.Second)
}

func Server(wg *sync.WaitGroup) {
	defer wg.Done()
	PrintLog("Server | Starting server")
	addr := net.UDPAddr{
		Port: 1201,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp4", &addr)
	if err != nil {
		PrintLog("Server | Error in listenUDP " + err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	StartTime = time.Now()

	numReaders := 10
	var wg_2 sync.WaitGroup
	for i := 0; i < numReaders; i++ {
		wg_2.Add(1)
		go ListenAndServe(conn, &wg_2)
	}
	wg_2.Wait()
	fmt.Println("Server | Server exiting")
}

func ListenAndServe(conn *net.UDPConn, wg *sync.WaitGroup) {
	defer wg.Done()
	buf := make([]byte, 2048)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)

		BrLock.Lock()
		Bytes += int64(n)
		BrLock.Unlock()

		if err != nil {
			PrintLog("ListenAndServe | server Error in Reading from UDP" + err.Error())
		}
		go HandleConnect(conn, remoteAddr, buf[:n])
	}
}

func HandleConnect(conn *net.UDPConn, remoteAddr *net.UDPAddr, data []byte) {
	request := utils.JSONBytesToRequest(data)
	if request.Req_type == types.Suspect {
		types.Su.Lock()
		types.Suspicion, _ = strconv.ParseBool(request.Data)
		types.Su.Unlock()
		return
	}
	//// fmt.Println(request)

	PrintLog("HandleConnect | Data recvd from :" + remoteAddr.String())
	data, err := json.Marshal(helpers.GetChangesMapMachines())
	if err != nil {
		PrintLog("HandleConnect | Error in server json Marshal:" + err.Error())
	}

	// rand.Seed(time.Now().UnixNano())
	// randomNumber := rand.Intn(100) + 1

	// if float64(randomNumber) <= DropProbability {
	// 	PrintLog("HandleConnect | Dropping packet")
	// 	return
	// }

	_, err = conn.WriteToUDP(
		utils.ResponseToJSONBytes(
			types.Response{
				Resp_type: types.Ack,
				Data:      string(data),
			},
		),
		remoteAddr,
	)

	AckCount++

	Machines, _ := utils.BytesToMachineMap([]byte(request.Data))
	PrintLog("HandleConnect | Data recvd from :" + fmt.Sprint(Machines))
	PrintLog("HandleConnect | Machines :" + fmt.Sprint(helpers.GetMachines()))

	helpers.CompareAndUpdateMap(Machines)

	PrintLog("HandleConnect | Updated Machines after Compare :" + fmt.Sprint(helpers.GetMachines()))

	if err != nil {
		PrintLog("Error in writing server response" + err.Error())
	}
}

func SelectStatus(machine string) {
	types.Su.RLock()
	if types.Suspicion {
		helpers.UpdateMachineStatus(machine, types.Suspect)
		fmt.Println("##########################################################")
		fmt.Println("Machine " + machine + " is a Suspected Machine!!!")
		fmt.Println("##########################################################")
	} else {
		helpers.UpdateMachineStatus(machine, types.Dead)
		fmt.Println("##########################################################")
		fmt.Println("Machine " + machine + " is a Dead Machine!!!")
		fmt.Println("##########################################################")
	}
	types.Su.RUnlock()

}

func Client(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		machineMap := helpers.GetMachines()

		PrintLog("Client | Updates Map: " + fmt.Sprint(helpers.GetChangesMapMachines()))

		// Get Keys from the map
		var machineMapKeys []string
		for k := range machineMap {
			machineMapKeys = append(machineMapKeys, k)
		}

		// Shuffle them randomly
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(machineMapKeys), func(i, j int) { machineMapKeys[i], machineMapKeys[j] = machineMapKeys[j], machineMapKeys[i] })

		for i := range machineMapKeys {
			// totalChecks += 1
			hostName := utils.GetHostName(machineMapKeys[i])
			if machineMapKeys[i] == disseminator.MachineId || machineMap[machineMapKeys[i]].Status == types.Dead {
				continue
			}

			hostName = hostName + ":1201"

			PrintLog("Client | Sending client data:" + hostName)

			udpAddr, err := net.ResolveUDPAddr("udp4", hostName)
			if err != nil {
				PrintLog("Client | Error in resolve udp addr:" + err.Error())
				SelectStatus(machineMapKeys[i])
			}
			conn, err := net.DialUDP("udp4", nil, udpAddr)
			if err != nil {
				PrintLog("Client | Error in dial UDP:" + err.Error())
			}

			PrintLog("Client | Sending data to:" + hostName)

			data, err := json.Marshal(helpers.GetChangesMapMachines())

			if err != nil {
				PrintLog("Client | Error in client json Marshal:" + err.Error())
			}
			_, err = conn.Write(
				utils.RequestToJSONBytes(
					types.Request{
						Req_type: "ping",
						Data:     string(data),
					},
				),
			)
			PingCount++

			if err != nil {
				PrintLog("Client | Error in UDP write in client:" + err.Error())
				SelectStatus(machineMapKeys[i])
			}
			helpers.ReduceTTLForChangesMap()
			err = conn.SetReadDeadline(time.Now().Add(1000 * time.Millisecond))
			if err != nil {
				PrintLog("Client | Error in setting read Deadline:" + err.Error())
			}
			buffer := make([]byte, 2048)
			n, _, err := conn.ReadFromUDP(buffer)

			if err != nil {
				PrintLog("Client | Error in Reading from UDP:" + err.Error())
				fmt.Println("Client | Error in Reading from UDP:" + err.Error())
				SelectStatus(machineMapKeys[i])
			}

			BrLock.Lock()
			Bytes += int64(n)
			BrLock.Unlock()

			request := utils.JSONBytesToRequest(buffer[:n])

			kMachines, _ := utils.BytesToMachineMap([]byte(request.Data))

			if len(kMachines) != 0 {
				PrintLog("Client | Data recvd from :" + fmt.Sprint(kMachines))
				PrintLog("Client | Machines :" + fmt.Sprint(helpers.GetMachines()))

				helpers.CompareAndUpdateMap(kMachines)

				PrintLog("Client | Updated Machines after Compare :" + fmt.Sprint(helpers.GetMachines()))
			}

		}
		time.Sleep(2 * time.Second)
	}
}

func UpdateSuspicion(suspect bool) {
	machineMap := helpers.GetMachines()

	// Get Keys from the map
	var machineMapKeys []string
	for k := range machineMap {
		machineMapKeys = append(machineMapKeys, k)
	}

	for i := range machineMapKeys {
		hostName := utils.GetHostName(machineMapKeys[i])
		if machineMap[machineMapKeys[i]].Status == types.Dead {
			continue
		}

		hostName = hostName + ":1201"
		// fmt.Println("----------Sending client data:", hostName)
		udpAddr, err := net.ResolveUDPAddr("udp4", hostName)
		if err != nil {
			PrintLog("UpdateSuspicion | Error in resolve udp addr:" + err.Error())
		}
		conn, err := net.DialUDP("udp4", nil, udpAddr)
		if err != nil {
			PrintLog("UpdateSuspicion | Error in dial UDP:" + err.Error())
		}
		data := strconv.FormatBool(suspect)
		if err != nil {
			PrintLog("UpdateSuspicion | Error in client json Marshal:" + err.Error())
		}
		_, err = conn.Write(
			utils.RequestToJSONBytes(
				types.Request{
					Req_type: types.Suspect,
					Data:     string(data),
				},
			),
		)
		if err != nil {
			PrintLog("UpdateSuspicion | Error in UDP write in client:" + err.Error())
		}
	}
}
