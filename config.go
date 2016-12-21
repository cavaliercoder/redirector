package main

type Config struct {
	DatabasePath   string     `json:"databasePath"`
	ListenAddr     string     `json:"listenAddr"`
	MgmtAddr       string     `json:"mgmtAddr"`
	LogFile        string     `json:"logFile"`
	KeyBuilderName string     `json:"keyBuilder"`
	KeyBuilder     KeyBuilder `json:"-"`
}

func GetConfig() (*Config, error) {
	return &Config{
		DatabasePath: "./redirector.db",
		ListenAddr:   ":8080",
		MgmtAddr:     "127.0.0.1:9321",
		LogFile:      "-",
		KeyBuilder:   RequestURIPathKeyBuilder(),
	}, nil
}
