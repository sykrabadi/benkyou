package server

import (
	"log"
	"net/http"

	"benkyou/service"

	"github.com/labstack/echo"
)

func RunServer(kanjiService *service.Service, port string) {
	e := echo.New()
	e.GET("/get-question/:level", getQuestion(kanjiService))

	if err := e.Start(port); err !=nil{
		log.Fatal(err)
	}
}

func getQuestion(kanjiSvc *service.Service) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		level := ctx.Param("level")
		question := kanjiSvc.GetQuestionByLevel(level)

		return ctx.JSON(http.StatusOK, question)
	}
}
