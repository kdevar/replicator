package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Env string

func (e *Env) Set(s string) error {
	*e = Env(s)
	return nil
}

func (e *Env) String() string {
	return string(*e)
}

type EnvConfig struct {
	Dev   Config
	Stage Config
	Prod  Config
}

func (ec *EnvConfig) GetConfig(e Env) Config {
	switch e {
	case DEVENV:
		return ec.Dev
	case STAGEENV:
		return ec.Stage
	case PRODENV:
		return ec.Prod
	default:
		return ec.Dev
	}
}

const (
	DEVENV   Env = "dev"
	STAGEENV Env = "stage"
	PRODENV  Env = "prod"
)

type Config struct {
	Session       *session.Session
	Kinesis       *kinesis.Kinesis
	IncludeTables []string `yaml:"IncludeTables"`
	ExcludeTables []string `yaml:"ExcludeTables"`
	Host          string   `yaml:"Host"`
	Port          int      `yaml:"Port"`
	User          string   `yaml:"User"`
	Password      string   `yaml:"Password"`
	Flavor        string   `yaml:"Flavor"`
	Env           Env      `yaml:"Env"`
}

func NewConfig() *Config {
	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	kinesis := kinesis.New(session, &aws.Config{})

	var env Env
	flag.Var(&env, "env", "what is the env")

	var configPath string
	flag.StringVar(&configPath, "cfg", "", "relative path to config file")

	flag.Parse()

	log.Printf("Path = %v", configPath)

	yamlFile, err := ioutil.ReadFile(configPath)

	fmt.Println(string(yamlFile))




	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	var envConfig EnvConfig
	err = yaml.Unmarshal(yamlFile, &envConfig)

	log.Println(envConfig)

	if err != nil {
		log.Printf("yamlFile unmarshal err   #%+v \n", err)
	}

	currentConfig := envConfig.GetConfig(env)
	log.Printf("Config constructed env=%v config=%+v \n", env, currentConfig)
	currentConfig.Session = session
	currentConfig.Kinesis = kinesis
	return &currentConfig
}
