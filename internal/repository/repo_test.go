package repository

import (
	"context"
	"testing"

	"bn-riftstats/internal/riotclient"
)

func newTestRepo() *Repository {
	dd := riotclient.NewDDragonManager("16.21.1")
	return NewRepository(dd)
}

func TestRegisterAndAuthenticate(t *testing.T) {
	r := newTestRepo()

	profile, err := r.RegisterUser("Tester@Example.COM", "strongpass1", "Tester")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if profile.Email != "tester@example.com" || profile.Role != "USER" {
		t.Fatalf("register normalization wrong: %+v", profile)
	}

	// Correct password authenticates (case-insensitive email).
	got, err := r.AuthenticateUser("TESTER@example.com", "strongpass1")
	if err != nil {
		t.Fatalf("authenticate failed: %v", err)
	}
	if got.UserID != profile.UserID {
		t.Fatalf("authenticated wrong user: %+v", got)
	}

	// Wrong password is rejected.
	if _, err := r.AuthenticateUser("tester@example.com", "nope-nope-no"); err == nil {
		t.Fatalf("wrong password was accepted")
	}
	// Unknown email is rejected (no auto-register).
	if _, err := r.AuthenticateUser("nobody@example.com", "whatever123"); err == nil {
		t.Fatalf("unknown email was accepted")
	}
	// Passwords are stored hashed, never plaintext.
	r.mu.RLock()
	stored := r.passwords["tester@example.com"]
	r.mu.RUnlock()
	if stored == "" || stored == "strongpass1" {
		t.Fatalf("password not hashed")
	}
}

func TestRegisterValidation(t *testing.T) {
	r := newTestRepo()
	if _, err := r.RegisterUser("not-an-email", "longenough1", ""); err == nil {
		t.Fatalf("invalid email accepted")
	}
	if _, err := r.RegisterUser("ok@example.com", "short", ""); err == nil {
		t.Fatalf("short password accepted")
	}
	if _, err := r.RegisterUser("dup@example.com", "longenough1", ""); err != nil {
		t.Fatalf("first register failed: %v", err)
	}
	if _, err := r.RegisterUser("dup@example.com", "longenough2", ""); err == nil {
		t.Fatalf("duplicate email accepted")
	}
}

func TestRecalculateTierGrades(t *testing.T) {
	r := newTestRepo()
	processed, _ := r.RecalculateTierGrades()
	if processed == 0 {
		t.Fatalf("expected slices to be re-scored")
	}
	resp, err := r.GetTierList(context.Background(), "GLOBAL", "EMERALD_PLUS", "16.21")
	if err != nil {
		t.Fatalf("tier list failed: %v", err)
	}
	if len(resp.Entries) == 0 {
		t.Fatalf("tier list empty after recalculation")
	}
	for _, e := range resp.Entries {
		if e.TierGrade == "" || e.SampleQuality == "" {
			t.Fatalf("entry missing grade/quality after recalculation: %+v", e)
		}
	}
}

func TestHealthStatusStandby(t *testing.T) {
	r := newTestRepo()
	pg, redis := r.HealthStatus(context.Background())
	if pg != "STANDBY" || redis != "STANDBY" {
		t.Fatalf("expected STANDBY without connections, got pg=%s redis=%s", pg, redis)
	}
}
