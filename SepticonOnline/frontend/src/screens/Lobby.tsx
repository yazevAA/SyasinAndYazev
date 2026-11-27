// src/screens/Lobby.tsx

import { useCallback, useEffect, useState } from "react";
import { createRoom, joinRoom, listRooms } from "../api/rooms";
import { ApiError, checkHealth } from "../api/client";
import { usePlayerId } from "../hooks/usePlayerId";
import type { Room } from "../types/room";

interface LobbyProps {
  onEnterRoom: (roomId: string) => void;
}

export function Lobby({ onEnterRoom }: LobbyProps) {
  const { playerId, playerName, setPlayerName } = usePlayerId();
  const [rooms, setRooms] = useState<Room[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [online, setOnline] = useState<boolean | null>(null);

  const refresh = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await listRooms();
      setRooms(data);
      setOnline(true);
    } catch (e) {
      setOnline(false);
      setError(e instanceof ApiError ? e.message : "Не удалось загрузить комнаты");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    checkHealth().then(setOnline);
    void refresh();
    // polling списка комнат раз в 4 секунды (контракт разрешает)
    const id = setInterval(() => void refresh(), 4000);
    return () => clearInterval(id);
  }, [refresh]);

  const nameReady = playerName.trim().length > 0;

  const handleCreate = async () => {
    if (!nameReady) {
      setError("Введи имя игрока");
      return;
    }
    setError(null);
    try {
      const created = await createRoom(playerId, playerName.trim());
      onEnterRoom(created.id);
    } catch (e) {
      setError(e instanceof ApiError ? e.message : "Не удалось создать комнату");
    }
  };

  const handleJoin = async (roomId: string) => {
    if (!nameReady) {
      setError("Введи имя игрока");
      return;
    }
    setError(null);
    try {
      await joinRoom(roomId, playerId, playerName.trim());
      onEnterRoom(roomId);
    } catch (e) {
      setError(e instanceof ApiError ? e.message : "Не удалось войти в комнату");
    }
  };

  return (
    <div style={styles.page}>
      <h1>Septikon Online — Лобби</h1>

      <div style={styles.status}>
        Статус сервера:{" "}
        {online === null ? "проверка…" : online ? "🟢 онлайн" : "🔴 недоступен"}
      </div>

      <div style={styles.block}>
        <label style={styles.label}>
          Имя игрока
          <input
            style={styles.input}
            value={playerName}
            onChange={(e) => setPlayerName(e.target.value)}
            placeholder="Например, Игрок 1"
          />
        </label>
        <small style={styles.hint}>playerId: {playerId}</small>
      </div>

      <div style={styles.block}>
        <button style={styles.primary} onClick={handleCreate} disabled={!nameReady}>
          Создать комнату
        </button>
        <button style={styles.secondary} onClick={() => void refresh()} disabled={loading}>
          Обновить список
        </button>
      </div>

      {error && <div style={styles.error}>{error}</div>}

      <h2>Открытые комнаты</h2>
      {rooms.length === 0 && !loading && (
        <p style={styles.hint}>
          Пока нет открытых комнат. Создай свою — она появится здесь для второго игрока.
        </p>
      )}

      <ul style={styles.list}>
        {rooms.map((room) => {
          const alreadyIn = room.players.some((p) => p.id === playerId);
          return (
            <li key={room.id} style={styles.listItem}>
              <div>
                <strong>{room.id}</strong>{" "}
                <span style={styles.statusBadge}>
                  {room.status === "waiting" ? "ожидает игрока" : "идёт игра"}
                </span>
                <div style={styles.hint}>
                  Игроки:{" "}
                  {room.players.map((p) => p.name).join(", ") || "—"}
                </div>
              </div>
              <button
                style={styles.primary}
                onClick={() => void handleJoin(room.id)}
                disabled={!nameReady || alreadyIn}
              >
                {alreadyIn ? "Ты уже здесь" : "Войти"}
              </button>
            </li>
          );
        })}
      </ul>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  page: {
    maxWidth: 720,
    margin: "0 auto",
    padding: 24,
    fontFamily: "system-ui, sans-serif",
  },
  status: { marginBottom: 16, fontSize: 14 },
  block: { marginBottom: 16, display: "flex", gap: 12, alignItems: "center", flexWrap: "wrap" },
  label: { display: "flex", flexDirection: "column", gap: 4, flex: 1, minWidth: 240 },
  input: { padding: "8px 10px", fontSize: 16, borderRadius: 6, border: "1px solid #ccc" },
  hint: { color: "#666", fontSize: 12 },
  primary: {
    padding: "8px 16px",
    fontSize: 15,
    borderRadius: 6,
    border: "none",
    background: "#2563eb",
    color: "#fff",
    cursor: "pointer",
  },
  secondary: {
    padding: "8px 16px",
    fontSize: 15,
    borderRadius: 6,
    border: "1px solid #2563eb",
    background: "#fff",
    color: "#2563eb",
    cursor: "pointer",
  },
  error: {
    padding: "10px 12px",
    borderRadius: 6,
    background: "#fee2e2",
    color: "#991b1b",
    marginBottom: 16,
  },
  list: { listStyle: "none", padding: 0, margin: 0, display: "grid", gap: 8 },
  listItem: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
    gap: 12,
    padding: 12,
    border: "1px solid #e5e7eb",
    borderRadius: 8,
  },
  statusBadge: {
    fontSize: 12,
    padding: "2px 8px",
    borderRadius: 999,
    background: "#eef2ff",
    color: "#3730a3",
  },
};