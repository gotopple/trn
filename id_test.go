package id

import (
	"strings"
	"testing"
)

// TOOD test for invalid input
func TestNewTRN(t *testing.T) {
	c := []struct {
		p      string
		s      string
		r      string
		a      string
		pre    string
		prefix string
		eo     error
	}{
		{``, ``, ``, ``, ``, `trn:::::/`, nil},
		{`topple`, ``, ``, ``, ``, `trn:topple::::/`, nil},
		{``, `content`, ``, ``, ``, `trn::content:::/`, nil},
		{``, ``, `us-west`, ``, ``, `trn:::us-west::/`, nil},
		{``, ``, ``, `1234`, ``, `trn::::1234:/`, nil},
		{``, ``, ``, ``, `prefix`, `trn:::::prefix/`, nil},
		{`topple`, `content`, `us-west`, `1234`, `prefix`, `trn:topple:content:us-west:1234:prefix/`, nil},
	}
	for _, e := range c {
		o := NewTRN(e.p, e.s, e.r, e.a, e.pre)
		if !strings.HasPrefix(string(o), e.prefix) {
			t.Errorf("ContentID for (%v, %v, %v, %v) was incorrect, got: %v wanted w/ prefix: %v", e.p, e.r, e.a, e.pre, o, e.prefix)
		}
	}
}
