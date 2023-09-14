package config

type AppManagerServiceClientConfig struct {
	ServerAddr  string `json:"serverAddr"`
	DialTimeout int64  `json:"dialTimeout"` // in milliseconds
	CallTimeout int64  `json:"callTimeout"` // in milliseconds
}

type LoggingManagerServiceClientConfig struct {
	ServerAddr  string `json:"serverAddr"`
	DialTimeout int64  `json:"dialTimeout"` // in milliseconds
	CallTimeout int64  `json:"callTimeout"` // in milliseconds
}
