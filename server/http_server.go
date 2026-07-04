package server

import (
	"log"
	"net/http"

	"benkyou/model"
	"benkyou/service"

	"github.com/labstack/echo"
)

func RunServer(kanjiService *service.Service, port string) {
	e := echo.New()
	e.GET("/get-question", getQuestion(kanjiService))

	if err := e.Start(port); err != nil {
		log.Fatal(err)
	}
}

func getQuestion(kanjiSvc *service.Service) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		level := ctx.QueryParam("level")
		if level == "" {
			level = "n5"
		}

		wordType := ctx.QueryParam("wordType")
		if wordType == "" {
			wordType = model.WordTypeKeiyoushi
		}

		question, err := kanjiSvc.GetQuestionByLevel(level, wordType)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, question)
	}
}
