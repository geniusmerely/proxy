package main

type Config struct {
	BindIp     string `toml:"bind_ip"`
	ListenAddr string `toml:"listen_addr"`
}

func DefaultConfig() *Config {
	return &Config{
		BindIp:     "",
		ListenAddr: ":8888",
	}
}
