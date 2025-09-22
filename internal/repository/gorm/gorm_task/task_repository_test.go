package gorm_task

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
		inputUser task.Task
		wantErr   bool
	}{
		{
			name: "successfull creation",
			inputUser: task.Task{
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
			gotErr := r.Create(&tt.inputUser)

			if tt.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
				assert.NotZero(t, tt.inputUser.ID)
				assert.NotZero(t, tt.inputUser.CreatedAt)

				var dbTask task.Task
				result := db.First(&dbTask, tt.inputUser.ID)
				assert.NoError(t, result.Error)
				assert.Equal(t, tt.inputUser.Task, dbTask.Task)
				assert.Equal(t, tt.inputUser.Status, dbTask.Status)

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
			UserID: 1,
		}
		err := db.Create(&taskToCreate).Error
		require.NoError(t, err)

		r := NewTaskRepository(db)
		got, gotErr := r.GetByID(1, taskToCreate.ID)

		assert.NoError(t, gotErr)
		assert.NotNil(t, got)
		assert.Equal(t, taskToCreate.ID, got.ID)
		assert.Equal(t, taskToCreate.Task, got.Task)
		assert.Equal(t, taskToCreate.Status, got.Status)
		assert.Equal(t, taskToCreate.UserID, got.UserID)
	})

	t.Run("non-existing id", func(t *testing.T) {
		db := setupTestDB(t)
		r := NewTaskRepository(db)
		got, gotErr := r.GetByID(1, 9999)

		assert.Error(t, gotErr)
		assert.Nil(t, got)
	})
}

func TestTaskRepository_List(t *testing.T) {
	tasks := []task.Task{
		{Task: "Buy Milk", Status: "pending", UserID: 1},
		{Task: "Buy Milk 2", Status: "pending", UserID: 1},
		{Task: "Buy Milk 3", Status: "pending", UserID: 2}, // different user
	}

	t.Run("successful response with userID filter", func(t *testing.T) {
		db := setupTestDB(t)
		r := NewTaskRepository(db)

		for i := range tasks {
			require.NoError(t, db.Create(&tasks[i]).Error)
		}

		got, err := r.List(1)
		assert.NoError(t, err)
		assert.Len(t, got, 2) // only UserID 1 tasks

		for _, tsk := range got {
			assert.Equal(t, 1, tsk.UserID)
			assert.NotZero(t, tsk.ID)
			assert.NotZero(t, tsk.CreatedAt)
		}
	})

	t.Run("empty response", func(t *testing.T) {
		db := setupTestDB(t)
		r := NewTaskRepository(db)

		got, err := r.List(99) // userID with no tasks
		assert.NoError(t, err)
		assert.Len(t, got, 0)
	})

	t.Run("database error", func(t *testing.T) {
		db := setupTestDB(t)
		// drop table to simulate failure
		db.Migrator().DropTable(&task.Task{})
		r := NewTaskRepository(db)

		got, err := r.List(1)
		assert.Error(t, err)
		assert.Nil(t, got)
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
		taskToCreate := task.Task{UserID: 1, Task: "Buy Milk", Status: "pending"}
		require.NoError(t, db.Create(&taskToCreate).Error)

		r := NewTaskRepository(db)
		err := r.Delete(1, taskToCreate.ID)
		assert.NoError(t, err)

		var fetched task.Task
		err = db.First(&fetched, taskToCreate.ID).Error
		assert.Error(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	})

	t.Run("delete non-existing task", func(t *testing.T) {
		db := setupTestDB(t)
		r := NewTaskRepository(db)
		err := r.Delete(1, 9999)
		assert.NoError(t, err) // GORM does not error if record not found
	})
}

func TestTaskRepository_UpdateStatus(t *testing.T) {
	t.Run("successful status update", func(t *testing.T) {
		db := setupTestDB(t)

		taskToCreate := task.Task{UserID: 1, Task: "Buy Milk", Status: "pending"}
		require.NoError(t, db.Create(&taskToCreate).Error)

		r := NewTaskRepository(db)
		err := r.UpdateStatus(1, taskToCreate.ID, "completed")
		assert.NoError(t, err)

		var updated task.Task
		require.NoError(t, db.First(&updated, taskToCreate.ID).Error)
		assert.Equal(t, "completed", updated.Status)
	})

	t.Run("update status non-existing task", func(t *testing.T) {
		db := setupTestDB(t)
		r := NewTaskRepository(db)

		// Non-existing userID + taskID
		err := r.UpdateStatus(1, 9999, "completed")
		assert.NoError(t, err) // GORM returns nil error if no rows affected
	})
}
