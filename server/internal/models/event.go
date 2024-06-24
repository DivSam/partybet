package models

type Event struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Date      string  `json:"date"`
	Duration  string  `json:"duration"`
	Outcome   bool    `json:"outcome"`
	Total     float32 `json:"total"`
	HandleYes float32 `json:"handle_yes"`
	HandleNo  float32 `json:"handle_no"`
	YesName   string  `json:"yes_name"`
	NoName    string  `json:"no_name"`
}

// return new event with params
func NewEvent(id int, name string, date string, duration string, outcome bool, total float32, handleYes float32, handleNo float32, yesName string, noName string) *Event {
	return &Event{
		ID:        id,
		Name:      name,
		Date:      date,
		Duration:  duration,
		Outcome:   outcome,
		Total:     total,
		HandleYes: handleYes,
		HandleNo:  handleNo,
		YesName:   yesName,
		NoName:    noName,
	}
}

func (e *Event) UpdateHandle(outcome bool, amount float32) {
	if outcome {
		e.HandleYes += amount
	} else {
		e.HandleNo += amount
	}

	e.Total += amount
}

func (e *Event) UpdateOutcome(outcome bool) {
	e.Outcome = outcome
}
