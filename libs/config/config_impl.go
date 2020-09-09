package config

import (
	"strings"

	"github.com/spf13/viper"
)

type configuration struct {
	v *viper.Viper
}

var conf *configuration

func init() {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	conf = &configuration{v}
}

// Get configuration for key 'string'
func Get(key string) string {
	return conf.v.GetString(key)
}

// GetInt configuration for key 'int'
func GetInt(key string) int {
	return conf.v.GetInt(key)
}
