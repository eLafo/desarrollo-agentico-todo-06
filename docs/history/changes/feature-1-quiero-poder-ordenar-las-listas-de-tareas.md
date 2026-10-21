# Quiero poder ordenar las listas de tareas

**Issue:** #1 · **Tipo:** feature

## Resumen

La lista de tareas se puede reordenar a mano: arrastrando cada tarea por un asa y, para quien no usa el ratón, con botones «Subir» y «Bajar» que funcionan con el teclado y sin JavaScript. El orden elegido se guarda en la base de datos y sobrevive a recargar la página.

## Qué se construyó

- Arrastrar y soltar una tarea en otro puesto de la lista (con JavaScript, ratón y pantallas de más de 760 px).
- Dos botones de icono por tarea, «Subir» y «Bajar», junto a editar y eliminar. El de subir de la primera tarea y el de bajar de la última salen desactivados.
- Un orden persistente: las tareas nuevas salen al final y completar, reabrir o editar una tarea no la mueve.
- Una ruta nueva, `PUT /tasks/{task_id}/move`, que usan tanto los botones como el arrastre.

## Implementación técnica

| Fichero | Papel |
| --- | --- |
| `migrations/20261020000000_add_position_to_tasks.{up,down}.fizz` | Añade `tasks.position` (entero, `0` por defecto) |
| `migrations/20261020000001_number_task_positions.{up,down}.sql` | Numera las tareas existentes por orden de creación; la inversa las deja en `0` |
| `models/task.go` | Campo `Position`, `AllTasks`, `BeforeCreate`, `Task.Move`, `Tasks.CanMoveUp` y `Tasks.CanMoveDown`, y el error `ErrMoveOutOfRange` |
| `actions/tasks.go`, `actions/app.go` | Acción `TasksMove` y su ruta; el índice y el alta con error leen la lista con `models.AllTasks` |
| `templates/tasks/index.plush.html` | El asa y los dos formularios de mover por tarea |
| `public/assets/app.js`, `app.css` | El arrastre, el asa, el estilo de la tarea arrastrada y de los botones desactivados |
| `templates/application.plush.html` | Sube la versión de los dos ficheros de `public/assets/` para saltarse la caché |
| `e2e/fixtures.ts` | La fixture `tasks` deja también el puesto de cada tarea |

Cambios técnicos que importan:

1. **El orden vive en `position`**, de 1 en adelante. La lista se lee por `position`, luego `created_at` y luego `id`, así que los empates se resuelven de forma estable.
2. **Una tarea nueva toma el puesto siguiente al mayor** (`BeforeCreate`), por lo que sale al final. Completar, reabrir o editar no tocan `position`.
3. **`Task.Move(tx, offset)` renumera la lista**: saca la tarea de su puesto, la inserta `offset` puestos más allá y solo escribe las filas cuyo puesto cambia, con `UpdateColumns`. Un destino fuera de la lista devuelve `ErrMoveOutOfRange`.
4. **Una sola ruta para todo.** Los botones envían `offset` de `-1` o `1`; el arrastre envía los puestos que se ha desplazado la tarea. `TasksMove` responde `303` a `/`, `404` si la tarea no existe y `422` si `offset` no es un número o saca la tarea de la lista.
5. **El JavaScript solo mejora.** El servidor decide si el movimiento es válido: tras soltar, `app.js` repinta la lista con la respuesta del servidor y, si este rechaza el movimiento, recarga la página. Con Escape o soltando fuera de la lista, la tarea vuelve a su sitio y no se envía nada.

## Cómo se usa

- **Con el ratón (pantallas de más de 760 px):** se agarra el asa de seis puntos a la izquierda de la tarea, que aparece al pasar por encima o al enfocarla, y se suelta en el puesto nuevo.
- **Con el teclado o sin JavaScript:** con el tabulador se llega a «Subir tarea: …» o «Bajar tarea: …» y se pulsa Intro o espacio. Cada pulsación mueve la tarea un puesto y recarga la lista; el foco vuelve al principio de la página.
- **En pantallas estrechas** no hay asa y se ordena con los botones.
- Mover no muestra ningún aviso.

## Configuración

Dos migraciones nuevas, que `cmd/app` aplica al arrancar y que en local y en los tests se aplican con `soda migrate` (`-e test` para la base de test):

- `20261020000000_add_position_to_tasks` añade la columna `position`.
- `20261020000001_number_task_positions` da a las tareas que ya existen el puesto que les corresponde por su orden de creación, de modo que el orden anterior no cambia.

Ninguna variable de entorno nueva.

## Cómo se probó

La verificación del cambio quedó **validada en el intento 2 de 3**. El acta, con los vídeos, capturas y comprobaciones, está en https://github.com/eLafo/desarrollo-agentico-todo-06/issues/1#issuecomment-6310568658.

| Criterio de la issue | Resultado |
| --- | --- |
| 1. Arrastrar y soltar una tarea en otro puesto | Cumplido. Arrastrar la última al principio y la primera al final deja el orden y lo mantiene al recargar. Soltar en el mismo sitio no cambia nada. El servidor rechaza con 404 una tarea inexistente y con 422 un desplazamiento imposible. |
| 2. «Subir» y «Bajar» con el teclado y sin JavaScript | Cumplido. Con tabulador e Intro, y con el JavaScript bloqueado, cada botón mueve la tarea un puesto. Subir la primera y bajar la última salen desactivados y el servidor tampoco los admite (422). Sin token o con `GET`, no se mueve nada. |
| 3. El orden se guarda | Cumplido. El mismo orden en un navegador nuevo, sin sesión, en el HTML del servidor y en la base. |
| 4. Las tareas nuevas salen al final | Cumplido, también en una lista ya reordenada. |
| 5. Completar o reabrir no mueve la tarea | Cumplido. Hay un test que reordena la lista antes y comprueba que la tarea conserva su puesto. En el intento 1 faltaba esta prueba y la verificación se rechazó. |
| Sin regresiones | La suite completa en verde: 107 tests de Go, ninguno saltado, y 7 tests de navegador, 2 de ellos nuevos. |

Los tests nuevos de Go y de navegador fallan sobre la base del cambio y pasan con él. Mover una tarea solo cambia su puesto: título, estado y fecha de creación quedan igual.

Pruebas nuevas: `actions/tasks_test.go`, `models/task_test.go`, `e2e/tareas/ordenar-arrastrando.spec.ts` y `e2e/tareas/ordenar-con-subir-y-bajar.spec.ts`.

Para repetirlas (con el entorno del worktree cargado y la base de test migrada, ver `docs/procedures/development.md`):

```bash
GO_ENV=test go test -p 1 ./...
npx playwright test
```

Quedaron preguntas para el humano que no bloquean, y que la verificación dejó planteadas en la propuesta: si se ve claro por dónde se agarra una tarea y dónde caerá, si los iconos se entienden en el ordenador y en el móvil, y si es cómodo mover una tarea varios puestos seguidos solo con el teclado.

## Notas

- **Decisión del humano:** ante la pregunta de cómo reordena quien no usa el ratón, eligió la opción B (arrastrar y botones «Subir» y «Bajar»). Está en [docs/history/decisions/feature-1-quiero-poder-ordenar-las-listas-de-tareas.md](../decisions/feature-1-quiero-poder-ordenar-las-listas-de-tareas.md).
- **Sin JavaScript, el foco vuelve al principio** tras cada pulsación, así que mover la misma tarea varios puestos obliga a tabular de nuevo hasta ella. No se ha corregido.
- **El asa solo existe con JavaScript, ratón y más de 760 px.** Con pantalla táctil estrecha se ordena con los botones.
- **Sin cubrir con test de navegador:** soltar una tarea fuera de la lista (vuelve a su sitio) y que el servidor rechace el movimiento (la página se recarga). El arrastre principal sí tiene test.
- **Pequeños retoques que la verificación sugirió y no se hicieron:** el nombre de `Test_Tasks_CanMoveUpAndDown` en `models/task_test.go` no dice qué caso prueba, y el mensaje de `ErrMoveOutOfRange` empieza en minúscula y no acaba en punto (solo se ve en la página de error de desarrollo).
- **Fuera de alcance:** ordenar por un criterio (fecha, nombre), separar las completadas y mover varias tareas a la vez.
- La documentación viva (`docs/knowledge/`) ya describe el orden, la ruta y la columna; no se ha visto nada desfasado.
