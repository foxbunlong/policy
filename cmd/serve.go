package cmd

import (
	"fmt"
	"io"

	"github.com/oeoen/policy/driver"
	"github.com/oeoen/policy/driver/config"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/uber/jaeger-client-go"
	jConfig "github.com/uber/jaeger-client-go/config"
)

// migrateCmd represents the migrate command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve service",
	Run:   serveCommand,
}

func init() {

	RootCmd.AddCommand(serveCmd)
}

func serveCommand(cmd *cobra.Command, args []string) {
	// var d driver.Driver = driver.NewDefaultDriver(logrusx.New(), false)
	d := driver.NewDefaultDriver(logrus.New(), false)
	closer, err := initJaeger(d.Configuration())
	if err != nil {
		d.Configuration().Logger().Println(logrus.Fields{"error": err}, "init_tracer_failed")
	}
	defer func() {
		if closer != nil {
			closer.Close() //nolint: errcheck
		}
	}()
	errC := make(chan error)
	go func() {
		err := d.Registry().Handler().Serve()
		d.Configuration().Logger().Error(err)
		errC <- err
	}()
	<-errC
	return
}

func initJaeger(cn config.Provider) (closer io.Closer, err error) {
	tracer, closer, err := cn.TracingJaegerConfig().NewTracer(jConfig.Logger(jaeger.StdLogger))

	if err != nil {
		err = fmt.Errorf("Could not initialize jaeger tracer: %+v\n", err)
		return
	}
	opentracing.SetGlobalTracer(tracer)
	return
}
