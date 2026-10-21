# Cambios documentados

Cada caso que cierra el workflow `adw` deja aquí su documento de cierre: qué se construyó, cómo se usa, cómo se probó y qué conviene saber para trabajar sobre ello. Lo escribe la skill `document` en la rama del caso, y viaja en el PR con el código.

Este índice es condicional: cada entrada dice cuándo conviene cargar su documento. Se llega a él desde [`docs/README.md`](../../README.md): quien planifica, como las skills `feature`, `bug` y `chore`, lee solo los documentos cuyas condiciones encajan con el encargo.

Un documento de cierre cuenta qué pasó en su caso, y no se reescribe cuando otro caso cambia el mismo código. Lo vigente está en `docs/knowledge/` y `docs/procedures/`.

## Formato de las entradas

```markdown
- [<tipo>-<N>-<slug>.md](<tipo>-<N>-<slug>.md): <qué cambió, en una línea>
  - Cuando se trabaje con <rutas, acciones o plantillas concretas>
  - Cuando se modifique <modelo, campo o regla relacionados>
  - Cuando se investigue <un problema concreto de esa área>
```

Entre tres y cinco condiciones por documento, y las entradas en orden por número de issue. Si un documento ya tiene entrada, se actualizan sus condiciones en lugar de añadir otra.

## Índice

- [feature-1-quiero-poder-ordenar-las-listas-de-tareas.md](feature-1-quiero-poder-ordenar-las-listas-de-tareas.md): la lista de tareas se reordena arrastrando o con botones «Subir» y «Bajar», y el orden se guarda en `tasks.position`
  - Cuando se trabaje con `PUT /tasks/{task_id}/move`, `TasksMove` o `Task.Move` en `models/task.go`
  - Cuando se modifique el orden de la lista, la columna `position` de `tasks` o cómo se crea una tarea nueva
  - Cuando se toque `templates/tasks/index.plush.html`, el asa de arrastre o los botones de cada tarea en `public/assets/app.js` y `app.css`
  - Cuando se escriban migraciones o fixtures de e2e que inserten tareas (`e2e/fixtures.ts`)
  - Cuando se investigue que una tarea cambie de sitio al completarla, o que el orden no se conserve al recargar
