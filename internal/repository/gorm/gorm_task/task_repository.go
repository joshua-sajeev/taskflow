package gorm_task

import (
	"taskflow/internal/domain/task"

	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Compile-time check
var _ TaskRepositoryInterface = (*TaskRepository)(nil)

func (r *TaskRepository) Create(t *task.Task) error {
	return r.db.Create(t).Error
}

func (r *TaskRepository) GetByID(userID int, id int) (*task.Task, error) {
	var t task.Task
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepository) List(userID int) ([]task.Task, error) {
	var tasks []task.Task
	if err := r.db.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepository) Update(t *task.Task) error {
	return r.db.Save(t).Error
}

func (r *TaskRepository) Delete(id int) error {
	return r.db.Delete(&task.Task{}, id).Error
}

func (r *TaskRepository) UpdateStatus(id int, status string) error {
	return r.db.Model(&task.Task{}).Where("id = ?", id).Update("status", status).Error
}
