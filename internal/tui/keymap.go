package tui

// KeyMap captures planned vim-like key bindings.
type KeyMap struct {
	Down     string
	Up       string
	Add      string
	Delete   string
	Complete string
	Edit     string
	Quit     string
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Down:     "j",
		Up:       "k",
		Add:      "a",
		Delete:   "d",
		Complete: "c",
		Edit:     "e",
		Quit:     "q",
	}
}
