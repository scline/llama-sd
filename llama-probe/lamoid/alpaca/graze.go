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
	"syscall"
	"time"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v2"
)

//GrazeAnatomy - A method called on LamoidEnv which registers the current running LLAMA configuration
//to the LLAMA-SERVER
func (g *LamoidEnv) GrazeAnatomy() {

	//Build the registration Payload
	lamoidAnatomy := &PayLoad{
		Port:      g.Port,
		Keepalive: g.KeepAlive,
		Ip:        g.SourceIP,
		Group:     g.Group,
	}

	//TODO: Read in version number
	lamoidAnatomy.Tags.Version = "0.1.0"
	lamoidAnatomy.Tags.ProbeName = g.ProbeName
	lamoidAnatomy.Tags.ProbeShortname = g.ProbeShortName

	byteArray, err := json.Marshal(lamoidAnatomy)
	if err != nil {
		log.Println(err)
	}

	//Build and Validate the LLAMA-SERVER url
	serverURL, err := url.ParseRequestURI(fmt.Sprintf("%sapi/v1/register", g.Server))
	if err != nil {
		log.Printf("[LAMOID-REGISTER]: The url constructed was not a valid URI, check LLAMA_SERVER, %s", err)
	}

	//Build the HTTP Post request
	request, err := http.NewRequest("POST", serverURL.String(), bytes.NewBuffer(byteArray))

	if err != nil {
		log.Printf("[LAMOID-REGISTER]: There was a problem creating a new request object, %s", err)
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	//HTTP Client
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	//Process HTTP request and log the status
	response, err := client.Do(request)

	if err != nil {
		log.Printf("[LAMOID-REGISTER]: There was a problem making a request, %s", err)
	}

	defer func() {
		err := response.Body.Close()

		if err != nil {
			log.Printf("[LAMOID-REGISTER]: There was a problem closing the response from LLAMA Server, %s", err)
		}
	}()

	log.Printf("[LAMOID-REGISTER]: Response Status: %s", response.Status)
}

//StartReflector - A method called on LamoidEnv which starts the LLAMA Reflector application and updates LamoidEnv with
//its OS process identification (PID)
func (g *LamoidEnv) StartReflector() {

	// Build os exec command to launch reflector with a given param
	reflector := exec.Command("reflector", fmt.Sprintf("-port %v", g.Port))

	// Set the process to output to Standard Out
	reflector.Stdout = os.Stdout

	// Execute the exec command to start reflector, panic on error.
	log.Print("[LAMOID]: Starting Reflector")
	err := reflector.Start()

	if err != nil {
		log.Printf("[LAMOID-REFLECTOR]: There was an error starting the reflector, %s", err)
	}

	//Set PID
	g.ReflectorPID = reflector.Process.Pid
}

//StartCollector - A method called on LamoidEnv which starts the LLAMA Collector application and updates LamoidEnv with
//its OS process identification (PID)
func (g *LamoidEnv) StartCollector() {

	// Build os exec command to launch colelctor with a given param
	collector := exec.Command("collector", "-llama.config /opt/alpaca/config.yaml")

	// Set the process to output to Standard Out
	collector.Stdout = os.Stdout

	// Execute the exec command to start colelctor, panic on error.
	log.Print("[LAMOID]: Starting Collector")
	err := collector.Start()

	if err != nil {
		log.Printf("[LAMOID-COLLECTOR]: There was an error starting the collector, %s", err)
	}

	//Set PID
	g.CollectorPID = collector.Process.Pid
}

//GrazeConfig - A method called on LamoidEnv which retrieves the running configuration from the LLAMA Servers configuration
//API. Writes that config as a YAML to the local node for consumption by the collector process. Must be ran before the
//collector is started.
func (g *LamoidEnv) GrazeConfig() []byte {

	// Build and validate URL
	configReqURL, err := url.ParseRequestURI(fmt.Sprintf("%s/api/v1/config/%s", g.Server, g.Group))
	if err != nil {
		log.Printf("[LAMOID-URL]: The url constructed was not a valid URI, check LLAMA_SERVER & LLAMA_GROUP , %s", err)
	}

	// Configure request url params
	configReqParam := url.Values{}
	configReqParam.Add("llamaport", fmt.Sprintf("%v", g.Port))

	// Build request
	request, err := http.NewRequest("GET", configReqURL.String(), strings.NewReader(configReqParam.Encode()))
	if err != nil {
		log.Printf("[LAMOID-CLIENT]: There was a problem creating a new request object, %s", err)
	}

	//HTTP Client
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	// Process HTTP request
	response, err := client.Do(request)
	if err != nil {
		log.Printf("[LAMOID-CLIENT]: There was a problem making a request to LLAMA Server, %s", err)
	}

	defer func() {
		err := response.Body.Close()

		if err != nil {
			log.Printf("[LAMOID-CLIENT]: There was a problem closing the config response from LLAMA Server, %s", err)
		}
	}()

	// Read response into bytes
	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("[LAMOID-CLIENT]: There was a problem reading the config response from LLAMA_SERVER, %s", err)
	}

	return respBytes

}

func (g *LamoidEnv) WriteConfig(respBytes []byte) {

	//Serialize YAML Data into a struct
	configRaw := LLamaConfig{}

	yamlErr := yaml.Unmarshal(respBytes, &configRaw)
	if yamlErr != nil {
		log.Printf("[LAMOID-ERR]: There was a problem reading the raw configuration, %s", yamlErr)
	}

	yamlData, err := yaml.Marshal(&configRaw)
	if err != nil {
		log.Printf("[LAMOID-ERR]: There was a problem searializing YAML data into bytes, %s", err)
	}

	//Write configuration to local node
	ioErr := os.WriteFile("/opt/alpaca/config.yaml", yamlData, 0644)

	if ioErr != nil {
		log.Printf("[LAMOID-ERR]: There was a problem writing data to the config file, %s", ioErr)
	}

}

func (g *LamoidEnv) ReadConfig() LLamaConfig {

	configRaw := LLamaConfig{}

	yamlFile, err := os.ReadFile("/opt/alpaca/config.yaml")
	if err != nil {
		log.Printf("[LAMOID-ERR]: There was a problem reading the config file, %s", err)
	}

	err = yaml.Unmarshal(yamlFile, &configRaw)
	if err != nil {
		log.Printf("[LAMOID-ERR]: There was a problem reading config into type struct, %s", err)
	}

	return configRaw

}

func (g *LamoidEnv) ValidateConfig() string {
	// Validate Running config Against Fetched config
	//Code duplication for now, will break out later after initial testing and refactoring.

	// Build and validate URL
	configReqURL, err := url.ParseRequestURI(fmt.Sprintf("%s/api/v1/config/%s", g.Server, g.Group))
	if err != nil {
		log.Fatalf("[LAMOID-URL]: The url constructed was not a valid URI, check LLAMA_SERVER & LLAMA_GROUP , %s", err)
	}

	// Configure request url params
	configReqParam := url.Values{}
	configReqParam.Add("llamaport", fmt.Sprintf("%v", g.Port))

	// Build request
	request, err := http.NewRequest("GET", configReqURL.String(), strings.NewReader(configReqParam.Encode()))
	if err != nil {
		log.Printf("[LAMOID-CLIENT]: There was a problem creating a new request object, %s", err)
	}

	//HTTP Client
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	// Process HTTP request
	response, err := client.Do(request)
	if err != nil {
		log.Printf("[LAMOID-CLIENT]: There was a problem making a request to LLAMA Server, %s", err)
	}

	defer func() {
		err := response.Body.Close()

		if err != nil {
			log.Printf("[LAMOID-CLIENT]: There was a problem closing the config response from LLAMA Server, %s", err)
		}
	}()

	// Read response into bytes
	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("[LAMOID-CLIENT]: There was a problem reading the config response from LLAMA_SERVER, %s", err)
	}

	//Note: Need to really play around with this to see if we can preserve the YAML formatting
	//returned from LLAMA

	//Serialize YAML Data into a struct
	newConfigRaw := LLamaConfig{}

	yamlErr := yaml.Unmarshal(respBytes, &newConfigRaw)
	if yamlErr != nil {
		log.Printf("[LAMOID-ERR]: There was a problem reading the raw configuration from llama server, %s", err)
	}

	//TODO: Read current running config. Validate
	currentConfigRaw := g.ReadConfig()

	//TODO: Compare both YAML files Validate
	var lamoidAction string

	if !cmp.Equal(currentConfigRaw, newConfigRaw) {
		log.Printf("[LAMOID-INFO]: New Configuration Detected")
		lamoidAction = "Update"
	} else {
		log.Printf("[LAMOID-INFO]: No Configuration Detected")
		lamoidAction = "Continue"
	}

	return lamoidAction
}

func (g *LamoidEnv) StartGrazing() {
	//Initial Run
	g.StartReflector()
	g.GrazeAnatomy()

	//Give the LLama sometime to eat....sheeeeeeshhhh
	time.Sleep(time.Second * 10)

	g.WriteConfig(g.GrazeConfig())

	g.StartCollector()
}

func (g *LamoidEnv) Graze() {
	// Main Loop for running the llama-probe

	g.StartGrazing()

Graze:
	for {
		switch g.ValidateConfig() {
		case "Continue":
			continue Graze
		case "Update":
			log.Printf("[LAMOID-INFO]: Stoping Collector")

			err := syscall.Kill(g.CollectorPID, syscall.SIGHUP)
			if err != nil {
				log.Printf("[LAMOID-ERR]: There was a problem trying to send SIGHUP to collector process, %s", err)
			}

			log.Printf("[LAMOID-INFO]: Grazing Registration")

			g.GrazeAnatomy()

			log.Printf("[LAMOID-INFO]: Writing New Config")
			time.Sleep(time.Second * 10)

			g.WriteConfig(g.GrazeConfig())

			log.Printf("[LAMOID-INFO]: Starting New Collector")
			g.StartCollector()

			time.Sleep(time.Second * 30)
			continue Graze
		}
	}

}
