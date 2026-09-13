package models

type CS2Player struct {
	Name       string               `json:"name"`
	Team       string               `json:"team"`
	State      CS2State             `json:"state"`
	Weapons    map[string]CS2Weapon `json:"weapons"`
	MatchStats CS2MatchStats        `json:"match_stats"`
}

type CS2MatchStats struct {
	Kills   int `json:"kills"`
	Deaths  int `json:"deaths"`
	Assists int `json:"assists"`
}

type CS2State struct {
	Health     int `json:"health"`
	Armor      int `json:"armor"`
	Flashed    int `json:"flashed"`
	Smoked     int `json:"smoked"`
	Burning    int `json:"burning"`
	Money      int `json:"money"`
	RoundKills int `json:"round_kills"`
}

type CS2Weapon struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CS2Map struct {
	Name   string   `json:"name"`
	Phase  string   `json:"phase"`
	Round  int      `json:"round"`
	TeamCT CS2Score `json:"team_ct"`
	TeamT  CS2Score `json:"team_t"`
}

type CS2Score struct {
	Score int `json:"score"`
}

type CS2Round struct {
	Phase   string `json:"phase"`
	Bomb    string `json:"bomb"`
	WinTeam string `json:"win_team"`
}

type CS2Payload struct {
	Map        CS2Map               `json:"map"`
	Player     CS2Player            `json:"player"`
	Round      CS2Round             `json:"round"`
	AllPlayers map[string]CS2Player `json:"allplayers"`
}
