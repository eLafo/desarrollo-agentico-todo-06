# Matriz de Eisenhower

Todas las tareas, repartidas en los cuatro cuadrantes de la matriz de Eisenhower según su importancia y su urgencia, con las que faltan por clasificar aparte. Es una pantalla aparte, `/matrix`. La importancia y la fecha límite de cada tarea se eligen al añadirla o al editarla (ver [tareas.md](tareas.md)).

## Cómo se llega y cómo se vuelve

- Desde la lista, con el enlace «Ver la matriz de Eisenhower» de la cabecera, debajo del número de tareas.
- La pantalla tiene el enlace «← Volver a la lista», el título «Matriz de Eisenhower» y la explicación «Cada tarea está en el cuadrante que le corresponde por su importancia y su urgencia. Es urgente si su fecha límite es hoy, mañana o pasado mañana, o si ya ha pasado.».

## Las zonas

Los cuatro cuadrantes forman una cuadrícula de 2×2, con lo importante arriba y lo urgente a la izquierda. En pantallas estrechas (hasta 760 px) quedan en una columna, en el mismo orden. Debajo de la matriz está «Sin clasificar».

| Zona | Texto de ayuda | Qué tareas tiene | Si está vacía |
| --- | --- | --- | --- |
| Hacer ya | Urgentes e importantes | Importantes y urgentes | «No hay tareas en este cuadrante.» |
| Planificar | Importantes, no urgentes | Importantes y no urgentes | «No hay tareas en este cuadrante.» |
| Delegar | Urgentes, no importantes | No importantes y urgentes | «No hay tareas en este cuadrante.» |
| Descartar | Ni urgentes ni importantes | No importantes y no urgentes | «No hay tareas en este cuadrante.» |
| Sin clasificar | Les falta la fecha límite o decidir su importancia. Edítalas para colocarlas en la matriz. | Con la importancia «Sin decidir» o sin fecha límite | «No hay tareas sin clasificar.» |

- Las cinco zonas se ven siempre, también vacías.
- Una tarea está clasificada si tiene la importancia decidida («Importante» o «No importante») y una fecha límite. Si le falta cualquiera de las dos, está en «Sin clasificar».
- Cada tarea sale una sola vez entre las cinco zonas.
- Dentro de cada zona, las tareas salen en el orden de la lista: por orden de creación, de la más antigua a la más reciente.
- Las tareas completadas siguen en su zona, marcadas como en la lista: tachadas, con la etiqueta «Completada» y el botón «Reabrir tarea: …».
- La matriz no muestra la importancia ni la fecha límite de cada tarea: se ven y se cambian al editar.

## Urgencia

- Una tarea es urgente si su fecha límite es hoy, mañana o pasado mañana, o si ya ha pasado. El margen son 2 días (`models.UrgentDays`).
- Se cuenta en días de calendario respecto al día en que se abre la matriz, con «hoy» en la zona Europe/Madrid (`models.TimeZone`), sea cual sea la del servidor.
- Una tarea pasa a urgente sola, sin tocarla, cuando su fecha límite entra en el margen.
- La regla solo está en `models/matrix.go`: ni la plantilla ni el JavaScript calculan nada.

## Acciones desde la matriz

Cada tarea tiene los mismos controles que en la lista, que aparecen igual al pasar el ratón:

| Acción | Qué pasa | Aviso |
| --- | --- | --- |
| Completar o reabrir | Con JavaScript, la tarea cambia sin recargar la página y el cambio se guarda en segundo plano; sin JavaScript, vuelve a la matriz | Ninguno |
| Editar | Abre la edición con «← Volver a la matriz»; «Guardar cambios» y «Cancelar» vuelven a la matriz | «Tarea actualizada.» |
| Eliminar | Pide confirmación («¿Seguro que quieres eliminar esta tarea?») y vuelve a la matriz | «Tarea eliminada.» |

Los formularios y el enlace de editar de la matriz llevan `return_to=/matrix`. Las acciones solo aceptan ese valor para volver a la matriz; cualquier otro vuelve a la lista.

## Código

- `TasksMatrix` (`actions/tasks.go`), en `GET /matrix`: carga las tareas en el orden de la lista y las reparte con `models.NewMatrix(tasks, models.Today(time.Now()))`.
- `models/matrix.go`: `TimeZone`, `UrgentDays`, `Today`, que da el día de hoy en Europe/Madrid, y `Matrix` con `NewMatrix`, que reparte las tareas en `Do`, `Schedule`, `Delegate`, `Eliminate` y `Unclassified`. Recibe el día de hoy como parámetro.
- `returnPath` (`actions/tasks.go`) decide adónde vuelven editar, guardar, completar y eliminar.
- Plantillas: `templates/tasks/matrix.plush.html` (la pantalla), `_quadrant.plush.html` (una zona) y `_task.plush.html` (una tarea, compartida con la lista).
