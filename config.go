package main

// Config contains runtime configuration for the redirector service.
type Config struct {
	// File path of the BoltDB database
	DatabasePath string `json:"databasePath"`

	// Interface and TCP port to bind to
	ListenAddr string `json:"listenAddr"`

	// Interface and TCP to bind the management service to
	MgmtAddr string `json:"mgmtAddr"`

	// Logfile to write to
	LogFile string `json:"logFile"`

	// The name of the key build to use when mapping URLs
	KeyBuilderName string `json:"keyBuilder"`

	// An instance of a KeyBuilder to use when mapping URLs
	KeyBuilder KeyBuilder `json:"-"`
}

// GetConfig returns a pointer to a singleton configuration struct.
func GetConfig() (*Config, error) {
	return &Config{
		DatabasePath: "./redirector.db", // current working directory
		ListenAddr:   ":8080",           // public on TCP 8080
		MgmtAddr:     "127.0.0.1:9321",  // local IPv4 only on 9321
		LogFile:      "-",               // stdout
		KeyBuilder:   RequestURIPathKeyBuilder(),
	}, nil
}
