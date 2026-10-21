import { test, expect } from '../fixtures';

test.describe('Matriz', () => {
  test('cada tarea en su cuadrante', async ({ page, tasks }) => {
    await tasks([
      { title: 'Entregar el informe', importance: 'important', dueInDays: -1 },
      { title: 'Preparar la reunión', importance: 'important', dueInDays: 2 },
      { title: 'Planificar el trimestre', importance: 'important', dueInDays: 3 },
      { title: 'Renovar el dominio', importance: 'important', dueInDays: 0, completed: true },
      { title: 'Responder al proveedor', importance: 'not_important', dueInDays: 1 },
      { title: 'Ordenar el escritorio', importance: 'not_important', dueInDays: 10 },
      { title: 'Leer el manual' },
    ]);

    // 1. Abre la matriz
    await page.getByRole('link', { name: 'Ver la matriz de Eisenhower' }).click();

    // 2. Cada zona tiene exactamente sus tareas, en el orden de la lista
    await expect(page.getByRole('region', { name: 'Hacer ya' }).locator('.task-title')).toHaveText([
      'Entregar el informe',
      'Preparar la reunión',
      'Renovar el dominio',
    ]);
    await expect(page.getByRole('region', { name: 'Planificar' }).locator('.task-title')).toHaveText([
      'Planificar el trimestre',
    ]);
    await expect(page.getByRole('region', { name: 'Delegar' }).locator('.task-title')).toHaveText([
      'Responder al proveedor',
    ]);
    await expect(page.getByRole('region', { name: 'Descartar' }).locator('.task-title')).toHaveText([
      'Ordenar el escritorio',
    ]);
    await expect(page.getByRole('region', { name: 'Sin clasificar' }).locator('.task-title')).toHaveText([
      'Leer el manual',
    ]);

    // 3. La completada está en su cuadrante, marcada como en la lista
    await expect(
      page.getByRole('region', { name: 'Hacer ya' }).getByRole('button', { name: 'Reabrir tarea: Renovar el dominio' }),
    ).toBeVisible();
  });
});
