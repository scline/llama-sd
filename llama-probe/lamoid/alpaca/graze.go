package alpaca

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/google/go-cmp/cmp"
)

//GrazeAnatomy - A method called on LamoidEnv which registers the current running LLAMA configuration
//to the LLAMA-SERVER
func (g *LamoidEnv) GrazeAnatomy() {

	log.Print("[LAMOID-REGISTER]: Regestiering With LLAMA Server....")

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
		log.Printf("[URL-ERROR]: The url constructed was not a valid URI, check LLAMA_SERVER, %s", err)
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

	log.Print("[LAMOID-REGISTER]: Regestiering Process Completed")
	log.Printf("[LAMOID-REGISTER]: Response Status: %s", response.Status)
}

//StartReflector - A method called on LamoidEnv which starts the LLAMA Reflector application and updates LamoidEnv with
//its OS process identification (PID)
func (g *LamoidEnv) StartReflector() {

	// Build os exec command to launch reflector with a given param
	reflector := exec.Command("reflector", "-port", fmt.Sprint(g.Port))

	// Set the process to output to Standard Out
	reflector.Stdout = os.Stdout
	reflector.Stderr = os.Stderr

	// Execute the exec command to start reflector, panic on error.
	log.Print("[LAMOID]: Starting Reflector")
	err := reflector.Start()

	if err != nil {
		log.Printf("[LAMOID-REFLECTOR]: There was an error starting the reflector, %s", err)
	}

	go func() {
		err = reflector.Wait()
		if err != nil {
			log.Printf("[ERROR]: Reflector processed didn't close gracfully")
		}
	}()

	//Set PID
	log.Printf("[REFLECTOR-PID]: %v", reflector.Process.Pid)
	g.ReflectorPID = reflector.Process.Pid
}

//StartCollector - A method called on LamoidEnv which starts the LLAMA Collector application and updates LamoidEnv with
//its OS process identification (PID)
func (g *LamoidEnv) StartCollector() {

	// Build os exec command to launch colelctor with a given param
	collector := exec.Command("collector", "-llama.config", "config.yaml")

	// Set the process to output to Standard Out
	collector.Stdout = os.Stdout
	collector.Stderr = os.Stderr

	// Execute the exec command to start colelctor, panic on error.
	log.Print("[LAMOID]: Starting Collector")
	err := collector.Start()

	if err != nil {
		log.Printf("[LAMOID-COLLECTOR]: There was an error starting the collector, %s", err)
	}

	go func() {
		err = collector.Wait()
		if err != nil {
			log.Printf("[ERROR]: Collector processed didn't close gracfully")
		}
	}()

	//Set PID
	log.Printf("[COLLECTOR-PID]: %v", collector.Process.Pid)
	g.CollectorPID = collector.Process.Pid
}

//GrazeConfig - A method called on LamoidEnv which retrieves the running configuration from the LLAMA Servers configuration
//API. Writes that config as a YAML to the local node for consumption by the collector process. Must be ran before the
//collector is started.
func (g *LamoidEnv) GrazeConfig() []byte {

	// Build and validate URL
	configReqURL, err := url.ParseRequestURI(fmt.Sprintf("%sapi/v1/config/%s", g.Server, g.Group))
	if err != nil {
		log.Printf("[LAMOID-URL]: The url constructed was not a valid URI, check LLAMA_SERVER & LLAMA_GROUP , %s", err)
	}

	// Build request
	request, err := http.NewRequest("GET", configReqURL.String(), nil)
	if err != nil {
		log.Printf("[LAMOID-CLIENT]: There was a problem creating a new request object, %s", err)
	}

	configReqQuery := request.URL.Query()
	configReqQuery.Add("llamaport", fmt.Sprint(g.Port))
	request.URL.RawQuery = configReqQuery.Encode()

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

	yamlFile, err := os.Create("config.yaml")
	if err != nil {
		return
	}

	defer func() {
		err = yamlFile.Close()
		if err != nil {
			log.Printf("[YAML-WRITE-ERROR]: %s", err)
		}
	}()

	_, writeErr := yamlFile.Write(respBytes)

	if writeErr != nil {
		log.Printf("[YAML-WRITE-ERROR]: %s", err)
	}

}

func (g *LamoidEnv) WriteTempConfig(respBytes []byte) {

	yamlFile, err := os.Create("tmp-config.yaml")
	if err != nil {
		return
	}

	defer func() {
		err = yamlFile.Close()
		if err != nil {
			log.Printf("[YAML-WRITE-ERROR]: %s", err)
		}
	}()

	_, writeErr := yamlFile.Write(respBytes)

	if writeErr != nil {
		log.Printf("[YAML-WRITE-ERROR]: %s", err)
	}

}

func (g *LamoidEnv) ReadConfig(configFile string) []byte {

	var configRawData []byte

	configReader, err := os.Open(configFile)
	if err != nil {
		log.Print("There was a problem reading config.yaml")
	}

	defer func() {
		err := configReader.Close()

		if err != nil {
			log.Print("There was a problem closing config.yaml")
		}
	}()

	_, readErr := configReader.Read(configRawData)
	if readErr != nil {
		log.Print("There was a problem reading config file to raw bytes.")
	}

	return configRawData

}

func (g *LamoidEnv) ValidateConfig() bool {

	g.WriteTempConfig(g.GrazeConfig())

	newConfig := md5.Sum(g.ReadConfig("tmp-config.yaml"))

	currentConfig := md5.Sum(g.ReadConfig("config.yaml"))

	log.Printf("[NEW-CONFIG]: Hash - %s", fmt.Sprint(newConfig))
	log.Printf("[OLD-CONFIG]: Hash - %s", fmt.Sprint(currentConfig))

	os.Remove("tmp-config.yaml")

	return cmp.Equal(newConfig, currentConfig)

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
	//TODO: Better Process Management for Collector and Reflector
	//TODO: CLI Flag to control config check interval
	//TODO: Unit Testing
	//TODO: Documentation
	g.StartGrazing()

Graze:
	for {
		time.Sleep(time.Second * 30)
		switch g.ValidateConfig() {
		case true:
			continue Graze
		case false:
			log.Printf("[LAMOID-INFO]: New Config - Stoping Collector")

			//TODO: This doesn't work like I expected it to.....  ¯\_(ツ)_/¯
			err := syscall.Kill(g.CollectorPID, syscall.SIGHUP)
			if err != nil {
				log.Printf("[LAMOID-ERR]: There was a problem trying to send SIGHUP to collector process, %s", err)
			}
			time.Sleep(time.Second * 10)

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
