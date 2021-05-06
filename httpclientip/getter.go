package httpclientip

import (
	"net"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/siddhant2408/golang-libraries/errors"
)

// Getter get the client IP for an HTTP request.
// It handles correctly the X-Forwarded-For header.
type Getter struct {
	trusted []interface{}
}

// NewGetter returns a new Getter.
//
// trusted is a list IPs/ranges for trusted proxies.
func NewGetter(trusted []string) (*Getter, error) {
	trustedParsed, err := parseTrusted(trusted)
	if err != nil {
		return nil, errors.Wrap(err, "parse trusted")
	}
	return &Getter{
		trusted: trustedParsed,
	}, nil
}

func parseTrusted(trusted []string) ([]interface{}, error) {
	parsed := make([]interface{}, len(trusted))
	for i, t := range trusted {
		p, err := parseTrustedValue(t)
		if err != nil {
			return nil, err
		}
		parsed[i] = p
	}
	return parsed, nil
}

func parseTrustedValue(t string) (interface{}, error) {
	if strings.Contains(t, "/") {
		_, ipNet, err := net.ParseCIDR(t)
		if err != nil {
			return nil, errors.Wrap(err, "parse CIDR")
		}
		return ipNet, nil
	}
	ip, err := parseIP(t)
	if err != nil {
		return nil, errors.Wrap(err, "parse IP")
	}
	return ip, nil
}

func parseIP(s string) (net.IP, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, errors.Newf("invalid IP address: %s", s)
	}
	return ip, nil
}

// GetClientIP returns the client IP for an HTTP request.
func (g *Getter) GetClientIP(req *http.Request) (net.IP, error) {
	ips, err := getIPs(req)
	if err != nil {
		return nil, errors.Wrap(err, "get IPs")
	}
	var ip net.IP
	for _, ipStr := range ips {
		ip, err = parseIP(ipStr)
		if err != nil {
			return nil, errors.Wrap(err, "parse IP")
		}
		if !g.IsTrusted(ip) {
			return ip, nil
		}
	}
	if ip != nil {
		return ip, nil
	}
	return nil, errors.New("no client IP")
}

// getIPs returns the IPs for an HTTP request.
// Low index means nearest to the server,
// high index means nearest to the client.
func getIPs(req *http.Request) ([]string, error) {
	var ips []string
	ips = append(ips, getIPsHeaderXFF(req)...)
	ip, err := getIPRemoteAddr(req)
	if err != nil {
		return nil, errors.Wrap(err, "remote addr")
	}
	if ip != "" {
		ips = append(ips, ip)
	}
	reverseSliceString(ips)
	return ips, nil
}

func getIPsHeaderXFF(req *http.Request) []string {
	var ips []string
	for _, xff := range req.Header[textproto.CanonicalMIMEHeaderKey("X-Forwarded-For")] {
		tmp := strings.Split(xff, ",")
		if len(ips) == 0 {
			ips = tmp // Optimization that prevents allocation of a new slice.
		} else {
			ips = append(ips, tmp...)
		}
	}
	for i, ip := range ips {
		ips[i] = strings.TrimSpace(ip)
	}
	return ips
}

func getIPRemoteAddr(req *http.Request) (string, error) {
	if req.RemoteAddr == "" {
		return "", nil
	}
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "", errors.Wrap(err, "split host port")
	}
	return ip, nil
}

func reverseSliceString(ss []string) {
	for i, j := 0, len(ss)-1; i < j; i, j = i+1, j-1 {
		ss[i], ss[j] = ss[j], ss[i]
	}
}

// IsTrusted returns true if the IP is trusted, false otherwise.
func (g *Getter) IsTrusted(ip net.IP) bool {
	for _, t := range g.trusted {
		switch t := t.(type) {
		case net.IP:
			if t.Equal(ip) {
				return true
			}
		case *net.IPNet:
			if t.Contains(ip) {
				return true
			}
		}
	}
	return false
}
