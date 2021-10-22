package alpaca

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
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
		log.Fatalf("[GRAZE-URL]: The url constructed was not a valid URI, check LLAMA_SERVER, %s", err)
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
	collector := exec.Command("collector", "-llama.config /opt/alpaca/config.yaml")

	collector.Stdout = os.Stdout

	err := collector.Start()

	if err != nil {
		log.Fatalf("[LLAMA-COLLECTOR]: There was an error starting the collector, %s", err)
	}

	g.CollectorPID = collector.Process.Pid
}

func (g *LamoidEnv) GrazeConfig() {
	// Fetch Config write to yaml on local host
	configReqURL, err := url.ParseRequestURI(fmt.Sprintf("%s/api/v1/config/%s", g.ServerURL, g.Group))
	if err != nil {
		log.Fatalf("[CONFIG-URL]: The url constructed was not a valid URI, check LLAMA_SERVER & LLAMA_GROUP , %s", err)
	}

	configReqParam := url.Values{}
	configReqParam.Add("llamaport", fmt.Sprintf("%v", g.Port))

	request, err := http.NewRequest("GET", configReqURL.String(), strings.NewReader(configReqParam.Encode()))
	if err != nil {
		log.Printf("[CONFIG-CLIENT]: There was a problem creating a new request object, %s", err)
	}

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	response, err := client.Do(request)
	if err != nil {
		log.Printf("[CONFIG-CLIENT]: There was a problem making a request to LLAMA Server, %s", err)
	}

	defer func() {
		err := response.Body.Close()

		if err != nil {
			log.Printf("[CONFIG-CLIENT]: There was a problem closing the config response from LLAMA Server, %s", err)
		}
	}()

	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("[CONFIG-CLIENT]: There was a problem reading the config response from LLAMA_SERVER, %s", err)
	}

	//Need to really play around with this to see if we can preserve the YAML formatting
	//returned from LLAMA
	configRaw := LLamaConfig{}

	yamlErr := yaml.Unmarshal(respBytes, &configRaw)
	if yamlErr != nil {
		log.Printf("[YAML-ERR]: There was a problem reading the raw configuration, %s", err)
	}

	yamlData, err := yaml.Marshal(&configRaw)
	if err != nil {
		log.Printf("[YAML-ERR]: There was a problem searializing YAML data into bytes, %s", err)
	}

	//Write configuration to local node
	ioErr := ioutil.WriteFile("/opt/alpaca/config.yaml", yamlData, 0644)

	if ioErr != nil {
		log.Fatalf("[IO-CONFIG]: There was a problem writing data to the config file, %s", ioErr)
	}

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
