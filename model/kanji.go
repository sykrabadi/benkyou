package model

type PopulatedKanji map[string][]Kanji

type PopulatedKanjiByLevel map[string]map[string][]Examples

// nested map with outer key is level and inner key is
// word type
type PopulatedKanjiByWordType map[string][]Examples

type Readings struct {
	On  []string `json:"on"`
	Kun []string `json:"kun"`
}

type Examples struct {
	Word    string `json:"word"`
	Reading string `json:"reading"`
	Meaning string `json:"meaning"`
	Type    string `json:"type"`
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

	start = (req.Page - 1) * req.PageSize

	end = start + req.PageSize

	if end > len(kanjiData) {
		end = len(kanjiData)
	}

	return start, end
}

var (
	LevelN5 Level = "n5"
	LevelN4 Level = "n4"
	LevelN3 Level = "n3"
	LevelN2 Level = "n2"
	LevelN1 Level = "n1"

	WordTypeKeiyoushi = "keiyoushi"
	WordTypeMeishi    = "meishi"
)

type Options struct {
	Option string `json:"option"`
	Answer bool   `json:"answer"`
}

type Question struct {
	Question string    `json:"question"`
	Meaning  string    `json:"meaning"`
	Furigana string    `json:"furigana"`
	Type     string    `json:"type"`
	Options  []Options `json:"options"`
}
