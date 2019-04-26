package domain

type Config struct {
	config map[string]interface{}
}

func (c *Config) SetFromBytes(data []byte) error {
	
}

func (c *Config) Get(serviceName string) (map[string]interface{}, error){
	
}