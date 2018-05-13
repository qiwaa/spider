package config

type Config struct {
	Addr       string		// dht address
	EntryNodes []string		// entry nodes
}

func NewFixedConfig() *Config {
	c := &Config{
		Addr: ":6881",
		EntryNodes: []string{
			"router.utorrent.com:6881",
		},
	}
	return c
}
