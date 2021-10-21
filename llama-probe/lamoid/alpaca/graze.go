package alpaca

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func (g *LamoidEnv) ValidateEnvironment() {
	// Validate Server Configuration
}

func (g *LamoidEnv) StartReflector() {
	// Start llama reflector and update the process id ref.
	reflector := exec.Command("reflector", fmt.Sprintf("-port %v", g.Port))

	reflector.Stdout = os.Stdout

	err := reflector.Start()

	if err != nil {
		log.Fatalf("[LLAMA-REFLECTOR]: There was an error starting the reflector, %s", err)
	}

	g.ReflectorPID = reflector.Process.Pid
}

func (g *LamoidEnv) StartCollector() {
	// Start llama collector and update the process id ref.
}

func (g *LamoidEnv) GetConfig() {
	// Fetch Config write to yaml on local host
}

func (g *LamoidEnv) ValidateConfig() {
	// Validate Running config Against Fetched config
}

func (g *LamoidEnv) NewServerUrl() {
	// Construct Server URL update Env ref.

}

func (g *LamoidEnv) Graze() {
	// Main Loop for running the llama-probe
}
