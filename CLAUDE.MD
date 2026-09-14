# CS HUD — Claude Instructions

## Stack
- Backend: Go
- Frontend: HTML, CSS and vanilla JavaScript
- Frontend is served as static files by Go
- Game data comes from Counter-Strike GSI
- Target use: OBS Browser Source
- Reference resolution: 1920x1080
- Minimum supported resolution: 1280x720

## Project goals
This is a competitive Counter-Strike broadcast HUD.

Priority:
1. Reliability
2. Information readability
3. Performance
4. Visual polish

## V1 philosophy
Keep the implementation simple.

Do NOT introduce:
- frontend frameworks
- build systems unless explicitly requested
- WebGL
- PixiJS
- complex Canvas effects
- large external dependencies

Advanced visual effects belong to V2/V3.

## Design language
The HUD has its own visual identity.

Prefer:
- dark neutral surfaces
- CT blue and T orange only as accents
- compact layouts
- strong information hierarchy
- subtle chamfered geometry
- thin accent lines
- broadcast-style information density

Avoid:
- SaaS/dashboard appearance
- large rounded cards
- glassmorphism
- decorative gradients
- excessive shadows/glows
- excessive animation
- emoji-like UI except temporary status indicators

## Existing design
The following components already have an approved visual direction:
- top scoreboard
- round history
- CT/T player panels
- observed player

Do NOT redesign these components unless explicitly requested.

## Development rules
- Do not modify backend behavior while doing visual tasks.
- Do not change GSI parsing unless the task requires it.
- Preserve existing selectors/events unless a change is necessary.
- Prefer small, isolated changes.
- Do not refactor unrelated code.
- Do not invent game data that GSI does not provide.
- Do not add comments to code (Go or frontend). Code must be self-explanatory through naming and structure alone.

When asked to make a significant change:
1. inspect the existing implementation;
2. explain the proposed change briefly;
3. wait for approval before implementing, unless explicitly told to implement directly.