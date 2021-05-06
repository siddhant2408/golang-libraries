package errors

import (
	"fmt"
	"io"

	"github.com/siddhant2408/golang-libraries/bufpool"
)

// Formattable represents a formattable error.
type Formattable interface {
	error
	WriteErrorMessage(w Writer, verbose bool) bool
}

// Writer represents a writer used by Formattable.
type Writer interface {
	io.Writer
	io.StringWriter
}

// Format formats an error.
func Format(err Formattable, s fmt.State, verb rune) {
	// This type assertion should never fail, because the internal type of the fmt package implements WriteString.
	// If it fails, it means there is something seriously wrong, and we should stop the application.
	// See https://github.com/golang/go/issues/20786
	w := s.(Writer) //nolint:errcheck
	switch {
	case verb == 'v' && s.Flag('+'):
		writeError(w, err, true)
	case verb == 'v' || verb == 's':
		writeError(w, err, false)
	case verb == 'q':
		_, _ = fmt.Fprintf(w, "%q", Error(err))
	}
}

// Error formats an error on a single line.
func Error(err Formattable) string {
	buf := bufpool.Get()
	defer bufpool.Put(buf)
	writeError(buf, err, false)
	return buf.String()
}

func writeError(w Writer, err Formattable, verbose bool) {
	var separator string
	if verbose {
		separator = "\n"
	} else {
		separator = ": "
	}
	for {
		ok := err.WriteErrorMessage(w, verbose)
		werr := Unwrap(err)
		if werr == nil {
			return
		}
		if ok {
			_, _ = w.WriteString(separator)
		}
		err, ok = werr.(Formattable)
		if !ok {
			_, _ = w.WriteString(werr.Error())
			return
		}
	}
}
