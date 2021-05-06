// Package debugutils provides debug utilities.
//
// It registers handlers on http.DefaultServeMux:
//  - from the net/http/pprof package
//  - from the expvar package
package debugutils

import (
	"context"
	"expvar"
	"net"
	"net/http"
	_ "net/http/pprof" // Register HTTP handlers.
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/ctxhttpsrv"
	"github.com/siddhant2408/golang-libraries/ctxutils"
	"github.com/siddhant2408/golang-libraries/errorhandle"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/timeutils"
)

func init() {
	http.HandleFunc("/debug", httpHandler)
	expvar.Publish("buildinfo", expvar.Func(expvarBuildInfo))
	expvar.Publish("envvars", expvar.Func(expvarEnvVars))
	expvar.Publish("goroutines", expvar.Func(expvarGoroutines))
	expvar.Publish("uptime", expvar.Func(expvarUptime))
	expvar.NewString("goarch").Set(runtime.GOARCH)
	expvar.NewString("goos").Set(runtime.GOOS)
	expvar.NewInt("gomaxprocs").Set(int64(runtime.GOMAXPROCS(0)))
	expvar.NewString("goversion").Set(runtime.Version())
	expvar.NewInt("numcpu").Set(int64(runtime.NumCPU()))
	expvar.NewInt("pid").Set(int64(os.Getpid()))
}

// StartHTTPServer starts an HTTP server that provides debug information.
func StartHTTPServer(addr string) (closeutils.F, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "listen")
	}
	ctx := context.Background()
	cl := ctxutils.Start(ctx, func(ctx context.Context) {
		err := ctxhttpsrv.Serve(ctx, l, http.DefaultServeMux)
		if err != nil {
			errorhandle.Handle(ctx, err, errorhandle.Wait())
		}
	})
	return cl, nil
}

var httpContent = []byte(`<!DOCTYPE html>
<html>
<head>
<title>/debug</title>
</head>
<body>
<p><a href="/debug/vars">vars</a></p>
<p><a href="/debug/pprof/">pprof</a></p>
</body>
</html>`)

func httpHandler(w http.ResponseWriter, req *http.Request) {
	_, _ = w.Write(httpContent)
}

func expvarGoroutines() interface{} {
	return runtime.NumGoroutine()
}

func expvarEnvVars() interface{} {
	envs := os.Environ()
	m := make(map[string]string, len(envs))
	for _, env := range envs {
		kv := strings.SplitN(env, "=", 2)
		if len(kv) == 2 {
			k, v := kv[0], kv[1]
			m[k] = v
		}
	}
	return m
}

var startTime = timeutils.Now()

func expvarUptime() interface{} {
	return timeutils.Since(startTime).Seconds()
}

func expvarBuildInfo() interface{} {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		return bi
	}
	return nil
}
