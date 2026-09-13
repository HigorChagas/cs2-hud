package services

import "cs2-hud/models"

func TranslateMatchState(payload models.CS2Payload) models.MatchState {
	teams := translateTeams(payload.Map)
	assignPlayers(teams, payload)

	return models.MatchState{
		Round:          payload.Map.Round,
		Map:            payload.Map.Name,
		Teams:          teams,
		C4:             translateC4(payload.Round),
		Spectating:     len(payload.AllPlayers) > 0,
		Warmup:         payload.Map.Phase == "warmup",
		ObservedPlayer: TranslatePlayer(payload.Player),
	}
}

// Espectando/observando, o CS2 manda "allplayers" com todo mundo; jogando
// ativamente, ele só manda o "player" local. Usa o que estiver disponível.
func assignPlayers(teams map[string]models.Team, payload models.CS2Payload) {
	if len(payload.AllPlayers) > 0 {
		for _, p := range payload.AllPlayers {
			assignPlayer(teams, p)
		}
		return
	}
	assignPlayer(teams, payload.Player)
}

func assignPlayer(teams map[string]models.Team, p models.CS2Player) {
	team, ok := teams[p.Team]
	if !ok {
		return
	}

	player := TranslatePlayer(p)
	team.Players[player.Name] = player
	teams[p.Team] = team
}

func translateTeams(m models.CS2Map) map[string]models.Team {
	return map[string]models.Team{
		"CT": {Name: "CT", Side: "CT", Score: m.TeamCT.Score, Players: map[string]models.Player{}},
		"T":  {Name: "T", Side: "T", Score: m.TeamT.Score, Players: map[string]models.Player{}},
	}
}

func TranslatePlayer(p models.CS2Player) models.Player {
	primary, primaryType, secondary, secondaryType, utilities := translateWeapons(p.Weapons)

	return models.Player{
		Name:                p.Name,
		Team:                p.Team,
		HP:                  p.State.Health,
		Vest:                p.State.Armor,
		PrimaryWeapon:       primary,
		PrimaryWeaponType:   primaryType,
		SecondaryWeapon:     secondary,
		SecondaryWeaponType: secondaryType,
		Utilities:           utilities,
		Alive:               p.State.Health > 0,
		Flashed:             p.State.Flashed,
		Smoked:              p.State.Smoked,
		Burning:             p.State.Burning,
		Money:               p.State.Money,
		RoundKills:          p.State.RoundKills,
		Kills:               p.MatchStats.Kills,
		Deaths:              p.MatchStats.Deaths,
		Assists:             p.MatchStats.Assists,
	}
}

func translateWeapons(weapons map[string]models.CS2Weapon) (primary, primaryType, secondary, secondaryType string, utilities []string) {
	for _, w := range weapons {
		switch {
		case isIgnored(w.Type):
			continue
		case isUtility(w.Type):
			utilities = append(utilities, w.Name)
		case isSecondary(w.Type):
			secondary = w.Name
			secondaryType = w.Type
		default:
			primary = w.Name
			primaryType = w.Type
		}
	}
	return primary, primaryType, secondary, secondaryType, utilities
}

// Faca e C4 não são armas/utilidades do HUD: a faca todo mundo sempre tem,
// e a C4 já é representada por MatchState.C4, não pelo inventário do player.
func isIgnored(weaponType string) bool {
	return weaponType == "Knife" || weaponType == "C4"
}

func isUtility(weaponType string) bool {
	return weaponType == "Grenade"
}

func isSecondary(weaponType string) bool {
	return weaponType == "Pistol"
}

func translateC4(round models.CS2Round) models.C4 {
	return models.C4{
		Planted: round.Bomb == "planted",
	}
}
