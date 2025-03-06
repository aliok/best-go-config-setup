package main

import (
	"encoding/json"
	"github.com/aliok/best-go-config-setup/util"
	"log"
	"os"

	"github.com/invopop/jsonschema"
	"sigs.k8s.io/yaml"

	"github.com/aliok/best-go-config-setup/pkg"
)

// this is the main function for the configbuilder, which would generate the configuration JSON schema and the reference configuration file.
func main() {
	//
	// CREATE THE JSON SCHEMA FOR THE CONFIGURATION
	//

	// we are going to generate the JSON schema for the configuration and write it to configuration-schema.gen.json
	reflector := new(jsonschema.Reflector)
	// treat code comments as JSON schema descriptions
	if err := reflector.AddGoComments("github.com/aliok/best-go-config-setup", "pkg"); err != nil {
		log.Fatalf("Failed to add comments: %v", err)
	}
	// generate the JSON schema
	schema := reflector.Reflect(&pkg.Config{})

	// fix the schema for arrays
	util.VisitSchema(schema, "array", util.FixArrayDefaultValues)

	// marshal the schema to JSON
	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal schema: %v", err)
	}

	// write the schema to a file
	if err := os.WriteFile("configuration-schema.gen.json", schemaJSON, 0644); err != nil {
		log.Fatalf("Failed to write schema to file: %v", err)
	}

	//
	// CREATE THE DEFAULT CONFIG FILE (reference config)
	//

	// create a blank Config instance, then set defaults.
	// this is the reference configuration.
	cfg := pkg.Config{}
	if err := pkg.HandleConfig(&cfg); err != nil {
		log.Fatalf("Error while defaulting or validating the blank config. Are you sure the default values for fields are good?: %v", err)
	}

	// write default config (reference config) to default-config.gen.yaml
	cfgYaml, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatalf("Failed to marshal config to yaml: %v", err)
	}
	// prepend the JSON schema header for IDE support
	cfgYaml = append([]byte("# yaml-language-server: $schema=./configuration-schema.gen.json \n"), cfgYaml...)

	// write to file
	if err := os.WriteFile("default-config.gen.yaml", cfgYaml, 0644); err != nil {
		log.Fatalf("Failed to write config to file: %v", err)
	}
}
