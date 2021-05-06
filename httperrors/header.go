package httperrors

import (
	"net/http"

	"github.com/siddhant2408/golang-libraries/errors"
)

func writeHeaderMessage(w errors.Writer, h http.Header) {
	if len(h) == 0 {
		return
	}
	_, _ = w.WriteString("	headers:\n")
	for k, v := range h {
		_, _ = w.WriteString("		")
		_, _ = w.WriteString(k)
		_, _ = w.WriteString(": ")
		for i, vv := range v {
			if i != 0 {
				_, _ = w.WriteString("|")
			}
			_, _ = w.WriteString(vv)
		}
		_, _ = w.WriteString("\n")
	}
}
