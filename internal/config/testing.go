package config

import "testing"

/*
	Метод возвращающаюй конфигурацию для тестирования сервера
*/
func InitTestConfig(t *testing.T) *Config {
	t.Helper()

	// Инициализирую конфигурацию
	return &Config{
		Services: ServicesConfigs{
			Server: ServerConfig{
				Host: "127.0.0.1",
				Port: ":4000",
			},

			DB: DatabaseConfig{
				DbUrl: "postgres://exchanger:qwerty@localhost:5432/exchanger_server_test?sslmode=disable",
			},

			Redis: RedisConfig{
				Host: "localhost",
				Port: 6379,
				DB:   1,
			},

			NSQ: NsqConfig{
				Host: "80.87.197.206",
				Port: 4150,
			},
		},

		Secrets: SecretsConfig{
			AccessSecret:  "accesssecret",
			RefreshSecret: "refreshsecret",
			TokenSecret:   "tokensecret",
		},

		Users: UsersConfig{
			Managers:   []string{},
			Developers: []string{"I0HuKc"},
			Admins:     []string{},
		},

		Plugins: PluginsConfig{
			AesKey: "fn5LyPGTnB18gl24nieHavsmKfKRmvLR",
		},
	}
}
