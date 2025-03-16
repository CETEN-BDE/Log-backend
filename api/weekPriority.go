package api

import (
	"log-backend/autogen"
	"log-backend/internal/models"

	"github.com/labstack/echo/v4"
)

func (s Server) GetPlanningWeekPriority(ctx echo.Context) error {
	userID := ctx.Get("user_id")
	var days []models.DayPriority
	result := s.db.Find(&days, "account_id = ?", userID)
	if len(days) == 0 {
		for i := 1; i <= 5; i++ {
			dayPriority := models.DayPriority{
				AccountID: userID.(uint), Day: i, Priority: i}
			result := s.db.Create(&dayPriority)
			if result.Error != nil {
				autogen.GetPlanningWeekPriority500JSONResponse{Message: "Can't access to DB"}.VisitGetPlanningWeekPriorityResponse(ctx.Response())
				return result.Error
			}
			days = append(days, dayPriority)
		}
	}
	if result.Error != nil {
		autogen.GetPlanningWeekPriority500JSONResponse{Message: "Can't access to DB"}.VisitGetPlanningWeekPriorityResponse(ctx.Response())
		return result.Error
	}
	var resp autogen.GetPlanningWeekPriority200JSONResponse
	for _, modelDay := range days {
		day := autogen.DayPriority{
			Day:      modelDay.Day,
			Priority: modelDay.Priority,
		}
		resp.WeekPriority = append(resp.WeekPriority, day)
	}
	resp.VisitGetPlanningWeekPriorityResponse(ctx.Response())
	return nil
}

func (s Server) PatchPlanningWeekPriority(ctx echo.Context) error {
	userID := ctx.Get("user_id")
	var req autogen.WeekPriority
	err := ctx.Bind(&req)
	if err != nil {
		autogen.PatchPlanningWeekPriority400JSONResponse{Message: "Wrong Request"}.VisitPatchPlanningWeekPriorityResponse(ctx.Response())
		return err
	}
	
	for _, day := range req.WeekPriority {
		var dayPriority models.DayPriority
		result := s.db.Where("day = ? AND account_id = ?", day.Day, userID).Find(&dayPriority)
		if result.Error != nil {
			autogen.PatchPlanningWeekPriority500JSONResponse{Message: "Can't access to DB"}.VisitPatchPlanningWeekPriorityResponse(ctx.Response())
			return result.Error
		}
		dayPriority.Day = day.Day
		dayPriority.Priority = day.Priority
		result = s.db.Save(&dayPriority)
		if result.Error != nil {
			autogen.PatchPlanningWeekPriority500JSONResponse{Message: "Can't access to DB"}.VisitPatchPlanningWeekPriorityResponse(ctx.Response())
			return result.Error
		}
	}
	autogen.PatchPlanningWeekPriority200JSONResponse{Message: "Success"}.VisitPatchPlanningWeekPriorityResponse(ctx.Response())
	return nil
}
