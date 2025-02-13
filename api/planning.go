package api

import (
	"errors"
	"log-backend/autogen"
	"log-backend/internal/models"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// func (s Server) GetPlanning(ctx echo.Context) error {
// 	autogen.GetPlanning200JSONResponse{Message: "Hello, World!"}.VisitGetPlanningResponse(ctx.Response())
// 	return nil
// }

// GetPlanningUserId implements autogen.ServerInterface.
func (s Server) GetPlanningAccountID(ctx echo.Context, accountID int) error {
	var perms []models.Permanence
	result := s.db.Find(&perms, "account_id = ?", accountID)
	if result.Error != nil {
		autogen.GetPlanningAccountID500JSONResponse{Message: "Can't access to DB"}.VisitGetPlanningAccountIDResponse(ctx.Response())
	}
	var resp autogen.GetPlanningAccountID200JSONResponse
	for _, modelPerm := range perms {
		perm := autogen.Permanence{
			Date:   types.Date{Time: modelPerm.Date},
			Status: modelPerm.Status,
		}
		resp = append(resp, perm)
	}
	logrus.Infof("GET /planning/%d", accountID)
	resp.VisitGetPlanningAccountIDResponse(ctx.Response())
	return nil
}

// PatchPlanningAccountID implements autogen.ServerInterface.
func (s Server) PatchPlanningAccountID(ctx echo.Context, accountID int) error {
	var req autogen.Permanence
	err := ctx.Bind(&req)
	if err != nil {
		autogen.PatchPlanningAccountID400JSONResponse{Message: "Wrong Request"}.VisitPatchPlanningAccountIDResponse(ctx.Response())
		return err
	}
	perm := models.Permanence{
		AccountID: uint(accountID),
		Date:      req.Date.Time,
		Status:    req.Status,
	}
	var permDb models.Permanence
	result := s.db.Where("date = ?", req.Date.Time).First(&permDb, "account_id = ?", accountID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		s.db.Save(&perm)
	}
	if result.Error != nil {
		autogen.PatchPlanningAccountID500JSONResponse{Message: "Can't access to DB"}.VisitPatchPlanningAccountIDResponse(ctx.Response())
		return result.Error
	}
	permDb.AccountID = perm.AccountID
	permDb.Date = perm.Date
	permDb.Status = perm.Status
	result = s.db.Save(&permDb)
	if result.Error != nil  {
		autogen.PatchPlanningAccountID500JSONResponse{Message: "DB Error"}.VisitPatchPlanningAccountIDResponse(ctx.Response())
		return result.Error
	}
	return nil
}
