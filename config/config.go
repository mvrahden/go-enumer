package config

import (
	"sort"
	"strings"

	env "github.com/ilyakaznacheev/cleanenv"
)

const (
	SupportUndefined    = "undefined"
	SupportIgnoreCase   = "ignore-case"
	SupportEntInterface = "ent"
)

type Args Options
type Options struct {
	TransformStrategy string     `yaml:"transform" env-default:"noop"`
	Serializers       stringList `yaml:"serializers"`
	SupportedFeatures stringList `yaml:"support"`
}

func (o *Options) Clone() *Options {
	if o == nil {
		return nil
	}
	args := Args(*o)
	return (*Options)(&args)
}

type stringList []string

func (sl stringList) Sort() {
	sort.Strings(sl)
}

func (sl stringList) Contains(s string) bool {
	for _, v := range sl {
		if s == v {
			return true
		}
	}
	return false
}

func (sl stringList) Unique() (u []string) {
	sl.Sort()
	for idx, v := range sl {
		if idx > 0 && sl[idx-1] == v {
			continue
		}
		u = append(u, v)
	}
	return u
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
	_ = env.ReadEnv(cfg)
	cfg.Serializers = cfg.Serializers.Unique()
	cfg.SupportedFeatures = cfg.SupportedFeatures.Unique()
}
