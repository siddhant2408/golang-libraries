package spewutils

import (
	"bytes"
	"testing"
)

func Test(t *testing.T) {
	buf := new(bytes.Buffer)
	writeValueWithoutNewline(buf, map[string]interface{}{
		"string":  "test",
		"int":     123,
		"float64": 123.456,
		"slice":   make([]int, 2, 3),
		"pointer": new(int),
	})
	s := buf.String()
	expected := `(map[string]interface {}) (len=5) {
	(string) (len=7) "float64": (float64) 123.456,
	(string) (len=3) "int": (int) 123,
	(string) (len=7) "pointer": (*int)(0),
	(string) (len=5) "slice": ([]int) (len=2) {
		(int) 0,
		(int) 0
	},
	(string) (len=6) "string": (string) (len=4) "test"
}`
	if s != expected {
		t.Fatalf("unexpected string:\ngot %q\nwant %q", s, expected)
	}
}
