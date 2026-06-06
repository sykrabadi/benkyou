package service_test

import (
	"testing"

	"benkyou/model"
	typeErrors "benkyou/types/errors"

	"benkyou/service"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	t.Run("it should return error because directory does not exists", func(t *testing.T) {
		_, err := service.NewKanjiService("abc")
		assert.Error(t, err)
		assert.ErrorIs(t, err, typeErrors.ErrDirNotFound)
	})

	t.Run("it should not return error and kanji data are exists", func(t *testing.T) {
		svc, err := service.NewKanjiService("../data")
		assert.NoError(t, err)
		assert.NotEmpty(t, svc.Data)
	})
}

func TestPopulateKanji(t *testing.T) {
	t.Run("it should return error because file does not exists", func(t *testing.T) {
		_, _, err := service.PopulateKanjiByLevel("", "")
		assert.ErrorIs(t, err, typeErrors.ErrFileNotFound)
	})

	t.Run("it should not return error and kanji data exists", func(t *testing.T) {
		n5Level := "N5"
		dirLocation := "../data/"
		data, n5Examples,err := service.PopulateKanjiByLevel(dirLocation, n5Level)

		assert.NoError(t, err)
		assert.NotEmpty(t, data)
		assert.NotEmpty(t, n5Examples)
	})

	t.Run("it should not return error and kanji for N5 and N4 are exists", func(t *testing.T) {
		dirLocation := "../data/"
		n5Level := "N5"
		n5Data, n5Examples, err := service.PopulateKanjiByLevel(dirLocation, n5Level)

		assert.NoError(t, err)
		assert.NotEmpty(t, n5Data)

		n4Level := "N4"
		n4Data, n4Examples, err := service.PopulateKanjiByLevel(dirLocation, n4Level)

		assert.NoError(t, err)
		assert.NotEmpty(t, n4Data)
		assert.NotEmpty(t, n4Examples)
		assert.NotEmpty(t, n5Data)
		assert.NotEmpty(t, n5Examples)
	})
}

func TestListKanjiByLevel(t *testing.T) {
	t.Run("it should return error because level does not exists", func(t *testing.T) {
		svc, err := service.NewKanjiService("../data")
		assert.NoError(t, err)

		req := model.ListKanjiByLevelRequest{
			Level: "N100",
		}

		_, err = svc.ListKanjiByLevel(req)
		assert.ErrorIs(t, err, typeErrors.ErrLevelDoesNotExists)
	})

	t.Run("it should not return error and N5 kanji are exists", func(t *testing.T) {
		svc, err := service.NewKanjiService("../data")
		assert.NoError(t, err)

		req := model.ListKanjiByLevelRequest{
			Level: "n5",
			PageSize: 2,
			Page: 1,
		}
		result, err :=svc.ListKanjiByLevel(req)
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})
}

func TestGetQuestionByLevel(t *testing.T){
	t.Run("it should not return empty question", func(t *testing.T) {
		svc, err := service.NewKanjiService("../data")
		assert.NoError(t, err)

		assert.NoError(t, err)
		question := svc.GetQuestionByLevel("n5")

		assert.NotEmpty(t, question.Question)
		assert.NotEmpty(t, question.Meaning)
		assert.NotEmpty(t, question.Furigana)
		assert.NotEmpty(t, question.Options)
	})
}
