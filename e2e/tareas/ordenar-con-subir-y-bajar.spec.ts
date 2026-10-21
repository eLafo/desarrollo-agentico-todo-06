import { test, expect } from '../fixtures';

test.use({ javaScriptEnabled: false });

test.describe('Tareas', () => {
  test('ordenar con Subir y Bajar, con el teclado y sin JavaScript', async ({ page, tasks }) => {
    await tasks([{ title: 'Comprar pan' }, { title: 'Llamar al banco', completed: true }, { title: 'Pagar la luz' }]);
    await page.reload();

    // 1. Sin JavaScript no hay asa; la primera no puede subir y la última no puede bajar
    await expect(page.locator('.task-handle').first()).toBeHidden();
    await expect(page.getByRole('button', { name: 'Subir tarea: Comprar pan' })).toBeDisabled();
    await expect(page.getByRole('button', { name: 'Bajar tarea: Pagar la luz' })).toBeDisabled();
    await expect(page.getByRole('button', { name: 'Bajar tarea: Comprar pan' })).toBeEnabled();

    // 2. Pulsa Intro en «Bajar» de «Comprar pan»
    await page.getByRole('button', { name: 'Bajar tarea: Comprar pan' }).press('Enter');
    await expect(page.locator('.task-title')).toHaveText(['Llamar al banco', 'Comprar pan', 'Pagar la luz']);

    // 3. Pulsa Intro en «Subir» de «Pagar la luz»
    await page.getByRole('button', { name: 'Subir tarea: Pagar la luz' }).press('Enter');
    await expect(page.locator('.task-title')).toHaveText(['Llamar al banco', 'Pagar la luz', 'Comprar pan']);
    await expect(page.getByRole('button', { name: 'Subir tarea: Llamar al banco' })).toBeDisabled();
    await expect(page.getByRole('button', { name: 'Bajar tarea: Comprar pan' })).toBeDisabled();

    // 4. Recarga: el orden se ha guardado
    await page.reload();
    await expect(page.locator('.task-title')).toHaveText(['Llamar al banco', 'Pagar la luz', 'Comprar pan']);
  });
});
