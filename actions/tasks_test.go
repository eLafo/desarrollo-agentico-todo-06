package actions

import (
	"net/url"
	"strings"
	"time"

	"todo/models"

	"github.com/gofrs/uuid"
)

func (as *ActionSuite) Test_TasksIndex_EmptyState() {
	res := as.HTML("/").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "No hay tareas todavía")
}

func (as *ActionSuite) Test_TasksIndex_ListsTasks() {
	as.LoadFixture("tres tareas")

	res := as.HTML("/").Get()
	as.Equal(200, res.Code)

	body := res.Body.String()
	as.Contains(body, "Comprar material de oficina")
	as.Contains(body, "Llamar al gestor")
	as.Contains(body, "Publicar la oferta de empleo")
	as.Contains(body, "Completada")
}

func (as *ActionSuite) Test_TasksCreate_Valid() {
	res := as.HTML("/tasks").Post(&models.Task{Title: "  Nueva tarea de prueba  "})
	as.Equal(303, res.Code)
	as.Equal("/", res.Header().Get("Location"))

	task := &models.Task{}
	as.NoError(as.DB.First(task))
	as.Equal("Nueva tarea de prueba", task.Title)
	as.False(task.Completed)
}

func (as *ActionSuite) Test_TasksCreate_EmptyTitle_Invalid() {
	res := as.HTML("/tasks").Post(&models.Task{Title: ""})
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "El título es obligatorio.")

	count, err := as.DB.Count(&models.Task{})
	as.NoError(err)
	as.Equal(0, count)
}

func (as *ActionSuite) Test_TasksCreate_OnlySpaces_Invalid() {
	res := as.HTML("/tasks").Post(&models.Task{Title: "     "})
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "El título es obligatorio.")

	count, err := as.DB.Count(&models.Task{})
	as.NoError(err)
	as.Equal(0, count)
}

func (as *ActionSuite) Test_TasksToggle_CompletesAndReopens() {
	task := &models.Task{Title: "Pendiente de completar"}
	verrs, err := as.DB.ValidateAndCreate(task)
	as.NoError(err)
	as.False(verrs.HasAny())

	res := as.HTML("/tasks/%s/toggle", task.ID).Put(url.Values{})
	as.Equal(303, res.Code)

	as.NoError(as.DB.Reload(task))
	as.True(task.Completed)

	res = as.HTML("/tasks/%s/toggle", task.ID).Put(url.Values{})
	as.Equal(303, res.Code)

	as.NoError(as.DB.Reload(task))
	as.False(task.Completed)
}

func (as *ActionSuite) Test_TasksToggle_NotFound() {
	id := uuid.Must(uuid.NewV4())
	res := as.HTML("/tasks/%s/toggle", id).Put(url.Values{})
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_TasksToggle_InvalidID() {
	res := as.HTML("/tasks/no-es-un-uuid/toggle").Put(url.Values{})
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_TasksDestroy() {
	task := &models.Task{Title: "Tarea que se elimina"}
	verrs, err := as.DB.ValidateAndCreate(task)
	as.NoError(err)
	as.False(verrs.HasAny())

	res := as.HTML("/tasks/%s", task.ID).Delete()
	as.Equal(303, res.Code)

	count, err := as.DB.Count(&models.Task{})
	as.NoError(err)
	as.Equal(0, count)
}

func (as *ActionSuite) Test_TasksDestroy_NotFound() {
	id := uuid.Must(uuid.NewV4())
	res := as.HTML("/tasks/%s", id).Delete()
	as.Equal(404, res.Code)
}

func (as *ActionSuite) taskForEditing(completed bool) *models.Task {
	task := &models.Task{Title: "Título original", Completed: completed}
	verrs, err := as.DB.ValidateAndCreate(task)
	as.NoError(err)
	as.False(verrs.HasAny())
	as.NoError(as.DB.Reload(task))
	return task
}

func (as *ActionSuite) Test_TasksIndex_EditLink() {
	task := as.taskForEditing(false)
	res := as.HTML("/").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "/tasks/"+task.ID.String()+"/edit")
	as.Contains(res.Body.String(), "Editar tarea: Título original")
}

func (as *ActionSuite) Test_TasksEdit_PrefillsTitle() {
	task := as.taskForEditing(true)
	res := as.HTML("/tasks/%s/edit", task.ID).Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), `value="Título original"`)
	as.Contains(res.Body.String(), `value="PUT"`)
	as.Contains(res.Body.String(), "Guardar cambios")
	as.Contains(res.Body.String(), "Cancelar")
	as.NoError(as.DB.Reload(task))
	as.True(task.Completed)
}

func (as *ActionSuite) Test_TasksEdit_NotFound() {
	res := as.HTML("/tasks/%s/edit", uuid.Must(uuid.NewV4())).Get()
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_TasksEdit_InvalidID() {
	res := as.HTML("/tasks/no-es-un-uuid/edit").Get()
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_TasksEdit_EscapesTitle() {
	task := &models.Task{Title: `<img src=x onerror="alert(1)">`}
	verrs, err := as.DB.ValidateAndCreate(task)
	as.NoError(err)
	as.False(verrs.HasAny())
	res := as.HTML("/tasks/%s/edit", task.ID).Get()
	as.Equal(200, res.Code)
	as.NotContains(res.Body.String(), `<img src=x`)
	as.Contains(res.Body.String(), "&lt;img")
}

func (as *ActionSuite) Test_TasksUpdate_PersistsTitleOnly() {
	task := as.taskForEditing(true)
	id, createdAt, updatedAt := task.ID, task.CreatedAt, task.UpdatedAt
	res := as.HTML("/tasks/%s", id).Put(url.Values{
		"title": {"  Título actualizado  "}, "completed": {"false"},
		"id": {uuid.Must(uuid.NewV4()).String()}, "created_at": {"2000-01-01T00:00:00Z"},
	})
	as.Equal(303, res.Code)
	as.Equal("/", res.Header().Get("Location"))
	as.NoError(as.DB.Reload(task))
	as.Equal("Título actualizado", task.Title)
	as.Equal(id, task.ID)
	as.True(task.Completed)
	as.True(task.CreatedAt.Equal(createdAt))
	as.False(task.UpdatedAt.Before(updatedAt))
	count, err := as.DB.Count(&models.Task{})
	as.NoError(err)
	as.Equal(1, count)
}

func (as *ActionSuite) Test_TasksUpdate_PendingTaskRemainsPending() {
	task := as.taskForEditing(false)
	res := as.HTML("/tasks/%s", task.ID).Put(url.Values{"title": {"Nuevo título"}, "completed": {"true"}})
	as.Equal(303, res.Code)
	as.NoError(as.DB.Reload(task))
	as.Equal("Nuevo título", task.Title)
	as.False(task.Completed)
}

func (as *ActionSuite) Test_TasksUpdate_EmptyTitle() {
	task := as.taskForEditing(false)
	res := as.HTML("/tasks/%s", task.ID).Put(url.Values{"title": {""}})
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "El título es obligatorio.")
	as.Contains(res.Body.String(), `aria-describedby="title-error"`)
	as.NoError(as.DB.Reload(task))
	as.Equal("Título original", task.Title)
}

func (as *ActionSuite) Test_TasksUpdate_OnlySpaces() {
	task := as.taskForEditing(true)
	res := as.HTML("/tasks/%s", task.ID).Put(url.Values{"title": {"   "}})
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "El título es obligatorio.")
	as.NoError(as.DB.Reload(task))
	as.Equal("Título original", task.Title)
	as.True(task.Completed)
}

func (as *ActionSuite) Test_TasksUpdate_TitleTooLong() {
	task := as.taskForEditing(false)
	title := strings.Repeat("á", models.TitleMaxLength+1)
	res := as.HTML("/tasks/%s", task.ID).Put(url.Values{"title": {title}})
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "El título no puede superar los 200 caracteres.")
	as.Contains(res.Body.String(), title)
	as.NoError(as.DB.Reload(task))
	as.Equal("Título original", task.Title)
}

func (as *ActionSuite) Test_TasksUpdate_TitleAtLimit() {
	task := as.taskForEditing(false)
	title := strings.Repeat("á", models.TitleMaxLength)
	res := as.HTML("/tasks/%s", task.ID).Put(url.Values{"title": {title}})
	as.Equal(303, res.Code)
	as.NoError(as.DB.Reload(task))
	as.Equal(title, task.Title)
}

func (as *ActionSuite) Test_TasksUpdate_NotFound() {
	res := as.HTML("/tasks/%s", uuid.Must(uuid.NewV4())).Put(url.Values{"title": {"Título nuevo"}})
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_TasksUpdate_InvalidID() {
	res := as.HTML("/tasks/no-es-un-uuid").Put(url.Values{"title": {"Título nuevo"}})
	as.Equal(404, res.Code)
}

// classifiedTask crea una tarea con su importancia y la fecha límite a dueInDays días de hoy.
func (as *ActionSuite) classifiedTask(title, importance string, dueInDays int) *models.Task {
	due := models.Today(time.Now()).AddDate(0, 0, dueInDays)
	task := &models.Task{Title: title, Importance: importance, DueDate: &due}
	verrs, err := as.DB.ValidateAndCreate(task)
	as.NoError(err)
	as.False(verrs.HasAny())
	return task
}

// dueInDays devuelve la fecha a esos días de hoy, en el formato del campo de fecha.
func dueInDays(days int) string {
	return models.Today(time.Now()).AddDate(0, 0, days).Format(models.DueDateLayout)
}

// matrixZone devuelve el HTML de una zona de la matriz: desde id="<id>" hasta el siguiente </section>.
func matrixZone(body, id string) string {
	start := strings.Index(body, `id="`+id+`"`)
	if start < 0 {
		return ""
	}
	end := strings.Index(body[start:], "</section>")
	if end < 0 {
		return body[start:]
	}
	return body[start : start+end]
}

// matrixZones son los ids de las cinco zonas de la matriz.
var matrixZones = []string{"quadrant-do", "quadrant-schedule", "quadrant-delegate", "quadrant-eliminate", "unclassified"}

func (as *ActionSuite) Test_TasksIndex_LinksToMatrix() {
	res := as.HTML("/").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), `href="/matrix"`)
	as.Contains(res.Body.String(), "Ver la matriz de Eisenhower")
}

// Guarda: la lista no muestra la fecha límite de las tareas.
func (as *ActionSuite) Test_TasksIndex_DoesNotShowDueDate() {
	task := as.classifiedTask("Con fecha límite", models.ImportanceImportant, 5)

	res := as.HTML("/").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Con fecha límite")
	as.NotContains(res.Body.String(), task.DueDateValue())
}

func (as *ActionSuite) Test_TasksMatrix_EmptyState() {
	res := as.HTML("/matrix").Get()
	as.Equal(200, res.Code)

	body := res.Body.String()
	as.Contains(body, "Matriz de Eisenhower")
	as.Contains(body, `<a href="/" class="back-link">`)
	as.Contains(body, "Volver a la lista")
	for _, heading := range []string{
		"Hacer ya", "Urgentes e importantes",
		"Planificar", "Importantes, no urgentes",
		"Delegar", "Urgentes, no importantes",
		"Descartar", "Ni urgentes ni importantes",
		"Sin clasificar", "Les falta la fecha límite o decidir su importancia. Edítalas para colocarlas en la matriz.",
	} {
		as.Contains(body, heading)
	}
	as.Equal(4, strings.Count(body, "No hay tareas en este cuadrante."))
	as.Contains(matrixZone(body, "unclassified"), "No hay tareas sin clasificar.")
}

func (as *ActionSuite) Test_TasksMatrix_PlacesTasksInQuadrants() {
	as.classifiedTask("Entregar el informe vencido", models.ImportanceImportant, -1)
	as.classifiedTask("Preparar la reunión", models.ImportanceImportant, 2)
	as.classifiedTask("Planificar el trimestre", models.ImportanceImportant, 3)
	as.classifiedTask("Responder al proveedor", models.ImportanceNotImportant, 0)
	as.classifiedTask("Ordenar el escritorio", models.ImportanceNotImportant, 3)
	as.classifiedTask("Revisar el contrato", models.ImportanceUndecided, 1)
	undated := &models.Task{Title: "Leer el manual", Importance: models.ImportanceImportant}
	verrs, err := as.DB.ValidateAndCreate(undated)
	as.NoError(err)
	as.False(verrs.HasAny())

	res := as.HTML("/matrix").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()

	expected := map[string]string{
		"Entregar el informe vencido": "quadrant-do",
		"Preparar la reunión":         "quadrant-do",
		"Planificar el trimestre":     "quadrant-schedule",
		"Responder al proveedor":      "quadrant-delegate",
		"Ordenar el escritorio":       "quadrant-eliminate",
		"Revisar el contrato":         "unclassified",
		"Leer el manual":              "unclassified",
	}
	for title, zone := range expected {
		for _, id := range matrixZones {
			if id == zone {
				as.Contains(matrixZone(body, id), title, id)
			} else {
				as.NotContains(matrixZone(body, id), title, id)
			}
		}
	}
}

func (as *ActionSuite) Test_TasksMatrix_ExistingTasksUnclassified() {
	as.LoadFixture("tres tareas")

	res := as.HTML("/matrix").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()

	unclassified := matrixZone(body, "unclassified")
	as.Contains(unclassified, "Comprar material de oficina")
	as.Contains(unclassified, "Llamar al gestor")
	as.Contains(unclassified, "Publicar la oferta de empleo")
	as.Contains(unclassified, "task-completed")
	as.Contains(unclassified, "Reabrir tarea: Publicar la oferta de empleo")
	as.Equal(4, strings.Count(body, "No hay tareas en este cuadrante."))
}

func (as *ActionSuite) Test_TasksMatrix_ActionsLikeTheList() {
	classified := as.classifiedTask("Preparar la demo", models.ImportanceImportant, 1)
	unclassified := &models.Task{Title: "Pensar en la demo"}
	verrs, err := as.DB.ValidateAndCreate(unclassified)
	as.NoError(err)
	as.False(verrs.HasAny())

	res := as.HTML("/matrix").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()

	for zone, task := range map[string]*models.Task{"quadrant-do": classified, "unclassified": unclassified} {
		html := matrixZone(body, zone)
		as.Contains(html, "data-task-toggle", zone)
		as.Contains(html, `data-confirm="¿Seguro que quieres eliminar esta tarea?"`, zone)
		as.Equal(2, strings.Count(html, `name="return_to" value="/matrix"`), zone)
		as.Contains(html, "/tasks/"+task.ID.String()+"/edit?return_to=/matrix", zone)
		as.Contains(html, "Completar tarea: "+task.Title, zone)
		as.Contains(html, "Editar tarea: "+task.Title, zone)
		as.Contains(html, "Eliminar tarea: "+task.Title, zone)
	}
}

func (as *ActionSuite) Test_TasksMatrix_EscapesTitle() {
	task := &models.Task{Title: `<img src=x onerror="alert(1)">`}
	verrs, err := as.DB.ValidateAndCreate(task)
	as.NoError(err)
	as.False(verrs.HasAny())

	res := as.HTML("/matrix").Get()
	as.Equal(200, res.Code)
	as.NotContains(res.Body.String(), `<img src=x`)
	as.Contains(res.Body.String(), "&lt;img")
}

func (as *ActionSuite) Test_TasksCreate_WithImportanceAndDueDate() {
	due := dueInDays(4)
	res := as.HTML("/tasks").Post(url.Values{
		"title": {"Preparar la demo"}, "importance": {"important"}, "due_date": {due},
	})
	as.Equal(303, res.Code)
	as.Equal("/", res.Header().Get("Location"))

	task := &models.Task{}
	as.NoError(as.DB.First(task))
	as.Equal("Preparar la demo", task.Title)
	as.Equal(models.ImportanceImportant, task.Importance)
	as.Equal(due, task.DueDateValue())
}

func (as *ActionSuite) Test_TasksCreate_DefaultsUnclassified() {
	res := as.HTML("/tasks").Post(url.Values{"title": {"Solo con el título"}})
	as.Equal(303, res.Code)

	task := &models.Task{}
	as.NoError(as.DB.First(task))
	as.Equal(models.ImportanceUndecided, task.Importance)
	as.Nil(task.DueDate)
}

func (as *ActionSuite) Test_TasksCreate_InvalidDueDate() {
	res := as.HTML("/tasks").Post(url.Values{
		"title": {"Con una fecha imposible"}, "importance": {"important"}, "due_date": {"2026-02-30"},
	})
	as.Equal(422, res.Code)

	body := res.Body.String()
	as.Contains(body, "La fecha límite no es válida.")
	as.Contains(body, `aria-describedby="due-date-error"`)
	as.NotContains(body, `aria-describedby="title-error"`)

	count, err := as.DB.Count(&models.Task{})
	as.NoError(err)
	as.Equal(0, count)
}

func (as *ActionSuite) Test_TasksCreate_InvalidImportance() {
	res := as.HTML("/tasks").Post(url.Values{"title": {"Con una importancia rara"}, "importance": {"muy"}})
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "La importancia no es válida.")
	as.Contains(res.Body.String(), `aria-describedby="importance-error"`)

	count, err := as.DB.Count(&models.Task{})
	as.NoError(err)
	as.Equal(0, count)
}

func (as *ActionSuite) Test_TasksCreate_InvalidTitle_KeepsClassification() {
	due := dueInDays(3)
	res := as.HTML("/tasks").Post(url.Values{"title": {""}, "importance": {"important"}, "due_date": {due}})
	as.Equal(422, res.Code)

	body := res.Body.String()
	as.Contains(body, "El título es obligatorio.")
	as.Contains(body, `<option value="important" selected>`)
	as.Contains(body, `value="`+due+`"`)

	count, err := as.DB.Count(&models.Task{})
	as.NoError(err)
	as.Equal(0, count)
}

func (as *ActionSuite) Test_TasksEdit_PrefillsClassification() {
	task := as.classifiedTask("Preparar la demo", models.ImportanceImportant, 6)

	res := as.HTML("/tasks/%s/edit", task.ID).Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), `<option value="important" selected>`)
	as.Contains(res.Body.String(), `value="`+task.DueDateValue()+`"`)
}

func (as *ActionSuite) Test_TasksEdit_FromMatrix() {
	task := as.taskForEditing(false)

	res := as.HTML("/tasks/%s/edit?return_to=/matrix", task.ID).Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "Volver a la matriz")
	as.Contains(body, `<a href="/matrix" class="back-link">`)
	as.Contains(body, `<a href="/matrix" class="button">Cancelar</a>`)
	as.Contains(body, `name="return_to" value="/matrix"`)

	res = as.HTML("/tasks/%s/edit", task.ID).Get()
	as.Equal(200, res.Code)
	body = res.Body.String()
	as.Contains(body, "Volver a la lista")
	as.NotContains(body, "Volver a la matriz")
	as.Contains(body, `<a href="/" class="back-link">`)
	as.Contains(body, `<a href="/" class="button">Cancelar</a>`)
	as.NotContains(body, `name="return_to"`)
}

func (as *ActionSuite) Test_TasksUpdate_ChangesClassification() {
	task := as.taskForEditing(true)
	due := dueInDays(8)

	res := as.HTML("/tasks/%s", task.ID).Put(url.Values{
		"title": {"Título original"}, "importance": {"not_important"}, "due_date": {due},
	})
	as.Equal(303, res.Code)

	as.NoError(as.DB.Reload(task))
	as.Equal(models.ImportanceNotImportant, task.Importance)
	as.Equal(due, task.DueDateValue())
	as.True(task.Completed)
}

func (as *ActionSuite) Test_TasksUpdate_ClearsDueDate() {
	task := as.classifiedTask("Con fecha que se quita", models.ImportanceImportant, 3)

	res := as.HTML("/tasks/%s", task.ID).Put(url.Values{
		"title": {"Con fecha que se quita"}, "importance": {"important"}, "due_date": {""},
	})
	as.Equal(303, res.Code)

	as.NoError(as.DB.Reload(task))
	as.Nil(task.DueDate)
	as.Equal(models.ImportanceImportant, task.Importance)
}

func (as *ActionSuite) Test_TasksUpdate_InvalidDueDate() {
	task := as.classifiedTask("Fecha que no cambia", models.ImportanceImportant, 3)
	due := task.DueDateValue()

	res := as.HTML("/tasks/%s", task.ID).Put(url.Values{
		"title": {"Otro título"}, "importance": {"not_important"}, "due_date": {"2026-02-30"},
		"return_to": {"/matrix"},
	})
	as.Equal(422, res.Code)

	body := res.Body.String()
	as.Contains(body, "La fecha límite no es válida.")
	as.Contains(body, "Volver a la matriz")
	as.Contains(body, `name="return_to" value="/matrix"`)

	as.NoError(as.DB.Reload(task))
	as.Equal("Fecha que no cambia", task.Title)
	as.Equal(models.ImportanceImportant, task.Importance)
	as.Equal(due, task.DueDateValue())
}

func (as *ActionSuite) Test_TasksUpdate_ReturnsToMatrix() {
	task := as.taskForEditing(false)

	res := as.HTML("/tasks/%s", task.ID).Put(url.Values{
		"title": {"Título desde la matriz"}, "importance": {"important"}, "due_date": {dueInDays(1)},
		"return_to": {"/matrix"},
	})
	as.Equal(303, res.Code)
	as.Equal("/matrix", res.Header().Get("Location"))

	res = as.HTML("/matrix").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Tarea actualizada.")
	as.Contains(matrixZone(res.Body.String(), "quadrant-do"), "Título desde la matriz")
}

func (as *ActionSuite) Test_TasksToggle_ReturnsToMatrix() {
	task := as.taskForEditing(false)

	res := as.HTML("/tasks/%s/toggle", task.ID).Put(url.Values{"return_to": {"/matrix"}})
	as.Equal(303, res.Code)
	as.Equal("/matrix", res.Header().Get("Location"))

	as.NoError(as.DB.Reload(task))
	as.True(task.Completed)
}

// Guarda contra la redirección abierta: solo /matrix cambia adónde se vuelve.
func (as *ActionSuite) Test_TasksToggle_IgnoresOtherReturnPaths() {
	task := as.taskForEditing(false)

	for _, returnTo := range []string{"https://example.com", "/tasks"} {
		res := as.HTML("/tasks/%s/toggle", task.ID).Put(url.Values{"return_to": {returnTo}})
		as.Equal(303, res.Code)
		as.Equal("/", res.Header().Get("Location"), returnTo)
	}
}

func (as *ActionSuite) Test_TasksDestroy_ReturnsToMatrix() {
	task := as.taskForEditing(false)

	res := as.HTML("/tasks/%s", task.ID).Post(url.Values{"_method": {"DELETE"}, "return_to": {"/matrix"}})
	as.Equal(303, res.Code)
	as.Equal("/matrix", res.Header().Get("Location"))

	count, err := as.DB.Count(&models.Task{})
	as.NoError(err)
	as.Equal(0, count)
}
