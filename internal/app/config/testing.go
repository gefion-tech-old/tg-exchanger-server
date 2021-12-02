package config

import "testing"

/*
	Метод возвращающаюй конфигурацию для тестирования сервера
*/
func InitTestConfig(t *testing.T) *Config {
	t.Helper()

	// Инициализирую конфигурацию
	return &Config{
		Server: ServerConfig{
			Host: "127.0.0.1",
			Port: ":4000",
		},
		Redis: RedisConfig{
			Host: "localhost",
			Port: 6379,
			DB:   1,
		},
		Secrets: SecretsConfig{
			AccessSecret:  "accesssecret",
			RefreshSecret: "refreshsecret",
			TokenSecret:   "tokensecret",
		},
		DB: DatabaseConfig{
			DbUrl: "postgres://exchanger:qwerty@localhost:5432/exchanger_server_test?sslmode=disable",
		},
		Environment: EnvironmentConfig{},
	}

}
