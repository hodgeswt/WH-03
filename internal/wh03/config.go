package wh03

type Config struct {
	ClockFreq           int    `json:"clockFreq"`
	LogLevel            string `json:"logLevel"`
	RamK                int    `json:"ramK"`
	LogFile             string `json:"logFile"`
	RomK                int    `json:"romK"`
	RomFile             string `json:"romFile"`
	ProbeEnabled        bool   `json:"probeEnabled"`
	StackStart          string `json:"stackStart"`
	StepsPerInstruction int    `json:"stepsPerInstruction"`
}
