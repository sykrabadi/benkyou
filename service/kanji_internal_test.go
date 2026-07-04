package service

import (
	"testing"

	"benkyou/model"

	"github.com/stretchr/testify/assert"
)

func TestGetMeaningQuestion(t *testing.T) {
	t.Run("it should not return error", func(t *testing.T) {
		svc, err := NewKanjiService("../data")
		assert.NoError(t, err)
		question, err := svc.getMeaningQuestion(svc.Examples["n5"]["keiyoushi"], model.Examples{
			Word:    "欲しい",
			Reading: "ほ.しい",
			Meaning: "Ingin mendapatkan",
			Type:    "keiyoushi",
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, question.Question)
		assert.NotEmpty(t, question.Options)
	})
}

func TestGetKanjiQuestion(t *testing.T) {
	t.Run("it should not return error", func(t *testing.T) {
		svc, err := NewKanjiService("../data")
		assert.NoError(t, err)
		question, err := svc.getKanjiQuestion(svc.Examples["n5"]["keiyoushi"], model.Examples{
			Word:    "欲しい",
			Reading: "ほ.しい",
			Meaning: "Ingin mendapatkan",
			Type:    "keiyoushi",
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, question.Question)
		assert.NotEmpty(t, question.Options)
	})
}
