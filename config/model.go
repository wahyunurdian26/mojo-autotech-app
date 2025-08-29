package config

type Config struct {
	Db  Database
	Srv Server
}

type Database struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Name string `json:"name"`
}

type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
}
