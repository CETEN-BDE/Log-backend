package api

import (
	"log-backend/autogen"
	"log-backend/internal/models"

	"github.com/labstack/echo/v4"
)

func (s Server) GetPlanningWeekPriorityAccountID(ctx echo.Context, accountID int) error {
	var days []models.DayPriority
	result := s.db.Find(&days, "account_id = ?", accountID)
	if len(days) == 0 {
		for i := 1; i <= 5; i++ {
			dayPriority := models.DayPriority{
				AccountID: uint(accountID), Day: i, Priority: i}
			result := s.db.Create(&dayPriority)
			if result.Error != nil {
				autogen.GetPlanningWeekPriorityAccountID500JSONResponse{Message: "Can't access to DB"}.VisitGetPlanningWeekPriorityAccountIDResponse(ctx.Response())
				return result.Error
			}
			days = append(days, dayPriority)
		}
	}
	if result.Error != nil {
		autogen.GetPlanningWeekPriorityAccountID500JSONResponse{Message: "Can't access to DB"}.VisitGetPlanningWeekPriorityAccountIDResponse(ctx.Response())
		return result.Error
	}
	var resp autogen.GetPlanningWeekPriorityAccountID200JSONResponse
	for _, modelDay := range days {
		day := autogen.DayPriority{
			Day:      modelDay.Day,
			Priority: modelDay.Priority,
		}
		resp.WeekPriority = append(resp.WeekPriority, day)
	}
	resp.VisitGetPlanningWeekPriorityAccountIDResponse(ctx.Response())
	return nil
}

func (s Server) PatchPlanningWeekPriorityAccountID(ctx echo.Context, accountID int) error {
	var req autogen.WeekPriority
	err := ctx.Bind(&req)
	if err != nil {
		autogen.PatchPlanningWeekPriorityAccountID400JSONResponse{Message: "Wrong Request"}.VisitPatchPlanningWeekPriorityAccountIDResponse(ctx.Response())
		return err
	}
	
	for _, day := range req.WeekPriority {
		var dayPriority models.DayPriority
		result := s.db.Where("day = ? AND account_id = ?", day.Day, accountID).Find(&dayPriority)
		if result.Error != nil {
			autogen.PatchPlanningWeekPriorityAccountID500JSONResponse{Message: "Can't access to DB"}.VisitPatchPlanningWeekPriorityAccountIDResponse(ctx.Response())
			return result.Error
		}
		dayPriority.Day = day.Day
		dayPriority.Priority = day.Priority
		result = s.db.Save(&dayPriority)
		if result.Error != nil {
			autogen.PatchPlanningAccountID500JSONResponse{Message: "Can't access to DB"}.VisitPatchPlanningAccountIDResponse(ctx.Response())
			return result.Error
		}
	}
	autogen.PatchPlanningAccountID200JSONResponse{Message: "Success"}.VisitPatchPlanningAccountIDResponse(ctx.Response())
	return nil
}
