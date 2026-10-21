package models

import (
	"time"
	_ "time/tzdata" // la zona de TimeZone, también en la imagen de Alpine, que no trae zoneinfo
)

// TimeZone es la zona horaria con la que se decide qué día es hoy.
const TimeZone = "Europe/Madrid"

// UrgentDays es cuántos días antes de su fecha límite pasa una tarea a ser urgente.
const UrgentDays = 2

// location es TimeZone cargada. Con time/tzdata en el binario, LoadLocation no falla.
var location = func() *time.Location {
	loc, err := time.LoadLocation(TimeZone)
	if err != nil {
		panic(err)
	}
	return loc
}()

// Today devuelve el día de hoy en TimeZone para el instante now, a las 00:00 UTC, como las fechas límite.
func Today(now time.Time) time.Time {
	y, m, d := now.In(location).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// Matrix reparte las tareas en los cuatro cuadrantes de la matriz de Eisenhower y en «Sin clasificar».
type Matrix struct {
	Do           Tasks // Urgentes e importantes: «Hacer ya».
	Schedule     Tasks // Importantes y no urgentes: «Planificar».
	Delegate     Tasks // Urgentes y no importantes: «Delegar».
	Eliminate    Tasks // Ni urgentes ni importantes: «Descartar».
	Unclassified Tasks // Con la importancia sin decidir o sin fecha límite.
}

// NewMatrix reparte las tareas, en el orden en que llegan, según su importancia y su urgencia el día today.
func NewMatrix(tasks Tasks, today time.Time) Matrix {
	m := Matrix{}
	for _, t := range tasks {
		important := t.Importance == ImportanceImportant
		switch {
		case !t.classified():
			m.Unclassified = append(m.Unclassified, t)
		case important && t.urgent(today):
			m.Do = append(m.Do, t)
		case important:
			m.Schedule = append(m.Schedule, t)
		case t.urgent(today):
			m.Delegate = append(m.Delegate, t)
		default:
			m.Eliminate = append(m.Eliminate, t)
		}
	}
	return m
}

// classified dice si la tarea tiene la importancia decidida y fecha límite, y por tanto un cuadrante.
func (t Task) classified() bool {
	decided := t.Importance == ImportanceImportant || t.Importance == ImportanceNotImportant
	return decided && t.DueDate != nil
}

// urgent dice si la fecha límite de la tarea cae como mucho UrgentDays días después de today,
// también si ya ha pasado.
func (t Task) urgent(today time.Time) bool {
	return t.DueDate != nil && !t.DueDate.After(today.AddDate(0, 0, UrgentDays))
}
