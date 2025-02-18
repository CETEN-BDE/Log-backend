package api

import (
	"errors"
	"log-backend/autogen"
	"log-backend/internal/models"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["AccountID"] = account.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Generate encoded token
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		autogen.PostLogin500JSONResponse{Message: "Can't generate the token"}.VisitPostLoginResponse(ctx.Response())
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = t
	cookie.Expires = time.Now().Add(24 * time.Hour) // Cookie valid for 24 hours
	// cookie.HttpOnly = true                          // Prevents JavaScript from accessing the cookie
	// cookie.Secure = true                            // Ensure the cookie is sent over HTTPS (in production)
	// cookie.Path = "/"

	// Set the cookie in the response
	ctx.SetCookie(cookie)
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

func (s Server) GetAccount(ctx echo.Context) error {
	AccountID, err := GetAccountIDFromJWT(ctx)
	if err != nil {
		logrus.Error(err)
		autogen.GetAccount401JSONResponse{Message: "Unauthorized"}.VisitGetAccountResponse(ctx.Response())
		return err
	}
	var account models.Account
	result_account := s.db.Where("id = ?", AccountID).First(&account)
	if result_account.Error != nil {
		autogen.GetAccount500JSONResponse{Message: "Can't access to the DB"}.VisitGetAccountResponse(ctx.Response())
		return result_account.Error
	}
	resp := autogen.GetAccount200JSONResponse{
		UserId: uint(AccountID),
		FirstName: account.FirstName,
		LastName: account.LastName,
	}
	resp.VisitGetAccountResponse(ctx.Response())
	logrus.Info(AccountID)
	return nil
}