package alpaca

func (g *LamoidEnv) ValidateEnvironment() {
	// Validate Server Configuration
}

func (g *LamoidEnv) StartReflector() {
	// Start llama reflector and update the process id ref.
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
