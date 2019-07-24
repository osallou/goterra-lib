package terraconfig

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	consul "github.com/hashicorp/consul/api"
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

type ACL struct {
	AllowUserCreateNS bool `json:"acl_user_createns"`
}

// Config contains goterra configuration
type Config struct {
	loaded bool
	Redis  RedisConfig
	Mongo  MongoConfig `json:"mongo"`
	URL    string      `json:"url"`
	// Secret string
	Fernet []string `json:"fernet"`
	Web    WebConfig
	Deploy Deploy
	Amqp   string `json:"amqp"` // rabbitmq url connection "amqp://guest:guest@localhost:5672/"

	Git string `json:"git"` // public repo for recipe and templates

	Support string `json:"support"` // Support email address

	ACL ACL `json:"acl"` // ACL config

	DefaultImage string `json:"default_image"` // base image name to use if no recipe used
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
	if os.Getenv("GOT_SECRET") != "" {
		config.Fernet = []string{os.Getenv("GOT_SECRET")}
	}
	if os.Getenv("GOT_URL") != "" {
		config.URL = os.Getenv("GOT_URL")
	}
	if os.Getenv("GOT_AMQP") != "" {
		config.Amqp = os.Getenv("GOT_AMQP")
	}

	if os.Getenv("GOT_GIT") != "" {
		config.Git = os.Getenv("GOT_GIT")
	}

	if os.Getenv("GOT_SUPPORT") != "" {
		config.Support = os.Getenv("GOT_SUPPORT")
	}

	if os.Getenv("GOT_ACL_USER_CREATENS") == "1" {
		config.ACL.AllowUserCreateNS = true
	}
	cfg = config
	return config
}

// ConsulDeclare declare current service to consul
func ConsulDeclare(serviceName string, path string) error {
	// cfg := LoadConfig()
	if os.Getenv("GOT_CONSUL") != "" {
		fmt.Printf("Declare service to consul at %s\n", os.Getenv("GOT_CONSUL"))
		consulCfg := consul.DefaultConfig()
		consulCfg.Address = os.Getenv("GOT_CONSUL")
		client, err := consul.NewClient(consulCfg)
		if err != nil {
			return err
		}
		hostname, _ := os.Hostname()
		tags := []string{
			"got",
			"api",
			"traefik.backend=" + serviceName,
			"traefik.frontend.rule=PathPrefix:" + path,
			"traefik.enable=true",
		}
		check := &consul.AgentServiceCheck{
			CheckID:  hostname,
			HTTP:     fmt.Sprintf("http://%s:%d%s", hostname, cfg.Web.Port, path),
			Interval: "30s",
		}
		service := &consul.AgentServiceRegistration{
			ID:      hostname,
			Address: hostname,
			Name:    serviceName,
			Port:    int(cfg.Web.Port),
			Tags:    tags,
			Check:   check,
		}
		regerr := client.Agent().ServiceRegister(service)
		if regerr != nil {
			return regerr
		}
		fmt.Printf("service register in consul to handle api calls to %s\n", path)
	}
	return nil

}
