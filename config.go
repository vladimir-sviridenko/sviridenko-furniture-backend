package main

type Config struct {
	Server ServerConfig `json:"server"`
	Email  EmailConfig  `json:"email"`
}

type ServerConfig struct {
	Addr string `json:"addr"`
}

type EmailConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}
