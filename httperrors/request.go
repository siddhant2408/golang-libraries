package httperrors

import (
	"net/http"

	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/httpclientip"
	"github.com/siddhant2408/golang-libraries/httpurl"
)

const (
	requestTypeClient = "client"
	requestTypeServer = "server"
)

func writeRequestMessage(w errors.Writer, verbose bool, req *http.Request, typ string) {
	u := httpurl.Get(req)
	if !verbose {
		_, _ = w.WriteString("HTTP ")
		_, _ = w.WriteString(typ)
		_, _ = w.WriteString(" request ")
		_, _ = w.WriteString(req.Method)
		_, _ = w.WriteString(" ")
		_, _ = w.WriteString(u.String())
		return
	}
	_, _ = w.WriteString("HTTP ")
	_, _ = w.WriteString(typ)
	_, _ = w.WriteString(" request\n")
	_, _ = w.WriteString("	method: ")
	_, _ = w.WriteString(req.Method)
	_, _ = w.WriteString("\n")
	_, _ = w.WriteString("	URL: ")
	_, _ = w.WriteString(u.String())
	_, _ = w.WriteString("\n")
	if typ == requestTypeServer { // This is a kind of a hack, but I don't know any other way to differentiate client/server requests.
		clientIP, err := httpclientip.GetFromRequest(req)
		if err == nil {
			_, _ = w.WriteString("	client IP: ")
			_, _ = w.WriteString(clientIP.String())
			_, _ = w.WriteString("\n")
		}
	}
	ua := req.UserAgent()
	if ua != "" {
		_, _ = w.WriteString("	user-agent: ")
		_, _ = w.WriteString(ua)
		_, _ = w.WriteString("\n")
	}
	writeHeaderMessage(w, req.Header)
}
