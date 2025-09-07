package gorm

import (
	"errors"
	"taskflow/internal/domain/task"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Disable logs during tests
	})
	require.NoError(t, err)

	err = db.AutoMigrate(&task.Task{})
	require.NoError(t, err)

	return db
}
func TestTaskRepository_Create(t *testing.T) {
	tests := []struct {
		name      string
		inputTask task.Task
		wantErr   bool
	}{
		{
			name: "successfull creation",
			inputTask: task.Task{
				Task:   "Buy Milk",
				Status: "pending",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			r := NewTaskRepository(db)
			gotErr := r.Create(&tt.inputTask)

			if tt.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
				assert.NotZero(t, tt.inputTask.ID)
				assert.NotZero(t, tt.inputTask.CreatedAt)

				var dbTask task.Task
				result := db.First(&dbTask, tt.inputTask.ID)
				assert.NoError(t, result.Error)
				assert.Equal(t, tt.inputTask.Task, dbTask.Task)
				assert.Equal(t, tt.inputTask.Status, dbTask.Status)

			}

		})
	}
}

func TestTaskRepository_GetByID(t *testing.T) {
	t.Run("existing id", func(t *testing.T) {
		db := setupTestDB(t)

		taskToCreate := task.Task{
			Task:   "Buy Milk",
			Status: "pending",
		}
		err := db.Create(&taskToCreate).Error
		require.NoError(t, err) // now ID is set

		r := NewTaskRepository(db)
		got, gotErr := r.GetByID(taskToCreate.ID)

		assert.NoError(t, gotErr)
		assert.NotNil(t, got)
		assert.Equal(t, taskToCreate.ID, got.ID)
		assert.Equal(t, taskToCreate.Task, got.Task)
		assert.Equal(t, taskToCreate.Status, got.Status)
	})

	t.Run("non-existing id", func(t *testing.T) {
		db := setupTestDB(t)
		r := NewTaskRepository(db)
		got, gotErr := r.GetByID(9999)

		assert.Error(t, gotErr)
		assert.Nil(t, got)
	})
}
func TestTaskRepository_List(t *testing.T) {

	tasks := []task.Task{
		{Task: "Buy Milk", Status: "pending"},
		{Task: "Buy Milk 2", Status: "pending"},
		{Task: "Buy Milk 3", Status: "pending"},
	}

	t.Run("successful response", func(t *testing.T) {
		db := setupTestDB(t)
		r := NewTaskRepository(db)

		for i := range tasks {
			err := db.Create(&tasks[i]).Error
			require.NoError(t, err)
		}

		got, gotErr := r.List()
		assert.NoError(t, gotErr)
		assert.Len(t, got, len(tasks))

		for i, taskItem := range tasks {
			assert.Equal(t, taskItem.Task, got[i].Task)
			assert.Equal(t, taskItem.Status, got[i].Status)
			assert.NotZero(t, got[i].ID)
			assert.NotZero(t, got[i].CreatedAt)
		}
	})
	t.Run("Empyt response", func(t *testing.T) {
		db := setupTestDB(t)
		r := NewTaskRepository(db)

		got, gotErr := r.List()
		assert.NoError(t, gotErr)
		assert.Len(t, got, 0)

	})
}

func TestTaskRepository_Update(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		db := setupTestDB(t)

		taskToCreate := task.Task{Task: "Buy Milk", Status: "pending"}
		require.NoError(t, db.Create(&taskToCreate).Error)

		r := NewTaskRepository(db)
		taskToCreate.Task = "Buy Bread"
		taskToCreate.Status = "completed"
		err := r.Update(&taskToCreate)

		assert.NoError(t, err)

		var updated task.Task
		require.NoError(t, db.First(&updated, taskToCreate.ID).Error)
		assert.Equal(t, "Buy Bread", updated.Task)
		assert.Equal(t, "completed", updated.Status)
	})

	t.Run("update non-existing task", func(t *testing.T) {
		db := setupTestDB(t)
		r := NewTaskRepository(db)
		nonExistentTask := task.Task{ID: 9999, Task: "Nothing", Status: "pending"}
		err := r.Update(&nonExistentTask)

		// GORM's Save creates new record if ID doesn't exist, so error is nil
		assert.NoError(t, err)

		var fetched task.Task
		err = db.First(&fetched, 9999).Error
		assert.NoError(t, err) // record is created
		assert.Equal(t, "Nothing", fetched.Task)
	})
}

func TestTaskRepository_Delete(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		db := setupTestDB(t)
		taskToCreate := task.Task{Task: "Buy Milk", Status: "pending"}
		require.NoError(t, db.Create(&taskToCreate).Error)

		r := NewTaskRepository(db)
		err := r.Delete(taskToCreate.ID)
		assert.NoError(t, err)

		var fetched task.Task
		err = db.First(&fetched, taskToCreate.ID).Error
		assert.Error(t, err) // should not be found
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	})

	t.Run("delete non-existing id", func(t *testing.T) {
		db := setupTestDB(t)
		r := NewTaskRepository(db)
		err := r.Delete(9999)
		assert.NoError(t, err) // GORM Delete does not return error if record not found
	})
}

func TestTaskRepository_UpdateStatus(t *testing.T) {
	t.Run("successful status update", func(t *testing.T) {
		db := setupTestDB(t)
		taskToCreate := task.Task{Task: "Buy Milk", Status: "pending"}
		require.NoError(t, db.Create(&taskToCreate).Error)

		r := NewTaskRepository(db)
		err := r.UpdateStatus(taskToCreate.ID, "completed")
		assert.NoError(t, err)

		var updated task.Task
		require.NoError(t, db.First(&updated, taskToCreate.ID).Error)
		assert.Equal(t, "completed", updated.Status)
	})

	t.Run("update status non-existing task", func(t *testing.T) {
		db := setupTestDB(t)
		r := NewTaskRepository(db)
		err := r.UpdateStatus(9999, "completed")
		assert.NoError(t, err) // GORM does nothing but does not error
	})
}
