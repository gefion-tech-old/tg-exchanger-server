package config

type Config struct {
	Server  ServerConfig   `toml:"server"`
	Secrets SecretsConfig  `toml:"secrets"`
	Redis   RedisConfig    `toml:"redis"`
	NSQ     NsqConfig      `toml:"nsq"`
	DB      DatabaseConfig `toml:"database"`
	Users   UsersConfig    `toml:"users"`
}

type ServerConfig struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	LogLevel string `toml:"log_level"`
	Tmp      string `toml:"tmp"`
}

type SecretsConfig struct {
	AccessSecret  string `toml:"access_secret"`
	RefreshSecret string `toml:"refresh_secret"`
	TokenSecret   string `toml:"token_secret"`
}

type RedisConfig struct {
	Host string `toml:"host"`
	Port int64  `toml:"port"`
	DB   int8   `toml:"db"`
}

type DatabaseConfig struct {
	DbUrl string `toml:"db_url"`
}

type NsqConfig struct {
	Host string `toml:"host"`
	Port uint16 `toml:"port"`
}

type UsersConfig struct {
	Managers   []string `toml:"managers"`
	Developers []string `toml:"developers"`
	Admins     []string `toml:"admins"`
}

func Init() *Config {
	return &Config{}
}
