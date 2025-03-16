package api

import (
	"errors"
	"log-backend/autogen"
	"log-backend/internal/models"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (s Server) PostLogin(ctx echo.Context) error {
	var req autogen.LoginRequest
	err := ctx.Bind(&req)
	if err != nil {
		autogen.PostLogin400JSONResponse{Message: "missing or malformed request login"}.VisitPostLoginResponse(ctx.Response())
		return err
	}

	var account models.Account
	result_account := s.db.Where("email = ?", req.Email).First(&account)
	if errors.Is(result_account.Error, gorm.ErrRecordNotFound) {
		autogen.PostLogin401JSONResponse{Message: "Invalid email or password"}.VisitPostLoginResponse(ctx.Response())
		return result_account.Error
	}
	if result_account.Error != nil {
		autogen.PostLogin500JSONResponse{Message: "Can't access to the DB"}.VisitPostLoginResponse(ctx.Response())
		return result_account.Error
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(req.Password))
	if err != nil {
		autogen.PostLogin401JSONResponse{Message: "Invalid email or password"}.VisitPostLoginResponse(ctx.Response())
		return nil
	}

	// Create new session
	session, err := store.Get(ctx.Request(), SessionName)
	if err != nil {
		autogen.PostLogin500JSONResponse{Message: "Unable to create session"}.VisitPostLoginResponse(ctx.Response())
		return err
	}

	// Set session values
	session.Values["authenticated"] = true
	session.Values["user_id"] = account.ID
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400, // 24 hours
		HttpOnly: true,  // Prevents JavaScript from accessing the cookie
	}

	// Save session
	if err := session.Save(ctx.Request(), ctx.Response().Writer); err != nil {
		autogen.PostLogin500JSONResponse{Message: "Unable to save session"}.VisitPostLoginResponse(ctx.Response())
		return err
	}
	autogen.PostLogin200JSONResponse{Message: "Login Success"}.VisitPostLoginResponse(ctx.Response())
	return nil
}

func (s Server) PostRegister(ctx echo.Context) error {
	var req autogen.RegisterRequest
	err := ctx.Bind(&req)
	if err != nil {
		autogen.PostRegister400JSONResponse{Message: "Mal formed request"}.VisitPostRegisterResponse(ctx.Response())
		return err
	}
	var account models.Account
	result := s.db.Where("email = ?", req.Email).First(&account)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			autogen.PostRegister500JSONResponse{Message: "Can't hash the password"}.VisitPostRegisterResponse(ctx.Response())
			return err
		}

		account := models.Account{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Group:     req.Group,
			ClubId:    int(req.ClubId),
			Email:     req.Email,
			Password:  string(hashedPassword),
		}
		result_create := s.db.Create(&account)
		if result_create.Error != nil {
			autogen.PostRegister500JSONResponse{Message: "Can't access to the DB"}.VisitPostRegisterResponse(ctx.Response())
			return result_create.Error
		}
	} else {
		autogen.PostRegister403JSONResponse{Message: "Email already exists"}.VisitPostRegisterResponse(ctx.Response())
		return nil
	}

	autogen.PostRegister200JSONResponse{Message: "Register Success"}.VisitPostRegisterResponse(ctx.Response())
	return nil
}

func (s Server) GetAccountsClubId(ctx echo.Context, clubId int) error {
	accountList := []models.Account{}
	result := s.db.Where("club_id = ?", clubId).Find(&accountList)
	if result.Error != nil {
		autogen.GetAccountsClubId500JSONResponse{Message: "Can't access to the DB"}.VisitGetAccountsClubIdResponse(ctx.Response())
		return result.Error
	}
	resp := autogen.GetAccountsClubId200JSONResponse{}
	for _, account := range accountList {
		acc := autogen.Account{FirstName: account.FirstName, LastName: account.LastName, Group: &account.Group, UserId: account.ID}
		resp = append(resp, acc)
	}
	logrus.Infof("GET /account/%d: Account Club Id", clubId)
	resp.VisitGetAccountsClubIdResponse(ctx.Response())
	return nil
}

func (s Server) PostLogout(ctx echo.Context) error {
	session, err := store.Get(ctx.Request(), SessionName)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Could not get session"})
	}

	// Clear session values
	session.Values["authenticated"] = false
	session.Values["user_id"] = nil
	session.Options.MaxAge = -1 // Delete the cookie

	err = session.Save(ctx.Request(), ctx.Response().Writer)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Could not save session"})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Successfully logged out"})
}

func (s Server) GetAuthGoogle(ctx echo.Context) error {
    return s.HandleGoogleLogin(ctx)
}

func (s Server) GetAuthGoogleCallback(ctx echo.Context, params autogen.GetAuthGoogleCallbackParams) error {
    return s.HandleGoogleCallback(ctx)
}

func (s Server) GetAccount(ctx echo.Context) error {
	userID := ctx.Get("user_id")
	if userID == nil {
		logrus.Error("No user ID in context")
		autogen.GetAccount401JSONResponse{Message: "Unauthorized"}.VisitGetAccountResponse(ctx.Response())
		return errors.New("unauthorized")
	}
	var account models.Account
	result_account := s.db.Where("id = ?", userID).First(&account)
	if result_account.Error != nil {
		autogen.GetAccount500JSONResponse{Message: "Can't access to the DB"}.VisitGetAccountResponse(ctx.Response())
		return result_account.Error
	}
	resp := autogen.GetAccount200JSONResponse{
		UserId:    uint(userID.(uint)),
		FirstName: account.FirstName,
		LastName:  account.LastName,
		Group:     &account.Group,
		ClubId:    &account.ClubId,
	}
	resp.VisitGetAccountResponse(ctx.Response())
	logrus.Info("Retrieved account for user: ", userID)
	return nil
}
