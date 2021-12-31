package config

import (
	"strings"

	env "github.com/ilyakaznacheev/cleanenv"
)

type Args Options
type Options struct {
	TypeAliasName     string   `yaml:"typeAlias"`
	Output            string   `yaml:"output"`
	TransformStrategy string   `yaml:"transform" env-default:"noop"`
	AddPrefix         string   `yaml:"addPrefix"`
	Serializers       []string `yaml:"serializers"`
}

type Strings []string

func (s Strings) String() string {
	return strings.Join(s, ",")
}

func (s *Strings) Set(v string) error {
	v = strings.ReplaceAll(v, ", ", ",")
	v = strings.TrimSpace(v)
	v = strings.TrimSuffix(v, ",")
	*s = strings.Split(v, ",")
	return nil
}

func LoadWith(args *Args) *Options {
	cfg := (*Options)(args)
	loadFromFile("", cfg)
	return cfg
}

func LoadFrom(file string) *Options {
	var cfg Options
	loadFromFile(file, &cfg)
	return &cfg
}

func loadFromFile(file string, cfg *Options) {
	_ = env.ReadConfig(file, cfg)
}
