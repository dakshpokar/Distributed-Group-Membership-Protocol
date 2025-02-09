package types

// MachineInfo is a struct that stores the information of a machine
type MachineInfo struct {
	Status      string
	Incarnation int
	Timestamp   int64
	TTL         int
}
