package models

type Series struct {
	ID           int           `json:"id"`
	Title        string        `json:"title"`
	Participants []Participant `json:"participants"`
}

type Participant struct {
	Roster Roster `json:"roster"`
}

type Roster struct {
	ID     int    `json:"id"`
	TeamId TeamId `json:"team"`
	LineUp LineUp `json:"line_up"`
}

type TeamId struct {
	ID int `json:"id"`
}

type PlayerId struct {
	ID int `json:"id"`
}

type Player struct {
	ID       int    `json:"id"`
	Nickname string `json:"nick_name"`
}

type LineUp struct {
	Players []PlayerId `json:"players"`
}

type SeriesDetails struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
