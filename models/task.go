package models

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// TitleMaxLength es el número máximo de caracteres permitido en el título.
const TitleMaxLength = 200

// ErrMoveOutOfRange es el error que devuelve Move si el destino de la tarea cae fuera de la lista.
var ErrMoveOutOfRange = errors.New("la tarea no se puede mover fuera de la lista")

// Task representa una tarea de la lista.
type Task struct {
	ID        uuid.UUID `json:"id" db:"id" form:"-"`
	Title     string    `json:"title" db:"title" form:"title"`
	Completed bool      `json:"completed" db:"completed" form:"-"`
	Position  int       `json:"position" db:"position" form:"-"`
	CreatedAt time.Time `json:"created_at" db:"created_at" form:"-"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" form:"-"`
}

// Tasks es una colección de Task.
type Tasks []Task

// Validate se ejecuta en cada llamada a los métodos "pop.Validate*".
func (t *Task) Validate(tx *pop.Connection) (*validate.Errors, error) {
	verrs := validate.NewErrors()

	title := strings.TrimSpace(t.Title)
	if title == "" {
		verrs.Add("title", "El título es obligatorio.")
	} else if len([]rune(title)) > TitleMaxLength {
		verrs.Add("title", fmt.Sprintf("El título no puede superar los %d caracteres.", TitleMaxLength))
	}

	return verrs, nil
}

// AllTasks devuelve todas las tareas en el orden de la lista: por puesto y, a igual puesto, por creación.
func AllTasks(tx *pop.Connection) (Tasks, error) {
	tasks := Tasks{}
	if err := tx.Order("position asc, created_at asc, id asc").All(&tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// BeforeCreate pone la tarea nueva al final de la lista, un puesto después del mayor.
func (t *Task) BeforeCreate(tx *pop.Connection) error {
	last := Tasks{}
	if err := tx.Order("position desc").Limit(1).All(&last); err != nil {
		return err
	}

	t.Position = 1
	if len(last) > 0 {
		t.Position = last[0].Position + 1
	}
	return nil
}

// Move desplaza la tarea offset puestos en la lista: con un número negativo sube y con uno positivo baja.
func (t *Task) Move(tx *pop.Connection, offset int) error {
	tasks, err := AllTasks(tx)
	if err != nil {
		return err
	}

	from := slices.IndexFunc(tasks, func(other Task) bool { return other.ID == t.ID })
	if !tasks.canMove(from, offset) {
		return ErrMoveOutOfRange
	}

	ordered := slices.Insert(slices.Delete(slices.Clone(tasks), from, from+1), from+offset, tasks[from])
	for i := range ordered {
		if ordered[i].Position == i+1 {
			continue
		}
		ordered[i].Position = i + 1
		if err := tx.UpdateColumns(&ordered[i], "position"); err != nil {
			return err
		}
	}

	t.Position = from + offset + 1
	return nil
}

// CanMoveUp dice si la tarea con índice i en la lista puede subir un puesto.
func (ts Tasks) CanMoveUp(i int) bool {
	return ts.canMove(i, -1)
}

// CanMoveDown dice si la tarea con índice i en la lista puede bajar un puesto.
func (ts Tasks) CanMoveDown(i int) bool {
	return ts.canMove(i, 1)
}

// canMove dice si la tarea con índice i en la lista puede desplazarse offset puestos sin salir de ella.
func (ts Tasks) canMove(i, offset int) bool {
	to := i + offset
	return i >= 0 && i < len(ts) && to >= 0 && to < len(ts)
}
