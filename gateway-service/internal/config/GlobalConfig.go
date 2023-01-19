package config

const GlobalConfigKey = "globalconfig"
const DefaultConfigPath = "./gateway-service/config.yaml"

type GlobalConfig struct {
	Environment Environment `koanf:"Environment" envDefault:"dev"`
	TelegramEnv TelegramEnv
}

func (g *GlobalConfig) IsDev() bool {
	return g.Environment == DevEnvironment
}

func (g *GlobalConfig) IsStaging() bool {
	return g.Environment == StagingEnvironment
}

func (g *GlobalConfig) IsProd() bool {
	return g.Environment == ProdEnvironment
}
