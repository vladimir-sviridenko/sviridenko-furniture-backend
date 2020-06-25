package main

type Config struct {
	Server Server `json:""`
}

type Server struct {
	Addr string `json:""`
}
