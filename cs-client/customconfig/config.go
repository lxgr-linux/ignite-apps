package customconfig

type Config struct {
	Client *ClientConfig `yaml:"client"`
}

type ClientConfig struct {
	CsClient *CsClientConfig `yaml:"cs-client,omitempty"`
}

type CsClientConfig struct {
	Path string `yaml:"path"`
}
