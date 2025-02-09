package executor

import (
	"MP2/disseminator"
	"MP2/helpers"
	"MP2/logger"
	"MP2/pingAck"
	"MP2/types"
	"MP2/utils"
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

/*
*

	StartCommandOutput() display text to indicate the start of the command output

*
*/
func StartCommandOutput() {
	fmt.Println("################# EXECUTING ####################")
}

/*
*

	EndCommandOutput() display text to indicate the end of the command output

*
*/
func EndCommandOutput() {
	fmt.Println("################# COMPLETED ####################")
}

/*
*

	BreakLine() display text to indicate the separation of commands

*
*/
func BreakLine() {
	fmt.Println("------------------------------------------------")
}

/*
*

	HelpMenu() displays the help menu

*
*/
func HelpMenu() {
	fmt.Println("Welcome to the MP2 CLI")
	fmt.Println("Commands:")
	fmt.Println("1. Join Group")
	fmt.Println("2. Leave Group")
	fmt.Println("3. List all Machines in this Group")
	fmt.Println("4. Display this Machine's ID")
	fmt.Println("5. Enable Suspect")
	fmt.Println("6. Disable Suspect")
	fmt.Println("7. Display Suspect Status")
	fmt.Println("8. Display Suspected nodes")
	fmt.Println("9. Help Menu")
	fmt.Println("10. Network Stats")
	fmt.Println("11. Introduce Packet Drop")
	fmt.Println("Enter Command:")
}

/*
*

	ListAllMachines() lists all the machines in the group

*
*/
func ListAllMachines() {
	StartCommandOutput()
	fmt.Println("List of Machines in this Group:")
	BreakLine()
	machines := helpers.GetMachines()

	// Sort the machines by their machine ID
	utils.PrintMachinesMap(machines)
	EndCommandOutput()
}

/*
*

	DisplaySuspectedNodes() lists all the suspected machines in the group

*
*/
func DisplaySuspectedNodes() {
	StartCommandOutput()
	fmt.Println("List of Suspected Machines in this Group:")
	BreakLine()
	machines := helpers.GetMachines()

	// Sort the machines by their machine ID
	utils.PrintSuspectMachinesMap(machines)
	EndCommandOutput()
}

/*
*

	DisplayMachineID() displays the machine ID, which is the hostname of the machine

*
*/
func DisplayMachineID() {
	StartCommandOutput()
	fmt.Println("Machine ID: " + disseminator.MachineId)
	EndCommandOutput()
}

/*
*

	JoinGroup() join the group by sending a request to the introducer

*
*/
func JoinGroup() {
	fmt.Println("Joining Group...")
	go disseminator.RequestToJoinGroup()
	go pingAck.PingAck()
}

/*
*

	DisplaySuspectFlag() displays the suspect flag

*
*/
func DisplaySuspectFlag() {
	StartCommandOutput()
	fmt.Println("Suspect Flag: ", types.Suspicion)
	EndCommandOutput()
}

/*
*

	toggleSuspect() toggles the suspect flag

*
*/
func toggleSuspect(suspect bool) {
	pingAck.UpdateSuspicion(suspect)
}

func IntroducePacketDrop() {
	var cmd string
	fmt.Print("Enter Drop Probability: ")
	fmt.Scanln(&cmd)
	dropProbability, err := strconv.ParseFloat(cmd, 64)
	if err != nil {
		fmt.Println("Invalid Probability")
		return
	}
	pingAck.DropProbability = dropProbability
}

/*
*

	ExecuteCommand() takes in a command and executes it

*
*/
func ExecuteCommand(cmd string) string {
	switch cmd {
	case "1": // Join Group
		JoinGroup()
	case "2": // Leave Group
		fmt.Println("Leaving Group...")
		fmt.Println("Goodbye! at " + strconv.FormatInt(time.Now().UnixMilli(), 10))
		return "exit"
	case "3": // List all Machines in this Group
		ListAllMachines()
	case "4": // Display this Machine's ID
		DisplayMachineID()
	case "5": // Enable Suspect
		toggleSuspect(true)
		fmt.Println("Enabled Suspect")
	case "6": // Disable Suspect
		toggleSuspect(false)
		fmt.Println("Disabled Suspect")
	case "7": // Display Suspect Status
		DisplaySuspectFlag()
	case "8": // Display Suspected Nodes
		DisplaySuspectedNodes()
	case "9": // Help Menu
		HelpMenu()
	case "10": // Network Stats
		pingAck.NetworkStats()
	case "11":
		IntroducePacketDrop()
	case "c":
		exec.Command("clear")
	default:
		fmt.Println("Invalid Command")
	}
	return "ok"
}

/*
*

	InputCommand() takes in user input and executes the command

*
*/
func InputCommand() {
	HelpMenu()
	logger.InitializeLog()
	for {
		var cmd string
		fmt.Scanln(&cmd)
		status := ExecuteCommand(cmd)
		if status == "exit" {
			break
		}
	}
	defer logger.File.Close()
}
