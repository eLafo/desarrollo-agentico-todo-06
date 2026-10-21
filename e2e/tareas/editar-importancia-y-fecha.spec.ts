import { test, expect, dateInDays } from '../fixtures';

test.describe('Tareas', () => {
  test('editar importancia y fecha', async ({ page, tasks }) => {
    await tasks([{ title: 'Preparar la demo', importance: 'important', dueInDays: 1 }]);
    await page.reload();

    // 1. Abre la edición desde la lista: los campos vienen con sus valores
    await page.getByRole('link', { name: 'Editar tarea: Preparar la demo' }).click();
    await expect(page.getByRole('link', { name: 'Volver a la lista' })).toBeVisible();
    await expect(page.getByLabel('Importancia')).toHaveValue('important');
    await expect(page.getByRole('textbox', { name: 'Fecha límite' })).toHaveValue(dateInDays(1));

    // 2. Cambia la fecha límite a dentro de diez días y guarda
    await page.getByRole('textbox', { name: 'Fecha límite' }).fill(dateInDays(10));
    await page.getByRole('button', { name: 'Guardar cambios' }).click();

    // 3. Vuelve a la lista, con el aviso
    await expect(page.getByRole('status')).toHaveText('Tarea actualizada.');
    await expect(page.getByRole('heading', { name: 'Tareas', level: 1 })).toBeVisible();

    // 4. En la matriz, está en «Planificar»
    await page.getByRole('link', { name: 'Ver la matriz de Eisenhower' }).click();
    await expect(page.getByRole('region', { name: 'Planificar' }).locator('.task-title')).toHaveText(['Preparar la demo']);

    // 5. Edita otra vez, vacía la fecha y guarda: queda en «Sin clasificar»
    await page.getByRole('link', { name: 'Editar tarea: Preparar la demo' }).click();
    await page.getByRole('textbox', { name: 'Fecha límite' }).fill('');
    await page.getByRole('button', { name: 'Guardar cambios' }).click();
    await expect(page.getByRole('status')).toHaveText('Tarea actualizada.');
    await expect(page.getByRole('region', { name: 'Sin clasificar' }).locator('.task-title')).toHaveText([
      'Preparar la demo',
    ]);
    await expect(page.getByRole('region', { name: 'Planificar' }).locator('.task-title')).toHaveCount(0);
  });
});
