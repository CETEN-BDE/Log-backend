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
		autogen.PatchPlanningAccountID404JSONResponse{Message: "Not Found"}.VisitPatchPlanningAccountIDResponse(ctx.Response())
		return result.Error
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
	autogen.PatchPlanningAccountID200JSONResponse{Message: "Success"}.VisitPatchPlanningAccountIDResponse(ctx.Response())
	return nil
}

func (s Server) PostPlanningAccountID(ctx echo.Context, accountID int) error {
	var req autogen.Permanence
	err := ctx.Bind(&req)
	if err != nil {
		autogen.PostPlanningAccountID400JSONResponse{Message: "Wrong Request"}.VisitPostPlanningAccountIDResponse(ctx.Response())
		return err
	}
	perm := models.Permanence{
		AccountID: uint(accountID),
		Date:      req.Date.Time,
		Status:    req.Status,
	}

	var existingPerm models.Permanence
	result := s.db.Where("date = ? AND account_id = ?", req.Date.Time, accountID).First(&existingPerm)
	if result.Error == nil {
		autogen.PostPlanningAccountID409JSONResponse{Message: "Conflict: Date already exists"}.VisitPostPlanningAccountIDResponse(ctx.Response())
		return errors.New("date already exists")
	}
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		autogen.PostPlanningAccountID500JSONResponse{Message: "Can't access to DB"}.VisitPostPlanningAccountIDResponse(ctx.Response())
		return result.Error
	}

	result = s.db.Save(&perm)
	if result.Error != nil {
		autogen.PostPlanningAccountID500JSONResponse{Message: "Can't access to DB"}.VisitPostPlanningAccountIDResponse(ctx.Response())
		return result.Error
	}
	autogen.PostPlanningAccountID200JSONResponse{Message: "Success"}.VisitPostPlanningAccountIDResponse(ctx.Response())
	return nil
}

// DeletePlanningAccountID implements autogen.ServerInterface.
func (s Server) DeletePlanningAccountID(ctx echo.Context, accountID int) error {
	var req autogen.Permanence
	err := ctx.Bind(&req)
	if err != nil {
		autogen.DeletePlanningAccountID400JSONResponse{Message: "Wrong Request"}.VisitDeletePlanningAccountIDResponse(ctx.Response())
		return err
	}

	var permDb models.Permanence
	result := s.db.Where("date = ?", req.Date.Time).First(&permDb, "account_id = ?", accountID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		autogen.DeletePlanningAccountID404JSONResponse{Message: "Not Found"}.VisitDeletePlanningAccountIDResponse(ctx.Response())
		return result.Error
	}
	if result.Error != nil {
		autogen.DeletePlanningAccountID500JSONResponse{Message: "Can't access to DB"}.VisitDeletePlanningAccountIDResponse(ctx.Response())
		return result.Error
	}
	result = s.db.Delete(&permDb)
	if result.Error != nil {
		autogen.DeletePlanningAccountID500JSONResponse{Message: "Can't access to DB"}.VisitDeletePlanningAccountIDResponse(ctx.Response())
		return result.Error
	}
	autogen.DeletePlanningAccountID200JSONResponse{Message: "Success"}.VisitDeletePlanningAccountIDResponse(ctx.Response())
	return nil
}