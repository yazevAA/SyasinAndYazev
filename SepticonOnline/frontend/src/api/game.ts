// src/api/game.ts
import { request } from "./client";
import type { Position } from "../types/room";

export const rollDice = (roomId: string, playerId: string) =>
  request<{ result: number }>(`/rooms/${roomId}/game/roll`, {
    method: "POST",
    body: JSON.stringify({ playerId }),
  });

export const endTurn = (roomId: string, playerId: string) =>
  request<{ message: string }>(`/rooms/${roomId}/game/end-turn`, {
    method: "POST",
    body: JSON.stringify({ playerId }),
  });export const moveRobot = (
  roomId: string, 
  playerId: string, 
  robotId: string, 
  path: Position[]
) =>
  request<{ message: string }>(`/rooms/${roomId}/game/move-robot`, {
    method: "POST",
    body: JSON.stringify({ playerId, robotId, path }),
  });