package config

import (
	"sort"
	"strings"

	env "github.com/ilyakaznacheev/cleanenv"
	"github.com/mvrahden/go-enumer/pkg/utils/slices"
)

const (
	SerializerBinary = "binary"
	SerializerBSON   = "bson"
	SerializerGQL    = "graphql"
	SerializerJSON   = "json"
	SerializerSQL    = "sql"
	SerializerText   = "text"
	SerializerYaml   = "yaml"
	SerializerYamlV3 = "yaml.v3"
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

func (sl stringList) Contains(s string) bool {
	return slices.Any(sl, func(v string, idx int) bool {
		return s == v
	})
}

func (sl stringList) ensureUnique() []string {
	sort.Strings(sl)
	// collect unique values
	return slices.Reduce(sl, func(v string, o []string) []string {
		if len(o) == 0 {
			return append(o, v)
		}
		if last := o[len(o)-1]; last != v {
			return append(o, v)
		}
		return o
	})
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
	cfg.Serializers = cfg.Serializers.ensureUnique()
	cfg.SupportedFeatures = cfg.SupportedFeatures.ensureUnique()
}
