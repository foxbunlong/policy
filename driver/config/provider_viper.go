package config

import (
	"strings"

	"github.com/ory/viper"
	"github.com/ory/x/viperx"
	"github.com/sirupsen/logrus"
	jConfig "github.com/uber/jaeger-client-go/config"
)

type ViperProvider struct {
	l               logrus.FieldLogger
	ss              [][]byte
	generatedSecret []byte
	forcedHTTP      bool
}

const (
	ViperKeyPublicURL        = "urls.public"
	ViperKeyDSN              = "dsn"
	ViperKeyHost             = "serve.host"
	ViperKeyPort             = "serve.port"
	ViperKeyGetCookieSecrets = "secrets.cookie"
	ViperKeyService          = "service"
)

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

}

func NewViperProvider(l logrus.FieldLogger, forcedHTTP bool) Provider {
	return &ViperProvider{
		l:          l,
		forcedHTTP: forcedHTTP,
	}
}

func (v *ViperProvider) DSN() string {
	return viperx.GetString(v.l, ViperKeyDSN, "", "DATABASE_URL")
}

func (v *ViperProvider) TracingJaegerConfig() *jConfig.Configuration {
	c := &jConfig.Configuration{
		Reporter: &jConfig.ReporterConfig{
			CollectorEndpoint: viperx.GetString(v.l, "tracing.providers.jaeger.endpoint", "", "TRACING_PROVIDER_JAEGER_ENDPOINT"),
		},
		ServiceName: "Police",
		Sampler: &jConfig.SamplerConfig{
			SamplingServerURL: viperx.GetString(v.l, "tracing.providers.jaeger.sampling.server", "", "TRACING_PROVIDER_JAEGER_SAMPLING_SERVER_URL"),
			Param:             viperx.GetFloat64(v.l, "tracing.providers.jaeger.sampling.value", float64(1), "TRACING_PROVIDER_JAEGER_SAMPLING_VALUE"),
			Type:              viperx.GetString(v.l, "tracing.providers.jaeger.sampling.type", "const", "TRACING_PROVIDER_JAEGER_SAMPLING_TYPE"),
		},
	}

	return c
}

func (v *ViperProvider) GetCookieSecrets() [][]byte {
	return [][]byte{
		[]byte(viperx.GetString(v.l, ViperKeyGetCookieSecrets, "", "COOKIE_SECRET")),
	}
}

func (v *ViperProvider) ListenHost() string {
	return viperx.GetString(v.l, ViperKeyHost, "", "HOST")
}

func (v *ViperProvider) ListenPort() string {
	return viperx.GetString(v.l, ViperKeyPort, "", "PORT")
}
func (v *ViperProvider) Service() string {
	return viperx.GetString(v.l, ViperKeyService, "", "SERVICE")
}

func (v *ViperProvider) ServeHTTPS() bool {
	return !v.forcedHTTP
}
func (v *ViperProvider) Logger() logrus.FieldLogger {
	return v.l
}

func (v *ViperProvider) AppName() string {
	return "Policy"
}
