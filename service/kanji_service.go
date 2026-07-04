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

const (
	totalOptions          = 4
	maxAttemptsMultiplier = 3
)

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
		exampleByLevel[level] = populateExampleByWordType(examplesByLevel)
	}

	return &Service{
		Data:       data,
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

	kanjiDecoder := json.NewDecoder(kanjiFile)

	if err := kanjiDecoder.Decode(&data); err != nil {
		return []model.Kanji{}, []model.Examples{}, errors.Wrap(err, "fail decode JSON kanji file")
	}

	kotobaPath := fmt.Sprintf("%s/%s/kotoba", dataDir, level)

	kotobaDir, err := os.ReadDir(kotobaPath)
	if err != nil {
		return []model.Kanji{}, []model.Examples{}, errors.Wrap(err, "fail open kotoba directory")
	}

	exampleByLevel := make([]model.Examples, 0)

	for _, entry := range kotobaDir {
		wordTypePath := kotobaPath + "/" + entry.Name()
		examples, err := loadExamplesFromFile(wordTypePath)
		if err != nil {
			return []model.Kanji{}, []model.Examples{}, err
		}

		exampleByLevel = append(exampleByLevel, examples...)
	}

	return data, exampleByLevel, nil
}

func loadExamplesFromFile(path string) ([]model.Examples, error) {
	f, err := os.Open(path)
	if err != nil {
		return []model.Examples{}, errors.Wrap(typeErrors.ErrFileNotFound, "fail kotoba file")
	}

	defer f.Close()

	kotobaDecoder := json.NewDecoder(f)

	examples := make([]model.Examples, 0)

	if err := kotobaDecoder.Decode(&examples); err != nil {
		return []model.Examples{}, errors.Wrap(err, "fail decode JSON kanji file")
	}

	return examples, nil
}

func populateExampleByWordType(examples []model.Examples) map[string][]model.Examples {
	examplesByWordType := make(map[string][]model.Examples)
	for _, v := range examples {
		examplesByWordType[v.Type] = append(examplesByWordType[v.Type], v)
	}

	return examplesByWordType
}

func (s *Service) ListKanjiByLevel(req model.ListKanjiByLevelRequest) ([]model.Kanji, error) {
	if _, ok := s.Data[req.Level]; !ok {
		return []model.Kanji{}, typeErrors.ErrLevelDoesNotExists
	}

	data := s.Data[req.Level]

	startPaging, endPaging := model.HandleKanjiByRequestPaging(req, data)

	return data[startPaging:endPaging], nil
}

func (s *Service) GetQuestionByLevel(level string, wordType string) (model.Question, error) {
	pool, ok := s.Examples[level][wordType]
	if !ok || len(pool) == 0 {
		return model.Question{}, errors.Errorf("data untuk level %v dan jenis kata %v tidak ditemukan", level, wordType)
	}

	question := pool[rand.IntN(len(pool))]

	rng := rand.IntN(2)
	if rng < 1 {
		return s.getMeaningQuestion(pool, question)
	}

	return s.getKanjiQuestion(pool, question)
}

func (s *Service) getKanjiQuestion(pool []model.Examples, question model.Examples) (model.Question, error) {
	seen := make(map[string]bool)

	seen[stripDot(question.Reading)] = true

	options := make([]model.Options, 0)

	options = append(options, model.Options{
		Option: stripDot(question.Reading),
		Answer: true,
	})

	maxAttempts := len(pool) * maxAttemptsMultiplier
	attempts := 0

	for len(options) < totalOptions && attempts < maxAttempts {
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
		attempts++
	}

	if len(options) < totalOptions {
		return model.Question{}, typeErrors.ErrInsufficientOptions
	}

	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	return model.Question{
		Question: question.Word,
		Meaning:  question.Meaning,
		Furigana: question.Reading,
		Type:     question.Type,
		Options:  options,
	}, nil
}

func (s *Service) getMeaningQuestion(pool []model.Examples, question model.Examples) (model.Question, error) {
	seen := make(map[string]bool)

	seen[lowerWord(question.Meaning)] = true

	options := make([]model.Options, 0)

	options = append(options, model.Options{
		Option: capitalizeFirstLetter(question.Meaning),
		Answer: true,
	})

	maxAttempts := len(pool) * maxAttemptsMultiplier
	attempts := 0

	for len(options) < totalOptions && attempts < maxAttempts {
		candidate := pool[rand.IntN(len(pool))]

		if !seen[lowerWord(candidate.Meaning)] {
			seen[lowerWord(candidate.Meaning)] = true
			option := model.Options{
				Option: capitalizeFirstLetter(candidate.Meaning),
				Answer: false,
			}
			options = append(options, option)
		}

		attempts++
	}

	if len(options) < totalOptions {
		return model.Question{}, typeErrors.ErrInsufficientOptions
	}

	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	return model.Question{
		Question: question.Word,
		Meaning:  capitalizeFirstLetter(question.Meaning),
		Furigana: question.Reading,
		Type:     question.Type,
		Options:  options,
	}, nil
}

func capitalizeFirstLetter(word string) string {
	return strings.ToUpper(word[:1]) + word[1:]
}

func lowerWord(word string) string {
	return strings.ToLower(word)
}

func stripDot(reading string) string {
	return strings.ReplaceAll(reading, ".", "")
}

// TODO
func ReadKanjiByID(
	id string,
) {
}
