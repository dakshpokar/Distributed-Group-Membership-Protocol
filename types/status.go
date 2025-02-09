package types

// Alive, Dead, and Suspect are the possible statuses of a machine
// Do not directly have the statuses as strings in the code, use these constants instead
const (
	Alive   = "alive"
	Dead    = "dead"
	Suspect = "suspect"
	Ack     = "ack"
)
