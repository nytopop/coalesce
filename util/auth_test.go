// coalesce/util/auth_test.go

package util

import "testing"

func Test_GenerateSalt(t *testing.T) {
	salt, err := GenerateSalt()
	if err != nil {
		t.Error(err)
	}
	t.Log(salt)
}

func Test_ComputeToken(t *testing.T) {
	salt, err := GenerateSalt()
	if err != nil {
		t.Error(err)
	}

	pw := "correcthorsebatterystaple"

	token, err := ComputeToken(salt, pw)
	if err != nil {
		t.Error(err)
	}

	t.Log(token)
}

func Test_CheckToken(t *testing.T) {
	salt1, err := GenerateSalt()
	if err != nil {
		t.Error(err)
	}

	salt2, err := GenerateSalt()
	if err != nil {
		t.Error(err)
	}

	pw1 := "correcthorsebatterystaple"
	pw2 := "correctthehorse1234$$$$$$"

	token1, err := ComputeToken(salt1, pw1)
	if err != nil {
		t.Error(err)
	}

	token2, err := ComputeToken(salt2, pw2)
	if err != nil {
		t.Error(err)
	}

	// This should match, correct password && correct salt
	check1 := CheckToken(salt1, pw1, token1)
	if check1 != nil {
		t.Error("Correct salt/password, but did not match.")
	}

	// This should NOT match, correct password && incorrect salt
	check2 := CheckToken(salt1, pw2, token2)
	if check2 == nil {
		t.Error("Incorrect salt, but still matched.")
	}

	// This should NOT match, incorrect password && correct salt
	check3 := CheckToken(salt1, pw2, token1)
	if check3 == nil {
		t.Error("Incorrect password, but still matched.")
	}

	// This should NOT match, incorrect password && incorrect salt
	check4 := CheckToken(salt1, pw1, token2)
	if check4 == nil {
		t.Error("Incorrect salt/password, but still matched.")
	}
}
