package webcontext

// DatabaseConfig

func NewDatabaseConfig() *DatabaseConfig { // {{{
	return &DatabaseConfig{}
} // }}}

type DatabaseConfig struct {
	config

	Driver        string
	DSN           string
	MaxIdleConns  int
	MaxOpenConns  int
	SingularTable bool
}

func (this *DatabaseConfig) Name() string {
	return CONFIG_DATABASE
}

func (this *DatabaseConfig) Filepath(dir string) string { // {{{
	return filepath(dir, this)
} // }}}

func (this *DatabaseConfig) Load(dir string) error { // {{{
	err := load(this.Filepath(dir), this)
	this.setLoaded(err == nil)

	return err
} // }}}
