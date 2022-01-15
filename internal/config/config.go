package config

type Config struct {
	Services ServicesConfigs `toml:"services"`
	Secrets  SecretsConfig   `toml:"secrets"`
	Users    UsersConfig     `toml:"users"`
	Plugins  PluginsConfig   `toml:"plugins"`
}

type ServicesConfigs struct {
	Server ServerConfig   `toml:"server"`
	DB     DatabaseConfig `toml:"database"`
	Redis  RedisConfig    `toml:"redis"`
	NSQ    NsqConfig      `toml:"nsq"`
}

type PluginsConfig struct {
	Whitebit WhitebitConfig `toml:"whitebit"`
}

type ServerConfig struct {
	Host     string `toml:"HOST"`
	Port     string `toml:"PORT"`
	LogLevel string `toml:"LOG_LEVEL"`
	Tmp      string `toml:"TMP"`
}

type SecretsConfig struct {
	AccessSecret  string `toml:"ACCESS_SECRET"`
	RefreshSecret string `toml:"REFRESH_SECRET"`
	TokenSecret   string `toml:"TOKEN_SECRET"`
}

type RedisConfig struct {
	Host string `toml:"HOST"`
	Port int64  `toml:"PORT"`
	DB   int8   `toml:"DB"`
}

type DatabaseConfig struct {
	DbUrl string `toml:"DB_URL"`
}

type NsqConfig struct {
	Host string `toml:"HOST"`
	Port uint16 `toml:"PORT"`
}

type UsersConfig struct {
	Managers   []string `toml:"MANAGERS"`
	Developers []string `toml:"DEVELOPERS"`
	Admins     []string `toml:"ADMINS"`
}

type WhitebitConfig struct {
	PublicKey string `toml:"PUBLIC_KEY"`
	SecretKey string `toml:"PRIVATE_KEY"`
	URL       string `toml:"URL"`
}

func Init() *Config {
	return &Config{}
}
