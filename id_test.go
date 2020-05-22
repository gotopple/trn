package trn

import (
	"fmt"
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

func TestGenerateServiceInvite(t *testing.T) {
	o := NewTRN(`test`, `account`, `local`, ``, `/serviceinvite`)
	printTRN("Random ServiceInvite", o)
	o = NewTRN(`test`, `account`, `local`, ``, `/invite`)
	printTRN("Random Invite", o)
	o = NewTRN(`test`, `account`, `local`, ``, `/discount`)
	printTRN("Random Discount", o)
	o = NewTRN(`test`, `account`, `local`, `1`, `/user`)
	printTRN("Random User", o)
}

func printTRN(note string, o TRN) {
	fmt.Printf("%s: %s\n\t%s\t%s\t%s\t%s\t%s\n", note, o.Encode(), o.Partition(), o.Service(), o.Region(), o.Account(), o.Resource())
}

func TestEncoding(t *testing.T) {
	i := NewTRN(`topple`, `content`, `us-west`, `1234`, `prefix`)
	encoded := i.Encode()
	o, err := Decode(encoded)
	if err != nil {
		t.Error(err)
	}
	if o != i {
		t.Error(`decoded TRN does not match original`)
	}
}

func TestDecodeFailure(t *testing.T) {
	i := `not a b32 trn`
	_, err := Decode(i)
	if err == nil {
		t.Error(`apparently decoded a non-base32 non-trn string into a TRN`)
	}
}

func TestComponentReads(t *testing.T) {
	i := NewTRN(`topple`, `content`, `us-west`, `1234`, `prefix`)
	id, p, s, r, a, _ := i.Components()
	if id != `trn` {
		t.Error(`incorrect ID from Components, expected trn`)
	}
	if p != `topple` {
		t.Errorf(`incorrect partition from Components %v`, p)
	}
	if s != `content` {
		t.Errorf(`incorrect service from Components %v`, s)
	}
	if r != `us-west` {
		t.Errorf(`incorrect region from Components %v`, r)
	}
	if a != `1234` {
		t.Errorf(`incorrect account from Components %v`, a)
	}

	id = i.ID()
	p = i.Partition()
	s = i.Service()
	r = i.Region()
	a = i.Account()

	if id != `trn` {
		t.Error(`incorrect ID, expected trn`)
	}
	if p != `topple` {
		t.Errorf(`incorrect partition %v`, p)
	}
	if s != `content` {
		t.Errorf(`incorrect service %v`, s)
	}
	if r != `us-west` {
		t.Errorf(`incorrect region %v`, r)
	}
	if a != `1234` {
		t.Errorf(`incorrect account %v`, a)
	}
}

func TestParseServiceIdentifier(t *testing.T) {
	o, err := ParseServiceIdentifier(`metadata`)
	if err != nil || o != Metadata {
		t.Errorf(`incorrectly parsed the SI for "metadata" err: %v o: %v`, err, o)
	}
	o, err = ParseServiceIdentifier(`content`)
	if err != nil || o != Content {
		t.Errorf(`incorrectly parsed the SI for "content" err: %v o: %v`, err, o)
	}
	o, err = ParseServiceIdentifier(``)
	if err == nil {
		t.Errorf(`failed to return error on parsing empty string o: %v`, o)
	}
}

func TestSIString(t *testing.T) {
	if Metadata.String() != `metadata` {
		t.Errorf(`incorrect name for metadata service String(): %v`, Metadata.String())
	}
}
