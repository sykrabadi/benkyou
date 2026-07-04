package service

import "benkyou/model"

type Service struct {
	Data       model.PopulatedKanji
	// Examples   model.PopulatedKanjiByLevel
	Examples model.PopulatedKanjiByLevel
}

func NewService() *Service {
	return &Service{}
}
