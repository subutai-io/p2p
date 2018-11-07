package ptp

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type conf struct {
	IPTool  string `yaml:"iptool"`
	TAPTool string `yaml:"taptool"`
	INFFile string `yaml:"inf_file"`
	MTU     int    `yaml:"mtu"`
	PMTU    bool   `yaml:"pmtu"`
}

func (c *conf) readConf(filepath string) error {
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		Log(Debug, "Failed to load config: %v", err)
		c.IPTool = DefaultIPTool
		c.TAPTool = DefaultTAPTool
		c.INFFile = DefaultINFFile
		c.MTU = DefaultMTU
		c.PMTU = DefaultPMTU
		return fmt.Errorf("Failed to load configuration file from %s: %s", filepath, err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return fmt.Errorf("Failed to parse config: %v", err)
	}
	return nil
}

func (c *conf) getIPTool(preset string) string {
	if preset != "" {
		return preset
	}
	return c.IPTool
}

func (c *conf) getTAPTool(preset string) string {
	if preset != "" {
		return preset
	}
	return c.TAPTool
}

func (c *conf) getINFFile(preset string) string {
	if preset != "" {
		return preset
	}
	return c.INFFile
}

func (c *conf) getMTU(preset int) int {
	if preset != 0 {
		return preset
	}
	return c.MTU
}

func (c *conf) getPMTU(preset bool) bool {
	if preset {
		return preset
	}
	return c.PMTU
}
