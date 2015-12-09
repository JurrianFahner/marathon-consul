package config

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/allegro/marathon-consul/consul"
	"github.com/allegro/marathon-consul/marathon"
	"github.com/allegro/marathon-consul/metrics"
	"github.com/allegro/marathon-consul/sync"
	flag "github.com/ogier/pflag"
	"io/ioutil"
	"time"
)

type Config struct {
	Consul consul.ConsulConfig
	Web    struct {
		Listen string
	}
	Sync       sync.Config
	Marathon   marathon.Config
	Metrics    metrics.Config
	LogLevel   string
	configFile string
}

var config = &Config{Marathon: marathon.Config{}}

func New() (*Config, error) {
	if !flag.Parsed() {
		config.parseFlags()
	}
	flag.Parse()
	err := config.loadConfigFromFile()
	config.setLogLevel()

	return config, err
}

func (config *Config) parseFlags() {
	// Consul
	flag.StringVar(&config.Consul.Port, "consul-port", "8500", "Consul port")
	flag.BoolVar(&config.Consul.Auth.Enabled, "consul-auth", false, "Use Consul with authentication")
	flag.StringVar(&config.Consul.Auth.Username, "consul-auth-username", "", "The basic authentication username")
	flag.StringVar(&config.Consul.Auth.Password, "consul-auth-password", "", "The basic authentication password")
	flag.BoolVar(&config.Consul.SslEnabled, "consul-ssl", false, "Use HTTPS when talking to Consul")
	flag.BoolVar(&config.Consul.SslVerify, "consul-ssl-verify", true, "Verify certificates when connecting via SSL")
	flag.StringVar(&config.Consul.SslCert, "consul-ssl-cert", "", "Path to an SSL client certificate to use to authenticate to the Consul server")
	flag.StringVar(&config.Consul.SslCaCert, "consul-ssl-ca-cert", "", "Path to a CA certificate file, containing one or more CA certificates to use to validate the certificate sent by the Consul server to us")
	flag.StringVar(&config.Consul.Token, "consul-token", "", "The Consul ACL token")

	// Web
	flag.StringVar(&config.Web.Listen, "listen", ":4000", "accept connections at this address")

	// Sync
	flag.DurationVar(&config.Sync.Interval, "sync-interval", 15*time.Minute, "Marathon-consul sync interval")

	// Marathon
	flag.StringVar(&config.Marathon.Location, "marathon-location", "localhost:8080", "Marathon URL")
	flag.StringVar(&config.Marathon.Protocol, "marathon-protocol", "http", "Marathon protocol (http or https)")
	flag.StringVar(&config.Marathon.Username, "marathon-username", "", "Marathon username for basic auth")
	flag.StringVar(&config.Marathon.Password, "marathon-password", "", "Marathon password for basic auth")
	flag.BoolVar(&config.Marathon.VerifySsl, "marathon-ssl-verify", true, "Verify certificates when connecting via SSL")

	// Metrics
	flag.StringVar(&config.Metrics.Target, "metrics-target", "stdout", "Metrics destination stdout or graphite")
	flag.StringVar(&config.Metrics.Prefix, "metrics-prefix", "default", "Metrics prefix (default is resolved to <hostname>.<app_name>")
	flag.DurationVar(&config.Metrics.Interval, "metrics-interval", 30*time.Second, "Metrics reporting interval")
	flag.StringVar(&config.Metrics.Addr, "metrics-location", "", "Graphite URL (used when metrics-target is set to graphite)")

	// General
	flag.StringVar(&config.LogLevel, "log-level", "info", "Log level: panic, fatal, error, warn, info, or debug")
	flag.StringVar(&config.configFile, "config-file", "", "Path to a JSON file to read configuration from. Note: Will override options set earlier on the command line")
}

func (config *Config) loadConfigFromFile() error {
	if config.configFile == "" {
		return nil
	}
	jsonBlob, err := ioutil.ReadFile(config.configFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBlob, config)
}

func (config *Config) setLogLevel() {
	level, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.WithField("level", config.LogLevel).Fatal("bad level")
	}
	log.SetLevel(level)
}
