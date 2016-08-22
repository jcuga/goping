package main

import (
    "flag"
    "fmt"
    "strings"
    "os"
    "os/exec"
)

// helper functions from: http://www.darrencoxall.com/golang/executing-commands-in-go/
func printCommand(cmd *exec.Cmd) {
  fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
  if err != nil {
    os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
  }
}

func printOutput(outs []byte) {
  if len(outs) > 0 {
    fmt.Printf("==> Output: %s\n", string(outs))
  }
}

struct PingPlace {
  name string `json:"name"`
  address string `json:"address"`
}

struct PingPlaceList {
  pingFrequencySec int `json:"ping_frequency_sec"`
  pingTimeoutMs int `json:"ping_timeout_ms"`
  addresses []PingPlace `json:"addresses"`
}

func main() {
  addressListFilename := flag.String("a", "address_list.json", "File of addresses to ping.")
  flag.Parse()
  fmt.Printf("Loading addresses to ping from: %s\n", *addressListFilename)

  // Create an *exec.Cmd
  cmd := exec.Command("ping", "-n", "1", "www.google.com")

  // Combine stdout and stderr
  printCommand(cmd)
  output, err := cmd.CombinedOutput()
  printError(err)
  printOutput(output) // => go version go1.3 darwin/amd64

}
