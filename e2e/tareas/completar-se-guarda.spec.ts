import { test, expect } from '../fixtures';

test.describe('Tareas', () => {
  test('completar se guarda', async ({ page, tasks }) => {
    await tasks([{ title: 'Llamar al proveedor' }, { title: 'Revisar el presupuesto' }]);
    await page.reload();

    // 1. Pulsa «Completar tarea: Llamar al proveedor» y espera a que se envíe
    const toggled = page.waitForResponse((response) => response.url().includes('/toggle'));
    await page.getByRole('button', { name: 'Completar tarea: Llamar al proveedor' }).click();
    await toggled;
    await expect(page.getByRole('button', { name: 'Reabrir tarea: Llamar al proveedor' })).toBeVisible();

    // 2. Al recargar, sigue completada
    await page.reload();
    await expect(page.getByRole('button', { name: 'Reabrir tarea: Llamar al proveedor' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Completar tarea: Revisar el presupuesto' })).toBeVisible();
  });
});
