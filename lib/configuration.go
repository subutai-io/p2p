package ptp

import (
	"io/ioutil"
)

type YamlConfiguration struct {
	IPTool  string `yaml:"iptool"`  // Network interface configuration tool
	AddTap  string `yaml:"addtap"`  // Path to addtap.bat
	InfFile string `yaml:"inffile"` // Path to deltap.bat
}

func (y *YamlConfiguration) GetIPTool() string {
	return y.IPTool
}

func (y *YamlConfiguration) GetAddTap() string {
	return y.AddTap
}

func (y *YamlConfiguration) GetInfFile() string {
	return y.InfFile
}

func (y *YamlConfiguration) Assign(_IPTool, _AddTap, _InfFile string) error {
	y.IPTool = _IPTool
	y.AddTap = _AddTap
	y.InfFile = _InfFile
	return nil
}

func (y *YamlConfiguration) Read() []byte {
	// TODO: Remove hard-coded path
	yamlFile, err := ioutil.ReadFile(ConfigDir + "/p2p/config.yaml")
	if err != nil {
		Log(Debug, "Failed to load config: %v", err)
		y.IPTool = "/sbin/ip"
		y.AddTap = "C:\\Program Files\\TAP-Windows\\bin\\tapinstall.exe"
		y.InfFile = "C:\\Program Files\\TAP-Windows\\driver\\OemVista.inf"
	}
	return yamlFile
}
