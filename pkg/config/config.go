package config

type Config interface {
	Sanitize()
	Validate() error
}

type InitConfig interface {
	Init() error
}
