package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// TitleMaxLength es el número máximo de caracteres permitido en el título.
const TitleMaxLength = 200

// Valores de la importancia de una tarea.
const (
	ImportanceUndecided    = "undecided"
	ImportanceImportant    = "important"
	ImportanceNotImportant = "not_important"
)

// DueDateLayout es el formato de la fecha límite en los formularios: el de un campo de fecha (AAAA-MM-DD).
const DueDateLayout = "2006-01-02"

// Task representa una tarea de la lista.
type Task struct {
	ID         uuid.UUID  `json:"id" db:"id" form:"-"`
	Title      string     `json:"title" db:"title" form:"title"`
	Completed  bool       `json:"completed" db:"completed" form:"-"`
	Importance string     `json:"importance" db:"importance" form:"importance"`
	DueDate    *time.Time `json:"due_date" db:"due_date" form:"-"`
	// dueDateInvalid indica que SetDueDate recibió algo que no es una fecha. Lo comprueba Validate.
	dueDateInvalid bool      `db:"-" form:"-"`
	CreatedAt      time.Time `json:"created_at" db:"created_at" form:"-"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at" form:"-"`
}

// Tasks es una colección de Task.
type Tasks []Task

// BeforeValidate deja la importancia sin decidir si llega vacía, como el valor por defecto de la base.
func (t *Task) BeforeValidate(tx *pop.Connection) error {
	if t.Importance == "" {
		t.Importance = ImportanceUndecided
	}
	return nil
}

// SetDueDate pone la fecha límite a partir del valor de un campo de fecha (AAAA-MM-DD). Un valor
// vacío quita la fecha; uno que no es un día del calendario la quita y lo anota para Validate.
func (t *Task) SetDueDate(value string) {
	t.DueDate = nil
	t.dueDateInvalid = false

	value = strings.TrimSpace(value)
	if value == "" {
		return
	}

	date, err := time.Parse(DueDateLayout, value)
	if err != nil {
		t.dueDateInvalid = true
		return
	}
	t.DueDate = &date
}

// DueDateValue devuelve la fecha límite con DueDateLayout, o "" si la tarea no tiene fecha.
func (t Task) DueDateValue() string {
	if t.DueDate == nil {
		return ""
	}
	return t.DueDate.Format(DueDateLayout)
}

// Validate se ejecuta en cada llamada a los métodos "pop.Validate*".
func (t *Task) Validate(tx *pop.Connection) (*validate.Errors, error) {
	verrs := validate.NewErrors()

	title := strings.TrimSpace(t.Title)
	if title == "" {
		verrs.Add("title", "El título es obligatorio.")
	} else if len([]rune(title)) > TitleMaxLength {
		verrs.Add("title", fmt.Sprintf("El título no puede superar los %d caracteres.", TitleMaxLength))
	}

	switch t.Importance {
	case ImportanceUndecided, ImportanceImportant, ImportanceNotImportant:
	default:
		verrs.Add("importance", "La importancia no es válida.")
	}

	if t.dueDateInvalid {
		verrs.Add("due_date", "La fecha límite no es válida.")
	}

	return verrs, nil
}
