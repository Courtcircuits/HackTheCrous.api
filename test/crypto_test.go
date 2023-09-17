package test

import (
	"testing"

	"github.com/Courtcircuits/HackTheCrous.api/util"
)

func TestGenNonce(t *testing.T) {
	nonce := util.GenNonce(10)

	if len(nonce) != 10 {
		t.Errorf("nonce should be 10 characters long but is %q characters long", len(nonce))
	}
}

func TestGenHash(t *testing.T) {
	s := "test123"
	sp := "123test"

	s_hashed := util.Hash(s)
	s_hashed_bis := util.Hash(s)
	sp_hashed := util.Hash(sp)

	if s == s_hashed {
		t.Errorf("hash must be different than what it used to be")
	}

	if s == sp_hashed {
		t.Errorf("no collision possible")
	}

	if s_hashed != s_hashed_bis {
		t.Errorf("hash must be equal always equal to the same seed through (CANT EXPLAIN IT CORRECTLY) %q != %q", s_hashed, s_hashed_bis)
	}

}

func TestGenHashAndSaltedAndSalted(t *testing.T) {
	s := "test123"
	sp := "123test"

	s_hashed, salt := util.HashAndSalted(s)
	sp_hashed, _ := util.HashAndSalted(sp)

	if len(s_hashed) > 64 {
		t.Errorf("hash shouldn't be longer than 64 characters")
	}

	if s == s_hashed {
		t.Errorf("hash must be different than what it used to be")
	}

	if s == sp_hashed {
		t.Errorf("no collision possible")
	}

	if !util.CompareHash(s_hashed, s, salt) {
		t.Errorf("%q + %q doesn't return %q", salt, s, s_hashed)
	}
}
