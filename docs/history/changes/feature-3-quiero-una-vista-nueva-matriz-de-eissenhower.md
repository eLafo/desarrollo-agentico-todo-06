# Quiero una vista nueva: matriz de eissenhower

**Issue:** #3 · **Tipo:** feature

## Resumen
Una pantalla nueva, `/matrix`, reparte todas las tareas en los cuatro cuadrantes de la matriz de Eisenhower, con una zona «Sin clasificar» aparte. Para poder repartirlas, cada tarea gana una importancia y una fecha límite, que se eligen al añadirla y al editarla; la urgencia se calcula a partir de la fecha.

## Qué se construyó
- La matriz en `/matrix`, con un enlace desde la cabecera de la lista («Ver la matriz de Eisenhower») y otro de vuelta («← Volver a la lista»).
- Cinco zonas siempre visibles: «Hacer ya», «Planificar», «Delegar», «Descartar» y «Sin clasificar».
- Importancia (`Sin decidir`, `Importante`, `No importante`) y fecha límite opcional en los formularios de añadir y de editar.
- Las mismas acciones que en la lista (completar o reabrir, editar, eliminar), y todas devuelven a la matriz. Las tareas completadas aparecen en su zona, marcadas.

## Implementación técnica
| Fichero | Papel |
| --- | --- |
| `models/matrix.go` | `Matrix` y `NewMatrix`: reparte las tareas; `Today`, `TimeZone`, `UrgentDays`. Única fuente de la regla de urgencia. |
| `models/task.go` | Campos `Importance` y `DueDate`, `SetDueDate`, `DueDateValue` y su validación. |
| `migrations/20261020000000_add_importance_and_due_date_to_tasks.{up,down}.fizz` | Columnas `importance` (por defecto `undecided`) y `due_date` (nula). |
| `actions/tasks.go`, `actions/app.go` | `TasksMatrix` (`GET /matrix`), `returnPath`, y guardado de importancia y fecha al crear y editar. |
| `templates/tasks/matrix.plush.html`, `_quadrant`, `_task`, `_classification` | La pantalla, una zona, una tarea (compartida con la lista) y los campos del formulario. |

Cambios que importan:
- **Una tarea está clasificada** solo si tiene la importancia decidida y fecha límite; si le falta cualquiera, va a «Sin clasificar».
- **Urgente** es vencer hoy, mañana o pasado mañana, o haber vencido (`UrgentDays = 2`), en días de calendario. Se calcula en la zona Europe/Madrid (`models.Today`), con `time/tzdata` incluido en el binario para que funcione en la imagen de Alpine.
- **Volver a la matriz:** los formularios y el enlace de editar llevan `return_to=/matrix`; `returnPath` solo acepta ese valor y, si no, vuelve a la lista.
- **La tarjeta de tarea** se sacó de la lista a `_task.plush.html`, para usarla igual en las dos pantallas.
- **Tests:** fixture `tasks` de `e2e/fixtures.ts` con `importance` y `dueInDays`, y `dateInDays`.

## Cómo se usa
1. En la lista, al añadir una tarea (o al editarla), elige su importancia y, si quieres, su fecha límite.
2. Pulsa «Ver la matriz de Eisenhower» en la cabecera de la lista.
3. Cada tarea está en su zona; las que no tienen importancia decidida o fecha límite están en «Sin clasificar». Edítalas para colocarlas.
4. Completa, edita o elimina desde la matriz; se vuelve a ella.

Detalle vigente en [`docs/knowledge/features/matriz.md`](../../knowledge/features/matriz.md).

## Configuración
Una migración nueva (`add_importance_and_due_date_to_tasks`). Las tareas que ya existen quedan sin fecha, con la importancia «Sin decidir» y en «Sin clasificar», sin perder datos; la migración se puede deshacer y volver a aplicar. No hay variables de entorno nuevas.

## Cómo se probó
Acta de la verificación, validada en el intento 1 de 3: https://github.com/eLafo/desarrollo-agentico-todo-06/issues/3#issuecomment-6309894425

| Afirmación | Resultado |
| --- | --- |
| Criterio 1: se llega a la matriz desde la lista, tiene las cinco zonas y se vuelve | ✅ |
| Criterio 2: importancia y fecha al añadir y al editar, guardadas sin desfase de día | ✅ |
| Criterio 3: el reparto; pasado mañana es urgente y dentro de 3 días no | ✅ |
| Criterio 4: todas las tareas salen una vez, también las completadas | ✅ |
| Criterio 5: completar, editar y eliminar desde la matriz, también sin JavaScript | ✅ |
| Criterio 6: las tareas existentes quedan sin clasificar; la migración se deshace | ✅ |
| Sin regresiones: migraciones, 45 tests de acciones y 19 del modelo, 14 de 14 en navegador | ✅ |

Para repetirlo:

```bash
bash .claude/skills/check/scripts/check.sh
GO_ENV=test go test -p 1 ./...
npx playwright test e2e/matriz
```

## Notas
- Las decisiones del humano al aclarar la issue (cómo se clasifica, dónde van las tareas sin clasificar, qué se puede hacer desde la matriz, el margen de urgencia y cuándo sale una tarea de «Sin clasificar») están en [`docs/history/decisions/feature-3-quiero-una-vista-nueva-matriz-de-eissenhower.md`](../decisions/feature-3-quiero-una-vista-nueva-matriz-de-eissenhower.md).
- La lista no cambia salvo el enlace a la matriz; la matriz no muestra la importancia ni la fecha de cada tarea, que se ven y cambian al editar.
- Quedaron abiertas en el acta, sin bloquear: la ayuda de la pantalla de editar («Los demás datos de la tarea no cambian.») ya no es del todo cierta; «hoy» es el día en Madrid y no el del reloj del servidor; la consulta que carga las tareas está repetida en tres acciones; ningún test cubre volver a «Sin decidir» al editar; y la documentación de añadir dice que un error de fecha conserva la fecha, cuando el campo vuelve vacío (`docs/knowledge/features/tareas.md`, desfasada en ese punto).
