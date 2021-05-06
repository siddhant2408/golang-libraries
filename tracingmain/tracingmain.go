// Package tracingmain provides tracing related utilities for a main package.
//
// This package exists because there is a cyclic dependency issue.
package tracingmain

import (
	"os"
	"strconv"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/envutils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/httptracing"
	ddtrace_opentracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	ddtrace_tracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	envEnabled      = "TRACING_ENABLED"
	envDDAgentAddr  = "TRACING_DATADOG_AGENT_ADDR"
	envDDAgentDebug = "TRACING_DATADOG_AGENT_DEBUG"
)

// Init initializes tracing.
func Init(service string, version string, env envutils.Env) (closeutils.F, error) {
	enabled, err := isEnabled(env)
	if err != nil {
		return nil, errors.Wrap(err, "enabled")
	}
	if !enabled {
		return func() {}, nil
	}
	tr, err := newTracer(service, version, env)
	if err != nil {
		return nil, errors.Wrap(err, "new tracer")
	}
	opentracing.SetGlobalTracer(tr)
	httptracing.WrapDefaultTransport()
	cl := func() {
		ddtrace_tracer.Stop()
	}
	return cl, nil
}

func isEnabled(env envutils.Env) (bool, error) {
	enabledStr, ok := os.LookupEnv(envEnabled)
	if ok {
		enabled, err := strconv.ParseBool(enabledStr)
		if err != nil {
			return false, errors.Wrapf(err, "parse environment variable %q", envEnabled)
		}
		return enabled, err
	}
	enabled := env == envutils.Staging || env == envutils.Production
	return enabled, nil
}

func newTracer(service string, version string, env envutils.Env) (opentracing.Tracer, error) {
	opts := []ddtrace_tracer.StartOption{
		ddtrace_tracer.WithServiceName(service),
		ddtrace_tracer.WithGlobalTag("version", version),
		ddtrace_tracer.WithGlobalTag("env", env.String()),
	}
	agentAddr := GetDatadogAgentAddr()
	if agentAddr != "" {
		opts = append(opts, ddtrace_tracer.WithAgentAddr(agentAddr))
	}
	agentDebug, err := getAgentDebug()
	if err != nil {
		return nil, errors.Wrap(err, "debug")
	}
	if agentDebug {
		opts = append(opts, ddtrace_tracer.WithDebugMode(true))
	}
	return ddtrace_opentracer.New(opts...), nil
}

// GetDatadogAgentAddr returns the address used by the DatDog agent, if defined.
func GetDatadogAgentAddr() string {
	return os.Getenv(envDDAgentAddr)
}

func getAgentDebug() (bool, error) {
	debugStr, ok := os.LookupEnv(envDDAgentDebug)
	if !ok {
		return false, nil
	}
	enabled, err := strconv.ParseBool(debugStr)
	if err != nil {
		return false, errors.Wrap(err, "parse bool")
	}
	return enabled, nil
}
