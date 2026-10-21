# Decisiones

Lo que decidió el humano al responder una pregunta con opciones de un caso del workflow `adw`, normalmente al aclarar la issue. Lo escribe la skill `document` al cerrar el caso: un fichero por caso, `<tipo>-<N>-<slug>.md`, con el mismo nombre que su documento de cierre en [`../changes/`](../changes/README.md).

Las elecciones técnicas del plan no van aquí: las recoge el documento de cierre.

## Formato de un fichero

````markdown
# Decisiones de <título del caso>

**Issue:** #<N> · **Documento de cierre:** [<tipo>-<N>-<slug>.md](../changes/<tipo>-<N>-<slug>.md)

## <La pregunta, en una línea>

- **Contexto:** <por qué hubo que preguntar>
- **Opciones:** <las que se ofrecieron>
- **Elección:** <la respuesta del humano, con sus palabras> · [comentario](<enlace>)
- **Consecuencia:** <qué hace hoy el sistema por esta decisión>
````

Una sección por pregunta.

## Formato de las entradas

```markdown
- [<tipo>-<N>-<slug>.md](<tipo>-<N>-<slug>.md): <lo decidido, en una línea>
  - Cuando <situación concreta en la que importa>
```

Entre una y tres condiciones por fichero, y las entradas en orden por número de issue.

## Índice

- [feature-3-quiero-una-vista-nueva-matriz-de-eissenhower.md](feature-3-quiero-una-vista-nueva-matriz-de-eissenhower.md): la importancia y la fecha límite salen de la tarea; urgente a 2 días; sin clasificar si falta cualquiera de las dos
  - Cuando se cambie cómo se decide la urgencia o el margen de `models.UrgentDays`
  - Cuando se decida qué tareas van a «Sin clasificar» o qué se puede hacer desde la matriz
