package xservice

import (
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/xinpianchang/xservice/core"
	"github.com/xinpianchang/xservice/pkg/config"
	"github.com/xinpianchang/xservice/pkg/netx"
)

type Options struct {
	Name               string
	Version            string
	Build              string
	Description        string
	Config             *viper.Viper
	GrpcOptions        []grpc.ServerOption
	SentryOptions      sentry.ClientOptions
	EchoTracingSkipper middleware.Skipper
}

type Option func(*Options)

func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func Version(version string) Option {
	return func(o *Options) {
		o.Version = version
	}
}

func Build(build string) Option {
	return func(o *Options) {
		o.Build = build
	}
}

func Description(description string) Option {
	return func(o *Options) {
		o.Description = description
	}
}

func Config(config *viper.Viper) Option {
	return func(o *Options) {
		o.Config = config
	}
}

func WithGrpcOptions(options ...grpc.ServerOption) Option {
	return func(o *Options) {
		o.GrpcOptions = options
	}
}

func WithSentry(options sentry.ClientOptions) Option {
	return func(o *Options) {
		o.SentryOptions = options
	}
}

func WithEchoTracingSkipper(skipper middleware.Skipper) Option {
	return func(o *Options) {
		o.EchoTracingSkipper = skipper
	}
}

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	loadEnvOptions(opts)

	for _, option := range options {
		option(opts)
	}

	if opts.Name == "" {
		opts.Name = core.DefaultServiceName
	}

	nameexp := `^[a-zA-Z0-9\-\_\.]+$`
	if ok, _ := regexp.MatchString(nameexp, opts.Name); !ok {
		log.Fatal("invalid service name", zap.String("name", opts.Name), zap.String("suggest", nameexp))
	}
	os.Setenv(core.EnvServiceName, opts.Name)

	if opts.Version == "" {
		opts.Version = "v0.0.1"
	}
	os.Setenv(core.EnvServiceVersion, opts.Version)

	if opts.Build == "" {
		opts.Build = fmt.Sprint("dev-", time.Now().UnixNano())
	}

	if opts.Config == nil {
		opts.loadConfig()
	}

	if opts.Config.IsSet(core.ConfigServiceAddr) {
		addviceAddr := opts.Config.GetString(core.ConfigServiceAdviceAddr)
		if addviceAddr == "" {
			address := opts.Config.GetString(core.ConfigServiceAddr)
			_, port, err := net.SplitHostPort(address)
			if err != nil {
				log.Fatal("invalid address", zap.Error(err))
			}
			addviceAddr = net.JoinHostPort(netx.InternalIp(), port)

			opts.Config.SetDefault(core.ConfigServiceAdviceAddr, addviceAddr)
		}
	}

	return opts
}

func loadEnvOptions(opts *Options) {
	if opts.Name == "" {
		opts.Name = os.Getenv(core.EnvServiceName)
	}

	if opts.Version == "" {
		opts.Version = os.Getenv(core.EnvServiceVersion)
	}
}

func (t *Options) loadConfig() {
	if err := config.LoadGlobal(); err != nil {
		log.Fatal("load config", zap.Error(err))
	}
	t.Config = viper.GetViper()
}
