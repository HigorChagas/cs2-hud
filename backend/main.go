package main

import (
	"cs2-hud/models"
	"cs2-hud/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var currentMatch models.MatchState
var mu sync.RWMutex

var roundHistory []models.RoundResult

var spectatingState bool
var missingAllPlayersStreak int

// O CS2 nem sempre manda "allplayers" em todo heartbeat (só quando há
// mudança), então um único request sem esse campo não significa que parou
// de espectar. Só desliga o estado depois de algumas atualizações seguidas
// sem "allplayers" — assim o badge não fica piscando.
func updateSpectatingState(hasAllPlayers bool) bool {
	if hasAllPlayers {
		missingAllPlayersStreak = 0
		spectatingState = true
		return spectatingState
	}

	missingAllPlayersStreak++
	if missingAllPlayersStreak >= 3 {
		spectatingState = false
	}
	return spectatingState
}

// Se o placar voltou pra 0-0 mas ainda temos histórico de rounds guardado,
// é sinal de que uma partida nova começou (mapa reiniciado, troca de partida)
// sem o servidor ser reiniciado — limpa o histórico da partida anterior.
func resetRoundHistoryIfNewMatch(totalRoundsPlayed int) {
	if totalRoundsPlayed == 0 && len(roundHistory) > 0 {
		roundHistory = nil
	}
}

// O GSI manda "round.phase":"over" + "round.win_team" várias vezes seguidas
// enquanto o round está terminando, e o contador "map.round" às vezes já
// avança pro próximo round antes do phase sair de "over" — não dá pra confiar
// nele pra evitar duplicata. O placar (soma dos scores) só sobe quando um
// round é realmente fechado, então usamos ele como fonte da verdade.
func recordRoundResult(round models.CS2Round, teams map[string]models.Team, totalRoundsPlayed int) {
	if round.Phase != "over" || round.WinTeam == "" {
		return
	}
	if len(roundHistory) >= totalRoundsPlayed {
		return
	}

	result := models.RoundResult{
		Winner: round.WinTeam,
		Reason: roundOutcomeReason(round, teams, round.WinTeam),
	}
	roundHistory = append(roundHistory, result)
}

// O GSI não manda um "motivo" explícito pro fim do round, então a gente
// infere: bomba explodiu/foi defusada vêm direto de round.bomb; senão,
// olha se o time perdedor ficou com todos mortos (eliminação) ou não (tempo).
func roundOutcomeReason(round models.CS2Round, teams map[string]models.Team, winner string) string {
	switch round.Bomb {
	case "exploded":
		return "bomb"
	case "defused":
		return "defuse"
	}

	loser := "T"
	if winner == "T" {
		loser = "CT"
	}

	team, ok := teams[loser]
	if !ok || len(team.Players) == 0 {
		return "time"
	}

	for _, p := range team.Players {
		if p.Alive {
			return "time"
		}
	}
	return "elimination"
}

func createMatchStateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Método Não Permitido", http.StatusMethodNotAllowed)
		return
	}

	var match models.CS2Payload
	err := json.NewDecoder(r.Body).Decode(&match)

	if err != nil {
		http.Error(w, "JSON Inválido", http.StatusBadRequest)
		return
	}

	matchState := services.TranslateMatchState(match)

	mu.Lock()
	matchState.Spectating = updateSpectatingState(len(match.AllPlayers) > 0)
	totalRoundsPlayed := match.Map.TeamCT.Score + match.Map.TeamT.Score
	resetRoundHistoryIfNewMatch(totalRoundsPlayed)
	recordRoundResult(match.Round, matchState.Teams, totalRoundsPlayed)
	matchState.RoundHistory = roundHistory
	currentMatch = matchState
	fmt.Printf(
		"\r\033[K Round: %d | CT players:%d T players:%d | Bomb planted: %t",
		currentMatch.Round,
		len(currentMatch.Teams["CT"].Players),
		len(currentMatch.Teams["T"].Players),
		currentMatch.C4.Planted,
	)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := map[string]string{
		"status": "sucesso",
	}

	json.NewEncoder(w).Encode(resp)
}

func streamMatchStateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Método Não Permitido", http.StatusMethodNotAllowed)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming não suportado", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			mu.RLock()
			payload, err := json.Marshal(currentMatch)
			mu.RUnlock()

			if err != nil {
				continue
			}

			fmt.Fprintf(w, "data: %s\n\n", payload)
			flusher.Flush()
		}
	}
}

// Enquanto o frontend ainda está mudando bastante, evita que o navegador
// sirva uma versão antiga do index.html/JS do cache do disco.
func noCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/game-state", createMatchStateHandler)
	mux.HandleFunc("/game-state/stream", streamMatchStateHandler)
	mux.Handle("/", noCache(http.FileServer(http.Dir("static"))))

	log.Println("Servidor rodando na porta :8080...")
	httpErr := http.ListenAndServe(":8080", mux)
	if httpErr != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", httpErr)
	}

}
