import { test as base, expect } from '@playwright/test';
import pg from 'pg';

export { expect };

type Task = {
  title: string;
  completed?: boolean;
  importance?: 'undecided' | 'important' | 'not_important';
  // Días desde hoy hasta la fecha límite; negativos si ya ha vencido. Sin él, la tarea no tiene fecha.
  dueInDays?: number;
};

type Fixtures = {
  // Deja la tabla con exactamente estas tareas, en este orden de creación.
  tasks: (tasks: Task[]) => Promise<void>;
};

// La fecha a `days` días de hoy, como AAAA-MM-DD (el formato del campo de fecha). Cuenta «hoy» en
// Europe/Madrid, como la aplicación, sea cual sea la zona del equipo.
export function dateInDays(days: number): string {
  const today = new Intl.DateTimeFormat('en-CA', { timeZone: 'Europe/Madrid' }).format(new Date());
  const date = new Date(`${today}T00:00:00Z`);
  date.setUTCDate(date.getUTCDate() + days);
  return date.toISOString().slice(0, 10);
}

// La misma base de datos que usa la aplicación con GO_ENV=test (database.yml).
function client() {
  return new pg.Client({
    host: process.env.DB_HOST ?? '127.0.0.1',
    port: Number(process.env.DB_PORT ?? 5432),
    user: 'postgres',
    password: 'postgres',
    database: 'todo_test',
  });
}

// Cada test empieza con la lista vacía y la página abierta. Si necesita tareas, las crea con
// `tasks` y recarga la página.
export const test = base.extend<Fixtures>({
  tasks: async ({}, use) => {
    const db = client();
    await db.connect();
    await db.query('DELETE FROM tasks');
    await use(async (tasks) => {
      await db.query('DELETE FROM tasks');
      for (const [i, task] of tasks.entries()) {
        await db.query(
          `INSERT INTO tasks (id, title, completed, importance, due_date, created_at, updated_at)
           VALUES (gen_random_uuid(), $1, $2, $3, $4, now() + $5 * interval '1 millisecond', now())`,
          [
            task.title,
            task.completed ?? false,
            task.importance ?? 'undecided',
            task.dueInDays === undefined ? null : dateInDays(task.dueInDays),
            i,
          ],
        );
      }
    });
    await db.end();
  },
  page: async ({ page, tasks }, use) => {
    void tasks;
    await page.goto('/');
    await use(page);
  },
});
