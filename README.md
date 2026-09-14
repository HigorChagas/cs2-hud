# CS2 HUD

A broadcast HUD for Counter-Strike 2, built to run as an **OBS Browser Source**. It receives match data through the game's **Game State Integration (GSI)**, processes it in a Go backend, and renders it in a static page (plain HTML/CSS/JS, no frameworks).

## Stack

- **Backend:** Go — receives GSI payloads over HTTP and streams the match state via Server-Sent Events (SSE).
- **Frontend:** vanilla HTML, CSS and JavaScript, served as a static file by the backend itself.
- **Reference resolution:** 1920x1080 — supports down to 1280x720 (automatic scaling via `zoom`, works in Chromium/OBS).

## Features

The HUD tracks and reacts to the following match states, all derived directly from the GSI payload:

- CT/T score and current round, with side swaps (halftime) resolved automatically.
- Round history (reason: bomb, defuse, elimination, time).
- Warmup — hides round number, bomb status and round history while active.
- Freeze time — center badge shows "COMPRANDO" (buying).
- Match end (gameover) — badge shows "FIM DE JOGO" (game over).
- Bomb planted, with an indicator of who is currently carrying the C4.
- CT and T panels with HP, armor, weapon, utility, KDA and money per player.
- Alive/dead state, flash, smoke and burning per player.
- Spectator/observer mode, with a dedicated card for the observed player.

## Requirements

- Go 1.26 or later (see `backend/go.mod`).
- Counter-Strike 2.

## Setting up GSI in CS2

The HUD depends on a GSI config file installed inside the game folder, which is **not part of this repository** (it lives outside the project, inside your CS2 installation). Create the file:

```
<CS2 install path>/game/csgo/cfg/gamestate_integration_hud.cfg
```

With the following content (adjust `data` to whatever information you want to expose):

```cfg
"CS2 HUD"
{
    "uri" "http://127.0.0.1:8080/game-state"
    "timeout" "5.0"
    "buffer"  "0.1"
    "throttle" "0.1"
    "heartbeat" "30.0"
    "data"
    {
        "provider"            "1"
        "map"                 "1"
        "round"                "1"
        "player_id"           "1"
        "player_state"        "1"
        "player_weapons"      "1"
        "player_match_stats"  "1"
        "allplayers_id"       "1"
        "allplayers_state"    "1"
        "allplayers_weapons"  "1"
        "allplayers_match_stats" "1"
    }
}
```

Restart CS2 after creating the file.

## Running the server

```bash
cd backend
go run .
```

The server starts on `http://localhost:8080`:

- `POST /game-state` — endpoint that receives GSI payloads.
- `GET /game-state/stream` — SSE stream consumed by the frontend.
- `GET /` — serves the static HUD (`backend/static/index.html`).

## Using it in OBS

1. Add a **Browser Source**.
2. URL: `http://localhost:8080`.
3. Resolution: 1920x1080 (works down to 1280x720).
4. Leave "Shutdown source when not visible" unchecked, so the HUD keeps its state across scenes.

## Project structure

```
backend/
├── main.go              — HTTP handlers, state kept between requests (spectating, C4 carrier, round history)
├── models/
│   ├── cs2.go            — structs for the raw GSI payload
│   └── match.go          — structs for the translated state consumed by the frontend
├── services/
│   └── translator.go     — stateless translation from the GSI payload into HUD state
└── static/
    ├── index.html         — the HUD (HTML + CSS + JS)
    └── icons/             — SVGs used in round history and the bomb carrier indicator
```

## Known limitations (V1)

- Players are identified by name (no stable ID is available in the payload used today); two players with the same name on the same team can collide in the roster.
- GSI doesn't always send the `allplayers` block on every heartbeat; when it's missing, the HUD keeps the last known reliable state (spectating and bomb carrier) instead of assuming a change.
- No round/bomb timer indicator — the GSI data currently parsed doesn't capture that.

## Project philosophy

This is a deliberately simple V1: no frontend frameworks, no build system, no complex visual effects. Priorities, in order: reliability, information readability, performance, visual polish. Full guidelines in [CLAUDE.MD](CLAUDE.MD).
