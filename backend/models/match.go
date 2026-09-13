package models

type Team struct {
	Name    string            `json:"nome"`
	Side    string            `json:"side"`
	Score   int               `json:"score"`
	Players map[string]Player `json:"players"`
}

type Player struct {
	Name                string   `json:"nome"`
	Team                string   `json:"team"`
	HP                  int      `json:"hp"`
	Vest                int      `json:"vest"`
	PrimaryWeapon       string   `json:"primaryWeapon"`
	PrimaryWeaponType   string   `json:"primaryWeaponType"`
	SecondaryWeapon     string   `json:"secondaryWeapon"`
	SecondaryWeaponType string   `json:"secondaryWeaponType"`
	Utilities           []string `json:"utilities"`
	Alive               bool     `json:"alive"`
	Flashed             int      `json:"flashed"`
	Smoked              int      `json:"smoked"`
	Burning             int      `json:"burning"`
	Money               int      `json:"money"`
	RoundKills          int      `json:"roundKills"`
	Kills               int      `json:"kills"`
	Deaths              int      `json:"deaths"`
	Assists             int      `json:"assists"`
}

type MatchState struct {
	Round          int             `json:"round"`
	RoundTime      int             `json:"roundTime"`
	Map            string          `json:"map"`
	Teams          map[string]Team `json:"teams"`
	C4             C4              `json:"c4"`
	Spectating     bool            `json:"spectating"`
	Warmup         bool            `json:"warmup"`
	RoundHistory   []RoundResult   `json:"roundHistory"`
	ObservedPlayer Player          `json:"observedPlayer"`
}

type RoundResult struct {
	Winner string `json:"winner"`
	Reason string `json:"reason"`
}

type C4 struct {
	Carrier *string `json:"carrier"`
	Planted bool    `json:"planted"`
}
