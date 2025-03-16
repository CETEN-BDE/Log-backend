package api

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const (
	SessionName = "log-session"
)

var store *sessions.CookieStore

// Store for sessions
func (s Server) Init() {
    store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
    InitializeOAuth()
}

// SecurityHandler validates the security requirements defined in the OpenAPI spec
func (s *Server) SecurityHandler(ctx echo.Context, securityName string, scopes []string) error {
    logrus.Infof("Security check: %s, scopes: %v", securityName, scopes)

    // Handle the "Auth" security scheme
    if securityName == "Auth" {
        // First check for valid session (basic authentication)
        session, err := store.Get(ctx.Request(), SessionName)
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, "Could not get session")
        }

        // Check if user is authenticated
        if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
            return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
        }

        // Get user ID from session
        userID, ok := session.Values["user_id"].(uint)
        if !ok {
            return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in session")
        }

        // Store the user ID in context for later use
        ctx.Set("user_id", userID)

        // Check for specific scopes
        for _, scope := range scopes {
            switch scope {
            case "BDE":
                // Check if user is a BDE member
                isBDE, err := IsMemberOfBDE(s.db, userID)
                if err != nil {
                    return echo.NewHTTPError(http.StatusInternalServerError, "Could not verify BDE membership")
                }
                
                if !isBDE {
                    return echo.NewHTTPError(http.StatusForbidden, "BDE membership required")
                }
            // Add more scope checks as needed
            // case "Admin":
            //     ...
            }
        }

        // All checks passed
        return nil
    }

    // Unknown security scheme
    return echo.NewHTTPError(http.StatusInternalServerError, "Unknown security requirement")
}
