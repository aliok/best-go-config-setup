package main

import (
	"fmt"
	"log"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"sigs.k8s.io/yaml"

	"github.com/aliok/best-go-config-setup/pkg"
)

// this is the main function for the application, which would run some business logic with the loaded configuration.
func main() {
	// viper should use app-config.yaml file as the configuration file in the current directory
	viper.SetConfigName("app-config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// read the config file (the location of the file should be set by the caller)
	if err := viper.ReadInConfig(); err != nil {
		// ok to not have a config file
		log.Printf("Failed to read config file, going to use defaults: %v", err)
	}

	// optionally, override the config with environment variables
	// viper.AutomaticEnv()

	// configure viper to use the `json` tag
	viperOpt := func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "json"
	}
	// Unmarshal into struct using Viper
	var cfg pkg.Config
	if err := viper.Unmarshal(&cfg, viperOpt); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Set default values for the configuration and validate it
	if err := pkg.HandleConfig(&cfg); err != nil {
		log.Fatalf("Failed to handle config: %v", err)
	}

	// output the loaded configuration
	cfgYaml, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatalf("Failed to marshal config to yaml: %v", err)
	}
	fmt.Printf("Read config\n%s\n", string(cfgYaml))
	// Outputs as:
	// Read config
	// features:
	//  enabled_features:
	//  - feature3
	//  - feature4
	// http_server:
	//  bind_address: 0.0.0.0
	//  port: 12345
	// logging:
	//  log_format: json
	//  log_level: 2

	// note that `port` and `enabled_features` fields are set to what is in the configuration file `app-config.yaml`.
	// other fields are set to their default values.

	// you can change the configuration file and run the program again to see the changes.
	// try setting values that would fail the validation, like setting `port` to 0.

	// ...
	// run business logic with the loaded configuration
	// ...

}
