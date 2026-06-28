package model

type PopulatedKanji map[string][]Kanji

type PopulatedKanjiByLevel map[string][]Examples

type Readings struct {
	On  []string `json:"on"`
	Kun []string `json:"kun"`
}

type Examples struct {
	Word    string `json:"word"`
	Reading string `json:"reading"`
	Meaning string `json:"meaning"`
}

type Kanji struct {
	ID          string     `json:"id"`
	Character   string     `json:"character"`
	Readings    Readings   `json:"readings"`
	Meanings    []string   `json:"meanings"`
	JlptLevel   string     `json:"jlpt_level"`
	StrokeCount int        `json:"stroke_count"`
	Examples    []Examples `json:"examples"`
}

type Level string

type KanjiLevel map[string]Level

type ListKanjiByLevelRequest struct {
	Level    string `json:"level"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

func HandleKanjiByRequestPaging(req ListKanjiByLevelRequest, kanjiData []Kanji) (start, end int) {
	const minLimit = 5

	// if req.PageSize <= 0 {
	// 	req.PageSize = 5
	// }

	start = (req.Page - 1) * req.PageSize

	end = start + req.PageSize

	if end > len(kanjiData) {
		end = len(kanjiData)
	}

	return start, end
}

var (
	LevelN5 Level = "N5"
	LevelN4 Level = "N4"
	LevelN3 Level = "N3"
	LevelN2 Level = "N2"
	LevelN1 Level = "N1"
)

type Options struct {
	Option string `json:"option"`
	Answer bool   `json:"answer"`
}

type Question struct {
	Question string    `json:"question"`
	Meaning  string    `json:"meaning"`
	Furigana string    `json:"furigana"`
	Options  []Options `json:"options"`
}
