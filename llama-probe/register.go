package main


import (
	"encoding/json"
	"net/http"
	"strconv"
	"bytes"
	"time"
	"log"
	"os"
)


// Create struct for JSON we send to the server for registration
type PayLoad struct {
	Port 		int 	`json:"port"`
	Keepalive 	int   	`json:"keepalive,omitempty"`
	Ip 			string  `json:"ip,omitempty"`
	Tags 		struct {
		Version			string	`json:"version"`
		ProbeShortname	string	`json:"probe_shortname"`
		ProbeName		string	`json:"probe_name"`
	}  `json:"tags"`
	Group 		string  `json:"group,omitempty"`
}


// Load environment variables from OS
func initEnvVars() map[string]string {

	// List of OS environment variables to load and store
	envs := []string{
		"PROBE_NAME",
		"PROBE_SHORTNAME",
		"LLAMA_SERVER",
		"LLAMA_KEEPALIVE",
		"LLAMA_GROUP",
		"LLAMA_PORT",
		"LLAMA_SOURCE_IP"}

	m := make(map[string]string)

	for _, env := range envs {
		m[env] = os.Getenv(env)
	}

	// Error out if we are missing the following ENV variables | LLAMA_SERVER, PROBE_NAME, PROBE_SHORTNAME
	if ( m["LLAMA_SERVER"] == "" || m["PROBE_NAME"] == "" || m["PROBE_SHORTNAME"] == "") {

		// Log that we are missing required information
		log.Println("ERROR: Missing required environment variables!")
		log.Println("LLAMA_SERVER: ", m["LLAMA_SERVER"])
		log.Println("PROBE_NAME: ", m["PROBE_NAME"])
		log.Println("PROBE_SHORTNAME: ", m["PROBE_SHORTNAME"] )

		// Exit app with error code 1
		os.Exit(1)
	}

	// Pass map of environment variables back to system
	return m
}


// Main program
func main() {
	// Program version for tagging
	appversion := "0.1.0"
	log.Println("Starting up registration client version:", appversion)

	// Load and print environment variables
	envs := initEnvVars()
	log.Println("Environment Variables Loaded:", envs)

	var registration PayLoad

	// If no port infromation is given, auto set this to 8100
	if envs["LLAMA_PORT"] == "" {
		envs["LLAMA_PORT"] = "8100"
	}

	// LLAMA Port Setting, convert string to int | OPTIONAL
	if envs["LLAMA_PORT"] != "" {
		intPort, err := strconv.Atoi(envs["LLAMA_PORT"])
		if err != nil {
			log.Println(err)
		}
		registration.Port = intPort
	}

	// LLAMA Keepalive Settings, convert string to int | OPTIONAL
	if envs["LLAMA_KEEPALIVE"] != "" {
		intKeepalive, err := strconv.Atoi(envs["LLAMA_KEEPALIVE"])
		if err != nil {
			log.Println(err)
		}
		registration.Keepalive = intKeepalive
	}

	// LLAMA Group Settings | OPTIONAL
	if envs["LLAMA_GROUP"] != "" {
		registration.Group = envs["LLAMA_GROUP"]
	}

	// LLAMA Source IP | OPTIONAL
	if envs["LLAMA_SOURCE_IP"] != "" {
		registration.Ip = envs["LLAMA_SOURCE_IP"]
	}

	// Register tag values
	registration.Tags.Version = appversion
	registration.Tags.ProbeShortname = envs["PROBE_SHORTNAME"]
	registration.Tags.ProbeName = envs["PROBE_NAME"]

	// Convert PayLoad struct to JSON now that we have all our values.
    byteArray, err := json.Marshal(registration)
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
