// src/App.tsx

import { useState } from "react";
import { Lobby } from "./screens/Lobby";
import { RoomScreen } from "./screens/Room";
import { usePlayerId } from "./hooks/usePlayerId";

type Screen = { name: "lobby" } | { name: "room"; roomId: string };

export default function App() {
  const [screen, setScreen] = useState<Screen>({ name: "lobby" });
  const { playerId } = usePlayerId(); // Достаем playerId

  if (screen.name === "lobby") {
    return <Lobby onEnterRoom={(roomId) => setScreen({ name: "room", roomId })} />;
  }

  return (
    <RoomScreen
      roomId={screen.roomId}
      playerId={playerId}
      onBack={() => setScreen({ name: "lobby" })}
    />
  );
}