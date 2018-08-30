package config

type Config struct {
	Addr       string   // dht address
	EntryNodes []string // entry nodes
}

func NewFixedConfig() *Config {
	c := &Config{
		Addr: ":6881",
		EntryNodes: []string{
			"router.bittorrent.com:6881",
			"router.utorrent.com:6881",
			"dht.transmissionbt.com:6881",
		},
	}
	return c
}
