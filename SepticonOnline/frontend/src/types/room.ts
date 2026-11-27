// src/types/room.ts

export type RoomStatus = "waiting" | "playing";

export interface RoomPlayer {
  id: string;
  name: string;
}

export interface Room {
  id: string;
  status: RoomStatus;
  players: RoomPlayer[];
  game?: Game; // отсутствует, пока status === "waiting" (omitempty)
}

export type GameStatus = "waiting" | "playing" | "finished";

export interface Game {
  id: string;
  status: GameStatus;
  players: Player[];
  rockets: Rocket[];
  turn: TurnManager;
  winnerId?: string;
  repairs: RepairTask[];
}

export interface Player {
  id: string;
  name: string;
  ready: boolean;
  station: Board;
  inventory: Inventory;
  robots: Robot[];
  productionTasks: ProductionTask[];
}

export interface Board {
  levels: Level[];
  surface: Surface;
}

export interface Level {
  number: number;
  compartments: Compartment[];
}

export interface Surface {
  level: number; // 9
  present: boolean;
}

export type CompartmentType =
  | "unknown"
  | "production"
  | "combat"
  | "energy"
  | "resource"
  | "rocket"
  | "command";

export interface Compartment {
  id: string; // "L1-C1"
  level: number;
  position: number;
  type: CompartmentType;
  damaged: boolean;
}

export interface Position {
  level: number; // 1..9
  position: number; // 1..21
}

export interface Robot {
  id: string;
  ownerId: string;
  position: Position;
  active: boolean;
}

export type ResourceType = "metal" | "uranium" | "oxygen" | "biomass" | "energy";

export interface Inventory {
  metal: number;
  uranium: number;
  oxygen: number;
  biomass: number;
  energy: number;
}

export interface ProductionTask {
  id: string;
  compartmentId: string;
  resource: ResourceType;
  amount: number;
  requiredTurns: number;
  remainingTurns: number;
  completed: boolean;
  collected: boolean;
}

export interface RepairTask {
  id: string;
  ownerId: string;
  compartmentId: string;
  remainingTurns: number;
  completed: boolean;
}

export type RocketStatus = "ready" | "flying" | "arrived";

export interface Rocket {
  id: string;
  ownerId: string;
  sourceStation: string;
  targetStation: string;
  position: number;
  totalDistance: number;
  speed: number;
  status: RocketStatus;
}

export interface TurnState {
  currentPlayerId: string;
  number: number;
  diceResult: number;
  diceRolled: boolean;
  actionUsed: boolean;
}

export interface TurnManager {
  playerIds: string[];
  turn: TurnState; // да, дважды turn (Manager -> Turn)
}