package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

const (
	// environment possible values
	EnvDevelopment = "development"
	EnvStaging     = "staging"
	EnvUat         = "uat"
	EnvProduction  = "production"

	// env vars name
	ENV    = "APP_ENV"
	REMOTE = "APP_REMOTE"
	DEBUG  = "APP_DEBUG"
)

// Bootstrap will load and set runtime configuration for a service based on given environment
func Bootstrap(name string) error {
	setEnv()
	viper.SetConfigType("json")
	errLocal := local(name)
	errRemote := remote(name)
	if errLocal != nil && errRemote != nil {
		return errors.Errorf("no valid configuration available for %s service", name)
	}
	if env := os.Getenv(ENV); env != EnvDevelopment && errRemote != nil {
		return errRemote
	}
	return nil
}

// remote will resolve configuration from remote consul storage
// it first resolved from value of env variable named `REMOTE`
// when none defined in env var, sensible default given `consul:8500`
// the consul path here must be formatted as `/service/{name}/config.json`
func remote(name string) error {
	host := os.Getenv(REMOTE)
	if len(host) == 0 {
		host = "consul:8500" // give sensible default
	}
	cpath := fmt.Sprintf("/service/%s/config.json", name)
	viper.AddRemoteProvider("consul", host, cpath)
	if err := viper.ReadRemoteConfig(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed retrieving config from %s with path %s", host, cpath))
	}
	return nil
}

// local will resolve config from local file in current executable directory
// with format `service_name.config.json`
func local(name string) error {
	viper.AddConfigPath(".")
	viper.SetConfigName(name + ".config")
	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed retrieving config from %s.config.json", name))
	}
	return nil
}

// setEnv validate and enforce value of environment var
// when no env name defined, `production` will be used by default
func setEnv() {
	env := EnvProduction
	if configEnv := os.Getenv(ENV); configEnv != "" {
		environments := []string{EnvDevelopment, EnvStaging, EnvUat, EnvProduction}
		for _, e := range environments {
			if configEnv == e {
				env = configEnv
				break
			}
		}
	}
	os.Setenv(ENV, env)
	setLogLevel(env)
}

// setLogLevel enforce log level based on env
func setLogLevel(env string) {
	var level int8
	switch env {
	case EnvDevelopment:
		level = 5 // debug
	case EnvStaging:
		level = 4 // info
	case EnvUat:
		level = 3 // warning
	default:
		level = 2 // error
	}
	if b, _ := strconv.ParseBool(os.Getenv(DEBUG)); b {
		level = 5 // debug
	}
	viper.Set("log_level", level)
}
