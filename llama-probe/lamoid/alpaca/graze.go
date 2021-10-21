package alpaca

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"
)

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

	//TODO: Read in version number at compile, add to LamoidEnv struct
	//so its available to all methods.
	lamoidAnatomy.Tags.Version = "0.1.0"
	lamoidAnatomy.Tags.ProbeName = g.ProbeName
	lamoidAnatomy.Tags.ProbeShortname = g.ProbeShortName

	byteArray, err := json.Marshal(lamoidAnatomy)
	if err != nil {
		log.Println(err)
	}

	//Might as well validate that the string is a URL since its comming from the
	//Deployment Environment
	serverURL, err := url.ParseRequestURI(fmt.Sprintf("%s/api/v1/register", g.ServerURL))
	if err != nil {
		log.Fatalf("[GRAZE-URL]: The url constructed for the http client was not a valid URI, check LLAMA_SERVER env, %s", err)
	}

	request, err := http.NewRequest("POST", serverURL.String(), bytes.NewBuffer(byteArray))

	if err != nil {
		log.Printf("[GRAZE-CLIENT]: There was a problem creating a new request object, %s", err)
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	response, err := client.Do(request)

	if err != nil {
		log.Printf("[GRAZE-CLIENT]: There was a problem making a request, %s", err)
	}

	defer func() {
		err := response.Body.Close()

		if err != nil {
			log.Printf("[GRAZE-CLIENT]: There was a problem closing the response from LLAMA Server, %s", err)
		}
	}()

	log.Printf("Response Status: %s", response.Status)
}

func (g *LamoidEnv) StartCollector() {
	// Start llama collector and update the process id ref.
}

func (g *LamoidEnv) GrazeConfig() {
	// Fetch Config write to yaml on local host
}

func (g *LamoidEnv) ValidateConfig() {
	// Validate Running config Against Fetched config
}

func (g *LamoidEnv) Graze() {
	// Main Loop for running the llama-probe

	//Initial Run
	g.StartReflector()
	g.GrazeAnatomy()

	//Give the LLama sometime to eat....sheeeeeeshhhh
	time.Sleep(time.Second * 10)

	g.GrazeConfig()
	g.StartCollector()

	// Do Loop Stuff Later

}
