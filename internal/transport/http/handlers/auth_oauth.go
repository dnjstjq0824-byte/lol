package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"bn-riftstats/internal/domain"
	"bn-riftstats/internal/transport/http/middleware"
)

// 22_Auth: POST /api/v1/auth/register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeProblem(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}

	user, err := h.repo.RegisterUser(req.Email, req.Password, req.DisplayName)
	if err != nil {
		msg := err.Error()
		status := http.StatusBadRequest
		title := "Registration Failed"
		if strings.Contains(msg, "already registered") {
			status = http.StatusConflict
			title = "Conflict"
		}
		writeProblem(w, status, title, msg)
		return
	}

	token, err := middleware.GenerateToken(user.UserID, user.Email, user.Role, h.cfg.JWTSecret)
	if err != nil {
		writeProblem(w, http.StatusInternalServerError, "Token Generation Failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, domain.AuthTokenResponse{
		AccessToken: token,
		ExpiresIn:   86400,
	})
}

// oauthState creates a signed, timestamped CSRF token for the OAuth round-trip.
func (h *Handler) oauthState() string {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	mac := hmac.New(sha256.New, []byte(h.cfg.JWTSecret))
	mac.Write([]byte("riftstats-oauth:" + ts))
	return base64.RawURLEncoding.EncodeToString([]byte(ts + "." + hex.EncodeToString(mac.Sum(nil))))
}

func (h *Handler) verifyOAuthState(state string) bool {
	raw, err := base64.RawURLEncoding.DecodeString(state)
	if err != nil {
		return false
	}
	parts := strings.SplitN(string(raw), ".", 2)
	if len(parts) != 2 {
		return false
	}
	ts, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || time.Since(time.Unix(ts, 0)) > 10*time.Minute {
		return false
	}
	mac := hmac.New(sha256.New, []byte(h.cfg.JWTSecret))
	mac.Write([]byte("riftstats-oauth:" + parts[0]))
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(parts[1]), []byte(expected))
}

// 22_Auth: GET /api/v1/auth/oauth/google — redirects to Google consent screen.
// Requires GOOGLE_CLIENT_ID / GOOGLE_CLIENT_SECRET / GOOGLE_REDIRECT_URL.
func (h *Handler) GoogleOAuthStart(w http.ResponseWriter, r *http.Request) {
	if h.cfg.GoogleClientID == "" || h.cfg.GoogleRedirectURL == "" {
		writeProblem(w, http.StatusNotImplemented, "Social Login Disabled", "Google OAuth is not configured on this server (GOOGLE_CLIENT_ID / GOOGLE_REDIRECT_URL).")
		return
	}
	authURL := "https://accounts.google.com/o/oauth2/v2/auth?" + url.Values{
		"client_id":     {h.cfg.GoogleClientID},
		"redirect_uri":  {h.cfg.GoogleRedirectURL},
		"response_type": {"code"},
		"scope":         {"openid email profile"},
		"state":         {h.oauthState()},
		"access_type":   {"online"},
		"prompt":        {"select_account"},
	}.Encode()
	http.Redirect(w, r, authURL, http.StatusFound)
}

// 22_Auth: GET /api/v1/auth/oauth/google/callback — exchanges the code,
// provisions the user, and redirects to the SPA with a fresh JWT.
func (h *Handler) GoogleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	if h.cfg.GoogleClientID == "" || h.cfg.GoogleClientSecret == "" {
		writeProblem(w, http.StatusNotImplemented, "Social Login Disabled", "Google OAuth is not configured on this server.")
		return
	}
	if !h.verifyOAuthState(r.URL.Query().Get("state")) {
		writeProblem(w, http.StatusBadRequest, "Invalid State", "OAuth state token is missing, expired, or forged.")
		return
	}
	code := r.URL.Query().Get("code")
	if code == "" {
		writeProblem(w, http.StatusBadRequest, "Missing Code", "Authorization code is missing.")
		return
	}

	tokenResp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"code":          {code},
		"client_id":     {h.cfg.GoogleClientID},
		"client_secret": {h.cfg.GoogleClientSecret},
		"redirect_uri":  {h.cfg.GoogleRedirectURL},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		writeProblem(w, http.StatusBadGateway, "OAuth Exchange Failed", err.Error())
		return
	}
	defer tokenResp.Body.Close()
	var tokenBody struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(tokenResp.Body).Decode(&tokenBody); err != nil || tokenBody.AccessToken == "" {
		writeProblem(w, http.StatusBadGateway, "OAuth Exchange Failed", "Google did not return an access token.")
		return
	}

	userReq, _ := http.NewRequestWithContext(r.Context(), "GET", "https://openidconnect.googleapis.com/v1/userinfo", nil)
	userReq.Header.Set("Authorization", "Bearer "+tokenBody.AccessToken)
	userResp, err := http.DefaultClient.Do(userReq)
	if err != nil {
		writeProblem(w, http.StatusBadGateway, "OAuth Userinfo Failed", err.Error())
		return
	}
	defer userResp.Body.Close()
	var userinfo struct {
		Email         string `json:"email"`
		Name          string `json:"name"`
		EmailVerified bool   `json:"email_verified"`
	}
	if err := json.NewDecoder(userResp.Body).Decode(&userinfo); err != nil || userinfo.Email == "" {
		writeProblem(w, http.StatusBadGateway, "OAuth Userinfo Failed", "Google did not return a verified email.")
		return
	}
	if !userinfo.EmailVerified {
		writeProblem(w, http.StatusForbidden, "Email Not Verified", "Only Google-verified emails can sign in.")
		return
	}

	user, err := h.repo.UpsertOAuthUser(userinfo.Email, userinfo.Name)
	if err != nil {
		writeProblem(w, http.StatusInternalServerError, "Provisioning Failed", err.Error())
		return
	}
	token, err := middleware.GenerateToken(user.UserID, user.Email, user.Role, h.cfg.JWTSecret)
	if err != nil {
		writeProblem(w, http.StatusInternalServerError, "Token Generation Failed", err.Error())
		return
	}
	http.Redirect(w, r, "/oauth/callback?token="+url.QueryEscape(token), http.StatusFound)
}

// validateFavoriteRequest rejects malformed bookmark payloads early.
func validateFavoriteRequest(req *domain.CreateFavoriteRequest) error {
	switch strings.ToUpper(req.TargetType) {
	case "SUMMONER":
		if strings.TrimSpace(req.TargetPuuid) == "" {
			return fmt.Errorf("targetPuuid is required for SUMMONER favorites")
		}
	case "CHAMPION":
		if req.TargetChampionID <= 0 {
			return fmt.Errorf("targetChampionId is required for CHAMPION favorites")
		}
	default:
		return fmt.Errorf("targetType must be SUMMONER or CHAMPION")
	}
	if strings.TrimSpace(req.PlatformRegion) == "" {
		return fmt.Errorf("platformRegion is required")
	}
	return nil
}
