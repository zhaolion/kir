package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/copier"
	"github.com/spf13/afero"
	"github.com/zhaolion/kir/pkg/config/env"
	"gopkg.in/yaml.v2"
)

var configer *Configer

// defaultLookupPaths find config file in those paths
var defaultLookupPaths = []string{"", "./conf", "/conf", "../", "../conf", "../..", "../../conf"}

// SupportedExts are universally supported extensions.
var SupportedExts = []string{"json", "yml", "toml"}

func init() {
	configer = New()
}

// AddConfigPath into lookup path list
func AddConfigPath(in string) {
	configer.AddPath(in)
}

// Initialize your config
func Initialize(configFile string, cfgStructPtr interface{}) error {
	if err := set(configFile, cfgStructPtr); err != nil {
		return err
	}

	if err := configer.loadConfig(); err != nil {
		return err
	}

	return copier.Copy(cfgStructPtr, configer.container)
}

func set(configFile string, cfgStructPtr interface{}) error {
	configName, configType := fileInfo(configFile)
	configer.configFile = ""
	configer.configName = configName

	if stringInSlice(configType, SupportedExts) {
		configer.configType = configType
	} else {
		return InvalidConfigTypeError(configType)
	}

	if err := checkObject(cfgStructPtr); err != nil {
		return err
	}

	configer.container = cfgStructPtr
	return nil
}

// Get return your config
func Get() interface{} {
	return configer.container
}

// Reset Intended for testing, will reset all to default settings.
// In the public interface for the viper package so applications
// can use it in their testing as well.
func Reset() {
	configer = New()
}

// New return config manager
func New() *Configer {
	return &Configer{
		container:   nil,
		configPaths: defaultLookupPaths,
		configName:  "config",
		configType:  "json",
		fs:          afero.NewOsFs(),
	}
}

// Configer config manager
type Configer struct {
	container interface{}

	configPaths []string
	configName  string
	configType  string
	configFile  string

	// The filesystem to read config from.
	fs afero.Fs
}

// AddPath into lookup path list
func (c *Configer) AddPath(in string) {
	in = absPath(in)
	if in != "" {
		c.configPaths = append(c.configPaths, in)
	}
}

func (c *Configer) getConfigType() string {
	if c.configType != "" {
		return c.configType
	}

	cf, err := c.getConfigFile()
	if err != nil {
		return ""
	}

	_, ext := fileInfo(cf)
	return ext
}

func (c *Configer) getConfigFile() (string, error) {
	if c.configFile == "" {
		cf, err := c.findConfigFile()
		if err != nil {
			return "", err
		}
		c.configFile = cf
	}
	return c.configFile, nil
}

// Search all configPaths for any config file.
// Returns the first path that exists (and is a config file).
func (c *Configer) findConfigFile() (string, error) {
	for _, cp := range c.configPaths {
		file := c.searchInPath(cp)
		if file != "" {
			return file, nil
		}
	}
	return "", FileNotFoundError{c.configName, fmt.Sprintf("%s", c.configPaths)}
}

func (c *Configer) searchInPath(in string) (filename string) {
	for _, ext := range SupportedExts {
		if b, _ := exists(c.fs, filepath.Join(in, c.configName+"."+ext)); b {
			return filepath.Join(in, c.configName+"."+ext)
		}
	}

	return ""
}

func (c *Configer) loadConfig() error {
	if err := c.processFile(); err != nil {
		return err
	}

	if err := c.processEnv(); err != nil {
		return err
	}

	return nil
}

func (c *Configer) processFile() error {
	filename, err := c.getConfigFile()
	if err != nil {
		return err
	}

	if !stringInSlice(c.getConfigType(), SupportedExts) {
		return UnsupportedConfigError(c.getConfigType())
	}

	file, err := afero.ReadFile(c.fs, filename)
	if err != nil {
		return err
	}

	cfg, err := c.unmarshalReader(bytes.NewReader(file))
	if err != nil {
		return err
	}

	c.container = cfg
	return nil
}

func (c *Configer) unmarshalReader(in io.Reader) (interface{}, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(in); err != nil {
		return nil, err
	}
	cfg := copyObject(c.container)

	switch strings.ToLower(c.getConfigType()) {
	case "json":
		if err := json.Unmarshal(buf.Bytes(), cfg); err != nil {
			return nil, ParseError{err}
		}
	case "yml":
		if err := yaml.Unmarshal(buf.Bytes(), cfg); err != nil {
			return nil, ParseError{err}
		}
	case "toml":
		if err := toml.Unmarshal(buf.Bytes(), cfg); err != nil {
			return nil, ParseError{err}
		}
	}

	return cfg, nil
}

func (c *Configer) processEnv() error {
	return env.Encode(c.container)
}
