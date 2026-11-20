package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"septikon-online/internal/robot"
	"septikon-online/internal/room"
	"strings"
)

type Handler struct {
	rooms *room.Manager
}

func NewHandler(manager *room.Manager) *Handler {
	return &Handler{
		rooms: manager,
	}
}

type CreateRoomRequest struct {
	PlayerID   string `json:"playerId"`
	PlayerName string `json:"playerName"`
}

type JoinRoomRequest struct {
	PlayerID   string `json:"playerId"`
	PlayerName string `json:"playerName"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{
		Error: message,
	})
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		return err
	}

	return nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/rooms" {
		h.handleRooms(w, r)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/rooms/")

	if path == r.URL.Path || path == "" {
		http.NotFound(w, r)
		return
	}

	parts := strings.Split(path, "/")

	if len(parts) == 1 {
		h.handleRoom(w, r, parts[0])
		return
	}

	if len(parts) == 2 && parts[1] == "join" {
		h.handleJoinRoom(w, r, parts[0])
		return
	}

	if len(parts) == 3 && parts[1] == "game" {
		h.handleGameAction(w, r, parts[0], parts[2])
		return
	}

	http.NotFound(w, r)
}

// GET /api/rooms
// POST /api/rooms
func (h *Handler) handleRooms(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, h.rooms.ListRooms())

	case http.MethodPost:
		var request CreateRoomRequest

		if err := decodeJSON(w, r, &request); err != nil {
			writeError(w, http.StatusBadRequest, "некорректный JSON")
			return
		}

		createdRoom, err := h.rooms.CreateRoom(
			request.PlayerID,
			request.PlayerName,
		)

		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeJSON(w, http.StatusCreated, createdRoom)

	default:
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
	}
}

// GET /api/rooms/{id}
func (h *Handler) handleRoom(w http.ResponseWriter, r *http.Request, roomID string) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}

	foundRoom, err := h.rooms.GetRoom(roomID)

	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, foundRoom)
}

// POST /api/rooms/{id}/join
func (h *Handler) handleJoinRoom(w http.ResponseWriter, r *http.Request, roomID string) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}

	var request JoinRoomRequest

	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}

	joinedRoom, err := h.rooms.JoinRoom(
		roomID,
		request.PlayerID,
		request.PlayerName,
	)

	if err != nil {
		status := http.StatusBadRequest

		if errors.Is(err, room.ErrRoomNotFound) {
			status = http.StatusNotFound
		}

		writeError(w, status, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, joinedRoom)
}

type GameActionRequest struct {
	PlayerID      string           `json:"playerId"`
	RobotID       string           `json:"robotId"`
	Path          []robot.Position `json:"path"`
	TaskID        string           `json:"taskId"`
	CompartmentID string           `json:"compartmentId"`
}

type DiceResponse struct {
	Result int `json:"result"`
}

func (h *Handler) handleGameAction(
	w http.ResponseWriter,
	r *http.Request,
	roomID string,
	action string,
) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}

	var request GameActionRequest

	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}

	switch action {
	case "roll":
		result, err := h.rooms.RollDice(roomID, request.PlayerID)

		if err != nil {
			status := http.StatusBadRequest

			if errors.Is(err, room.ErrRoomNotFound) {
				status = http.StatusNotFound
			}

			writeError(w, status, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, DiceResponse{
			Result: result,
		})

	case "end-turn":
		err := h.rooms.EndTurn(roomID, request.PlayerID)

		if err != nil {
			status := http.StatusBadRequest

			if errors.Is(err, room.ErrRoomNotFound) {
				status = http.StatusNotFound
			}

			writeError(w, status, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"message": "ход завершен",
		})

	default:
		http.NotFound(w, r)

	case "move-robot":
		err := h.rooms.MoveRobot(
			roomID,
			request.PlayerID,
			request.RobotID,
			request.Path,
		)

		if err != nil {
			status := http.StatusBadRequest

			if errors.Is(err, room.ErrRoomNotFound) {
				status = http.StatusNotFound
			}

			writeError(w, status, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"message": "робот перемещен",
		})

	case "repair":
		err := h.rooms.StartRepair(
			roomID,
			request.PlayerID,
			request.TaskID,
			request.CompartmentID,
		)

		if err != nil {
			status := http.StatusBadRequest

			if errors.Is(err, room.ErrRoomNotFound) {
				status = http.StatusNotFound
			}

			writeError(w, status, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"message": "ремонт начат",
		})
	}
}
