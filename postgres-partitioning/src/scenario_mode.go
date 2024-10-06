package src

import "fmt"

type ScenarioMode int

const (
	Simple ScenarioMode = iota
	Partition
	DropPartition
)

func (m *ScenarioMode) String() string {
	switch *m {
	case Simple:
		return "simple"
	case Partition:
		return "partition"
	case DropPartition:
		return "drop-partition"
	default:
		return fmt.Sprintf("Mode(%d)", m)
	}
}

func (m *ScenarioMode) Set(s string) error {
	switch s {
	case "simple":
		*m = Simple
	case "partition":
		*m = Partition
	case "drop-partition":
		*m = DropPartition
	default:
		return fmt.Errorf("invalid mode: %q", s)
	}
	return nil
}
