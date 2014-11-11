package initialize

import "os"

type EnvReader interface {
	ReadPort() string
}

type EnvVarReader struct{}

func (e EnvVarReader) ReadPort() string {
	return os.Getenv("PORT")
}

func ReadPort(envReader EnvReader) string {
	port := envReader.ReadPort()
	if port == "" {
		port = "8080"
	}
	return port
}
