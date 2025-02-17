package api

import (
	"log-backend/autogen"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
	"github.com/sirupsen/logrus"
)

func autoFill(start types.Date, end types.Date) error {
	logrus.Infof("AutoFill is not implemented yet for %s to %s", start.String(), end.String())
	return nil
}

func (s Server) PostPlanningAutofill(ctx echo.Context) error {
	var req autogen.Autofill
	err := ctx.Bind(&req)
	if err != nil {
		autogen.PostPlanningAutofill400JSONResponse{Message: "Wrong Request"}.VisitPostPlanningAutofillResponse(ctx.Response())
		return err
	}
	err = autoFill(req.StartDate, req.StartDate)
	if err != nil {
		autogen.PostPlanningAutofill500JSONResponse{Message: "AutoFill failed because :"}.VisitPostPlanningAutofillResponse(ctx.Response())
		return err
	}
	autogen.PostPlanningAutofill200JSONResponse{Message: "Success"}.VisitPostPlanningAutofillResponse(ctx.Response())
	return nil
}