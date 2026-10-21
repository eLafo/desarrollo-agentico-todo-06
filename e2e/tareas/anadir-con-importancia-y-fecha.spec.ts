import { test, expect, dateInDays } from '../fixtures';

test.describe('Tareas', () => {
  test('añadir con importancia y fecha', async ({ page }) => {
    // 1. Escribe el título, elige «Importante» y pon la fecha límite de mañana
    await page.getByRole('textbox', { name: 'Nueva tarea' }).fill('Preparar la demo');
    await page.getByLabel('Importancia').selectOption('Importante');
    await page.getByRole('textbox', { name: 'Fecha límite' }).fill(dateInDays(1));

    // 2. Pulsa «Añadir»: la tarea se crea
    await page.getByRole('button', { name: 'Añadir' }).click();
    await expect(page.getByRole('status')).toHaveText('Tarea creada.');
    await expect(page.locator('.task-title')).toHaveText(['Preparar la demo']);

    // 3. En la matriz, está en «Hacer ya»
    await page.getByRole('link', { name: 'Ver la matriz de Eisenhower' }).click();
    await expect(page.getByRole('region', { name: 'Hacer ya' }).locator('.task-title')).toHaveText(['Preparar la demo']);
  });
});
