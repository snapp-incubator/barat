// Package config
package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// C is the global singleton variable for the config.
var C *Config

// Config is the struct for the config file.
type Config struct {
	ProjectPath  string        `yaml:"project-path"`
	TomlPaths    []string      `yaml:"toml-path"`
	Exclude      Exclude       `yaml:"exclude"`
	MessageFuncs []MessageFunc `yaml:"message-functions"`
	Options      Opts          `yaml:"options"`
}

// Exclude is the struct for the exclude section in the config file. It contains the folders
// and regex keys that should be excluded.
type Exclude struct {
	Folders   []string `yaml:"folders"`
	RegexKeys []string `yaml:"regex-keys"`
}

// MessageFunc is the struct for the message-func section in the config file. It contains the function name
// that returns the message by i18n and the number of MessageID in arguments.
type MessageFunc struct {
	Name        string `yaml:"name"`
	MessageIDNo int    `yaml:"message-id-no"`
}

// ToMessageFuncs converts the map of the function's name to MessageID in arguments to a slice of MessageFunc.
func ToMessageFuncs(args map[string]int) []MessageFunc {
	var tmp []MessageFunc
	for k, v := range args {
		tmp = append(tmp, MessageFunc{Name: k, MessageIDNo: v})
	}
	return tmp
}

// Enable is the struct for the enable section in the config file. It contains the items
// that can be enabled or disabled.
type Enable struct {
	TomlCheck        bool `yaml:"toml-check"`
	DescriptionCheck bool `yaml:"description-check"`
	OtherKeyCheck    bool `yaml:"other-key-check"`
	CodeCheck        bool `yaml:"code-check"`
}

// Opts is the struct for the options section in the config file. It contains the items that can be enabled or disabled.
type Opts struct {
	Enable Enable `yaml:"enable"`
}

// LoadConfig loads the config file.
func LoadConfig(filepath string) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &C)
	if err != nil {
		return err
	}

	return nil
}
