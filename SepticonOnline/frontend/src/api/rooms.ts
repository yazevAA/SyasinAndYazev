// src/api/rooms.ts

import { request } from "./client";
import type { Room } from "../types/room";

export const listRooms = () => request<Room[]>("/rooms");

export const createRoom = (playerId: string, playerName: string) =>
  request<Room>("/rooms", {
    method: "POST",
    body: JSON.stringify({ playerId, playerName }),
  });

export const getRoom = (roomId: string) => request<Room>(`/rooms/${roomId}`);

export const joinRoom = (
  roomId: string,
  playerId: string,
  playerName: string,
) =>
  request<Room>(`/rooms/${roomId}/join`, {
    method: "POST",
    body: JSON.stringify({ playerId, playerName }),
  });