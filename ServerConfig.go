package webcontext

import (
	"fmt"
)

// [ServerConfig]

func NewServerConfig() *ServerConfig { // {{{
	return &ServerConfig{}
} // }}}

type ServerConfig struct {
	config

	Host             string
	Port             int
	NumWorkerJobs    int // Number of worker Jobs (per worker)
	NumServerWorkers int // Number of server workers
	LogMode          bool
}

func (this *ServerConfig) Name() string { // {{{
	return CONFIG_SERVER
} // }}}

func (this *ServerConfig) Addr() string { // {{{
	return fmt.Sprintf("%s:%d", this.Host, this.Port)
} // }}}

func (this *ServerConfig) Filepath(dir string) string { // {{{
	return filepath(dir, this)
} // }}}

func (this *ServerConfig) Load(dir string) error { // {{{
	err := load(this.Filepath(dir), this)
	this.setLoaded(err == nil)

	return err
} // }}}
