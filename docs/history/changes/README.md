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

- [feature-3-quiero-una-vista-nueva-matriz-de-eissenhower.md](feature-3-quiero-una-vista-nueva-matriz-de-eissenhower.md): matriz de Eisenhower en `/matrix`, con importancia y fecha límite en las tareas
  - Cuando se trabaje con `GET /matrix`, `TasksMatrix` o las plantillas `templates/tasks/matrix.plush.html`, `_quadrant`, `_task` o `_classification`
  - Cuando se modifique la importancia o la fecha límite de `models.Task`, o la regla de urgencia de `models/matrix.go`
  - Cuando se cambien las acciones de completar, editar o eliminar y adónde vuelven (`return_to`, `returnPath`)
  - Cuando se investigue una tarea en el cuadrante o en «Sin clasificar» equivocado, o un desfase de un día por la zona horaria
