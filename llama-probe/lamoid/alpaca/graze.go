package alpaca

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
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

func (g *LamoidEnv) GrazeAnatomy() {

	lamoidAnatomy := &PayLoad{
		Port:      g.Port,
		Keepalive: g.KeepAlive,
		Ip:        g.SourceIP,
		Group:     g.Group,
	}

	lamoidAnatomy.Tags.Version = "0.1.0"
	lamoidAnatomy.Tags.ProbeName = g.ProbeName
	lamoidAnatomy.Tags.ProbeShortname = g.ProbeShortName

	// Convert PayLoad struct to JSON now that we have all our values.
	byteArray, err := json.Marshal(lamoidAnatomy)
	if err != nil {
		log.Println(err)
	}

	// POST JSON to API Server
	url := envs["LLAMA_SERVER"] + "/api/v1/register"
	log.Println("Server URL:", url)

	// Print registration JSON we send to API server
	log.Println("JSON Post:", string(byteArray))

	// Loop the registration every 60 seconds
	request, error := http.NewRequest("POST", url, bytes.NewBuffer(byteArray))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{Timeout: 5 * time.Second}
	response, error := client.Do(request)
	if error != nil {
		log.Println(error)
	}
	defer response.Body.Close()

	// Log responce
	log.Println("Response Status:", response.Status)
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
