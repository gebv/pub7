package chats

const (
	START = "start"
	MENU  = "menu"
)

type Script struct {
	ScruptID string `toml:"id"`
	Title    string `toml:"title"`

	// question id
	StartQID string `toml:"start_qid"`

	Options []OptionOfScript `toml:"opts"`
}

type Workspace struct {
	Questions []Question `toml:"q"`
	Scripts   []Script   `toml:"s"`
}

func (w Workspace) FindQuestion(qid string) *Question {
	for _, q := range w.Questions {
		if q.QID == qid {
			return &q
		}
	}
	return nil
}

func (w Workspace) FindScript(sid string) *Script {
	for _, s := range w.Scripts {
		if s.ScruptID == sid {
			return &s
		}
	}
	return nil
}

type OptionOfScript struct {
	Key     string `toml:"key"`
	NextSID string `toml:"sid"`
}
