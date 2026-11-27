// src/hooks/usePlayerId.ts

import { useState } from "react";

const PLAYER_ID_KEY = "septikon.playerId";
const PLAYER_NAME_KEY = "septikon.playerName";

function generateId(): string {
  // crypto.randomUUID есть во всех современных браузерах
  if (typeof crypto !== "undefined" && "randomUUID" in crypto) {
    return crypto.randomUUID();
  }
  return `player-${Date.now()}-${Math.random().toString(36).slice(2, 10)}`;
}

function loadOrCreateId(): string {
  const existing = localStorage.getItem(PLAYER_ID_KEY);
  if (existing) return existing;
  const created = generateId();
  localStorage.setItem(PLAYER_ID_KEY, created);
  return created;
}

export function usePlayerId() {
  const [playerId] = useState(loadOrCreateId);
  const [playerName, setPlayerNameState] = useState(
    () => localStorage.getItem(PLAYER_NAME_KEY) ?? "",
  );

  const setPlayerName = (name: string) => {
    setPlayerNameState(name);
    localStorage.setItem(PLAYER_NAME_KEY, name);
  };

  return { playerId, playerName, setPlayerName };
}