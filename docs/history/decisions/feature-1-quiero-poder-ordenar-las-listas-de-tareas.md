# Decisiones de Quiero poder ordenar las listas de tareas

**Issue:** #1 · **Documento de cierre:** [feature-1-quiero-poder-ordenar-las-listas-de-tareas.md](../changes/feature-1-quiero-poder-ordenar-las-listas-de-tareas.md)

## ¿Cómo reordena las tareas quien no usa el ratón?

- **Contexto:** la issue pedía «ordenar con drag and drop», que necesita ratón o pantalla táctil. El resto de la lista (completar, reabrir, editar y eliminar) funciona también con el teclado y sin JavaScript, y las convenciones del proyecto piden que el JavaScript mejore lo que ya funciona sin él. Añadir controles que la issue no pedía ampliaba el alcance, y no añadirlos dejaba a parte de los usuarios sin poder reordenar. Se publicó la pregunta en [este comentario](https://github.com/eLafo/desarrollo-agentico-todo-06/issues/1#issuecomment-6300134630), con los supuestos de que el orden se guarda, las tareas nuevas salen al final y las completadas se ordenan igual que las pendientes.
- **Opciones:**
  - **A · Solo arrastrar y soltar.** Exactamente lo pedido; con teclado o sin JavaScript no se puede reordenar.
  - **B · Arrastrar y soltar, y botones «Subir» y «Bajar» en cada tarea.** Funcionan con teclado y sin JavaScript. Cada tarea muestra dos botones de icono más, junto a editar y eliminar; la primera no puede subir y la última no puede bajar.
  - **C · Arrastrar y soltar, y mover la tarea con el teclado.** Sin botones nuevos: se enfoca el asa y se mueve con las flechas arriba y abajo. Sin JavaScript no se puede reordenar.
- **Elección:** «La B» · [comentario](https://github.com/eLafo/desarrollo-agentico-todo-06/issues/1#issuecomment-6300226479)
- **Consecuencia:** cada tarea tiene un asa de arrastre (con JavaScript, ratón y más de 760 px) y los botones «Subir» y «Bajar», que funcionan con el teclado y sin JavaScript; el de subir de la primera y el de bajar de la última salen desactivados. Los supuestos publicados con la pregunta, que el humano no discutió, quedaron como parte de lo pedido.
