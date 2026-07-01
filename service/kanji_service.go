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

const totalOptions = 4

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
		Data:     data,
		Examples: exampleByLevel,
	}, nil
}

func PopulateKanjiByLevel(dataDir, level string) ([]model.Kanji, []model.Examples, error) {
	kanjiPath := fmt.Sprintf("%s/%s/kanji.json", dataDir, level)

	kanjiFile, err := os.Open(kanjiPath)
	if err != nil {
		return []model.Kanji{}, []model.Examples{}, errors.Wrap(typeErrors.ErrFileNotFound, "fail read kanji file")
	}

	defer kanjiFile.Close()

	data := make([]model.Kanji, 0)

	// exampleByLevel := make([]model.Examples, 0)

	kanjiDecoder := json.NewDecoder(kanjiFile)

	if err := kanjiDecoder.Decode(&data); err != nil {
		return []model.Kanji{}, []model.Examples{}, errors.Wrap(err, "fail decode JSON kanji file")
	}

	kotobaPath := fmt.Sprintf("%s/%s/kotoba.json", dataDir, level)

	kotobaFile, err := os.Open(kotobaPath)
	if err != nil {
		return []model.Kanji{}, []model.Examples{}, errors.Wrap(typeErrors.ErrFileNotFound, "fail kotoba file")
	}

	defer kotobaFile.Close()

	exampleByLevel := make([]model.Examples, 0)

	kotobaDecoder := json.NewDecoder(kotobaFile)

	if err := kotobaDecoder.Decode(&exampleByLevel); err != nil {
		return []model.Kanji{}, []model.Examples{}, errors.Wrap(err, "fail decode JSON kanji file")
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

	rng := rand.IntN(2)
	if rng < 1 {
		return s.getMeaningQuestion(pool, question)
	}

	return s.getKanjiQuestion(pool, question)
}

func (s *Service) getKanjiQuestion(pool []model.Examples, question model.Examples) model.Question {
	seen := make(map[string]bool)

	seen[stripDot(question.Reading)] = true

	options := make([]model.Options, 0)

	options = append(options, model.Options{
		Option: stripDot(question.Reading),
		Answer: true,
	})

	for len(options) < totalOptions {
		candidate := pool[rand.IntN(len(pool))]
		reading := stripDot(candidate.Reading)

		if !seen[reading] {
			seen[reading] = true
			option := model.Options{
				Option: reading,
				Answer: false,
			}
			options = append(options, option)
		}
	}

	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	return model.Question{
		Question: question.Word,
		Meaning:  question.Meaning,
		Furigana: question.Reading,
		Options:  options,
	}
}

func (s *Service) getMeaningQuestion(pool []model.Examples, question model.Examples) model.Question {
	seen := make(map[string]bool)

	seen[lowerWord(question.Meaning)] = true

	options := make([]model.Options, 0)

	options = append(options, model.Options{
		Option: capitalizeFirstLetter(question.Meaning),
		Answer: true,
	})

	for len(options) < totalOptions {
		candidate := pool[rand.IntN(len(pool))]

		if !seen[lowerWord(candidate.Meaning)] {
			seen[lowerWord(candidate.Meaning)] = true
			option := model.Options{
				Option: capitalizeFirstLetter(candidate.Meaning),
				Answer: false,
			}
			options = append(options, option)
		}
	}

	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	return model.Question{
		Question: question.Word,
		Meaning:  capitalizeFirstLetter(question.Meaning),
		Furigana: question.Reading,
		Options:  options,
	}
}

func capitalizeFirstLetter(word string) string {
	return strings.ToUpper(word[:1]) + word[1:]
}

func lowerWord(word string) string{
	return strings.ToLower(word)
}

func stripDot(reading string) string {
	return strings.ReplaceAll(reading, ".", "")
}

func ReadKanjiByID(
	id string,
) {
}
