package service

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"

	"benkyou/model"
	typeErrors "benkyou/types/errors"

	"github.com/pkg/errors"
)

const totalOptions = 3

func NewKanjiService(
	kanjiDir string,
) (*Service, error) {
	entries, err := os.ReadDir(kanjiDir)
	if err != nil {
		return nil, typeErrors.ErrDirNotFound
	}

	data := make(model.PopulatedKanji)

	exampleByLevel := make(model.PopulatedKanjiByLevel)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		level := entry.Name()
		kanji, examplesByLevel, err := PopulateKanjiByLevel(kanjiDir, level)
		if err != nil {
			return nil, errors.Wrapf(err, "fail to load level %s", level)
		}

		data[level] = kanji
		exampleByLevel[level] = examplesByLevel
	}

	return &Service{
		Data: data,
		Examples: exampleByLevel,
		}, nil
}

func PopulateKanjiByLevel(dataDir, level string) ([]model.Kanji, []model.Examples, error) {
	kanjiPath := fmt.Sprintf("%s/%s/kanji.json", dataDir, level)

	file, err := os.Open(kanjiPath)
	if err != nil {
		return []model.Kanji{}, []model.Examples{}, errors.Wrap(typeErrors.ErrFileNotFound, "fail read kanji file")
	}

	defer file.Close()

	data := make([]model.Kanji, 0)

	exampleByLevel := make([]model.Examples, 0)

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&data); err != nil {
		return []model.Kanji{}, []model.Examples{}, errors.Wrap(err, "fail decode JSON kanji file")
	}

	for i := range data {
		for j := range data[i].Examples {
			exampleByLevel = append(exampleByLevel, model.Examples{
				Word:    data[i].Examples[j].Word,
				Reading: data[i].Examples[j].Reading,
				Meaning: data[i].Examples[j].Meaning,
			})
		}
	}

	return data, exampleByLevel, nil
}

func (s *Service) ListKanjiByLevel(req model.ListKanjiByLevelRequest) ([]model.Kanji, error) {
	if _, ok := s.Data[req.Level]; !ok {
		return []model.Kanji{}, typeErrors.ErrLevelDoesNotExists
	}

	data := s.Data[req.Level]

	startPaging, endPaging := model.HandleKanjiByRequestPaging(req, data)

	return data[startPaging:endPaging], nil
}

func (s *Service) GetQuestionByLevel(level string) model.Question {
	pool := s.Examples[level]

	question := pool[rand.IntN(len(pool))]

	seen := make(map[string]bool)

	seen[question.Reading] = true

	options := make([]model.Options, 0)

	options = append(options, model.Options{
		Option: question.Reading,
		Answer: true,
	})

	for len(options) < totalOptions {
		candidate := pool[:rand.IntN(len(pool))]

		examples := candidate[:len(pool)]

		for _, example := range examples {
			reading := strings.ReplaceAll(example.Word, ".", "")
			if !seen[reading] {
				seen[reading] = true
				option := model.Options{
					Option: example.Word,
					Answer: false,
				}
				options = append(options, option)
			}
		}
	}

	return model.Question{
		Question: question.Word,
		Meaning:  question.Meaning,
		Furigana: question.Reading,
		Options:  options,
	}
}

func ReadKanjiByID(
	id string,
) {
}
