package models

import "strings"

func (ms *ModelSuite) Test_Task_Create_Valid() {
	task := &Task{Title: "Comprar café"}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.False(verrs.HasAny())

	ms.NotZero(task.ID)
	ms.False(task.Completed)

	count, err := ms.DB.Count(&Task{})
	ms.NoError(err)
	ms.Equal(1, count)
}

func (ms *ModelSuite) Test_Task_Title_Required() {
	task := &Task{Title: ""}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.True(verrs.HasAny())
	ms.Contains(verrs.Get("title"), "El título es obligatorio.")

	count, err := ms.DB.Count(&Task{})
	ms.NoError(err)
	ms.Equal(0, count)
}

func (ms *ModelSuite) Test_Task_Title_OnlySpaces_Invalid() {
	task := &Task{Title: "   \t  "}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.True(verrs.HasAny())
	ms.Contains(verrs.Get("title"), "El título es obligatorio.")
}

func (ms *ModelSuite) Test_Task_Title_MaxLength() {
	task := &Task{Title: strings.Repeat("a", TitleMaxLength)}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.False(verrs.HasAny())
}

func (ms *ModelSuite) Test_Task_Title_TooLong_Invalid() {
	task := &Task{Title: strings.Repeat("a", TitleMaxLength+1)}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.True(verrs.HasAny())
	ms.NotEmpty(verrs.Get("title"))
}

func (ms *ModelSuite) Test_Task_Toggle_Persists() {
	task := &Task{Title: "Revisar el correo"}
	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.False(verrs.HasAny())

	task.Completed = true
	ms.NoError(ms.DB.Update(task))

	reloaded := &Task{}
	ms.NoError(ms.DB.Find(reloaded, task.ID))
	ms.True(reloaded.Completed)
}

func (ms *ModelSuite) Test_Task_Importance_DefaultsToUndecided() {
	task := &Task{Title: "Sin clasificar todavía"}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.False(verrs.HasAny())

	reloaded := &Task{}
	ms.NoError(ms.DB.Find(reloaded, task.ID))
	ms.Equal(ImportanceUndecided, reloaded.Importance)
	ms.Nil(reloaded.DueDate)
}

func (ms *ModelSuite) Test_Task_Importance_Invalid() {
	task := &Task{Title: "Con una importancia que no existe", Importance: "muy"}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.Equal([]string{"La importancia no es válida."}, verrs.Get("importance"))

	count, err := ms.DB.Count(&Task{})
	ms.NoError(err)
	ms.Equal(0, count)
}

func (ms *ModelSuite) Test_Task_DueDate_Persists() {
	for _, value := range []string{"2026-11-02", "2020-01-15"} {
		task := &Task{Title: "Con fecha límite " + value, Importance: ImportanceImportant}
		task.SetDueDate(value)

		verrs, err := ms.DB.ValidateAndCreate(task)
		ms.NoError(err)
		ms.False(verrs.HasAny(), value)

		reloaded := &Task{}
		ms.NoError(ms.DB.Find(reloaded, task.ID))
		ms.Equal(value, reloaded.DueDateValue())
		ms.Equal(ImportanceImportant, reloaded.Importance)
	}
}

func (ms *ModelSuite) Test_Task_DueDate_EmptyMeansNone() {
	task := &Task{Title: "Sin fecha límite"}
	task.SetDueDate("  ")

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.False(verrs.HasAny())
	ms.Nil(task.DueDate)

	reloaded := &Task{}
	ms.NoError(ms.DB.Find(reloaded, task.ID))
	ms.Nil(reloaded.DueDate)
	ms.Equal("", reloaded.DueDateValue())
}

func (ms *ModelSuite) Test_Task_DueDate_Invalid() {
	for _, value := range []string{"2026-02-30", "mañana", "20/10/2026"} {
		task := &Task{Title: "Con una fecha que no es un día"}
		task.SetDueDate(value)

		verrs, err := ms.DB.ValidateAndCreate(task)
		ms.NoError(err)
		ms.Equal([]string{"La fecha límite no es válida."}, verrs.Get("due_date"), value)
		ms.Nil(task.DueDate, value)
	}

	count, err := ms.DB.Count(&Task{})
	ms.NoError(err)
	ms.Equal(0, count)
}

// Las filas de la fixture no dan valor a las columnas nuevas, como las de antes de la migración:
// prueba el valor por defecto de la base, no el de BeforeValidate.
func (ms *ModelSuite) Test_Task_ExistingRows_Unclassified() {
	ms.LoadFixture("tres tareas")

	tasks := Tasks{}
	ms.NoError(ms.DB.All(&tasks))
	ms.Len(tasks, 3)
	for _, task := range tasks {
		ms.Equal(ImportanceUndecided, task.Importance, task.Title)
		ms.Nil(task.DueDate, task.Title)
	}
}
