import { test, expect } from '../fixtures';

test.describe('Matriz', () => {
  test('editar vuelve a la matriz', async ({ page, tasks }) => {
    await tasks([{ title: 'Preparar la demo', importance: 'important', dueInDays: 1 }]);
    await page.getByRole('link', { name: 'Ver la matriz de Eisenhower' }).click();

    // 1. Abre la edición desde la matriz
    await page.getByRole('link', { name: 'Editar tarea: Preparar la demo' }).click();
    await expect(page.getByRole('heading', { name: 'Editar tarea' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Volver a la matriz' })).toBeVisible();

    // 2. Cambia la importancia a «No importante» y guarda
    await page.getByLabel('Importancia').selectOption('No importante');
    await page.getByRole('button', { name: 'Guardar cambios' }).click();

    // 3. Vuelve a la matriz, con el aviso, y la tarea está en «Delegar»
    await expect(page).toHaveURL(/\/matrix$/);
    await expect(page.getByRole('status')).toHaveText('Tarea actualizada.');
    await expect(page.getByRole('region', { name: 'Delegar' }).locator('.task-title')).toHaveText(['Preparar la demo']);
    await expect(page.getByRole('region', { name: 'Hacer ya' }).locator('.task-title')).toHaveCount(0);
  });
});
