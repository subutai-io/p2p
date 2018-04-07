package ptp

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	IPTool  string `yaml:"iptool"`  // Network interface configuration tool
	AddTap  string `yaml:"addtap"`  // Path to addtap.bat
	InfFile string `yaml:"inffile"` // Path to deltap.bat
}

func (y *Configuration) GetIPTool() string {
	return y.IPTool
}

func (y *Configuration) GetAddTap() string {
	return y.AddTap
}

func (y *Configuration) GetInfFile() string {
	return y.InfFile
}

func (y *Configuration) Read() error {
	// TODO: Remove hard-coded path
	yamlFile, err := ioutil.ReadFile(ConfigDir + "/p2p/config.yaml")
	if err != nil {
		Log(Debug, "Failed to load config: %v", err)
		y.IPTool = "/sbin/ip"
		y.AddTap = "C:\\Program Files\\TAP-Windows\\bin\\tapinstall.exe"
		y.InfFile = "C:\\Program Files\\TAP-Windows\\driver\\OemVista.inf"
	}
	err = yaml.Unmarshal(yamlFile, y)
	if err != nil {
		return fmt.Errorf("Failed to parse config: %v", err)
	}
	return nil
}
