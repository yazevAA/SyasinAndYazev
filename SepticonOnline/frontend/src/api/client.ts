// src/api/client.ts

/** Базовый путь. Благодаря Vite-прокси уходит на :8080. */
const BASE_URL = "/api";

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

export async function request<T>(
  path: string,
  init?: RequestInit,
): Promise<T> {
  let response: Response;

  try {
    response = await fetch(BASE_URL + path, {
      headers: { "Content-Type": "application/json" },
      ...init,
    });
  } catch {
    throw new ApiError(0, "Сервер недоступен. Проверь, запущен ли бэкенд.");
  }

  if (!response.ok) {
    let message = `Ошибка ${response.status}`;
    try {
      const body = await response.json();
      if (body && typeof body.error === "string") {
        message = body.error;
      }
    } catch {
      /* тело не JSON — оставляем дефолт */
    }
    throw new ApiError(response.status, message);
  }

  // 204 и пустое тело не парсим
  if (response.status === 204) {
    return undefined as T;
  }

  return (await response.json()) as T;
}

/** Проверка доступности бэкенда: GET /health */
export async function checkHealth(): Promise<boolean> {
  try {
    const response = await fetch("/health");
    if (!response.ok) return false;
    const body = await response.json();
    return body?.status === "ok";
  } catch {
    return false;
  }
}