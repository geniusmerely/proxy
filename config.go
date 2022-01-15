package main

type User struct {
	UserName string `toml:"username"`
	Password string `toml:"password"`
}
type Config struct {
	BindIp     string `toml:"bind_ip"`
	ListenAddr string `toml:"listen_addr"`
	Users      []User `toml:"users"`
}

func DefaultConfig() *Config {
	return &Config{
		BindIp:     "",
		ListenAddr: ":8888",
	}
}
