package ravenerrors

import (
	"fmt"
	"io"
	"testing"

	raven "github.com/getsentry/raven-go"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestCaptureAndWait(t *testing.T) {
	myerr := errors.New("error")
	_, err := CaptureAndWait(myerr, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestCaptureAndWaitWithClient(t *testing.T) {
	myerr := errors.New("error")
	_, err := CaptureAndWaitWithClient(raven.DefaultClient, myerr, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestCapture(t *testing.T) {
	myerr := errors.New("error")
	_, errch := Capture(myerr, nil)
	err := <-errch
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestCaptureWithClient(t *testing.T) {
	myerr := errors.New("error")
	_, errch := CaptureWithClient(raven.DefaultClient, myerr, nil)
	err := <-errch
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestNewPacket(t *testing.T) {
	myerr := errors.New("error")
	NewPacket(myerr)
}

func TestNewPacketTags(t *testing.T) {
	myerr := errors.New("error")
	myerr = errors.WithTag(myerr, "foo", "bar")
	pkt := NewPacket(myerr)
	expected := raven.Tags{
		{
			Key:   "foo",
			Value: "bar",
		},
	}
	testutils.Compare(t, "unexpected tags", pkt.Tags, expected)
}

func TestNewPacketValues(t *testing.T) {
	myerr := errors.New("error")
	myerr = errors.WithValue(myerr, "foo", "bar")
	pkt := NewPacket(myerr)
	if pkt.Extra["foo"] != "bar" {
		t.Fatalf("unexpected extra: got %v, want %v", pkt.Extra["foo"], "bar")
	}
}

func TestNewPacketSeverity(t *testing.T) {
	myerr := errors.New("error")
	myerr = WithSeverity(myerr, raven.FATAL)
	pkt := NewPacket(myerr)
	if pkt.Level != raven.FATAL {
		t.Fatalf("unexpected level: got %v, want %v", pkt.Level, raven.FATAL)
	}
}

func TestNewPacketWithClient(t *testing.T) {
	myerr := errors.New("error")
	NewPacketWithClient(raven.DefaultClient, myerr)
}

func TestNewExceptions(t *testing.T) {
	myerr := errors.New("error")
	NewExceptions(myerr)
}

func TestNewExceptionsWithClient(t *testing.T) {
	myerr := errors.New("error")
	NewExceptionsWithClient(raven.DefaultClient, myerr)
}

func TestNewStackTraces(t *testing.T) {
	myerr := errors.New("error")
	NewStackTraces(myerr, 3)
}

func TestNewStackTracesWithClient(t *testing.T) {
	myerr := errors.New("error")
	NewStackTracesWithClient(raven.DefaultClient, myerr, 3)
}

func TestNewStackTracesDefault(t *testing.T) {
	myerr := &testError{
		s: "error",
	}
	sts := NewStackTraces(myerr, 3)
	if len(sts) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(sts), 1)
	}
	st := sts[0]
	if len(st.Frames) == 0 {
		t.Fatal("no frames")
	}
	f := st.Frames[len(st.Frames)-1].Function
	if f != "TestNewStackTracesDefault" {
		t.Fatalf("unexpected function: got %q, want %q", f, "TestNewStackTracesDefault")
	}
}

func TestGetErrorStackTraces(t *testing.T) {
	myerr := errors.New("error")
	sts := getErrorStackTraces(myerr, 3, nil)
	if len(sts) != 1 {
		t.Fatalf("unexpected length: got %d, want %d", len(sts), 1)
	}
	st := sts[0]
	if len(st.Frames) == 0 {
		t.Fatal("no frames")
	}
	for _, f := range st.Frames {
		if f.Filename == "" {
			t.Fatal("Filename is empty")
		}
		if f.Lineno == 0 {
			t.Fatal("Lineno is equal to 0")
		}
		if f.Function == "" {
			t.Fatal("Function is empty")
		}
		if f.Module == "" {
			t.Fatal("Module is empty")
		}
		if f.ContextLine == "" {
			t.Fatal("ContextLine is empty")
		}
	}
}

func TestSeverity(t *testing.T) {
	err := errors.New("error")
	err = WithSeverity(err, raven.FATAL)
	sv := GetSeverity(err)
	if sv != raven.FATAL {
		t.Fatalf("unexpected severity: got %v, want %v", sv, raven.FATAL)
	}
}

func TestSeverityNil(t *testing.T) {
	err := WithSeverity(nil, raven.FATAL)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestSeverityFormat(t *testing.T) {
	err := errors.New("error")
	err = WithSeverity(err, raven.FATAL)
	s := fmt.Sprint(err)
	expected := "Raven/Sentry fatal: error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func BenchmarkSeverityFormat(b *testing.B) {
	err := errors.New("error")
	err = WithSeverity(err, raven.FATAL)
	for i := 0; i < b.N; i++ {
		_, _ = fmt.Fprint(io.Discard, err)
	}
}

func TestInterface(t *testing.T) {
	err := errors.New("error")
	itf := &raven.User{
		ID:       "test",
		Username: "test",
		Email:    "test@example.com",
		IP:       "1.2.3.4",
	}
	err = WithInterface(err, itf)
	itfs := getInterfaces(err)
	testutils.CompareFatal(t, "unexpected interfaces", itfs, []raven.Interface{itf})
}

func TestInterfaceNil(t *testing.T) {
	itf := &raven.User{
		ID:       "test",
		Username: "test",
		Email:    "test@example.com",
		IP:       "1.2.3.4",
	}
	err := WithInterface(nil, itf)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func TestInterfaceFormat(t *testing.T) {
	err := errors.New("error")
	itf := &raven.User{
		ID:       "test",
		Username: "test",
		Email:    "test@example.com",
		IP:       "1.2.3.4",
	}
	err = WithInterface(err, itf)
	s := fmt.Sprint(err)
	expected := "error"
	if s != expected {
		t.Fatalf("unexpected message: got %q, want %q", s, expected)
	}
}

func BenchmarkInterfaceFormat(b *testing.B) {
	err := errors.New("error")
	itf := &raven.User{
		ID:       "test",
		Username: "test",
		Email:    "test@example.com",
		IP:       "1.2.3.4",
	}
	err = WithInterface(err, itf)
	for i := 0; i < b.N; i++ {
		_, _ = fmt.Fprint(io.Discard, err)
	}
}

type testError struct {
	s string
}

func (e *testError) Error() string {
	return e.s
}
