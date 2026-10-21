import { test, expect } from '../fixtures';

declare global {
  interface Window {
    sinRecargar?: boolean;
  }
}

test.describe('Matriz', () => {
  test('completar y reabrir sin recargar', async ({ page, tasks }) => {
    await tasks([{ title: 'Preparar la demo', importance: 'important', dueInDays: 1 }]);
    await page.getByRole('link', { name: 'Ver la matriz de Eisenhower' }).click();
    const hacerYa = page.getByRole('region', { name: 'Hacer ya' });

    // 1. Marca la página, para saber si se recarga
    await page.evaluate(() => {
      window.sinRecargar = true;
    });

    // 2. Completa la tarea y espera a que se envíe
    let toggled = page.waitForResponse((response) => response.url().includes('/toggle'));
    await hacerYa.getByRole('button', { name: 'Completar tarea: Preparar la demo' }).click();
    await toggled;

    // 3. La tarea cambia sin recargar la página
    await expect(hacerYa.getByRole('button', { name: 'Reabrir tarea: Preparar la demo' })).toBeVisible();
    await expect(page).toHaveURL(/\/matrix$/);
    expect(await page.evaluate(() => window.sinRecargar)).toBe(true);

    // 4. Al recargar, sigue completada en el mismo cuadrante
    await page.reload();
    await expect(hacerYa.getByRole('button', { name: 'Reabrir tarea: Preparar la demo' })).toBeVisible();

    // 5. Reabre la tarea, también sin recargar
    await page.evaluate(() => {
      window.sinRecargar = true;
    });
    toggled = page.waitForResponse((response) => response.url().includes('/toggle'));
    await hacerYa.getByRole('button', { name: 'Reabrir tarea: Preparar la demo' }).click();
    await toggled;
    await expect(hacerYa.getByRole('button', { name: 'Completar tarea: Preparar la demo' })).toBeVisible();
    await expect(page).toHaveURL(/\/matrix$/);
    expect(await page.evaluate(() => window.sinRecargar)).toBe(true);

    // 6. Al recargar, sigue pendiente en el mismo cuadrante
    await page.reload();
    await expect(hacerYa.getByRole('button', { name: 'Completar tarea: Preparar la demo' })).toBeVisible();
  });
});
