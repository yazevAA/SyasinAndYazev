// src/components/StationBoard.tsx
import type { Board } from "../types/room";

interface Props {
  station: Board;
  title: string;
  color: string;
}

export function StationBoard({ station, title, color }: Props) {
  return (
    <div style={{ border: `2px solid ${color}`, padding: 16, borderRadius: 8 }}>
      <h3 style={{ marginTop: 0, color }}>{title}</h3>
      
      {/* Поверхность (уровень 9) */}
      <div style={{ marginBottom: 16, padding: 8, background: "#e5e7eb", borderRadius: 4, textAlign: "center", fontWeight: "bold" }}>
        Поверхность (Уровень 9)
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
                // Цвет в зависимости от типа и повреждения
                let bg = "#fff";
                let borderColor = "#ccc";
                
                if (comp.type === "production") bg = "#dcfce7"; // зеленый
                else if (comp.type === "combat") bg = "#fee2e2"; // красный
                else if (comp.type === "energy") bg = "#fef9c3"; // желтый
                else if (comp.type === "unknown") bg = "#f3f4f6"; // серый
                
                if (comp.damaged) {
                  bg = "#ef4444"; // красный при повреждении
                  borderColor = "#991b1b";
                }

                return (
                  <div
                    key={comp.id}
                    title={`${comp.id} (${comp.type})`}
                    style={{
                      width: 20,
                      height: 20,
                      backgroundColor: bg,
                      border: `1px solid ${borderColor}`,
                      borderRadius: 4,
                      cursor: "help",
                    }}
                  />
                );
              })}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}