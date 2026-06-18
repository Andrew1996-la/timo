package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Andrew1996-la/timo/internal/models"
	"github.com/Andrew1996-la/timo/internal/storage"
	"github.com/stretchr/testify/suite"

	_ "modernc.org/sqlite"
)

type TaskRepositorySuite struct {
	suite.Suite
	db   *sql.DB
	repo *TaskRepository
	ctx  context.Context
}

func (s *TaskRepositorySuite) SetupTest() {
	s.ctx = context.Background()

	db, err := sql.Open("sqlite", ":memory:")
	s.Require().NoError(err)

	err = storage.Migrate(db)
	s.Require().NoError(err)

	s.db = db
	s.repo = NewTaskRepository(db)
}

func (s *TaskRepositorySuite) TearDownTest() {
	if s.db != nil {
		s.Require().NoError(s.db.Close())
	}
}

func (s *TaskRepositorySuite) createTask(title string) *models.Task {
	task, err := s.repo.Create(s.ctx, title)
	s.Require().NoError(err)
	return task
}

func (s *TaskRepositorySuite) TestCreate() {
	s.Run("success", func() {
		task, err := s.repo.Create(s.ctx, "Test task")

		s.Require().NoError(err)
		s.Require().NotNil(task)
		s.NotZero(task.Id)
		s.Equal("Test task", task.Title)
		s.Zero(task.SpentSeconds)
		s.Nil(task.DeletedAt)
		s.False(task.CreatedAt.IsZero())
	})

	s.Run("with empty title", func() {
		task, err := s.repo.Create(s.ctx, "")

		s.Require().NoError(err)
		s.NotZero(task.Id)
		s.Empty(task.Title)
	})

	s.Run("with very long title", func() {
		longTitle := string(make([]byte, 1000))
		task, err := s.repo.Create(s.ctx, longTitle)

		s.Require().NoError(err)
		s.Equal(longTitle, task.Title)
	})
}

func (s *TaskRepositorySuite) TestGetByID() {
	s.Run("existing task", func() {
		created := s.createTask("Test task")

		task, err := s.repo.GetByID(s.ctx, created.Id)

		s.Require().NoError(err)
		s.Require().NotNil(task)
		s.Equal(created.Id, task.Id)
		s.Equal(created.Title, task.Title)
		s.Equal(created.SpentSeconds, task.SpentSeconds)
		s.Equal(created.CreatedAt.Unix(), task.CreatedAt.Unix())
		s.Equal(created.DeletedAt, task.DeletedAt)
	})

	s.Run("non-existent task", func() {
		task, err := s.repo.GetByID(s.ctx, 999)

		s.Require().Error(err)
		s.Nil(task)
		s.True(errors.Is(err, ErrTaskNotFound))
	})

	s.Run("deleted task", func() {
		created := s.createTask("To be deleted")
		s.Require().NoError(s.repo.Delete(s.ctx, created.Id))

		task, err := s.repo.GetByID(s.ctx, created.Id)

		s.Require().Error(err)
		s.Nil(task)
		s.True(errors.Is(err, ErrTaskNotFound))
	})
}

func (s *TaskRepositorySuite) TestGetAll() {
	s.Run("empty database", func() {
		tasks, err := s.repo.GetAll(s.ctx)

		s.Require().NoError(err)
		s.Empty(tasks)
	})

	s.Run("returns all tasks ordered by creation date", func() {
		task1 := s.createTask("First")
		time.Sleep(10 * time.Millisecond)
		task2 := s.createTask("Second")
		time.Sleep(10 * time.Millisecond)
		task3 := s.createTask("Third")

		tasks, err := s.repo.GetAll(s.ctx)

		s.Require().NoError(err)
		s.Len(tasks, 3)

		s.Equal(task3.Id, tasks[0].Id)
		s.Equal(task2.Id, tasks[1].Id)
		s.Equal(task1.Id, tasks[2].Id)

		for _, task := range tasks {
			s.NotZero(task.Id)
			s.NotEmpty(task.Title)
			s.False(task.CreatedAt.IsZero())
			s.Nil(task.DeletedAt)
			s.Zero(task.SpentSeconds)
		}
	})

	s.Run("excludes deleted tasks", func() {
		task1 := s.createTask("Keep")
		task2 := s.createTask("Delete")
		s.Require().NoError(s.repo.Delete(s.ctx, task2.Id))

		tasks, err := s.repo.GetAll(s.ctx)

		s.Require().NoError(err)
		s.Len(tasks, 1)
		s.Equal(task1.Id, tasks[0].Id)
	})

	s.Run("handles many tasks", func() {
		expectedCount := 100
		for i := 0; i < expectedCount; i++ {
			s.createTask("Task " + string(rune(i)))
		}

		tasks, err := s.repo.GetAll(s.ctx)

		s.Require().NoError(err)
		s.Len(tasks, expectedCount)
	})
}

func (s *TaskRepositorySuite) TestDelete() {
	s.Run("successful deletion", func() {
		task := s.createTask("To delete")

		err := s.repo.Delete(s.ctx, task.Id)

		s.Require().NoError(err)

		_, err = s.repo.GetByID(s.ctx, task.Id)
		s.True(errors.Is(err, ErrTaskNotFound))

		tasks, err := s.repo.GetAll(s.ctx)
		s.Require().NoError(err)
		s.Empty(tasks)
	})

	s.Run("non-existent task", func() {
		err := s.repo.Delete(s.ctx, 999)

		s.Require().Error(err)
		s.True(errors.Is(err, ErrTaskNotFound))
	})

	s.Run("already deleted task", func() {
		task := s.createTask("Delete twice")
		s.Require().NoError(s.repo.Delete(s.ctx, task.Id))

		err := s.repo.Delete(s.ctx, task.Id)

		s.Require().Error(err)
		s.True(errors.Is(err, ErrTaskNotFound))
	})

	s.Run("delete with negative id", func() {
		err := s.repo.Delete(s.ctx, -1)

		s.Require().Error(err)
		s.True(errors.Is(err, ErrTaskNotFound))
	})
}

func (s *TaskRepositorySuite) TestAddTime() {
	s.Run("successful time addition", func() {
		task := s.createTask("Time task")

		err := s.repo.AddTime(s.ctx, task.Id, 60)
		s.Require().NoError(err)

		task, err = s.repo.GetByID(s.ctx, task.Id)
		s.Require().NoError(err)
		s.Equal(60, task.SpentSeconds)

		err = s.repo.AddTime(s.ctx, task.Id, 30)
		s.Require().NoError(err)

		task, err = s.repo.GetByID(s.ctx, task.Id)
		s.Require().NoError(err)
		s.Equal(90, task.SpentSeconds)
	})

	s.Run("non-existent task", func() {
		err := s.repo.AddTime(s.ctx, 999, 60)

		s.Require().Error(err)
		s.True(errors.Is(err, ErrTaskNotFound))
	})

	s.Run("deleted task", func() {
		task := s.createTask("Delete before time")
		s.Require().NoError(s.repo.Delete(s.ctx, task.Id))

		err := s.repo.AddTime(s.ctx, task.Id, 60)

		s.Require().Error(err)
		s.True(errors.Is(err, ErrTaskNotFound))
	})

	s.Run("negative time", func() {
		task := s.createTask("Negative time")

		err := s.repo.AddTime(s.ctx, task.Id, -10)
		s.Require().NoError(err)

		task, err = s.repo.GetByID(s.ctx, task.Id)
		s.Require().NoError(err)
		s.Equal(-10, task.SpentSeconds)
	})

	s.Run("zero time", func() {
		task := s.createTask("Zero time")

		err := s.repo.AddTime(s.ctx, task.Id, 0)
		s.Require().NoError(err)

		task, err = s.repo.GetByID(s.ctx, task.Id)
		s.Require().NoError(err)
		s.Zero(task.SpentSeconds)
	})

	s.Run("large time", func() {
		task := s.createTask("Large time")

		err := s.repo.AddTime(s.ctx, task.Id, 1000000)
		s.Require().NoError(err)

		task, err = s.repo.GetByID(s.ctx, task.Id)
		s.Require().NoError(err)
		s.Equal(1000000, task.SpentSeconds)
	})
}

func (s *TaskRepositorySuite) TestIntegration() {
	s.Run("complete workflow", func() {
		task1 := s.createTask("Task 1")
		task2 := s.createTask("Task 2")

		s.Require().NoError(s.repo.AddTime(s.ctx, task1.Id, 120))

		task1, err := s.repo.GetByID(s.ctx, task1.Id)
		s.Require().NoError(err)
		s.Equal(120, task1.SpentSeconds)

		s.Require().NoError(s.repo.Delete(s.ctx, task2.Id))

		tasks, err := s.repo.GetAll(s.ctx)
		s.Require().NoError(err)
		s.Len(tasks, 1)
		s.Equal(task1.Id, tasks[0].Id)

		s.Require().NoError(s.repo.AddTime(s.ctx, task1.Id, 30))

		task1, err = s.repo.GetByID(s.ctx, task1.Id)
		s.Require().NoError(err)
		s.Equal(150, task1.SpentSeconds)

		s.Require().NoError(s.repo.Delete(s.ctx, task1.Id))

		tasks, err = s.repo.GetAll(s.ctx)
		s.Require().NoError(err)
		s.Empty(tasks)
	})

	s.Run("concurrent operations", func() {
		task := s.createTask("Concurrent")
		done := make(chan bool, 10)

		for i := 0; i < 10; i++ {
			go func() {
				err := s.repo.AddTime(s.ctx, task.Id, 10)
				s.Assert().NoError(err)
				done <- true
			}()
		}

		for i := 0; i < 10; i++ {
			<-done
		}

		task, err := s.repo.GetByID(s.ctx, task.Id)
		s.Require().NoError(err)
		s.Equal(100, task.SpentSeconds)
	})
}

func (s *TaskRepositorySuite) TestTimestamps() {
	s.Run("created_at is set correctly", func() {
		before := time.Now()
		task := s.createTask("Timestamp test")
		after := time.Now()

		s.True(task.CreatedAt.After(before) || task.CreatedAt.Equal(before))
		s.True(task.CreatedAt.Before(after) || task.CreatedAt.Equal(after))
	})

	s.Run("created_at doesn't change on update", func() {
		task := s.createTask("No change")
		originalCreated := task.CreatedAt

		s.Require().NoError(s.repo.AddTime(s.ctx, task.Id, 30))

		task, err := s.repo.GetByID(s.ctx, task.Id)
		s.Require().NoError(err)
		s.Equal(originalCreated.Unix(), task.CreatedAt.Unix())
	})

	s.Run("deleted_at is set on deletion", func() {
		task := s.createTask("Delete me")

		s.Require().NoError(s.repo.Delete(s.ctx, task.Id))

		_, err := s.repo.GetByID(s.ctx, task.Id)
		s.True(errors.Is(err, ErrTaskNotFound))
	})
}

func (s *TaskRepositorySuite) TestEdgeCases() {
	s.Run("duplicate titles", func() {
		s.createTask("Same title")
		s.createTask("Same title")
		s.createTask("Same title")

		tasks, err := s.repo.GetAll(s.ctx)
		s.Require().NoError(err)
		s.Len(tasks, 3)

		for _, task := range tasks {
			s.Equal("Same title", task.Title)
		}
	})

	s.Run("create after deletion", func() {
		task := s.createTask("First")
		s.Require().NoError(s.repo.Delete(s.ctx, task.Id))

		newTask := s.createTask("Second")
		s.NotEqual(task.Id, newTask.Id)

		tasks, err := s.repo.GetAll(s.ctx)
		s.Require().NoError(err)
		s.Len(tasks, 1)
		s.Equal(newTask.Id, tasks[0].Id)
	})

	s.Run("add time to multiple tasks", func() {
		task1 := s.createTask("Task 1")
		task2 := s.createTask("Task 2")

		s.Require().NoError(s.repo.AddTime(s.ctx, task1.Id, 10))
		s.Require().NoError(s.repo.AddTime(s.ctx, task2.Id, 20))

		task1, err := s.repo.GetByID(s.ctx, task1.Id)
		s.Require().NoError(err)
		s.Equal(10, task1.SpentSeconds)

		task2, err = s.repo.GetByID(s.ctx, task2.Id)
		s.Require().NoError(err)
		s.Equal(20, task2.SpentSeconds)
	})
}

func TestTaskRepositorySuite(t *testing.T) {
	suite.Run(t, new(TaskRepositorySuite))
}
