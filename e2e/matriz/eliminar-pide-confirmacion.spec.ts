import { test, expect } from '../fixtures';

test.describe('Matriz', () => {
  test('eliminar pide confirmación', async ({ page, tasks }) => {
    await tasks([
      { title: 'Enviar la factura', importance: 'not_important', dueInDays: 5 },
      { title: 'Revisar el presupuesto' },
    ]);
    await page.getByRole('link', { name: 'Ver la matriz de Eisenhower' }).click();
    const eliminar = page.getByRole('button', { name: 'Eliminar tarea: Enviar la factura' });

    // 1. Pulsa eliminar y cancela: la tarea sigue
    page.once('dialog', (dialog) => dialog.dismiss());
    await eliminar.click();
    await expect(page.getByRole('region', { name: 'Descartar' }).locator('.task-title')).toHaveText(['Enviar la factura']);

    // 2. Pulsa eliminar y acepta: la tarea desaparece y la página sigue en la matriz
    page.once('dialog', async (dialog) => {
      expect(dialog.message()).toBe('¿Seguro que quieres eliminar esta tarea?');
      await dialog.accept();
    });
    await eliminar.click();
    await expect(page.getByRole('status')).toHaveText('Tarea eliminada.');
    await expect(page).toHaveURL(/\/matrix$/);
    await expect(page.locator('.task-title')).toHaveText(['Revisar el presupuesto']);
  });
});
