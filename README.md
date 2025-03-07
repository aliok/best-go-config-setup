# A Robust Configuration Management Setup in Golang

Configuration management is a critical part of any application. It needs to be flexible, maintainable, and developer-friendly. This blueprint describes a powerful setup in Golang that allows for reading configuration files, setting defaults, validating inputs, and even generating JSON schemas for better user experience.

The code is available, with 2 entry points:

#### **[cmd/app/main.go](cmd/app/main.go)** 
  
This is the program that reads the configuration and runs the necessary business logic. It will read the config, set defaults, validate it, and then run the application.

#### **[cmd/configbuilder/main.go](cmd/configbuilder/main.go)** 

The entry point for the configuration builder, which generates the JSON schema for the configuration.

## The Evolution of This Configuration Setup

Initially, managing configuration in Go projects was straightforward but limited. I used environment variables and command-line flags for configuration, but this approach had several drawbacks:

- No structured configuration file support.
- Lack of validation meant that incorrect configurations could cause unexpected behavior.
- Default values had to be manually set in multiple places.
- Developers lacked auto-completion and validation support in their IDEs.

To address these issues, I gradually evolved my configuration setup to be more structured, using a combination of powerful libraries:

1. **Reading configuration from files with overrides from environment variables and flags** ([spf13/viper](https://github.com/spf13/viper))
2. **Automatically filling in default values from struct tags** ([go-defaultz](https://github.com/aliok/go-defaultz))
3. **Validating configuration fields with rules defined in struct tags** ([go-playground/validator](https://github.com/go-playground/validator/))
4. **Generating JSON schemas for IDE support and validation** ([invopop/jsonschema](https://github.com/invopop/jsonschema/))

## The Final Configuration Setup

This refined approach ensures:
- Users only need to define the values they want to override.
- Default values and validation rules are in a single place, reducing duplication and maintenance overhead as Ill as cognitive load.
- IDEs provide auto-completion and validation via JSON schema.

### Example Configuration

#### Default Configuration File (Generated or Provided)
```yaml
features:
  enabled_features:
  - feature1
  - feature2
http_server:
  bind_address: 0.0.0.0
  port: 8080
logging:
  log_format: json
  log_level: 2
```

#### User-Provided Configuration File (Overrides Defaults)
```yaml
http_server:
  port: 12345
features:
  enabled_features:
  - feature3
  - feature4
```

## How It Works and the Evolution

### 1. Reading Configuration

I use [**spf13/viper**](https://github.com/spf13/viper) to read configuration files while allowing environment variable and flag-based overrides.

Instead of using `GetXXX()` methods like `GetInt("http_server.port")`, I use a struct to define the configuration and bind it to Viper. 

This way, I can access configuration values directly through the struct.

```go
type Config struct {
    HTTPServerConfig HTTPServerConfig `mapstructure:"http_server"`
    // ...
}

type HTTPServerConfig struct {
    Port int `mapstructure:"port,omitempty"`
    BindAddress string `mapstructure:"bind_address,omitempty"`
    // ...
}

func main() {
    // ... <code to set up Viper>

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        log.Fatalf("Error unmarshaling config: %s", err)
    }
}
```

This is great for readability and maintainability, as the configuration struct serves as a single source of truth for the configuration schema.

Viper will fill in the configuration struct with values from the configuration file, environment variables, and flags.
If a field is not present in the configuration file, it will remain at its zero value.

However, Viper doesn't have a way to set up default values or do validation. I would have to manually go through the fields and set defaults, which is error-prone and tedious.
That's where the next steps come in.

One minor change before moving on, I change the `mapstructure` tag to `json` to be able to marshal the struct to JSON and YAML later.

```go
type Config struct {
    HTTPServerConfig HTTPServerConfig `json:"http_server"`
    // ...
}

type HTTPServerConfig struct {
    Port int `json:"port,omitempty"`
    BindAddress string `json:"bind_address,omitempty"`
    // ...
}

func main() {
    // ... <code to set up Viper>

    viperOpt := func(dc *mapstructure.DecoderConfig) {
        dc.TagName = "json"
    }
    var config Config
    if err := viper.Unmarshal(&cfg, viperOpt); err != nil {
        log.Fatalf("Error unmarshaling config: %s", err)
    }
}
```

### 2. Setting Default Values

Instead of manually handling defaults in multiple places, I use [**aliok/go-defaultz**](https://github.com/aliok/go-defaultz) to extract default values from struct tags and apply them automatically.

This would need setting up some tags to denote default values in the configuration struct.

```go
type HTTPServerConfig struct {
    Port int           `json:"port,omitempty"         default:"8080"`
    BindAddress string `json:"bind_address,omitempty" default:"0.0.0.0"`
    // ...
}

func main() {
    var cfg Config
    // ... <code to read config>

    if err := defaultz.ApplyDefaults(&cfg); err != nil {
        return err
    }
}
```

This would set the default values of fields if they are not provided in the configuration file.

For example, for this config file:
```yaml
http_server:
  port: 99999
```

the struct would be:
```text
HTTPServerConfig{
    Port: 99999,
    BindAddress: "0.0.0.0" 
}
```

It is great! I don't have to manually set defaults for each field, and the defaults are right there in the struct definition.
Plus, only the fields that are not provided in the configuration file will get the default values.

### 3. Validating Configuration
Using **go-playground/validator**, I enforce validation rules through struct tags. This ensures that invalid configurations are caught early.

To do that, let's add `validate` tags to the struct fields.

```go
type HTTPServerConfig struct {
    Port int           `json:"port,omitempty"         default:"8080"    validate:"required,min=1024,max=65535"`
    BindAddress string `json:"bind_address,omitempty" default:"0.0.0.0" validate:"required,ip4_addr"`
    // ...
}

func main() {
    var cfg Config
    // ... <code to read the config and set defaults>
	
    // validate the configuration using `validate` tags
    validate := validator.New()
    if err := validate.Struct(cfg); err != nil {
        return err
    }
}
```

This will ensure that the configuration is valid before the application starts. If a field is missing or doesn't meet the validation rules, an error will be returned.

### 4. Generating JSON Schemas

I leverage **invopop/jsonschema** to generate a JSON schema from my configuration struct, enabling better IDE support through auto-completion and validation.

The JSON schema is generated using the `jsonschema` tags.

```go
type HTTPServerConfig struct {
    Port int           `json:"port,omitempty"         default:"8080"    validate:"required,min=1024,max=65535" jsonschema:"default=8080,title=Port,description=Port number of the HTTP server"`
    BindAddress string `json:"bind_address,omitempty" default:"0.0.0.0" validate:"required,ip4_addr"           jsonschema:"default=0.0.0.0,title=Bind Address,description=IP address to bind the HTTP server to"`
    // ...
}

func main() {
    // NOTE: we don't read any config here, since we want the JSON schema for the configuration struct. we just use an empty struct.
    var cfg Config
    // ... <code to set defaults>

    reflector := new(jsonschema.Reflector)
    schema := reflector.Reflect(&pkg.Config{})

    // marshal the schema to JSON
    schemaJSON, err := json.MarshalIndent(schema, "", "  ")
    if err != nil {
        log.Fatalf("Failed to marshal schema: %v", err)
    }

    // write the schema to a file
    if err := os.WriteFile("configuration-schema.gen.json", schemaJSON, 0644); err != nil {
        log.Fatalf("Failed to write schema to file: %v", err)
    }
}
```

This will create a JSON schema file:

```json
{
  // ...
  "$defs": {
    "Config": {
      // ...
    },
    "HTTPServerConfig": {
      "properties": {
        "port": {
          "type": "integer",
          "description": "Port number of the HTTP server",
          "default": 8080
        },
        "bind_address": {
          "type": "string",
          "description": "IP address to bind the HTTP server to",
          "default": "0.0.0.0"
        }
      },
      // ...
    }
  }
}
```

This is nice!

However, I have an issue! The default value is defined twice, once of the JSON schema, once for the defaulting operation:

```go
    Port int `<...> default:"8080" validate:"..." jsonschema:"default=8080,..."`
```

To get rid of that, we can make `go-defaultz` to use the default value from the `jsonschema` tag:

```go

type HTTPServerConfig struct {
    Port int           `<...> jsonschema:"default=8080,title=Port,description=Port number of the HTTP server"`
    BindAddress string `<...> jsonschema:"default=0.0.0.0,title=Bind Address,description=IP address to bind the HTTP server to"`
    // ...
}

func main() {
    // NOTE: we don't read any config here, since we want the JSON schema for the configuration struct. we just use an empty struct.
    var cfg Config

    // use go-defaultz to apply defaults
    // reuse the `jsonschema` tag and the `default=` prefix
    defaulter := defaultz.NewDefaulterRegistry(
        defaultz.WithBasicDefaulters(),
        defaultz.WithDefaultExtractor(defaultz.NewDefaultzExtractor("jsonschema", "default=", ",")),
    )

    if err := defaulter.ApplyDefaults(&cfg); err != nil {
        return err
    }

    // ... <code to validate the configuration>
    // ... <code to generate the JSON schema>
    // ... <code to write the JSON schema to a file>
}
```

This way, the default value is defined only once, in the `jsonschema` tag.

The JSON schema generator, `invopop/jsonschema`, can actually read code comments and fill in the `description` field in the JSON schema. This is great for providing more context to users.

```go
type HTTPServerConfig struct {
    // Port number of the HTTP server
    Port int           `<...> jsonschema:"default=8080`
    
    // IP address to bind the HTTP server to
    BindAddress string `<...> jsonschema:"default=0.0.0.0"`
    // ...
}

func main(){
    // ... <code to create the config and set defaults>

    reflector := new(jsonschema.Reflector)

    // treat code comments as JSON schema descriptions
    if err := reflector.AddGoComments("github.com/aliok/best-go-config-setup", "pkg"); err != nil {
        log.Fatalf("Failed to add comments: %v", err)
    }

    // ... <code to generate the JSON schema and write to a file>
}
```

This will generate a JSON schema with descriptions, just like before, but with the descriptions filled in from the code comments:

```json
{
  // ...
  "$defs": {
    "Config": {
      // ...
    },
    "HTTPServerConfig": {
      "properties": {
        "port": {
          "type": "integer",
          "description": "Port number of the HTTP server",
          "default": 8080
        },
        "bind_address": {
          "type": "string",
          "description": "IP address to bind the HTTP server to",
          "default": "0.0.0.0"
        }
      },
      // ...
    }
  }
}
```

### 5. Final touches

Everything is good until now. However, there are 2 things I would like to add:
- The `default` field in JSON schema is messed up for slices.
- I would like to generate a reference configuration file with the default values filled in, so that users can see the default values.

To fix the first issue, we need to process the JSON schema and fix the `default` field for slices. See [./util/jsonschema.go](util/jsonschema.go) for that.

To get the second one done, I simply marshal the configuration struct to YAML and write it to a file.

I use [sigs.k8s.io/yaml](https://github.com/kubernetes-sigs/yaml) package for that person due to 2 reasons:
1. It can marshal the struct to YAML with the struct tags of `json`, unlike the `gopkg.in/yaml.v2` package that requires separate `yaml` tags. This way, I can reuse the `json` tags for both JSON and YAML.
2. It fixes some of the rendering/formatting issues with floats. For example the large floats are rendered in scientific notation with `gopkg.in/yaml.v2`.

### 6. Best Practices

See https://github.com/aliok/go-defaultz?tab=readme-ov-file#best-practices for some best practices around fields, their zero-values and defaulting them.

## Benefits of This Approach
- **Simplified configuration management**: Users only specify what they want to change.
- **Reduced duplication**: Defaults and validation live in a single place.
- **Better developer experience**: IDE support with JSON schemas.
- **More reliable applications**: Validation ensures correctness before the application starts.

If you’re working on a Go project, adopting a similar approach could save you a lot of hassle!

