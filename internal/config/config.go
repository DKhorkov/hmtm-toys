package config

import (
	"fmt"
	"time"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/loadenv"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func New() Config {
	return Config{
		Environment: loadenv.GetEnv("ENVIRONMENT", "local"),
		Version:     loadenv.GetEnv("VERSION", "latest"),
		HTTP: HTTPConfig{
			Host: loadenv.GetEnv("HOST", "0.0.0.0"),
			Port: loadenv.GetEnvAsInt("PORT", 8060),
		},
		Database: db.Config{
			Host:         loadenv.GetEnv("POSTGRES_HOST", "0.0.0.0"),
			Port:         loadenv.GetEnvAsInt("POSTGRES_PORT", 5432),
			User:         loadenv.GetEnv("POSTGRES_USER", "postgres"),
			Password:     loadenv.GetEnv("POSTGRES_PASSWORD", "postgres"),
			DatabaseName: loadenv.GetEnv("POSTGRES_DB", "postgres"),
			SSLMode:      loadenv.GetEnv("POSTGRES_SSL_MODE", "disable"),
			Driver:       loadenv.GetEnv("POSTGRES_DRIVER", "postgres"),
			Pool: db.PoolConfig{
				MaxIdleConnections: loadenv.GetEnvAsInt("MAX_IDLE_CONNECTIONS", 1),
				MaxOpenConnections: loadenv.GetEnvAsInt("MAX_OPEN_CONNECTIONS", 1),
				MaxConnectionLifetime: time.Second * time.Duration(
					loadenv.GetEnvAsInt("MAX_CONNECTION_LIFETIME", 20),
				),
				MaxConnectionIdleTime: time.Second * time.Duration(
					loadenv.GetEnvAsInt("MAX_CONNECTION_IDLE_TIME", 10),
				),
			},
		},
		Clients: ClientsConfig{
			SSO: ClientConfig{
				Host:         loadenv.GetEnv("SSO_CLIENT_HOST", "0.0.0.0"),
				Port:         loadenv.GetEnvAsInt("SSO_CLIENT_PORT", 8070),
				RetriesCount: loadenv.GetEnvAsInt("SSO_RETRIES_COUNT", 3),
				RetryTimeout: time.Second * time.Duration(
					loadenv.GetEnvAsInt("SSO_RETRIES_TIMEOUT", 1),
				),
			},
		},
		Logging: logging.Config{
			Level:       logging.Levels.DEBUG,
			LogFilePath: fmt.Sprintf("logs/%s.log", time.Now().UTC().Format("02-01-2006")),
		},
		Tracing: TracingConfig{
			Server: tracing.Config{
				ServiceName:    loadenv.GetEnv("TRACING_SERVICE_NAME", "hmtm-toys"),
				ServiceVersion: loadenv.GetEnv("VERSION", "latest"),
				JaegerURL: fmt.Sprintf(
					"http://%s:%d/api/traces",
					loadenv.GetEnv("TRACING_JAEGER_HOST", "0.0.0.0"),
					loadenv.GetEnvAsInt("TRACING_API_TRACES_PORT", 14268),
				),
			},
			Spans: SpansConfig{
				Root: tracing.SpanConfig{
					Opts: []trace.SpanStartOption{
						trace.WithAttributes(
							attribute.String("Environment", loadenv.GetEnv("ENVIRONMENT", "local")),
						),
					},
					Events: tracing.SpanEventsConfig{
						Start: tracing.SpanEventConfig{
							Name: "Calling handler",
							Opts: []trace.EventOption{
								trace.WithAttributes(
									attribute.String(
										"Environment",
										loadenv.GetEnv("ENVIRONMENT", "local"),
									),
								),
							},
						},
						End: tracing.SpanEventConfig{
							Name: "Received response from handler",
							Opts: []trace.EventOption{
								trace.WithAttributes(
									attribute.String(
										"Environment",
										loadenv.GetEnv("ENVIRONMENT", "local"),
									),
								),
							},
						},
					},
				},
				Repositories: SpanRepositories{
					Categories: tracing.SpanConfig{
						Opts: []trace.SpanStartOption{
							trace.WithAttributes(
								attribute.String(
									"Environment",
									loadenv.GetEnv("ENVIRONMENT", "local"),
								),
							),
						},
						Events: tracing.SpanEventsConfig{
							Start: tracing.SpanEventConfig{
								Name: "Calling database",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String(
											"Environment",
											loadenv.GetEnv("ENVIRONMENT", "local"),
										),
									),
								},
							},
							End: tracing.SpanEventConfig{
								Name: "Received response from database",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String(
											"Environment",
											loadenv.GetEnv("ENVIRONMENT", "local"),
										),
									),
								},
							},
						},
					},
					Tags: tracing.SpanConfig{
						Opts: []trace.SpanStartOption{
							trace.WithAttributes(
								attribute.String(
									"Environment",
									loadenv.GetEnv("ENVIRONMENT", "local"),
								),
							),
						},
						Events: tracing.SpanEventsConfig{
							Start: tracing.SpanEventConfig{
								Name: "Calling database",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String(
											"Environment",
											loadenv.GetEnv("ENVIRONMENT", "local"),
										),
									),
								},
							},
							End: tracing.SpanEventConfig{
								Name: "Received response from database",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String(
											"Environment",
											loadenv.GetEnv("ENVIRONMENT", "local"),
										),
									),
								},
							},
						},
					},
					Toys: tracing.SpanConfig{
						Opts: []trace.SpanStartOption{
							trace.WithAttributes(
								attribute.String(
									"Environment",
									loadenv.GetEnv("ENVIRONMENT", "local"),
								),
							),
						},
						Events: tracing.SpanEventsConfig{
							Start: tracing.SpanEventConfig{
								Name: "Calling database",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String(
											"Environment",
											loadenv.GetEnv("ENVIRONMENT", "local"),
										),
									),
								},
							},
							End: tracing.SpanEventConfig{
								Name: "Received response from database",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String(
											"Environment",
											loadenv.GetEnv("ENVIRONMENT", "local"),
										),
									),
								},
							},
						},
					},
					Masters: tracing.SpanConfig{
						Opts: []trace.SpanStartOption{
							trace.WithAttributes(
								attribute.String(
									"Environment",
									loadenv.GetEnv("ENVIRONMENT", "local"),
								),
							),
						},
						Events: tracing.SpanEventsConfig{
							Start: tracing.SpanEventConfig{
								Name: "Calling database",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String(
											"Environment",
											loadenv.GetEnv("ENVIRONMENT", "local"),
										),
									),
								},
							},
							End: tracing.SpanEventConfig{
								Name: "Received response from database",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String(
											"Environment",
											loadenv.GetEnv("ENVIRONMENT", "local"),
										),
									),
								},
							},
						},
					},
				},
				Clients: SpanClients{
					SSO: tracing.SpanConfig{
						Opts: []trace.SpanStartOption{
							trace.WithAttributes(
								attribute.String(
									"Environment",
									loadenv.GetEnv("ENVIRONMENT", "local"),
								),
							),
						},
						Events: tracing.SpanEventsConfig{
							Start: tracing.SpanEventConfig{
								Name: "Calling gRPC SSO client",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String(
											"Environment",
											loadenv.GetEnv("ENVIRONMENT", "local"),
										),
									),
								},
							},
							End: tracing.SpanEventConfig{
								Name: "Received response from gRPC SSO client",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String(
											"Environment",
											loadenv.GetEnv("ENVIRONMENT", "local"),
										),
									),
								},
							},
						},
					},
				},
			},
		},
		Validation: ValidationConfig{
			Master: MasterValidationConfig{
				Info: loadenv.GetEnvAsSlice(
					"MASTER_INFO_REGEXP",
					[]string{
						`^.{5,1000}$`,     // длина 5-1000 символов
						`^[А-Яа-яЁё\s]+$`, // только кириллица и пробелы
					},
					";",
				),
			},
			Toy: ToyValidationConfig{
				Name: loadenv.GetEnvAsSlice(
					"TOY_NAME_REGEXP",
					[]string{
						`^.{2,50}$`,       // длина 2-50 символов
						`^[А-Яа-яЁё\s]+$`, // только кириллица и пробелы
					},
					";",
				),
				Description: loadenv.GetEnvAsSlice(
					"TOY_DESCRIPTION_REGEXP",
					[]string{
						`^.{10,1000}$`,    // длина 10-1000 символов
						`^[А-Яа-яЁё\s]+$`, // только кириллица и пробелы
					},
					";",
				),
			},
			Tag: TagValidationConfig{
				Name: loadenv.GetEnvAsSlice(
					"TAG_NAME_REGEXP",
					[]string{
						`^.{2,20}$`,       // длина 2-20 символов
						`^[А-Яа-яЁё\s]+$`, // только кириллица и пробелы
					},
					";",
				),
			},
		},
	}
}

type HTTPConfig struct {
	Host string
	Port int
}

type TracingConfig struct {
	Server tracing.Config
	Spans  SpansConfig
}

type SpansConfig struct {
	Root         tracing.SpanConfig
	Repositories SpanRepositories
	Clients      SpanClients
}

type SpanClients struct {
	SSO tracing.SpanConfig
}

type SpanRepositories struct {
	Categories tracing.SpanConfig
	Tags       tracing.SpanConfig
	Masters    tracing.SpanConfig
	Toys       tracing.SpanConfig
}

type ClientsConfig struct {
	SSO ClientConfig
}

type ClientConfig struct {
	Host         string
	Port         int
	RetryTimeout time.Duration
	RetriesCount int
}

type ValidationConfig struct {
	Master MasterValidationConfig
	Toy    ToyValidationConfig
	Tag    TagValidationConfig
}

type MasterValidationConfig struct {
	Info []string // since Go's regex doesn't support backtracking.
}

type ToyValidationConfig struct {
	Name        []string // since Go's regex doesn't support backtracking.
	Description []string // since Go's regex doesn't support backtracking.
}

type TagValidationConfig struct {
	Name []string // since Go's regex doesn't support backtracking.
}

type Config struct {
	HTTP        HTTPConfig
	Clients     ClientsConfig
	Database    db.Config
	Logging     logging.Config
	Tracing     TracingConfig
	Validation  ValidationConfig
	Environment string
	Version     string
}
