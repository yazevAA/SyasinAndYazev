// src/components/RobotBoard.tsx
import { useState } from "react";
import type { Board, Robot, Position } from "../types/room";

interface Props {
  station: Board;
  robots: Robot[];
  title: string;
  color: string;
  onMove: (robotId: string, path: Position[]) => void;
  disabled: boolean; // Блокировка, если не наш ход или кубик не брошен
}

export function RobotBoard({ station, robots, title, color, onMove, disabled }: Props) {
  const [selectedRobotId, setSelectedRobotId] = useState<string | null>(null);
  const [path, setPath] = useState<Position[]>([]);

  // Ищем выбранного робота
  const selectedRobot = robots.find((r) => r.id === selectedRobotId);

  const handleCellClick = (level: number, position: number) => {
    if (disabled || !selectedRobotId) return;

    const newPos = { level, position };
    
    // Если кликнули на ту же клетку, где стоит робот - начинаем маршрут
    if (selectedRobot && selectedRobot.position.level === level && selectedRobot.position.position === position) {
      setPath([newPos]);
      return;
    }

    // Добавляем клетку в маршрут
    setPath((prev) => [...prev, newPos]);
  };

  const handleMove = () => {
    if (selectedRobotId && path.length > 0) {
      onMove(selectedRobotId, path);
      setPath([]);
      setSelectedRobotId(null);
    }
  };

  // Проверяем, стоит ли робот в этой клетке
  const getRobotAt = (level: number, position: number) => {
    return robots.find((r) => r.position.level === level && r.position.position === position);
  };

  // Проверяем, является ли клетка частью маршрута
  const isInPath = (level: number, position: number) => {
    return path.some((p) => p.level === level && p.position === position);
  };

  return (
    <div style={{ border: `2px solid ${color}`, padding: 16, borderRadius: 8 }}>
      <h3 style={{ marginTop: 0, color }}>{title}</h3>
      
      <div style={{ marginBottom: 16, minHeight: 40 }}>
        <strong>Выбранный робот:</strong> {selectedRobotId || "нет"} <br/>
        <strong>Маршрут:</strong> {path.map(p => `(L${p.level}-C${p.position})`).join(" → ")}
        <button onClick={handleMove} disabled={!selectedRobotId || path.length === 0} style={{ marginLeft: 10 }}>
          В путь!
        </button>
      </div>

      {/* Поверхность (Уровень 9) */}
      <div style={{ display: "flex", gap: 8, alignItems: "center", marginBottom: 16 }}>
        <div style={{ width: 60, fontSize: 12, fontWeight: "bold", color: "#666" }}>Поверхн.</div>
        <div style={{ display: "flex", gap: 2 }}>
          {Array.from({ length: 21 }).map((_, i) => {
            const pos = i + 1;
            const robot = getRobotAt(9, pos);
            return (
              <div 
                key={`surf-${pos}`}
                onClick={() => handleCellClick(9, pos)}
                title={`Поверхность-${pos}`}
                style={{
                  width: 20, height: 20, backgroundColor: "#e5e7eb", border: "1px solid #ccc",
                  borderRadius: 4, cursor: disabled ? "default" : "pointer",
                  display: "flex", justifyContent: "center", alignItems: "center",
                  outline: isInPath(9, pos) ? "2px solid blue" : "none"
                }}
              >
                {robot && (
                  <div 
                    onClick={(e) => { e.stopPropagation(); setSelectedRobotId(robot.id); }}
                    style={{ width: 14, height: 14, borderRadius: "50%", backgroundColor: color, cursor: "pointer" }}
                    title={`Робот ${robot.id}`}
                  />
                )}
              </div>
            );
          })}
        </div>
      </div>

      {/* Подземные уровни */}
      <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
        {station.levels.map((level) => (
          <div key={level.number} style={{ display: "flex", gap: 8, alignItems: "center" }}>
            <div style={{ width: 60, fontSize: 12, fontWeight: "bold", color: "#666" }}>
              Ур. {level.number}
            </div>
            <div style={{ display: "flex", gap: 2 }}>
              {level.compartments.map((comp) => {
                const robot = getRobotAt(level.number, comp.position);
                return (
                  <div
                    key={comp.id}
                    onClick={() => handleCellClick(level.number, comp.position)}
                    title={`${comp.id} (${comp.type})`}
                    style={{
                      width: 20, height: 20,
                      backgroundColor: comp.damaged ? "#ef4444" : (comp.type === "production" ? "#dcfce7" : "#f3f4f6"),
                      border: "1px solid #ccc", borderRadius: 4,
                      cursor: disabled ? "default" : "pointer",
                      display: "flex", justifyContent: "center", alignItems: "center",
                      outline: isInPath(level.number, comp.position) ? "2px solid blue" : "none"
                    }}
                  >
                    {robot && (
                      <div 
                        onClick={(e) => { e.stopPropagation(); setSelectedRobotId(robot.id); }}
                        style={{ width: 14, height: 14, borderRadius: "50%", backgroundColor: color, cursor: "pointer" }}
                        title={`Робот ${robot.id}`}
                      />
                    )}
                  </div>
                );
              })}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}