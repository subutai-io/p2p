package ptp

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Conf struct {
	IPTool  string `yaml:"iptool"`
	TAPTool string `yaml:"taptool"`
	INFFile string `yaml:"inf_file"`
	MTU     int    `yaml:"mtu"`
	PMTU    bool   `yaml:"pmtu"`
}

func (c *Conf) Load(filepath string) error {
	c.SetDefaults()
	if filepath == "" {
		return nil
	}
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return fmt.Errorf("config parse failed: %s", err.Error())
	}
	return nil
}

func (c *Conf) SetDefaults() {
	c.IPTool = DefaultIPTool
	c.TAPTool = DefaultTAPTool
	c.INFFile = DefaultINFFile
	c.MTU = DefaultMTU
	c.PMTU = DefaultPMTU
}

func (c *Conf) GetIPTool(preset string) string {
	if preset != "" {
		return preset
	}
	return c.IPTool
}

func (c *Conf) GetTAPTool(preset string) string {
	if preset != "" {
		return preset
	}
	return c.TAPTool
}

func (c *Conf) GetINFFile(preset string) string {
	if preset != "" {
		return preset
	}
	return c.INFFile
}

func (c *Conf) GetMTU(preset int) int {
	if preset != 0 {
		return preset
	}
	return c.MTU
}

func (c *Conf) GetPMTU() bool {
	return c.PMTU
}
