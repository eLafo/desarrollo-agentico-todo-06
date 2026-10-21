import { test, expect } from '../fixtures';

test.describe('Matriz', () => {
  test('ir a la matriz y volver', async ({ page, tasks }) => {
    await tasks([
      { title: 'Llamar al proveedor' },
      { title: 'Preparar la demo', importance: 'important', dueInDays: 1 },
    ]);
    await page.reload();

    // 1. Desde la lista, pulsa «Ver la matriz de Eisenhower»
    await page.getByRole('link', { name: 'Ver la matriz de Eisenhower' }).click();

    // 2. La matriz tiene su encabezado y sus cinco zonas
    await expect(page).toHaveURL(/\/matrix$/);
    await expect(page.getByRole('heading', { name: 'Matriz de Eisenhower', level: 1 })).toBeVisible();
    for (const name of ['Hacer ya', 'Planificar', 'Delegar', 'Descartar', 'Sin clasificar']) {
      await expect(page.getByRole('region', { name })).toBeVisible();
    }

    // 3. «Volver a la lista» vuelve a la lista, con las tareas como estaban
    await page.getByRole('link', { name: 'Volver a la lista' }).click();
    await expect(page.getByRole('heading', { name: 'Tareas', level: 1 })).toBeVisible();
    await expect(page.locator('.task-title')).toHaveText(['Llamar al proveedor', 'Preparar la demo']);
  });
});
