// Package ravenerrors adds support for github.com/siddhant2408/golang-libraries/errors's stacktrace to github.com/getsentry/raven-go.
//
// It works as a replacement for github.com/getsentry/raven-go.CaptureError*.
// Stacktraces are extracted from the error if available and replace raven's default behavior.
package ravenerrors

import (
	"fmt"
	"runtime"

	raven "github.com/getsentry/raven-go"
	"github.com/siddhant2408/golang-libraries/errors"
)

// CaptureAndWait is a replacement for github.com/getsentry/raven-go.CaptureErrorAndWait.
func CaptureAndWait(myerr error, tags map[string]string, interfaces ...raven.Interface) (eventID string, err error) {
	return captureAndWait(raven.DefaultClient, myerr, tags, interfaces...)
}

// CaptureAndWaitWithClient is a replacement for github.com/getsentry/raven-go.Client.CaptureErrorAndWait.
func CaptureAndWaitWithClient(client *raven.Client, myerr error, tags map[string]string, interfaces ...raven.Interface) (eventID string, err error) {
	return captureAndWait(client, myerr, tags, interfaces...)
}

func captureAndWait(client *raven.Client, myerr error, tags map[string]string, interfaces ...raven.Interface) (eventID string, err error) {
	eventID, errch := capture(client, myerr, 2, tags, interfaces...)
	err = <-errch
	err = errors.WithStack(err)
	return eventID, err
}

// Capture is a replacement for github.com/getsentry/raven-go.CaptureError.
func Capture(myerr error, tags map[string]string, interfaces ...raven.Interface) (eventID string, ch <-chan error) {
	return capture(raven.DefaultClient, myerr, 1, tags, interfaces...)
}

// CaptureWithClient is a replacement for github.com/getsentry/raven-go.Client.CaptureError.
func CaptureWithClient(client *raven.Client, myerr error, tags map[string]string, interfaces ...raven.Interface) (eventID string, ch <-chan error) {
	return capture(client, myerr, 1, tags, interfaces...)
}

func capture(client *raven.Client, myerr error, skip int, tags map[string]string, interfaces ...raven.Interface) (eventID string, ch <-chan error) {
	p := newPacket(client, myerr, skip+1, interfaces...)
	return client.Capture(p, tags)
}

// NewPacket is a replacement for github.com/getsentry/raven-go.Client.NewPacket.
func NewPacket(myerr error, interfaces ...raven.Interface) *raven.Packet {
	return newPacket(raven.DefaultClient, myerr, 1, interfaces...)
}

// NewPacketWithClient is a replacement for github.com/getsentry/raven-go.NewPacket.
func NewPacketWithClient(client *raven.Client, myerr error, interfaces ...raven.Interface) *raven.Packet {
	return newPacket(client, myerr, 1, interfaces...)
}

func newPacket(client *raven.Client, myerr error, skip int, interfaces ...raven.Interface) *raven.Packet {
	msg := fmt.Sprintf("%v\n%+v", myerr, myerr)
	interfaces = append(interfaces, newExceptions(client, myerr, skip+1))
	interfaces = append(interfaces, getInterfaces(myerr)...)
	pkt := raven.NewPacket(msg, interfaces...)
	pkt.Level = GetSeverity(myerr)
	pkt.AddTags(errors.Tags(myerr))
	for k, v := range errors.Values(myerr) {
		pkt.Extra[k] = v
	}
	return pkt
}

// NewExceptions is a replacement for github.com/getsentry/raven-go.NewException.
func NewExceptions(myerr error) raven.Exceptions {
	return newExceptions(raven.DefaultClient, myerr, 1)
}

// NewExceptionsWithClient is a replacement for github.com/getsentry/raven-go.NewException.
func NewExceptionsWithClient(client *raven.Client, myerr error) raven.Exceptions {
	return newExceptions(client, myerr, 1)
}

func newExceptions(client *raven.Client, myerr error, skip int) raven.Exceptions {
	sts := newStackTraces(client, myerr, skip+1, 3)
	vals := make([]*raven.Exception, len(sts))
	for i, st := range sts {
		vals[i] = raven.NewException(myerr, st)
	}
	return raven.Exceptions{
		Values: vals,
	}
}

// NewStackTraces is a replacement for github.com/getsentry/raven-go.NewStackTrace.
func NewStackTraces(myerr error, context int) []*raven.Stacktrace {
	return newStackTraces(raven.DefaultClient, myerr, 1, context)
}

// NewStackTracesWithClient is a replacement for github.com/getsentry/raven-go.NewStackTrace.
func NewStackTracesWithClient(client *raven.Client, myerr error, context int) []*raven.Stacktrace {
	return newStackTraces(client, myerr, 1, context)
}

func newStackTraces(client *raven.Client, myerr error, skip int, context int) []*raven.Stacktrace {
	sts := getErrorStackTraces(myerr, context, client.IncludePaths())
	if len(sts) == 0 {
		sts = []*raven.Stacktrace{
			raven.NewStacktrace(skip+1, context, client.IncludePaths()),
		}
	}
	return sts
}

func getErrorStackTraces(myerr error, context int, appPackagePrefixes []string) []*raven.Stacktrace {
	sfs := errors.StackFrames(myerr)
	sts := make([]*raven.Stacktrace, len(sfs))
	for i, sf := range sfs {
		sts[i] = convertFrames(sf, context, appPackagePrefixes)
	}
	return sts
}

func convertFrames(sf *runtime.Frames, context int, appPackagePrefixes []string) *raven.Stacktrace {
	var frames []*raven.StacktraceFrame
	for more := true; more; {
		var f runtime.Frame
		f, more = sf.Next()
		frame := convertFrame(f, context, appPackagePrefixes)
		if frame != nil {
			frames = append(frames, frame)
		}
	}
	for i, j := 0, len(frames)-1; i < j; i, j = i+1, j-1 {
		frames[i], frames[j] = frames[j], frames[i]
	}
	return &raven.Stacktrace{Frames: frames}
}

func convertFrame(f runtime.Frame, context int, appPackagePrefixes []string) *raven.StacktraceFrame {
	file := f.File
	if file == "" {
		file = "<unknown>"
	}
	return raven.NewStacktraceFrame(f.PC, f.Function, file, f.Line, context, appPackagePrefixes)
}

// WithSeverity adds a Severity to an error.
func WithSeverity(err error, sv raven.Severity) error {
	if err == nil {
		return nil
	}
	return &severity{
		err: err,
		sv:  sv,
	}
}

type severity struct {
	err error
	sv  raven.Severity
}

func (err *severity) Severity() raven.Severity {
	return err.sv
}

func (err *severity) WriteErrorMessage(w errors.Writer, verbose bool) bool {
	_, _ = w.WriteString("Raven/Sentry ")
	_, _ = w.WriteString(string(err.sv))
	return true
}

func (err *severity) Error() string                 { return errors.Error(err) }
func (err *severity) Format(s fmt.State, verb rune) { errors.Format(err, s, verb) }
func (err *severity) Unwrap() error                 { return err.err }

// GetSeverity returns the Severity for an error.
//
// It returns an empty string if there is no Severity.
func GetSeverity(err error) raven.Severity {
	var werr *severity
	ok := errors.As(err, &werr)
	if ok {
		return werr.Severity()
	}
	return ""
}

// WithInterface adds an Interface to an error.
func WithInterface(err error, i raven.Interface) error {
	if err == nil {
		return nil
	}
	return &itf{
		err: err,
		itf: i,
	}
}

type itf struct {
	err error
	itf raven.Interface
}

func (err *itf) Interface() raven.Interface {
	return err.itf
}

func (err *itf) WriteErrorMessage(w errors.Writer, verbose bool) bool {
	return false
}

func (err *itf) Error() string                 { return errors.Error(err) }
func (err *itf) Format(s fmt.State, verb rune) { errors.Format(err, s, verb) }
func (err *itf) Unwrap() error                 { return err.err }

func getInterfaces(err error) []raven.Interface {
	var itfs []raven.Interface
	for ; err != nil; err = errors.Unwrap(err) {
		if err, ok := err.(*itf); ok {
			itfs = append(itfs, err.Interface())
		}
	}
	return itfs
}
