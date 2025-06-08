package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/invopop/validation"
	"github.com/spf13/viper"
)

type Config struct {
	App           AppConfig
	Logger        LoggerConfig
	Postgres      PostgresConfig
	Observability ObservabilityConfig
	Redis         RedisConfig
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.App),
		validation.Field(&c.Logger),
		validation.Field(&c.Postgres),
		validation.Field(&c.Observability),
		validation.Field(&c.Redis),
	)
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Port int64  `mapstructure:"port"`
	Env  string `mapstructure:"env"`
	Key  string `mapstructure:"key"`
}

func (ac AppConfig) Validate() error {
	return validation.ValidateStruct(&ac,
		validation.Field(&ac.Name, validation.Required),
		validation.Field(&ac.Port, validation.Required, validation.Min(1), validation.Max(65535)),
		validation.Field(&ac.Env, validation.Required, validation.In("local", "production", "staging")),
		validation.Field(&ac.Key, validation.Required, validation.Length(32, 64)),
	)
}

type LoggerConfig struct {
	Mode   string `mapstructure:"mode"`
	Level  string `mapstructure:"level"`
	Enable *bool  `mapstructure:"enable"`
}

func (lc LoggerConfig) Validate() error {
	return validation.ValidateStruct(&lc,
		validation.Field(&lc.Mode, validation.Required),
		validation.Field(&lc.Level, validation.Required, validation.In("debug", "info", "warn", "error", "fatal")),
		validation.Field(&lc.Enable, validation.NotNil),
	)
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int64  `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
	SSLMode  string `mapstructure:"ssl_mode"`
	MaxConn  int64  `mapstructure:"max_conn"`
}

func (pc PostgresConfig) Validate() error {
	return validation.ValidateStruct(&pc,
		validation.Field(&pc.Host, validation.Required),
		validation.Field(&pc.Port, validation.Required, validation.Min(1), validation.Max(65535)),
		validation.Field(&pc.User, validation.Required),
		validation.Field(&pc.Password, validation.Required),
		validation.Field(&pc.DBName, validation.Required),
		validation.Field(&pc.SSLMode, validation.Required, validation.In("disable", "require", "verify-ca", "verify-full")),
		validation.Field(&pc.MaxConn, validation.Required, validation.Min(1), validation.Max(100)),
	)
}

type ObservabilityConfig struct {
	Enable       *bool  `mapstructure:"enable"`
	Mode         string `mapstructure:"mode"`
	OtlpEndpoint string `mapstructure:"otlp_endpoint"`
}

func (oc ObservabilityConfig) Validate() error {
	return validation.ValidateStruct(&oc,
		validation.Field(&oc.Enable, validation.NotNil),
		validation.Field(&oc.Mode, validation.Required, validation.In("otlp", "jaeger", "zipkin")),
		validation.Field(&oc.Enable, validation.NotNil),
	)
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int64  `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int64  `mapstructure:"database"`
}

func (rc RedisConfig) Validate() error {
	return validation.ValidateStruct(&rc,
		validation.Field(&rc.Host, validation.Required),
		validation.Field(&rc.Port, validation.Required),
		validation.Field(&rc.Database, validation.Min(0), validation.Max(15)),
	)
}

var (
	global *Config
)

func MustGet() Config {
	if global == nil {
		panic("must call LoadConfig/LoadConfigPath first")
	}

	return *global
}

func Get() *Config {
	return global
}

func LoadConfig(env string) (Config, error) {
	if global == nil {
		v := viper.New()
		v.SetConfigName(fmt.Sprintf("config/api/config-%s", env))
		v.AddConfigPath(".")
		v.AutomaticEnv()
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		bindEnvs(v, Config{})

		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Println("Config file not found, failing back to environment variables")
			} else {
				return Config{}, err
			}
		}

		var cfg Config
		err := v.Unmarshal(&cfg, func(dc *mapstructure.DecoderConfig) {
			dc.TagName = "mapstructure"
			dc.MatchName = looseMatch
		})
		if err != nil {
			log.Printf("unable to decode into struct, %v", err)
			return Config{}, err
		}
		err = cfg.Validate()
		if err != nil {
			return Config{}, err
		}
		global = &cfg

		os.Setenv("APP_ENV", cfg.App.Env)
		return cfg, nil
	}

	return *global, nil
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToUpper(snake)
}

func bindEnvs(v *viper.Viper, iface interface{}, parts ...string) {
	ifaceVal := reflect.ValueOf(iface)
	if ifaceVal.Kind() == reflect.Ptr {
		ifaceVal = ifaceVal.Elem()
	}

	ifaceType := ifaceVal.Type()
	for i := 0; i < ifaceVal.NumField(); i++ {
		field := ifaceVal.Field(i)
		typeField := ifaceType.Field(i)

		tag := toSnakeCase(typeField.Name)
		envKey := strings.ToUpper(strings.Join(append(parts, tag), "_"))
		viperKey := strings.ToLower(strings.Join(append(parts, tag), "."))

		_ = v.BindEnv(viperKey, envKey)
		if val, ok := os.LookupEnv(envKey); ok {
			v.Set(viperKey, val)
		}

		if field.Kind() == reflect.Struct && typeField.Type.Name() != "Time" {
			bindEnvs(v, field.Interface(), append(parts, tag)...)
		}
	}
}

func looseMatch(mapKey, fieldName string) bool {
	normalize := func(s string) string {
		return strings.ToLower(strings.ReplaceAll(s, "_", ""))
	}
	return normalize(mapKey) == normalize(fieldName)
}
