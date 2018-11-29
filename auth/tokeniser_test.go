package auth_test

import (
	"log"
	"testing"

	"github.com/LUSHDigital/microservice-core-golang/auth"
)

var (
	err       error
	tokeniser *auth.Tokeniser
)

func TestTokeniser_GenerateToken(t *testing.T) {
	// A test consumer.
	testConsumer := &auth.Consumer{
		ID:        999,
		FirstName: "Testy",
		LastName:  "McTest",
		Grants: []string{
			"testing.read",
			"testing.create",
		},
	}

	// Generate a test authToken.
	authToken, err := tokeniser.GenerateToken(testConsumer)
	if err != nil {
		t.Fatal(err)
	}

	parsedToken, err := tokeniser.ParseToken(authToken)
	if err != nil {
		t.Fatal(err)
	}

	deepEqual(t, testConsumer, &parsedToken.Claims.(*auth.Claims).Consumer)
}

func TestTokeniser_ValidateToken(t *testing.T) {
	// A test consumer.
	testConsumer := &auth.Consumer{
		ID:        999,
		FirstName: "Testy",
		LastName:  "McTest",
		Grants: []string{
			"testing.read",
			"testing.create",
		},
	}

	// Generate a test authToken.
	authToken, err := tokeniser.GenerateToken(testConsumer)
	if err != nil {
		log.Fatal(err)
	}

	ok, err := tokeniser.ValidateToken(authToken)
	if err != nil {
		log.Fatal(err)
	}

	deepEqual(t, true, ok)
}

func TestTokeniser_ValidateToken2(t *testing.T) {
	cases := []struct {
		name        string
		authToken   string
		expectedOk  bool
		expectedErr error
	}{
		{
			name:        "incorrect signing method",
			authToken:   "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJjb25zdW1lciI6eyJpZCI6OTk5LCJmaXJzdF9uYW1lIjoiVGVzdHkiLCJsYXN0X25hbWUiOiJNY1Rlc3QiLCJsYW5ndWFnZSI6IiIsImdyYW50cyI6WyJ0ZXN0aW5nLnJlYWQiLCJ0ZXN0aW5nLmNyZWF0ZSJdfSwiZXhwIjoxNTE4NjAzNzIwLCJqdGkiOiIyNTAwYjk3MS0wNTcxLTQ4Y2UtYmUzOS1jNWJhNGQwZmU0MGIiLCJpc3MiOiJ0ZXN0aW5nIn0.",
			expectedOk:  false,
			expectedErr: auth.ErrTokenInvalid,
		},
		{
			name:        "malformed token",
			authToken:   ".eyJjb25zdW1lciI6eyJpZCI6OTk5LCJmaXJzdF9uYW1lIjoiVGVzdHkiLCJsYXN0X25hbWUiOiJNY1Rlc3QiLCJsYW5ndWFnZSI6IiIsImdyYW50cyI6WyJ0ZXN0aW5nLnJlYWQiLCJ0ZXN0aW5nLmNyZWF0ZSJdfSwiZXhwIjoxNTE4NjAzNzIwLCJqdGkiOiIyNTAwYjk3MS0wNTcxLTQ4Y2UtYmUzOS1jNWJhNGQwZmU0MGIiLCJpc3MiOiJ0ZXN0aW5nIn0.",
			expectedOk:  false,
			expectedErr: auth.ErrTokenMalformed,
		},
		{
			name:        "missing token",
			expectedOk:  false,
			expectedErr: auth.ErrTokenMalformed,
		},
		{
			name: "expired token",

			authToken:   "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb25zdW1lciI6eyJpZCI6OTk5LCJmaXJzdF9uYW1lIjoiVGVzdHkiLCJsYXN0X25hbWUiOiJNY1Rlc3QiLCJsYW5ndWFnZSI6IiIsImdyYW50cyI6WyJ0ZXN0aW5nLnJlYWQiLCJ0ZXN0aW5nLmNyZWF0ZSJdfSwiZXhwIjoxNTE0NzY0ODAwLCJqdGkiOiI5MjJiNTJhNi0wYmRjLTQ5ZmEtOWM4NC0wNmRlZjc2YWM2MGMiLCJpc3MiOiJ0ZXN0aW5nIn0.qNFzO3UODL6W-r_tG7Bmc844Qg9clOdoM-mbAawAN6pTyhdcx888mag6zxyvxYiX4fHY__j1iCfxrrr0mYLtcaM3MhmOch_Nj5u0IyOHDjgtwCQT22pRR1Y878uq78LQ2ktY2pbqTAFZyRlTbzsiT2Zq1RCatPOlZpwORLfOUTA",
			expectedOk:  false,
			expectedErr: auth.ErrTokenExpired,
		},
		{
			name: "token not valid yet",

			authToken:   "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb25zdW1lciI6eyJpZCI6OTk5LCJmaXJzdF9uYW1lIjoiVGVzdHkiLCJsYXN0X25hbWUiOiJNY1Rlc3QiLCJsYW5ndWFnZSI6IiIsImdyYW50cyI6WyJ0ZXN0aW5nLnJlYWQiLCJ0ZXN0aW5nLmNyZWF0ZSJdfSwibmJmIjoyNTMzNzA3NjQ4MDAsImp0aSI6IjkyMmI1MmE2LTBiZGMtNDlmYS05Yzg0LTA2ZGVmNzZhYzYwYyIsImlzcyI6InRlc3RpbmcifQ.aKEg_6-7YVJgewm7-YL_8p4uFuSOzzq0DR-z0OMjIamlitZNyk4fY5YTyBuc0MFJT-dW-lrL8AMmCTqhFEOPYu-0uGKQPZUIGlBmc88fZb0yh5Pt-o3uSYncoU1Lx27P1GoFSQH_wVlhl_L3khTuTlshZR9p-Fe8wJOMUaTSUC8",
			expectedOk:  false,
			expectedErr: auth.ErrTokenExpired,
		},
		{
			name:        "invalid claims",
			authToken:   "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTg2MDQ3OTQsImp0aSI6IjQxNjU3NzNlLWQ0YzYtNGU1Ni04ZGJmLTU2YzM2YzlmMzA1OCIsImlzcyI6InRlc3RpbmcifQ.4jhNEfhCkUrweLT2T4lJBmHWTOjF8QHNQBBEQaxo3xnFl1ya0vnL0hWPHdJydnFuSJbrFSvi4iXQtdByuKEQg7ti3JCTKxHN68zQRayk_0c_M6bE_RqDnRnX-Qc65qNAiloRWwIdEvTy4LebClgU-0POWSdqhNnAGUw759tFah0",
			expectedOk:  false,
			expectedErr: auth.ErrTokenExpired,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ok, err := tokeniser.ValidateToken(c.authToken)
			if err != nil {
				deepEqual(t, c.expectedErr, err)
				return
			}
			deepEqual(t, c.expectedOk, ok)
		})
	}
}

func TestTokeniser_GetTokenConsumer(t *testing.T) {
	// A test consumer.
	testConsumer := &auth.Consumer{
		ID:        999,
		FirstName: "Testy",
		LastName:  "McTest",
		Grants: []string{
			"testing.read",
			"testing.create",
		},
	}

	// Generate a test authToken.
	authToken, err := tokeniser.GenerateToken(testConsumer)
	if err != nil {
		log.Fatal(err)
	}

	deepEqual(t, testConsumer, tokeniser.GetTokenConsumer(authToken))
}