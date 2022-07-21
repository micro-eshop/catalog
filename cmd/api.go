package cmd

import (
	"context"
	"flag"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/subcommands"
	"github.com/micro-eshop/catalog/internal/postgres"
	"github.com/micro-eshop/catalog/pkg/core/services"
	"github.com/micro-eshop/catalog/pkg/core/usecase"
	"github.com/micro-eshop/catalog/pkg/handlers"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	ginlogrus "github.com/toorop/gin-logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
)

type RunApiCmd struct {
	addr         string
	postgresConn string
	v            *viper.Viper
}

func NewRunApiCmd(v *viper.Viper) *RunApiCmd {
	return &RunApiCmd{v: v}
}

func (*RunApiCmd) Name() string     { return "run-api" }
func (*RunApiCmd) Synopsis() string { return "Run catalog api" }
func (*RunApiCmd) Usage() string {
	return `run-api`
}

func (p *RunApiCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.addr, "addr", ":8080", "address to listen")
	f.StringVar(&p.postgresConn, "postgresConn", p.v.GetString("POSTGRES_CONNECTION"), "postgresConn connection string")
}

func initLogger() *logrus.Logger {
	logger := log.New()
	// Log as JSON instead of the default ASCII formatter.
	logger.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logger.SetLevel(log.InfoLevel)
	logger.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	)))
	return logger
}

func (p *RunApiCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	shutdown := handlers.InitPrivder(ctx)
	defer shutdown(ctx)

	r := gin.New()
	log := initLogger()
	r.Use(ginlogrus.Logger(log), gin.Recovery())
	r.Use(otelgin.Middleware("catalog", otelgin.WithTracerProvider(otel.GetTracerProvider())))
	log.Infoln("Start import products")
	postgresClient, err := postgres.NewPostgresClient(ctx, p.postgresConn)
	if err != nil {
		log.WithError(err).Error("can't create postgres client")
		return subcommands.ExitFailure
	}
	defer postgresClient.Close(ctx)
	repo := postgres.NewPostgresCatalogRepository(postgresClient)
	service := services.NewCatalogService(repo)

	getById := usecase.NewGetProductByIdUseCase(service)
	getByIds := usecase.NewGetProductByIdsUseCase(service)
	catalog := handlers.NewCatalogHandler(getById, getByIds)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	catalog.Setup(r)
	if err := r.Run(p.addr); err != nil {
		log.WithError(err).WithContext(ctx).Errorln("failed to run api")
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
