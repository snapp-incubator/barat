package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var C *Config

type Config struct {
	ProjectPath  string        `yaml:"project-path"`
	TomlPaths    []string      `yaml:"toml-path"`
	Exclude      Exclude       `yaml:"exclude"`
	MessageFuncs []MessageFunc `yaml:"message-func"`
	Options      Opts          `yaml:"options"`
}

type Exclude struct {
	Folders   []string `yaml:"folders"`
	RegexKeys []string `yaml:"regex-key"`
}

type MessageFunc struct {
	Name        string `yaml:"name"`
	MessageIDNo int    `yaml:"message-id-no"`
}

func ToMessageFuncs(args map[string]int) []MessageFunc {
	var tmp []MessageFunc
	for k, v := range args {
		tmp = append(tmp, MessageFunc{Name: k, MessageIDNo: v})
	}
	return tmp
}

type Enable struct {
	TomlCheck        bool `yaml:"toml-check"`
	DescriptionCheck bool `yaml:"description-check"`
	OtherKeyCheck    bool `yaml:"other-key-check"`
	CodeCheck        bool `yaml:"code-check"`
}

type Opts struct {
	Enable Enable `yaml:"enable"`
}

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
