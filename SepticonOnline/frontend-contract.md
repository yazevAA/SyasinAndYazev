# Frontend Contract — Septikon Online

> Handoff snapshot: 2026-11-25. This document is a working contract reconstructed from prior project work. Verify every route and JSON field against the current Go handlers and `json` struct tags before relying on it in production.

## 1. Integration principles

- The Go backend is the source of truth for the board, players, turns, resources, robots, rockets, production, repairs, and victory conditions.
- The frontend renders received state and sends user actions. It must not independently simulate game rules.
- After an action, fetch the room state again. Handle server errors even if the UI has already disabled an unavailable action.
- Use the player ID returned/assigned by the backend. Do not assume that `players[0]` is always the local player.
- Initial synchronization can use polling of `GET /api/rooms/{id}` every 1–2 seconds while the game is active. WebSocket support is not confirmed.

## 2. Known routes

| Method | Path | Purpose |
|---|---|---|
| GET | `/health` | Health check |
| GET | `/api/rooms` | List rooms |
| POST | `/api/rooms` | Create room |
| GET | `/api/rooms/{id}` | Get room and game state |
| POST | `/api/rooms/{id}/join` | Join room |
| POST | `/api/rooms/{id}/game/roll` | Roll die |
| POST | `/api/rooms/{id}/game/end-turn` | End turn |
| POST | `/api/rooms/{id}/game/move-robot` | Move robot |
| POST | `/api/rooms/{id}/game/repair` | Repair action |

**Important:** these paths are the routes recorded in the previous handoff. Confirm the current router. In particular, verify whether the `/game` segment is present in every action route.

## 3. Request/response schema verification

The previous handoff discussed these data concepts:

- `Room`
- `Game`
- `Turn`
- `Player`
- `Board`
- `Compartment`
- `Inventory`
- `Robot`
- `Rocket`
- `ProductionTask`
- `RepairTask`

Do not create final TypeScript interfaces from names alone. Inspect the Go structs and their JSON tags. A previous note mentioned `playerID` in action bodies, but exact casing and required fields must be verified.

For each endpoint, document:
1. Exact URL and HTTP method.
2. Required request JSON.
3. Successful response JSON.
4. Error response format and HTTP status codes.
5. Whether the endpoint mutates state.
6. Whether the caller must be the current player.

## 4. Actions that need route verification

The full game interface may need endpoints for:
- adding/creating robots;
- production setup and resource collection;
- rocket launch and rocket status/actions;
- repair task status;
- any other action shown in the official rules.

These are not marked as available until confirmed in the current backend.

## 5. Polling strategy

Initial implementation:
- Fetch room state immediately after joining.
- While the game is active, poll `GET /api/rooms/{id}` every 1–2 seconds.
- Stop polling when leaving the room or unmounting the game screen.
- Avoid overlapping requests; if a request is still in flight, skip the next poll.
- Refresh state immediately after a successful action.
- Display a clear disconnected/loading state if requests fail.

When WebSocket support is implemented, replace polling with server-pushed updates or use polling only as a fallback.

## 6. CORS and local development

Frontend environment variable:

```env
VITE_API_BASE_URL=http://localhost:8080
```

The backend must allow the frontend development origin (commonly `http://localhost:5173`) through CORS. Confirm the actual configured origin and allowed methods/headers before integration.

## 7. Board rendering requirements

- Reproduce the board layout from the board image and official rules supplied by the project owner.
- Do not invent a different number or arrangement of levels/compartments.
- Prefer rendering cells from the board data returned by the backend.
- Keep visual coordinates/layout metadata separate from game state where possible.
- Show ownership, compartment type, robot position, damage, and other visual states only when those values exist in the backend response.
