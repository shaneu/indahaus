package auth_test

import (
	"testing"

	"github.com/shaneu/indahaus/pkg/auth"
)

// Success/Failure chars for nicer go test -v output
const (
	success = "\u2713"
	failure = "\u2717"
)

func TestAuthenticate(t *testing.T) {
	t.Log("Given the need to authenticate users.")

	tests := []struct {
		password string
		username string
		want     bool
	}{
		{
			username: "correct",
			password: "good",
			want:     true,
		},
		{
			username: "incorrect",
			password: "bad",
			want:     false,
		},
	}

	a := auth.New("correct", "good")

	for i, tt := range tests {
		t.Run("When inserting an IP result.", func(t *testing.T) {
			if got := a.Authenticate(tt.username, tt.password); got != tt.want {
				t.Fatalf("\t%s\tTest %d:\tShould be able to authenticate user : got=%t want=%t.", failure, i, got, tt.want)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to authenticate user.", success, i)
		})
	}
}
