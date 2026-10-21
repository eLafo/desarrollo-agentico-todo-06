package models

import (
	"errors"
	"io/fs"
	"strings"
	"time"

	"todo/migrations"
)

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

func (ms *ModelSuite) createTasks(titles ...string) Tasks {
	tasks := Tasks{}
	for _, title := range titles {
		task := Task{Title: title}
		verrs, err := ms.DB.ValidateAndCreate(&task)
		ms.NoError(err)
		ms.False(verrs.HasAny())
		tasks = append(tasks, task)
	}
	return tasks
}

func (ms *ModelSuite) listedTitles() []string {
	tasks, err := AllTasks(ms.DB)
	ms.NoError(err)
	titles := []string{}
	for _, task := range tasks {
		titles = append(titles, task.Title)
	}
	return titles
}

func (ms *ModelSuite) listedPositions() []int {
	tasks, err := AllTasks(ms.DB)
	ms.NoError(err)
	positions := []int{}
	for _, task := range tasks {
		positions = append(positions, task.Position)
	}
	return positions
}

func (ms *ModelSuite) Test_Task_Create_AppendsAtEnd() {
	tasks := ms.createTasks("A", "B", "C")
	ms.Equal([]int{1, 2, 3}, ms.listedPositions())

	ms.NoError(tasks[2].Move(ms.DB, -1))
	ms.NoError(ms.DB.Destroy(&tasks[0]))
	ms.createTasks("D")

	ms.Equal([]string{"C", "B", "D"}, ms.listedTitles())
}

func (ms *ModelSuite) Test_Task_Move_Up() {
	tasks := ms.createTasks("A", "B", "C")

	ms.NoError(tasks[1].Move(ms.DB, -1))
	ms.Equal([]string{"B", "A", "C"}, ms.listedTitles())
	ms.Equal([]int{1, 2, 3}, ms.listedPositions())
	ms.Equal(1, tasks[1].Position)
}

func (ms *ModelSuite) Test_Task_Move_Down() {
	tasks := ms.createTasks("A", "B", "C")

	ms.NoError(tasks[1].Move(ms.DB, 1))
	ms.Equal([]string{"A", "C", "B"}, ms.listedTitles())
	ms.Equal([]int{1, 2, 3}, ms.listedPositions())
	ms.Equal(3, tasks[1].Position)
}

func (ms *ModelSuite) Test_Task_Move_SeveralPlaces() {
	tasks := ms.createTasks("A", "B", "C")

	ms.NoError(tasks[2].Move(ms.DB, -2))
	ms.Equal([]string{"C", "A", "B"}, ms.listedTitles())
	ms.Equal([]int{1, 2, 3}, ms.listedPositions())

	ms.NoError(tasks[2].Move(ms.DB, 2))
	ms.Equal([]string{"A", "B", "C"}, ms.listedTitles())
	ms.Equal([]int{1, 2, 3}, ms.listedPositions())
}

func (ms *ModelSuite) Test_Task_Move_FirstCannotMoveUp() {
	tasks := ms.createTasks("A", "B", "C")

	for _, offset := range []int{-1, -4} {
		err := tasks[0].Move(ms.DB, offset)
		ms.True(errors.Is(err, ErrMoveOutOfRange), "offset %d", offset)
		ms.Equal([]string{"A", "B", "C"}, ms.listedTitles())
		ms.Equal([]int{1, 2, 3}, ms.listedPositions())
	}
}

func (ms *ModelSuite) Test_Task_Move_LastCannotMoveDown() {
	tasks := ms.createTasks("A", "B", "C")

	for _, offset := range []int{1, 4} {
		err := tasks[2].Move(ms.DB, offset)
		ms.True(errors.Is(err, ErrMoveOutOfRange), "offset %d", offset)
		ms.Equal([]string{"A", "B", "C"}, ms.listedTitles())
		ms.Equal([]int{1, 2, 3}, ms.listedPositions())
	}
}

func (ms *ModelSuite) Test_Task_Move_AfterDestroy() {
	tasks := ms.createTasks("A", "B", "C", "D")
	ms.NoError(ms.DB.Destroy(&tasks[1]))

	ms.NoError(tasks[2].Move(ms.DB, -1))
	ms.Equal([]string{"C", "A", "D"}, ms.listedTitles())
	ms.Equal([]int{1, 2, 3}, ms.listedPositions())

	err := tasks[3].Move(ms.DB, 1)
	ms.True(errors.Is(err, ErrMoveOutOfRange))
	ms.Equal([]string{"C", "A", "D"}, ms.listedTitles())
}

func (ms *ModelSuite) Test_Task_Move_TiedPositions() {
	tasks := ms.createTasks("A", "B", "C")
	for i := range tasks {
		tasks[i].Position = 0
		ms.NoError(ms.DB.UpdateColumns(&tasks[i], "position"))
	}
	ms.Equal([]int{0, 0, 0}, ms.listedPositions())
	ms.Equal([]string{"A", "B", "C"}, ms.listedTitles())

	ms.NoError(tasks[2].Move(ms.DB, -1))
	ms.Equal([]string{"A", "C", "B"}, ms.listedTitles())
	ms.Equal([]int{1, 2, 3}, ms.listedPositions())
}

func (ms *ModelSuite) Test_Tasks_CanMoveUpAndDown() {
	three := Tasks{{Title: "A"}, {Title: "B"}, {Title: "C"}}
	ms.False(three.CanMoveUp(0))
	ms.True(three.CanMoveUp(1))
	ms.True(three.CanMoveUp(2))
	ms.True(three.CanMoveDown(0))
	ms.True(three.CanMoveDown(1))
	ms.False(three.CanMoveDown(2))

	one := Tasks{{Title: "A"}}
	ms.False(one.CanMoveUp(0))
	ms.False(one.CanMoveDown(0))
}

func (ms *ModelSuite) Test_Task_Migration_NumbersByCreation() {
	day := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	for _, task := range []Task{
		{Title: "Tercera", CreatedAt: day.Add(9*time.Hour + 10*time.Minute)},
		{Title: "Primera", CreatedAt: day.Add(9 * time.Hour)},
		{Title: "Segunda", CreatedAt: day.Add(9*time.Hour + 5*time.Minute)},
	} {
		verrs, err := ms.DB.ValidateAndCreate(&task)
		ms.NoError(err)
		ms.False(verrs.HasAny())
	}
	ms.Equal([]string{"Tercera", "Primera", "Segunda"}, ms.listedTitles())

	sql, err := fs.ReadFile(migrations.FS(), "20261020000001_number_task_positions.up.sql")
	ms.NoError(err)
	ms.NoError(ms.DB.RawQuery(string(sql)).Exec())

	ms.Equal([]string{"Primera", "Segunda", "Tercera"}, ms.listedTitles())
	ms.Equal([]int{1, 2, 3}, ms.listedPositions())
}
