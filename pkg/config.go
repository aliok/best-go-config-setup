package pkg

import (
	"github.com/aliok/go-defaultz"
	"github.com/go-playground/validator/v10"
)

// `json`: Used for marshalling and unmarshalling JSON and YAML, plus used by Viper
// `jsonschema`: Used for generating JSON schema and defaulting
// `validate`: Used for validating the configuration

type Config struct {
	// HTTPServerConfig is the configuration for the HTTP server.
	HTTPServerConfig HTTPServerConfig `json:"http_server"`

	// FeatureConfig is the configuration for the features.
	FeatureConfig FeatureConfig `json:"features"`

	// LoggingConfig is the configuration for the logging.
	LoggingConfig LoggingConfig `json:"logging"`
}

type HTTPServerConfig struct {
	// Port is the port number for the HTTP server
	Port int `json:"port,omitempty" jsonschema:"default=8080" validate:"required,min=1,max=65535"`

	// BindAddress is the address to bind to
	BindAddress string `json:"bind_address,omitempty" jsonschema:"default=0.0.0.0" validate:"required,ip4_addr"`
}

type FeatureConfig struct {
	// EnabledFeatures is the list of enabled features
	EnabledFeatures []string `json:"enabled_features,omitempty" jsonschema:"omitempty,default=feature1 feature2"`
}

type LoggingConfig struct {
	// LogLevel is the log level for the application
	LogLevel *int8 `json:"log_level,omitempty" jsonschema:"default=2" validate:"required,min=-1,max=5"`
	// field above is a pointer to distinguish between zero value and default value

	// LogFormat is the format of the logs. Can be `json` or `pretty`.
	LogFormat string `json:"log_format,omitempty" jsonschema:"default=json,enum=json,enum=pretty" validate:"required,oneof=json pretty"`
}

func HandleConfig(cfg *Config) error {
	// use go-defaultz to apply defaults
	// reuse the `jsonschema` tag and the `default=` prefix
	defaulter := defaultz.NewDefaulterRegistry(
		defaultz.WithBasicDefaulters(),
		defaultz.WithDefaultExtractor(defaultz.NewDefaultzExtractor("jsonschema", "default=", ",")),
	)
	// apply defaults
	if err := defaulter.ApplyDefaults(cfg); err != nil {
		return err
	}

	// validate the configuration using `validate` tags
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return err
	}

	return nil
}
