package api

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log-backend/internal/models"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

var (
    googleOAuthConfig *oauth2.Config
)

type GoogleUser struct {
    ID            string `json:"id"`
    Email         string `json:"email"`
    VerifiedEmail bool   `json:"verified_email"`
    Name          string `json:"name"`
    GivenName     string `json:"given_name"`
    FamilyName    string `json:"family_name"`
    Picture       string `json:"picture"`
}

func InitializeOAuth() {
    googleOAuthConfig = &oauth2.Config{
        ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
        ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
        Scopes: []string{
            "https://www.googleapis.com/auth/userinfo.email",
            "https://www.googleapis.com/auth/userinfo.profile",
        },
        Endpoint: google.Endpoint,
    }
}

func (s *Server) HandleGoogleLogin(c echo.Context) error {
    state := generateState()
    session, err := store.Get(c.Request(), SessionName)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create session"})
    }
    
    session.Values["oauth_state"] = state
    err = session.Save(c.Request(), c.Response().Writer)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save session"})
    }

    url := googleOAuthConfig.AuthCodeURL(state)
    return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (s *Server) HandleGoogleCallback(c echo.Context) error {
    session, err := store.Get(c.Request(), SessionName)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get session"})
    }

    // Verify state
    expectedState := session.Values["oauth_state"]
    if expectedState != c.QueryParam("state") {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid OAuth state"})
    }

    code := c.QueryParam("code")
    token, err := googleOAuthConfig.Exchange(c.Request().Context(), code)
    if err != nil {
        logrus.Errorf("Failed to exchange token: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to exchange token"})
    }

    // Get user info
    client := googleOAuthConfig.Client(c.Request().Context(), token)
    resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
    if err != nil {
        logrus.Errorf("Failed to get user info: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get user info"})
    }
    defer resp.Body.Close()

    var googleUser GoogleUser
    if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
        logrus.Errorf("Failed to decode user info: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to decode user info"})
    }

    // Find or create user
    var account models.Account
    result := s.db.Where("google_id = ?", googleUser.ID).First(&account)
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        // Check if email exists
        emailResult := s.db.Where("email = ?", googleUser.Email).First(&account)
        if errors.Is(emailResult.Error, gorm.ErrRecordNotFound) {
            // Create new account
            account = models.Account{
                Email:     googleUser.Email,
                FirstName: googleUser.GivenName,
                LastName:  googleUser.FamilyName,
                GoogleID:  googleUser.ID,
                Group:     "member", // Default group
            }
            if err := s.db.Create(&account).Error; err != nil {
                logrus.Errorf("Failed to create account: %v", err)
                return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create account"})
            }
        } else if emailResult.Error != nil {
            logrus.Errorf("Database error: %v", emailResult.Error)
            return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
        } else {
            // Link existing account with Google
            account.GoogleID = googleUser.ID
            if err := s.db.Save(&account).Error; err != nil {
                logrus.Errorf("Failed to update account: %v", err)
                return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update account"})
            }
        }
    } else if result.Error != nil {
        logrus.Errorf("Database error: %v", result.Error)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
    }

    // Set session
    session.Values["authenticated"] = true
    session.Values["user_id"] = account.ID
    session.Options = &sessions.Options{
        Path:     "/",
        MaxAge:   86400,
        HttpOnly: true,
    }

    if err := session.Save(c.Request(), c.Response().Writer); err != nil {
        logrus.Errorf("Failed to save session: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save session"})
    }

    // Redirect to frontend
    return c.Redirect(http.StatusTemporaryRedirect, "http://localhost:5173/")
}

func generateState() string {
    b := make([]byte, 16)
    if _, err := rand.Read(b); err != nil {
        return fmt.Sprintf("%d", time.Now().UnixNano())
    }
    return fmt.Sprintf("%x", b)
}
