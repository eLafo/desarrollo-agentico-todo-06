package models

import "time"

// matrixToday es el día fijo con el que se reparten las tareas en estos tests.
var matrixToday = time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)

// dueTask crea en memoria una tarea con su importancia y la fecha límite a dueInDays días de matrixToday.
func dueTask(title, importance string, dueInDays int) Task {
	due := matrixToday.AddDate(0, 0, dueInDays)
	return Task{Title: title, Importance: importance, DueDate: &due}
}

// titles devuelve los títulos de las tareas, en su orden.
func titles(tasks Tasks) []string {
	out := []string{}
	for _, t := range tasks {
		out = append(out, t.Title)
	}
	return out
}

func (ms *ModelSuite) Test_Today_UsesMadridDate() {
	// 01:30 en Madrid, hora de verano.
	ms.Equal(time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC), Today(time.Date(2026, 10, 19, 23, 30, 0, 0, time.UTC)))
	// 23:30 en Madrid, el mismo día.
	ms.Equal(time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC), Today(time.Date(2026, 10, 20, 21, 30, 0, 0, time.UTC)))
	// 00:30 en Madrid, hora de invierno.
	ms.Equal(time.Date(2026, 1, 16, 0, 0, 0, 0, time.UTC), Today(time.Date(2026, 1, 15, 23, 30, 0, 0, time.UTC)))
}

func (ms *ModelSuite) Test_Matrix_ImportantTasks() {
	m := NewMatrix(Tasks{
		dueTask("Vencida", ImportanceImportant, -1),
		dueTask("Hoy", ImportanceImportant, 0),
		dueTask("Mañana", ImportanceImportant, 1),
		dueTask("Pasado mañana", ImportanceImportant, 2),
		dueTask("Dentro de tres días", ImportanceImportant, 3),
	}, matrixToday)

	ms.Equal([]string{"Vencida", "Hoy", "Mañana", "Pasado mañana"}, titles(m.Do))
	ms.Equal([]string{"Dentro de tres días"}, titles(m.Schedule))
	ms.Empty(m.Delegate)
	ms.Empty(m.Eliminate)
	ms.Empty(m.Unclassified)
}

func (ms *ModelSuite) Test_Matrix_NotImportantTasks() {
	m := NewMatrix(Tasks{
		dueTask("Vencida", ImportanceNotImportant, -1),
		dueTask("Pasado mañana", ImportanceNotImportant, 2),
		dueTask("Dentro de tres días", ImportanceNotImportant, 3),
		dueTask("Dentro de un mes", ImportanceNotImportant, 30),
	}, matrixToday)

	ms.Equal([]string{"Vencida", "Pasado mañana"}, titles(m.Delegate))
	ms.Equal([]string{"Dentro de tres días", "Dentro de un mes"}, titles(m.Eliminate))
	ms.Empty(m.Do)
	ms.Empty(m.Schedule)
	ms.Empty(m.Unclassified)
}

func (ms *ModelSuite) Test_Matrix_Unclassified() {
	m := NewMatrix(Tasks{
		dueTask("Sin decidir con fecha", ImportanceUndecided, 1),
		{Title: "Importante sin fecha", Importance: ImportanceImportant},
		{Title: "No importante sin fecha", Importance: ImportanceNotImportant},
		{Title: "Sin decidir sin fecha", Importance: ImportanceUndecided},
	}, matrixToday)

	ms.Equal([]string{
		"Sin decidir con fecha", "Importante sin fecha", "No importante sin fecha", "Sin decidir sin fecha",
	}, titles(m.Unclassified))
	ms.Empty(m.Do)
	ms.Empty(m.Schedule)
	ms.Empty(m.Delegate)
	ms.Empty(m.Eliminate)
}

func (ms *ModelSuite) Test_Matrix_EveryTaskOnce() {
	tasks := Tasks{
		dueTask("Importante vencida", ImportanceImportant, -1),
		dueTask("Importante hoy", ImportanceImportant, 0),
		dueTask("Importante mañana", ImportanceImportant, 1),
		dueTask("Importante pasado mañana", ImportanceImportant, 2),
		dueTask("Importante en tres días", ImportanceImportant, 3),
		dueTask("No importante vencida", ImportanceNotImportant, -1),
		dueTask("No importante pasado mañana", ImportanceNotImportant, 2),
		dueTask("No importante en tres días", ImportanceNotImportant, 3),
		dueTask("No importante en un mes", ImportanceNotImportant, 30),
		dueTask("Sin decidir con fecha", ImportanceUndecided, 1),
		{Title: "Importante sin fecha", Importance: ImportanceImportant},
		{Title: "No importante sin fecha", Importance: ImportanceNotImportant},
		{Title: "Sin decidir sin fecha", Importance: ImportanceUndecided},
	}
	m := NewMatrix(tasks, matrixToday)

	seen := map[string]bool{}
	total := 0
	for _, zone := range []Tasks{m.Do, m.Schedule, m.Delegate, m.Eliminate, m.Unclassified} {
		for _, t := range zone {
			ms.False(seen[t.Title], t.Title)
			seen[t.Title] = true
			total++
		}
	}
	ms.Equal(len(tasks), total)
}

func (ms *ModelSuite) Test_Matrix_BecomesUrgentByItself() {
	tasks := Tasks{dueTask("Vence el 13", ImportanceImportant, 3)}

	m := NewMatrix(tasks, matrixToday)
	ms.Equal([]string{"Vence el 13"}, titles(m.Schedule))
	ms.Empty(m.Do)

	m = NewMatrix(tasks, matrixToday.AddDate(0, 0, 1))
	ms.Equal([]string{"Vence el 13"}, titles(m.Do))
	ms.Empty(m.Schedule)
}

func (ms *ModelSuite) Test_Matrix_KeepsOrderAndCompleted() {
	completed := dueTask("Segunda, completada", ImportanceImportant, 1)
	completed.Completed = true

	m := NewMatrix(Tasks{
		dueTask("Primera", ImportanceImportant, 2),
		completed,
		dueTask("Tercera", ImportanceImportant, -5),
	}, matrixToday)

	ms.Equal([]string{"Primera", "Segunda, completada", "Tercera"}, titles(m.Do))
	ms.True(m.Do[1].Completed)
}
