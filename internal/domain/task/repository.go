package task

type TaskRepositoryInterface interface {
	Create(task *Task) error
	GetByID(id int) (*Task, error)
	List() ([]Task, error)
	Update(task *Task) error
	Delete(id int) error
	UpdateStatus(id int, status string) error
}
