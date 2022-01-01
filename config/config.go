package config

import (
	"strings"

	env "github.com/ilyakaznacheev/cleanenv"
)

const (
	SupportUndefined    = "undefined"
	SupportEntInterface = "ent"
)

type Args Options
type Options struct {
	TypeAliasName     string     `yaml:"typeAlias"`
	Output            string     `yaml:"output"`
	TransformStrategy string     `yaml:"transform" env-default:"noop"`
	AddPrefix         string     `yaml:"addPrefix"`
	Serializers       stringList `yaml:"serializers"`
	SupportedFeatures stringList `yaml:"support"`
}

type stringList []string

func (sl stringList) Contains(s string) bool {
	for _, v := range sl {
		if s == v {
			return true
		}
	}
	return false
}

func (sl stringList) String() string {
	return strings.Join(sl, ",")
}

func (sl *stringList) Set(v string) error {
	v = strings.ReplaceAll(v, ", ", ",")
	v = strings.TrimSpace(v)
	v = strings.TrimSuffix(v, ",")
	*sl = strings.Split(v, ",")
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
