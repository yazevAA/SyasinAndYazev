
# Frontend Contract — Septikon Online

Документ описывает взаимодействие Frontend с Backend проекта Septikon Online.

Backend написан на Go. Обмен данными осуществляется через HTTP API в формате JSON.

Базовый адрес:

```text
http://localhost:8080
```

---

# 1. Общие правила

- Все запросы с телом передают JSON.
- Ответы Backend возвращаются в формате JSON.
- Идентификаторы игроков и комнат передаются строками.
- Игровые действия выполняются только игроком, которому принадлежит текущий ход.
- Ошибки возвращаются в формате:

```json
{
  "error": "описание ошибки"
}
```

Для POST-запросов необходимо устанавливать заголовок:

```http
Content-Type: application/json
```

---

# 2. Основные TypeScript-интерфейсы

Ниже приведены интерфейсы для работы с игровыми объектами.

Важно: интерфейсы составлены по текущей структуре Backend. Перед использованием необходимо сверить их с фактическими JSON-ответами Go, особенно поля, которые могут быть изменены в ходе разработки.

## 2.1. Room

Игровая комната.

```ts
export interface Room {
  id: string;
  status: "waiting" | "playing";
  players: RoomPlayer[];
  game?: Game;
}

export interface RoomPlayer {
  id: string;
  name: string;
}
```

## 2.2. Game

Основная игровая информация.

```ts
export type GameStatus =
  | "waiting"
  | "playing"
  | "finished";

export interface Game {
  id: string;
  status: GameStatus;
  players: Player[];
  rockets: Rocket[];
  turn: TurnManager;
  winnerId?: string;
  repairs?: RepairTask[];
}
```

## 2.3. TurnManager

Менеджер очередности ходов.

```ts
export interface TurnManager {
  playerIDs: string[];
  turn: Turn;
}

export interface Turn {
  currentPlayerID: string;
  number: number;
  diceResult: number;
  diceRolled: boolean;
  actionUsed: boolean;
}
```

## 2.4. Player

Информация об игроке.

```ts
export interface Player {
  id: string;
  name: string;
  ready: boolean;
  station: Board;
  inventory: Inventory;
  robots: Robot[];
  productionTasks: ProductionTask[];
}
```

## 2.5. Board

Игровая станция.

```ts
export interface Board {
  levels: Level[];
  surface: Surface;
}

export interface Level {
  number: number;
  compartments: Compartment[];
}

export interface Surface {
  level: number;
  present: boolean;
}
```

## 2.6. Compartment

Отсек станции.

```ts
export type CompartmentType =
  | "unknown"
  | "production"
  | "combat"
  | "energy"
  | "resource"
  | "rocket"
  | "command";

export interface Compartment {
  id: string;
  level: number;
  position: number;
  type: CompartmentType;
  damaged: boolean;
}
```

## 2.7. Inventory

Ресурсы игрока.

```ts
export interface Inventory {
  metal: number;
  uranium: number;
  oxygen: number;
  biomass: number;
  energy: number;
}
```

## 2.8. Robot

Игровой робот.

```ts
export interface Position {
  level: number;
  position: number;
}

export interface Robot {
  id: string;
  ownerID: string;
  position: Position;
  active: boolean;
}
```

## 2.9. ProductionTask

Производственная задача.

```ts
export interface ProductionTask {
  id: string;
  resource: string;
  amount: number;
  requiredTurns: number;
  remainingTurns: number;
  completed: boolean;
}
```

## 2.10. Rocket

Ракета.

```ts
export type RocketStatus =
  | "ready"
  | "flying"
  | "arrived";

export interface Rocket {
  id: string;
  ownerID: string;
  sourceStation: string;
  targetStation: string;
  position: number;
  totalDistance: number;
  speed: number;
  status: RocketStatus;
}
```

## 2.11. RepairTask

Задача ремонта.

```ts
export interface RepairTask {
  id: string;
  ownerID: string;
  compartmentID: string;
  remainingTurns: number;
  completed: boolean;
}
```

---

# 3. HTTP API

## 3.1. Проверка сервера

```http
GET /health
```

Ответ:

```json
{
  "status": "ok"
}
```

## 3.2. Получение списка комнат

```http
GET /api/rooms
```

Получение списка доступных игровых комнат.

## 3.3. Создание комнаты

```http
POST /api/rooms
```

Тело запроса:

```json
{
  "playerID": "player-1",
  "playerName": "Player One"
}
```

В ответ Backend возвращает созданную комнату.

## 3.4. Получение комнаты

```http
GET /api/rooms/{roomID}
```

Пример:

```http
GET /api/rooms/room-1
```

Возвращает актуальное состояние комнаты и игры, если она уже создана.

## 3.5. Подключение второго игрока

```http
POST /api/rooms/{roomID}/join
```

Тело запроса:

```json
{
  "playerID": "player-2",
  "playerName": "Player Two"
}
```

После подключения второго игрока Backend создаёт игровую партию.

## 3.6. Бросок кубика

```http
POST /api/rooms/{roomID}/game/roll
```

Тело запроса:

```json
{
  "playerID": "player-1"
}
```

Пример ответа:

```json
{
  "result": 4
}
```

Результат броска используется для перемещения роботов.

## 3.7. Завершение хода

```http
POST /api/rooms/{roomID}/game/end-turn
```

Тело запроса:

```json
{
  "playerID": "player-1"
}
```

Backend проверяет, что игрок бросил кубик и использовал основное действие.

После завершения хода:

- Обновляются ракеты.
- Обновляется производство.
- Обновляется ремонт.
- Очередь переходит к следующему игроку.

## 3.8. Перемещение робота

```http
POST /api/rooms/{roomID}/game/move-robot
```

Тело запроса:

```json
{
  "playerID": "player-1",
  "robotID": "robot-1",
  "path": [
    {
      "level": 1,
      "position": 1
    },
    {
      "level": 1,
      "position": 2
    },
    {
      "level": 1,
      "position": 3
    },
    {
      "level": 1,
      "position": 4
    },
    {
      "level": 1,
      "position": 5
    }
  ]
}
```

Backend использует результат кубика для определения конечной позиции робота.

Важно: текущая реализация не гарантирует полноценную проверку связности маршрута. Frontend не должен считать собственную проверку маршрута заменой серверной проверке.

## 3.9. Начало ремонта

```http
POST /api/rooms/{roomID}/game/repair
```

Тело запроса:

```json
{
  "playerID": "player-1",
  "taskID": "repair-1",
  "compartmentID": "L1-C1"
}
```

Предварительные параметры ремонта:

- Стоимость: 2 металла.
- Продолжительность: 2 игровых хода.

Значения могут быть изменены.

---

# 4. Как обращаться к данным игры

После получения комнаты Frontend может обращаться к данным следующим образом:

```ts
const room: Room = await getRoom(roomID);

const game = room.game;

if (game) {
  const currentPlayerID =
    game.turn.turn.currentPlayerID;

  const players = game.players;

  const rockets = game.rockets;
}
```

## Полезные пути к данным

| Данные | Путь |
|---|---|
| ID комнаты | room.id |
| Статус комнаты | room.status |
| Список игроков комнаты | room.players |
| ID игры | room.game.id |
| Статус игры | room.game.status |
| Игроки партии | room.game.players |
| Текущий игрок | room.game.turn.turn.currentPlayerID |
| Номер хода | room.game.turn.turn.number |
| Результат кубика | room.game.turn.turn.diceResult |
| Кубик уже брошен | room.game.turn.turn.diceRolled |
| Действие использовано | room.game.turn.turn.actionUsed |
| Список ракет | room.game.rockets |
| Список уровней станции | room.game.players[0].station.levels |
| Отсеки уровня | room.game.players[0].station.levels[0].compartments |
| Ресурсы игрока | room.game.players[0].inventory |
| Роботы игрока | room.game.players[0].robots |
| Производственные задачи | room.game.players[0].productionTasks |
| Задачи ремонта | room.game.repairs |

Индекс `[0]` означает первого игрока. Для второго игрока используется `[1]`.

---

# 5. Пример API-клиента

Рекомендуется создать отдельный файл `src/api/client.ts`.

```ts
const API_URL = "http://localhost:8080";

export async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const response = await fetch(
    `${API_URL}${endpoint}`,
    {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
    }
  );

  const data = await response.json();

  if (!response.ok) {
    throw new Error(
      data.error || "Ошибка запроса к серверу"
    );
  }

  return data as T;
}
```

Пример использования:

```ts
const room = await apiRequest<Room>(
  "/api/rooms/room-1"
);
```

Создание комнаты:

```ts
const room = await apiRequest<Room>(
  "/api/rooms",
  {
    method: "POST",
    body: JSON.stringify({
      playerID: "player-1",
      playerName: "Player One",
    }),
  }
);
```

---

# 6. Проверка возможности действий

Frontend должен учитывать состояние текущего хода.

Пример:

```ts
function isMyTurn(
  game: Game,
  playerID: string
): boolean {
  return (
    game.turn.turn.currentPlayerID === playerID
  );
}
```

Можно ли бросить кубик:

```ts
function canRollDice(
  game: Game,
  playerID: string
): boolean {
  const turn = game.turn.turn;

  return (
    turn.currentPlayerID === playerID &&
    !turn.diceRolled
  );
}
```

Можно ли выполнить основное действие:

```ts
function canPerformAction(
  game: Game,
  playerID: string
): boolean {
  const turn = game.turn.turn;

  return (
    turn.currentPlayerID === playerID &&
    turn.diceRolled &&
    !turn.actionUsed
  );
}
```

Можно ли завершить ход:

```ts
function canEndTurn(
  game: Game,
  playerID: string
): boolean {
  const turn = game.turn.turn;

  return (
    turn.currentPlayerID === playerID &&
    turn.diceRolled &&
    turn.actionUsed
  );
}
```

Эти функции нужны для управления интерфейсом. Окончательную проверку прав и игровых правил всегда выполняет Backend.

---

# 7. Обновление состояния игры

Пока WebSocket не реализован, рекомендуется использовать периодические GET-запросы.

Например, обновлять комнату раз в 1–2 секунды:

```ts
async function updateRoom(roomID: string) {
  return await apiRequest<Room>(
    `/api/rooms/${roomID}`
  );
}
```

При этом:

- Frontend отправляет игровые действия через POST.
- Backend изменяет состояние игры.
- Frontend повторно получает комнату через GET.
- Интерфейс обновляет игровое поле и статистику.

В дальнейшем периодические запросы можно заменить на WebSocket.

---

# 8. Что пока не реализовано

На текущем этапе Backend предоставляет не все игровые действия через HTTP API.

Пока отсутствуют HTTP-маршруты для:

- Добавления роботов.
- Запуска производства.
- Получения произведённых ресурсов.
- Запуска ракет.
- Получения отдельного статуса ремонта.

Также пока необходимо доработать:

- Проверку допустимых маршрутов роботов.
- Полноценную синхронизацию игроков.
- Обработку отключения игроков.
- Сохранение состояния партий.
- Авторизацию.
- WebSocket.

---

# 9. Важные замечания для Frontend-разработчика

1. Не изменять состояние игры самостоятельно — все игровые изменения отправлять Backend.
2. Не считать локальные проверки достаточной защитой от неправильных действий.
3. Не хранить результат кубика как единственный источник истины — получать его с сервера.
4. После каждого действия обновлять состояние комнаты.
5. Не привязываться к фиксированному индексу игрока: определять своего игрока по ID.
6. Учитывать, что игра начинается после подключения второго игрока.
7. Обрабатывать ошибки сервера и показывать их пользователю.
8. Не предполагать, что все игровые действия уже доступны через API.

Документ будет обновляться по мере развития Backend.
