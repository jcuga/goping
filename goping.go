package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type PingPlace struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Config struct {
	PingFrequencySec int         `json:"ping_frequency_sec"`
	PingTimeoutSec   int         `json:"ping_timeout_sec"`
	Addresses        []PingPlace `json:"addresses"`
}

func getConfigFromFile(filename string) (*Config, error) {
	configFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	var config Config
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

// This works for both mac and linux output, not sure if for windows too...
func parseResults(cmd *exec.Cmd, name, address string, pattern *regexp.Regexp) {
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("event='ping_cmd_error' name='%s' addresss='%s' error='%s'\n", name, address, err)
	}
	if len(output) > 0 {
		for _, line := range strings.Split(string(output), "\n") {
			if matches := pattern.FindStringSubmatch(line); matches != nil && len(matches) >= 2 {
				log.Printf("event='ping_latency' name='%s' addresss='%s' latency_ms='%s'\n", name, address, matches[1])
				return
			}
		}
	}
	// guess we never found a ping latency in our response data
	log.Printf("event='missed_ping_latency' name='%s' addresss='%s'\n", name, address)
}

func pingLinux(address, name string, timeoutSec int, pattern *regexp.Regexp) {
	// -c 1 --> send one packet -w <sec> deadline/timeout in seconds before giving up
	cmd := exec.Command("ping", "-c", "1", "-w", strconv.Itoa(timeoutSec), address)
	parseResults(cmd, name, address, pattern)
}

func pingMac(address, name string, timeoutSec int, pattern *regexp.Regexp) {
	// -c 1 --> send one packet -t <sec> timeout in sec before ping exits
	// regardless of packets received
	cmd := exec.Command("ping", "-c", "1", "-t", strconv.Itoa(timeoutSec), address)
	parseResults(cmd, name, address, pattern)
}

func pingWindows(address, name string, timeoutSec int, pattern *regexp.Regexp) {
	log.Fatalf("TODO: support windows\n")
}

func ping(config *Config, address, name string, pattern *regexp.Regexp) {
	if config == nil {
		log.Fatalf("Config arg was nil.  Abort!\n")
	}
	switch os := runtime.GOOS; os {
	case "darwin":
		pingMac(address, name, config.PingTimeoutSec, pattern)
	case "linux":
		pingLinux(address, name, config.PingTimeoutSec, pattern)
	case "windows":
		pingWindows(address, name, config.PingTimeoutSec, pattern)
	default:
		log.Fatalf("Unsupported OS type: %s.  Can't establish ping cmd args.\n", os)
	}
}

func main() {
	// TODO: if this pattern is different for windows, make condition here
	// but this covers both mac/linux ping results
	LATENCY_PATTERN := regexp.MustCompile("time=(.*) ms")
	addressListFilename := flag.String("f", "address_list.json", "File of addresses to ping.")
	flag.Parse()
	log.Printf("event='program_args' config_filename='%s'\n", *addressListFilename)
	config, err := getConfigFromFile(*addressListFilename)
	if err != nil {
		log.Fatalf("Error parsing config file: %s\n", err)
	}
	log.Printf("event='config_values' timeout_sec='%d' ping_freq_sec='%d'\n", config.PingFrequencySec,
		config.PingTimeoutSec)
	if len(config.Addresses) == 0 {
		log.Fatalf("No addresses listed in config.  Nothing to ping.  Abort!\n")
	}
	for {
		select {
		case <-time.After(time.Second * time.Duration(config.PingFrequencySec)):
			for _, place := range config.Addresses {
				go ping(config, place.Address, place.Name, LATENCY_PATTERN)
			}
		}
	}
}
