package config

var (
	DebugMode = false
)

func SetDebugMode(mode bool) {
	DebugMode = mode
}

func GetDebugMode() bool{
	return DebugMode
}