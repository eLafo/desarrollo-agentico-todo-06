import { test, expect } from '../fixtures';

test.describe('Matriz', () => {
  test('sin clasificar', async ({ page, tasks }) => {
    await tasks([
      { title: 'Llamar al gestor' },
      { title: 'Revisar el contrato', dueInDays: 1 },
      { title: 'Leer el manual', importance: 'important' },
    ]);

    // 1. Abre la matriz
    await page.getByRole('link', { name: 'Ver la matriz de Eisenhower' }).click();

    // 2. Las tres están en «Sin clasificar»
    await expect(page.getByRole('region', { name: 'Sin clasificar' }).locator('.task-title')).toHaveText([
      'Llamar al gestor',
      'Revisar el contrato',
      'Leer el manual',
    ]);

    // 3. Y en ningún cuadrante
    for (const name of ['Hacer ya', 'Planificar', 'Delegar', 'Descartar']) {
      const quadrant = page.getByRole('region', { name });
      await expect(quadrant.locator('.task-title')).toHaveCount(0);
      await expect(quadrant.getByText('No hay tareas en este cuadrante.')).toBeVisible();
    }
  });
});
