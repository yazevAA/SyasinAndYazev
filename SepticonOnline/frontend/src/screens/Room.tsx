// src/screens/Room.tsx

import { useCallback, useEffect, useRef, useState } from "react";
import { getRoom } from "../api/rooms";
import { rollDice, endTurn } from "../api/game"; // Импортируем новые методы
import { ApiError } from "../api/client";
import { StationBoard } from "../components/StationBoard";
import type { Room } from "../types/room";

interface RoomScreenProps {
  roomId: string;
  playerId: string;
  onBack: () => void;
}

export function RoomScreen({ roomId, playerId, onBack }: RoomScreenProps) {
  const [room, setRoom] = useState<Room | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loadingAction, setLoadingAction] = useState(false);
  const inFlight = useRef(false);

  const fetchRoom = useCallback(async () => {
    if (inFlight.current) return;
    inFlight.current = true;
    try {
      const data = await getRoom(roomId);
      setRoom(data);
      setError(null);
    } catch (e) {
      setError(e instanceof ApiError ? e.message : "Ошибка загрузки");
    } finally {
      inFlight.current = false;
    }
  }, [roomId]);

  useEffect(() => {
    fetchRoom();
    const intervalId = setInterval(fetchRoom, 1500); // Поллинг 1.5 сек
    return () => clearInterval(intervalId);
  }, [fetchRoom]);

  // Обработчик броска
  const handleRoll = async () => {
    setLoadingAction(true);
    try {
      const res = await rollDice(roomId, playerId);
      alert(`Выпало: ${res.result}`);
      await fetchRoom(); // Сразу обновляем состояние
    } catch (e) {
      alert(e instanceof ApiError ? e.message : "Ошибка");
    } finally {
      setLoadingAction(false);
    }
  };

  // Обработчик конца хода
  const handleEndTurn = async () => {
    setLoadingAction(true);
    try {
      await endTurn(roomId, playerId);
      await fetchRoom();
    } catch (e) {
      alert(e instanceof ApiError ? e.message : "Ошибка");
    } finally {
      setLoadingAction(false);
    }
  };

  if (error) {
    return (
      <div style={{ padding: 24 }}>
        <p style={{ color: "red" }}>{error}</p>
        <button onClick={onBack}>← Назад в лобби</button>
      </div>
    );
  }

  if (!room) return <div style={{ padding: 24 }}>Загрузка...</div>;

  if (room.status === "waiting") {
    return (
      <div style={{ padding: 24 }}>
        <h1>Комната {room.id}</h1>
        <h2>Ожидание второго игрока...</h2>
        <ul>{room.players.map((p) => <li key={p.id}>{p.name} {p.id === playerId ? "(Ты)" : ""}</li>)}</ul>
        <button onClick={onBack}>← Назад в лобби</button>
      </div>
    );
  }

  const myData = room.game?.players.find((p) => p.id === playerId);
  const enemyData = room.game?.players.find((p) => p.id !== playerId);
  
  // Данные о ходе
  const turn = room.game?.turn.turn;
  const isMyTurn = turn?.currentPlayerId === playerId;
  const canRoll = isMyTurn && !turn?.diceRolled;
  const canEnd = isMyTurn && turn?.diceRolled && turn?.actionUsed;

  if (!myData || !enemyData || !turn) return <div>Ошибка данных игры.</div>;

  return (
    <div style={{ padding: 24, fontFamily: "sans-serif", maxWidth: 1200, margin: "0 auto" }}>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <h1>Игра: {room.id}</h1>
        <button onClick={onBack}>Выйти</button>
      </div>

      {/* Верхняя панель управления ходом */}
      <div style={{
        padding: 16, background: isMyTurn ? "#dcfce7" : "#f3f4f6",
        borderRadius: 8, marginBottom: 20, display: "flex", justifyContent: "space-between", alignItems: "center"
      }}>
        <div>
          <strong>Ход:</strong> {turn.currentPlayerId === playerId ? "Твой" : "Противника"} (Ход №{turn.number})
        </div>
        <div>
          {turn.diceRolled ? `Кубик: ${turn.diceResult}` : "Кубик не брошен"} 
          {turn.actionUsed ? " | Действие сделано" : " | Действие не сделано"}
        </div>
        <div style={{ display: "flex", gap: 10 }}>
          <button onClick={handleRoll} disabled={!canRoll || loadingAction}>
            🎲 Бросить кубик
          </button>
          <button onClick={handleEndTurn} disabled={!canEnd || loadingAction}>
            Завершить ход
          </button>
        </div>
      </div>

      {/* Игровые поля */}
      <div style={{ display: "flex", gap: 20, flexWrap: "wrap" }}>
        <StationBoard station={myData.station} title={`Твоя станция (${myData.name})`} color="#2563eb" />
        <StationBoard station={enemyData.station} title={`Враг: ${enemyData.name}`} color="#dc2626" />
      </div>
      
      <details style={{ marginTop: 20 }}>
        <summary>Сырой JSON хода (Turn)</summary>
        <pre>{JSON.stringify(turn, null, 2)}</pre>
      </details>
    </div>
  );
}