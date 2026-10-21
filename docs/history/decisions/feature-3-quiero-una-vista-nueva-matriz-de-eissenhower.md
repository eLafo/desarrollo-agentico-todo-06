# Decisiones de Quiero una vista nueva: matriz de eissenhower

**Issue:** #3 · **Documento de cierre:** [feature-3-quiero-una-vista-nueva-matriz-de-eissenhower.md](../changes/feature-3-quiero-una-vista-nueva-matriz-de-eissenhower.md)

## Primera ronda

### ¿Cómo se decide si una tarea es urgente y si es importante?

- **Contexto:** una tarea solo tenía título y si estaba completada; nada decía en qué cuadrante iba.
- **Opciones:** A, dos casillas «Urgente» e «Importante»; B, «Importante» como casilla y «Urgente» calculado a partir de una fecha límite nueva; C, mover las tareas entre cuadrantes desde la matriz.
- **Elección:** «1. B» · [comentario](https://github.com/eLafo/desarrollo-agentico-todo-06/issues/3#issuecomment-6300157552)
- **Consecuencia:** cada tarea tiene una importancia y una fecha límite, y la urgencia se calcula con la fecha; una tarea pasa a urgente sola con el tiempo.

### ¿Dónde van las tareas sin clasificar?

- **Contexto:** las tareas que ya existen y las que se añaden sin marcar nada no tienen clasificación.
- **Opciones:** A, en «ni urgente ni importante»; B, en una zona aparte, «Sin clasificar».
- **Elección:** «2. B» · [comentario](https://github.com/eLafo/desarrollo-agentico-todo-06/issues/3#issuecomment-6300157552)
- **Consecuencia:** la matriz tiene una quinta zona, «Sin clasificar», debajo de la cuadrícula.

### ¿Qué se puede hacer con una tarea desde la matriz?

- **Contexto:** había que fijar el alcance de la pantalla.
- **Opciones:** A, solo verla, con un enlace para editarla; B, lo mismo que en la lista.
- **Elección:** «3. B» · [comentario](https://github.com/eLafo/desarrollo-agentico-todo-06/issues/3#issuecomment-6300157552)
- **Consecuencia:** se puede completar, reabrir, editar y eliminar desde la matriz, y cada acción vuelve a ella.

## Segunda ronda

### ¿Cuántos días antes de su fecha límite pasa una tarea a ser urgente?

- **Contexto:** la elección de una fecha límite no fijaba el margen de la urgencia.
- **Opciones:** A, el mismo día; B, 2 días antes (hoy, mañana o pasado mañana); C, 7 días antes; o otro número.
- **Elección:** «1. B» · [comentario](https://github.com/eLafo/desarrollo-agentico-todo-06/issues/3#issuecomment-6300589210)
- **Consecuencia:** es urgente la tarea que vence hoy, mañana o pasado mañana, o que ya ha vencido (`models.UrgentDays = 2`).

### ¿Cuándo sale una tarea de «Sin clasificar»?

- **Contexto:** con una sola casilla, «sin marcar» no distingue «no importante» de «sin decidir».
- **Opciones:** A, en cuanto tiene fecha límite o está marcada como importante; B, cuando se decide su importancia, con tres valores («Sin decidir», «Importante», «No importante») y la fecha opcional.
- **Elección:** «2. Necesita tener importancia y fecha límite» · [comentario](https://github.com/eLafo/desarrollo-agentico-todo-06/issues/3#issuecomment-6300589210)
- **Consecuencia:** una tarea solo entra en un cuadrante si tiene la importancia decidida y una fecha límite; si le falta cualquiera de las dos, está en «Sin clasificar». Se eligió la importancia de tres valores, con «Sin decidir» por defecto, y la respuesta añade que la fecha es necesaria para clasificar.
