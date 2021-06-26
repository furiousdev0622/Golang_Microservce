package xservice

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/cloudflare/tableflip"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	gwrt "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	echomd "github.com/labstack/echo/v4/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/soheilhy/cmux"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/xinpianchang/xservice/core"
	"github.com/xinpianchang/xservice/core/middleware"
	"github.com/xinpianchang/xservice/pkg/echox"
	"github.com/xinpianchang/xservice/pkg/grpcx"
	"github.com/xinpianchang/xservice/pkg/log"
	"github.com/xinpianchang/xservice/pkg/signalx"
	"github.com/xinpianchang/xservice/pkg/tracingx"
)

type Server interface {
	Echo() *echo.Echo
	Serve() error
	GrpcRegister(desc *grpc.ServiceDesc, impl interface{}, handler ...GrpcRegisterHandler)
}

type grpcService struct {
	Desc    *grpc.ServiceDesc
	Impl    interface{}
	Handler GrpcRegisterHandler
}

type GrpcRegisterHandler func(ctx context.Context, mux *gwrt.ServeMux, conn *grpc.ClientConn) error

type serverImpl struct {
	options      *Options
	echo         *echo.Echo
	grpc         *grpc.Server
	grpcServices []*grpcService
}

func newServer(opts *Options) Server {
	server := &serverImpl{
		grpcServices: make([]*grpcService, 0, 128),
	}
	server.options = opts

	server.initGrpc()
	server.initEcho()

	return server
}

func (t *serverImpl) Echo() *echo.Echo {
	return t.echo
}

func (t *serverImpl) Serve() error {
	address := t.getHttpAddress()

	log.Info("server start",
		zap.String("name", t.options.Name),
		zap.Int("pid", os.Getpid()),
		zap.String("version", t.options.Version),
		zap.String("build", t.options.Build),
		zap.String("address", address),
		zap.String("runtime", runtime.Version()),
	)

	upg, err := tableflip.New(tableflip.Options{
		UpgradeTimeout: time.Minute,
	})
	if err != nil {
		log.Fatal("tableflip init", zap.Error(err))
	}
	defer upg.Stop()

	t.waitSignalForTableflip(upg)

	ln, err := upg.Fds.Listen("tcp", address)
	if err != nil {
		log.Fatal("listen", zap.Error(err))
	}
	defer ln.Close()

	mux := cmux.New(ln)
	defer mux.Close()

	grpcL := mux.Match(cmux.HTTP2())
	defer grpcL.Close()

	httpL := mux.Match(cmux.HTTP1Fast())
	defer httpL.Close()

	if len(t.grpcServices) > 0 {
		go t.serveGrpc(grpcL)
	}

	server := http.Server{
		Handler:           t.echo,
		ReadHeaderTimeout: time.Second * 30,
		IdleTimeout:       time.Minute * 1,
	}

	go func() {
		if err := server.Serve(httpL); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal("start http server", zap.Error(err))
			}
		}
	}()

	go func() {
		_ = mux.Serve()
	}()

	if err = upg.Ready(); err != nil {
		log.Fatal("ready", zap.Error(err))
	}

	// all ready
	t.registerGrpcServiceEtcd()

	signalx.AddShutdownHook(func(os.Signal) {
		_ = server.Shutdown(context.Background())
		sentry.Flush(time.Second * 2)
		log.Info("shutdown", zap.Int("pid", os.Getpid()))
	})

	<-upg.Exit()

	signalx.Shutdown()

	return nil
}

func (t *serverImpl) GrpcRegister(desc *grpc.ServiceDesc, impl interface{}, hs ...GrpcRegisterHandler) {
	var handler GrpcRegisterHandler
	if len(hs) > 0 {
		handler = hs[0]
	}
	t.grpcServices = append(t.grpcServices, &grpcService{desc, impl, handler})
}

func (t *serverImpl) getHttpAddress() string {
	return t.options.Config.GetString(core.ConfigServiceAddr)
}

func (t *serverImpl) waitSignalForTableflip(upg *tableflip.Upgrader) {
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGUSR2, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		for s := range sig {
			switch s {
			case syscall.SIGUSR2, syscall.SIGHUP:
				err := upg.Upgrade()
				if err != nil {
					log.Error("upgrade failed", zap.Error(err))
					continue
				}
				log.Info("upgrade succeeded", zap.Int("pid", os.Getpid()))
				return
			default:
				upg.Stop()
			}
		}
	}()
}

func (t *serverImpl) initEcho() {
	e := echo.New()
	t.echo = e

	e.Logger = log.NewEchoLogger()
	e.IPExtractor = echo.ExtractIPFromXFFHeader(echo.TrustPrivateNet(true))
	echox.ConfigValidator(e)
	e.HTTPErrorHandler = echox.HTTPErrorHandler

	// recover
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if x := recover(); x != nil {
					log.For(c.Request().Context()).Error("server panic error", zap.Any("error", x))
					_ = c.String(http.StatusInternalServerError, fmt.Sprint("internal server error, ", x))
				}
			}()
			return next(&echoContext{c})
		}
	})

	e.Use(echomd.RequestID())
	e.Use(middleware.Trace(t.options.Config.GetBool("jaeger.body_dump"), t.options.EchoTracingSkipper))
	e.Use(sentryecho.New(sentryecho.Options{Repanic: true}))

	// logger id & traceId & server-info
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("X-Service", fmt.Sprint(t.options.Name, "/", t.options.Version, "/", t.options.Build))
			id := c.Request().Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = c.Response().Header().Get(echo.HeaderXRequestID)
			}
			c.Set(echo.HeaderXRequestID, id)
			ctx := context.WithValue(c.Request().Context(), core.ContextHeaderXRequestID, id)
			c.SetRequest(c.Request().WithContext(ctx))

			traceId := tracingx.GetTraceID(c.Request().Context())
			if traceId != "" {
				c.Response().Header().Set("X-Trace-Id", traceId)
			}

			if span := opentracing.SpanFromContext(c.Request().Context()); span != nil {
				span.SetTag("requestId", id)
				span.SetTag("ip", c.RealIP())
			}

			if hub := sentryecho.GetHubFromContext(c); hub != nil {
				scope := hub.Scope()
				scope.SetTag("ip", c.RealIP())
				scope.SetTag("X-Forwarded-For", c.Request().Header.Get("X-Forwarded-For"))
				if traceId != "" {
					scope.SetTag("traceId", traceId)
				}
			}

			return next(c)
		}
	})

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.Group("/debug/*", middleware.Pprof())
}

// init grpc
// add middleware https://github.com/grpc-ecosystem/go-grpc-middleware
func (t *serverImpl) initGrpc() {
	options := make([]grpc.ServerOption, 0, 8)
	options = append(options,
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
			grpc_opentracing.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			grpcx.EnvoyproxyValidatorStreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpcx.EnvoyproxyValidatorUnaryServerInterceptor(),
		)),
	)
	options = append(options, t.options.GrpcOptions...)
	g := grpc.NewServer(options...)
	t.grpc = g
}

func (t *serverImpl) serveGrpc(ln net.Listener) {
	for _, service := range t.grpcServices {
		t.grpc.RegisterService(service.Desc, service.Impl)
		// log.Debug("register grpc service", zap.String("impl", reflect.TypeOf(service.Impl).String()))
	}

	go func() {
		_ = t.grpc.Serve(ln)
	}()

	address := t.getHttpAddress()
	grpcClientConn, err := grpc.DialContext(
		context.Background(),
		address,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			grpc_opentracing.UnaryClientInterceptor(
				grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
			),
		),
	)

	if err != nil {
		log.Fatal("grpc gateway client conn", zap.Error(err))
	}

	grpcGateway := gwrt.NewServeMux()

	for _, service := range t.grpcServices {
		if service.Handler == nil {
			continue
		}
		err := service.Handler(context.Background(), grpcGateway, grpcClientConn)
		if err != nil {
			log.Fatal("grpc register handler", zap.Error(err))
		}
		// log.Debug("register grpc gateway", zap.String("handler", runtime.FuncForPC(reflect.ValueOf(service.Handler).Pointer()).Name()))
	}

	t.echo.Group("/rpc/*", echo.WrapMiddleware(func(handler http.Handler) http.Handler {
		return grpcGateway
	}))
}

// registerGrpcServiceEtcd
// refer: https://etcd.io/docs/v3.5/dev-guide/grpc_naming/
func (t *serverImpl) registerGrpcServiceEtcd() {
	if len(t.grpcServices) == 0 {
		return
	}

	if os.Getenv(core.EnvEtcd) == "" {
		log.Warn("etcd not configured, service register ignored")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	go t.doRegisterGrpcServiceEtcd(ctx)

	signalx.AddShutdownHook(func(s os.Signal) {
		cancel()
		// deregister
		log.Debug("deregister service")
		client, err := serviceEtcdClient()
		if err != nil {
			return
		}

		em, _ := endpoints.NewManager(client, core.ServiceRegisterKeyPrefix)

		for _, service := range t.grpcServices {
			_ = em.DeleteEndpoint(context.Background(), serviceKey(os.Getenv(core.EnvServiceName), service.Desc))
		}
	})
}

func (t *serverImpl) doRegisterGrpcServiceEtcd(ctx context.Context) {
	l := log.Named("registerGrpcServiceEtcd")
	defer func() {
		if x := recover(); x != nil {
			l.Error("recover", zap.Any("err", x))
			sentry.CaptureException(errors.WithStack(errors.New(fmt.Sprint(x))))

			time.Sleep(time.Second * 10)
			go t.doRegisterGrpcServiceEtcd(ctx)
		}
	}()

	client, err := serviceEtcdClient()
	if err != nil {
		l.Error("get client", zap.Error(err))
		return
	}

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	ttl := int64(10) // seconds

	type serviceLease struct {
		id       clientv3.LeaseID
		lease    clientv3.Lease
		key      string
		endpoint endpoints.Endpoint
	}

	leaseMap := make(map[string]*serviceLease, len(t.grpcServices))
	getLease := func(desc *grpc.ServiceDesc) *serviceLease {
		if sl, ok := leaseMap[desc.ServiceName]; ok {
			return sl
		}
		addr := t.options.Config.GetString(core.ConfigServiceAdviceAddr)
		sl := &serviceLease{
			id:    0,
			lease: clientv3.NewLease(client),
			key:   serviceKey(os.Getenv(core.EnvServiceName), desc),
			endpoint: endpoints.Endpoint{
				Addr:     addr,
				Metadata: desc.Metadata,
			},
		}
		leaseMap[desc.ServiceName] = sl
		return sl
	}

	em, _ := endpoints.NewManager(client, core.ServiceRegisterKeyPrefix)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			for _, service := range t.grpcServices {
				sl := getLease(service.Desc)
				ll := l.With(zap.String("service", sl.key))
				// ll.Debug("register")
				if sl.id == 0 {
					leaseRsp, err := sl.lease.Grant(context.Background(), ttl)
					if err != nil {
						ll.Error("lease.Grant", zap.Error(err))
						continue
					}

					err = em.AddEndpoint(context.Background(), sl.key, sl.endpoint, clientv3.WithLease(leaseRsp.ID))
					if err != nil {
						ll.Error("kv.Put", zap.Error(err))
					}
					sl.id = leaseRsp.ID
				} else {
					_, err = sl.lease.KeepAliveOnce(context.Background(), sl.id)
					if err != nil {
						sl.id = 0
					}
				}
			}
		}

		// wait next loop
		<-ticker.C
	}
}

type echoContext struct {
	echo.Context
}

func (t *echoContext) Logger() echo.Logger {
	logger := t.Context.Logger()
	if l, ok := logger.(*log.EchoLogger); ok {
		return l.For(t.Request().Context())
	}
	return logger
}

func (t *echoContext) Path() string {
	return t.Request().URL.Path
}
