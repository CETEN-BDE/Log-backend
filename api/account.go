package api

import (
	"errors"
	"log-backend/autogen"
	"log-backend/internal/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
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

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = req.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Generate encoded token
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		autogen.PostLogin500JSONResponse{Message: "Can't generate the token"}.VisitPostLoginResponse(ctx.Response())
		return err
	}

	autogen.PostLogin200Response{Headers: autogen.PostLogin200ResponseHeaders{SetCookie: t}}.VisitPostLoginResponse(ctx.Response())
	return nil
}

func (s Server) PostRegister(ctx echo.Context) error {
	var req autogen.LoginRequest
	err := ctx.Bind(&req)
	if err != nil {
		autogen.PostRegister400JSONResponse{Message: "Mal formed request"}.VisitPostRegisterResponse(ctx.Response())
		return err
	}
	result := s.db.Where("email = ?", req.Email).First(&req)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			autogen.PostRegister500JSONResponse{Message: "Can't hash the password"}.VisitPostRegisterResponse(ctx.Response())
			return err
		}

		account := models.Account{
			Email:    req.Email,
			Password: string(hashedPassword),	
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