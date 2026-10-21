import { test, expect } from '../fixtures';

test.describe('Tareas', () => {
  test('ordenar arrastrando una tarea', async ({ page, tasks }) => {
    await tasks([{ title: 'Comprar pan' }, { title: 'Llamar al banco' }, { title: 'Pagar la luz', completed: true }]);
    await page.reload();
    const tarea = (title: string) => page.getByRole('listitem').filter({ hasText: title });

    // 1. Arrastra «Pagar la luz» por su asa hasta la mitad superior de «Comprar pan»
    await tarea('Pagar la luz').locator('.task-handle').dragTo(tarea('Comprar pan'), { targetPosition: { x: 200, y: 5 } });
    await expect(page.locator('.task-title')).toHaveText(['Pagar la luz', 'Comprar pan', 'Llamar al banco']);
    await expect(page.getByRole('button', { name: 'Subir tarea: Pagar la luz' })).toBeDisabled();

    // 2. Arrastra «Comprar pan» hasta la mitad inferior de «Llamar al banco»
    const destino = await tarea('Llamar al banco').boundingBox();
    await tarea('Comprar pan').locator('.task-handle').dragTo(tarea('Llamar al banco'), {
      targetPosition: { x: 200, y: destino!.height - 5 },
    });
    await expect(page.locator('.task-title')).toHaveText(['Pagar la luz', 'Llamar al banco', 'Comprar pan']);
    await expect(page.getByRole('button', { name: 'Bajar tarea: Comprar pan' })).toBeDisabled();

    // 3. Recarga: el orden se ha guardado
    await page.reload();
    await expect(page.locator('.task-title')).toHaveText(['Pagar la luz', 'Llamar al banco', 'Comprar pan']);
  });
});
