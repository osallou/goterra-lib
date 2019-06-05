package terraconfig

import (
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// RedisConfig contains redis server connection info
type RedisConfig struct {
	Host   string
	Port   uint
	Prefix string
}

// MongoConfig contains mongodb server connection info
type MongoConfig struct {
	URL string `json:"url"`
	DB  string `json:"db"`
}

// WebConfig defines web server
type WebConfig struct {
	Listen string
	Port   uint
}

// Deploy contains information to deploy scripts
type Deploy struct {
	Path string
}

// Config contains goterra configuration
type Config struct {
	loaded bool
	Redis  RedisConfig
	Mongo  MongoConfig `json:"mongo"`
	URL    string      `json:"url"`
	Secret string
	Web    WebConfig
	Deploy Deploy
}

// Singleton config
var cfg Config

// ConfigFile config file path
var ConfigFile string

// LoadConfig returns the singleton config object
func LoadConfig() Config {
	if cfg.loaded {
		return cfg
	}

	cfgFile := os.Getenv("GOT_CONFIG")
	if cfgFile != "" {
		ConfigFile = os.Getenv("GOT_CONFIG")
	}
	if ConfigFile == "" {
		ConfigFile = "goterra.yml"
	}
	log.Printf("Using config file %s\n", ConfigFile)

	cfgfile, _ := ioutil.ReadFile(ConfigFile)
	config := Config{loaded: true}
	yaml.Unmarshal([]byte(cfgfile), &config)
	// fmt.Printf("Config: %+v\n", config)
	cfg = config
	return config
}
