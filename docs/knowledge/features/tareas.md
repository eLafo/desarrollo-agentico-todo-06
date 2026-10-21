# Lista de tareas

Lo que ve y puede hacer el usuario con sus tareas, en la lista (`/`) y en la edición de una tarea (`/tasks/{task_id}/edit`). Las mismas tareas se ven también en la matriz de Eisenhower (`/matrix`), descrita en [matriz.md](matriz.md). El código está en `actions/tasks.go`, `models/task.go`, `templates/tasks/` y `public/assets/app.js`.

## La lista

- La pestaña y el enlace de la navegación lateral dicen «Tareas». El título de la página también.
- Debajo del título, el número de tareas: «3 tareas en tu lista», y en singular «1 tarea en tu lista». Debajo, el enlace «Ver la matriz de Eisenhower», que lleva a la matriz.
- Las tareas salen por orden de creación, de la más antigua a la más reciente.
- Cada tarea muestra:
  - un botón redondo para completarla o reabrirla;
  - su título;
  - la etiqueta «Completada», si lo está;
  - un enlace para editarla y un botón para eliminarla.
- La lista no muestra la importancia ni la fecha límite de las tareas: se ven y se cambian al editar.
- Los botones solo tienen icono, y su nombre accesible incluye el título: «Completar tarea: Comprar pan», «Reabrir tarea: …», «Editar tarea: …», «Eliminar tarea: …».
- Sin tareas, la lista muestra «No hay tareas todavía. Añade la primera con el formulario.», y el número dice «0 tareas en tu lista».

## Reglas del título

Las comprueba `Task.Validate` al crear y al editar:

| Regla | Mensaje si no se cumple |
| --- | --- |
| Es obligatorio. Un título solo con espacios cuenta como vacío | «El título es obligatorio.» |
| Tiene como mucho 200 caracteres (`models.TitleMaxLength`), contados sin los espacios del principio y del final | «El título no puede superar los 200 caracteres.» |

Antes de guardar, la acción quita los espacios del principio y del final. Los campos de título del navegador tienen también `maxlength="200"`.

## Importancia y fecha límite

- La importancia se elige en el campo «Importancia» entre «Sin decidir», «Importante» y «No importante». Por defecto es «Sin decidir» (`undecided`).
- La fecha límite es un día, sin hora, y es opcional. Se escribe en el campo «Fecha límite» con el selector de fecha del navegador. Se acepta una fecha ya pasada: la tarea cuenta como vencida.
- Una tarea necesita las dos, la importancia decidida y la fecha límite, para entrar en uno de los cuadrantes de la matriz; si le falta alguna, está en «Sin clasificar» (ver [matriz.md](matriz.md)).

Las comprueba `Task.Validate` al crear y al editar:

| Regla | Mensaje si no se cumple |
| --- | --- |
| La importancia es una de las tres opciones | «La importancia no es válida.» |
| La fecha límite, si la hay, es un día del calendario (AAAA-MM-DD) | «La fecha límite no es válida.» |

## Añadir

- El formulario está debajo de la lista: el campo «Nueva tarea», con el texto de ayuda «Añadir una tarea…», y el botón «Añadir». Debajo, los campos «Importancia», con «Sin decidir» elegido, y «Fecha límite», vacío.
- Si todo es válido, la tarea se crea pendiente, con la importancia y la fecha elegidas, aparece al final de la lista y sale el aviso «Tarea creada.».
- Si no, la lista se vuelve a mostrar con `422`: el formulario conserva el título, la importancia y la fecha, solo el campo que falla queda marcado como inválido y su mensaje sale debajo de él.

## Editar

- Se llega desde el enlace de editar de cada tarea. La pantalla tiene el enlace «← Volver a la lista», el título «Editar tarea» y el campo «Título» con el título actual. Debajo, la ayuda: «Hasta 200 caracteres. Los demás datos de la tarea no cambian.», y los campos «Importancia» y «Fecha límite», rellenos con los valores de la tarea.
- «Guardar cambios» cambia el título, la importancia y la fecha límite; dejar la fecha vacía la quita. Si la tarea estaba completada, lo sigue estando. Después vuelve a la lista con el aviso «Tarea actualizada.».
- Con un valor inválido, la pantalla se vuelve a mostrar con `422` y el mensaje bajo su campo, y la tarea no cambia. «Cancelar» vuelve a la lista sin cambiar nada.
- Desde la matriz, el enlace dice «← Volver a la matriz», y «Cancelar» y «Guardar cambios» vuelven a la matriz.
- Si la tarea no existe, `404`.

## Completar y reabrir

- El botón de cada tarea la alterna entre pendiente y completada. Una tarea completada lleva la etiqueta «Completada», y su botón pasa a reabrirla.
- Sin JavaScript, el formulario envía `PUT /tasks/{task_id}/toggle`, la acción cambia el estado y vuelve a la pantalla desde la que se completó: la lista o la matriz. No hay aviso.
- Con JavaScript, `app.js` cambia la tarea en la página sin recargarla (la etiqueta, el estado del botón y su nombre accesible) y envía el formulario en segundo plano (`POST` con sus campos, `_method=PUT` y el token incluidos), así que el cambio se guarda.

## Eliminar

- El botón de eliminar pide confirmación: «¿Seguro que quieres eliminar esta tarea?». Si se cancela, no pasa nada. Sin JavaScript, elimina sin preguntar.
- Al eliminar, la tarea desaparece, se vuelve a la pantalla desde la que se eliminó (la lista o la matriz) y sale el aviso «Tarea eliminada.».
- Si la tarea no existe, `404`.

## Avisos

Los avisos salen arriba, en la primera pantalla que se muestra después de la acción, y se anuncian a los lectores de pantalla (`role="status"`):

| Acción | Aviso |
| --- | --- |
| Añadir | «Tarea creada.» |
| Editar | «Tarea actualizada.» |
| Eliminar | «Tarea eliminada.» |
| Completar o reabrir | Ninguno |
