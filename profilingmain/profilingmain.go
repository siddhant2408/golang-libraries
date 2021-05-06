// Package profilingmain provides profiling related utilities for a main package.
package profilingmain

import (
	"fmt"

	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/envutils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/tracingmain"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

// Init initializes the profiling.
func Init(service string, version string, env envutils.Env) (closeutils.F, error) {
	if !isEnabled(env) {
		return func() {}, nil
	}
	opts := []profiler.Option{
		profiler.WithService(service),
		profiler.WithEnv(env.String()),
		profiler.WithTags(fmt.Sprintf("version:%s", version)),
		profiler.WithProfileTypes(profiler.CPUProfile, profiler.HeapProfile, profiler.GoroutineProfile),
	}
	agentAddr := tracingmain.GetDatadogAgentAddr()
	if agentAddr != "" {
		opts = append(opts, profiler.WithAgentAddr(agentAddr))
	}
	err := profiler.Start(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "start")
	}
	cl := func() {
		profiler.Stop()
	}
	return cl, nil
}

func isEnabled(env envutils.Env) bool {
	return env == envutils.Staging || env == envutils.Production
}
