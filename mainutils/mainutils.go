// Package mainutils provides main package related utilities.
package mainutils

import (
	"context"
	"expvar"
	"fmt"
	"log"

	_ "github.com/siddhant2408/golang-libraries/ballast" // Initializes ballast.
	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/ctxsignal"
	"github.com/siddhant2408/golang-libraries/debugutils"
	"github.com/siddhant2408/golang-libraries/envutils"
	"github.com/siddhant2408/golang-libraries/errorhandle"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/logmain"
	"github.com/siddhant2408/golang-libraries/panichandle"
	"github.com/siddhant2408/golang-libraries/profilingmain"
	_ "github.com/siddhant2408/golang-libraries/randutils" // Initializes random seed.
	"github.com/siddhant2408/golang-libraries/ravenmain"
	"github.com/siddhant2408/golang-libraries/sibutils/sibhttpua"
	_ "github.com/siddhant2408/golang-libraries/spewutils" // Initializes spew config.
	"github.com/siddhant2408/golang-libraries/tracingmain"
)

func init() {
	panichandle.Handler = handlePanic
}

// Run is a helper to run the main function.
//
// Features:
//  - create a "background" context
//  - initialize log (+exit)
//  - catch returned error
//    - send to Sentry/Raven (severity=fatal)
//    - print to log
func Run(f func(context.Context) error) {
	defer panichandle.Recover()
	ctx := context.Background()
	ctx, unregisterCtx := ctxsignal.RegisterCancel(ctx)
	defer unregisterCtx()
	err := f(ctx)
	if err != nil {
		errorhandle.Handle(ctx, err, errorhandle.Fatal())
	}
	log.Println("Exit")
}

// Init initializes the common services for the main package.
//
// It initializes:
//  - log start
//  - Raven / Sentry
//  - tracing
func Init(cfg Config) (closeutils.F, error) {
	err := cfg.validate()
	if err != nil {
		return nil, errors.Wrap(err, "validate config")
	}
	logmain.Start(cfg.Version, cfg.Env)
	sibhttpua.WrapDefaultTransport(cfg.AppName, cfg.Version)
	closeRaven, err := initRaven(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "Raven / Sentry")
	}
	closeProfiling, err := initProfiling(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "profiling")
	}
	closeTracing, err := initTracing(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "tracing")
	}
	expvar.NewString("appname").Set(cfg.AppName)
	expvar.NewString("environment").Set(cfg.Env.String())
	expvar.NewString("version").Set(cfg.Version)
	closeDebug, err := initDebug(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "debug")
	}
	cl := func() {
		closeDebug()
		closeTracing()
		closeProfiling()
		closeRaven()
	}
	return cl, nil
}

func initRaven(cfg Config) (closeutils.F, error) {
	if cfg.SentryDSN == "" {
		return func() {}, nil
	}
	return ravenmain.Init(cfg.SentryDSN, cfg.Version, cfg.Env)
}

func initProfiling(cfg Config) (closeutils.F, error) {
	if cfg.ProfilingDisabled {
		return func() {}, nil
	}
	return profilingmain.Init(cfg.AppName, cfg.Version, cfg.Env)
}

func initTracing(cfg Config) (closeutils.F, error) {
	if cfg.TracingDisabled {
		return func() {}, nil
	}
	return tracingmain.Init(cfg.AppName, cfg.Version, cfg.Env)
}

func initDebug(cfg Config) (closeutils.F, error) {
	if cfg.Debug == "" {
		return func() {}, nil
	}
	return debugutils.StartHTTPServer(cfg.Debug)
}

func handlePanic(r interface{}) {
	var err error
	switch r := r.(type) {
	case error:
		err = r
	default:
		err = errors.New(fmt.Sprint(r))
	}
	err = errors.Wrap(err, "panic")
	errorhandle.Handle(context.Background(), err, errorhandle.Fatal())
}

// Config is the configuration for Init.
type Config struct {
	AppName           string
	Version           string
	Env               envutils.Env
	SentryDSN         string
	ProfilingDisabled bool
	TracingDisabled   bool
	Debug             string
}

func (c Config) validate() error {
	if c.AppName == "" {
		return errors.New("AppName must not be empty")
	}
	if c.Version == "" {
		return errors.New("Version must not be empty")
	}
	err := envutils.Check(c.Env)
	if err != nil {
		return errors.Wrap(err, "Env")
	}
	return nil
}
