package webcontext

import (
	"fmt"
	"io/ioutil"
	"path"

	"encoding/json"    // JSON (http://www.json.org/)
	"gopkg.in/yaml.v2" // YAML (http://www.yaml.org/)
)

const (
	CONFIG_TYPE_JSON = "json"
	CONFIG_TYPE_YAML = "yaml"

	CONFIG_MAIN     = "main"
	CONFIG_SERVER   = "server"
	CONFIG_DATABASE = "database"

	SERVER_TYPE_HTML = "HTML"
	SERVER_TYPE_FCGI = "FCGI"
)

var (
	configType = CONFIG_TYPE_JSON
)

func UseYaml() { // {{{
	configType = CONFIG_TYPE_YAML
} // }}}

func UseJson() { // {{{
	configType = CONFIG_TYPE_JSON
} // }}}

func load(file string, config interface{}) error { // {{{
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	switch configType {
	case CONFIG_TYPE_YAML:
		return yaml.Unmarshal(data, config)
	case CONFIG_TYPE_JSON:
		return json.Unmarshal(data, config)
	default:
		return fmt.Errorf("Unknown parser: '%s'", configType)
	}
} // }}}

func filepath(dir string, config iConfig) string { // {{{
	return path.Join(dir, fmt.Sprintf("%s.%s", config.Name(), configType))
} // }}}

// [config]

type config struct {
	loaded bool
}

func (this *config) setLoaded(loaded bool) { // {{{
	this.loaded = loaded
} // }}}

func (this *config) IsLoaded() bool { // {{{
	return this.loaded
} // }}}

// [iConfig]

type iConfig interface {
	Name() string
	Filepath(dir string) string
	Load(dir string) error
}

// [Config]

func LoadConfig(path string) (*Config, *ConfigErrors) { // {{{
	config := NewConfig()
	return config, config.Load(path)
} // }}}

func NewConfig() *Config { // {{{
	return &Config{
		Main:     NewMainConfig(),
		Server:   NewServerConfig(),
		Database: NewDatabaseConfig(),
	}
} // }}}

type Config struct {
	config

	Main     *MainConfig
	Server   *ServerConfig
	Database *DatabaseConfig
}

func (this *Config) HasMain() bool { // {{{
	return this.Main != nil && this.Main.IsLoaded()
} // }}}

func (this *Config) HasServer() bool { // {{{
	return this.Server != nil && this.Server.IsLoaded()
} // }}}

func (this *Config) HasDatabase() bool { // {{{
	return this.Database != nil && this.Database.IsLoaded()
} // }}}

func (this *Config) Load(path string) *ConfigErrors { // {{{
	var err error
	errs := make(ConfigErrors, 0)

	if err = this.Main.Load(path); err != nil {
		errs[this.Main.Name()] = err
	}

	if err = this.Server.Load(path); err != nil {
		errs[this.Server.Name()] = err
	}

	if err = this.Database.Load(path); err != nil {
		errs[this.Database.Name()] = err
	}

	this.setLoaded(true)

	return &errs
} // }}}
