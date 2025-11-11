package config

type Config struct {
	Umami    *UmamiConfig    `validate:"required" mapstructure:"umami" yaml:"umami"`
	Database *DatabaseConfig `validate:"required" mapstructure:"database" yaml:"database"`
	Imports  []*ImportConfig `validate:"required,min=1,dive" mapstructure:"imports" yaml:"imports"`
}

type UmamiConfig struct {
	CollectionURL     string            `validate:"required,url" mapstructure:"collection_url" yaml:"collection_url"`
	Timeout           int               `validate:"min=1" mapstructure:"timeout" yaml:"timeout"`
	Retries           int               `validate:"required,min=1" mapstructure:"retries" yaml:"retries"`
	MaxRequests       int               `validate:"required,min=1" mapstructure:"max_clients" yaml:"max_requests"`
	IgnoreTLS         bool              `mapstructure:"ignore_tls" yaml:"ignore_tls"`
	CustomHTTPHeaders map[string]string `validate:"dive,keys,required,endkeys,required" mapstructure:"custom_http_headers" yaml:"custom_http_headers"`
}

type DatabaseConfig struct {
	Path string `validate:"required" mapstructure:"path" yaml:"path"`
}

type ImportConfig struct {
	Website WebsiteConfig `validate:"required" mapstructure:"website" yaml:"website"`
	Logs    LogsConfig    `validate:"required" mapstructure:"logs" yaml:"logs"`
}

type WebsiteConfig struct {
	ID      string `validate:"required,uuid4" mapstructure:"id" yaml:"id"`
	BaseURL string `validate:"required,url" mapstructure:"base_url" yaml:"base_url"`
}

type LogsConfig struct {
	Paths               []string `validate:"required,min=1,dive,required" mapstructure:"paths" yaml:"paths"`
	Type                string   `validate:"required,oneof=apache nginx custom" mapstructure:"type" yaml:"type"`
	TypeCustomRegex     string   `validate:"required_if=Type custom,regex" mapstructure:"type_custom_regex" yaml:"type_custom_regex"`
	TypeCustomTimestamp string   `validate:"required_if=Type custom" mapstructure:"type_custom_timestamp" yaml:"type_custom_timestamp"`
	IncludeExtension    string   `mapstructure:"include_extension" yaml:"include_extension"`
	Recursive           bool     `mapstructure:"recursive" yaml:"recursive"`
}
