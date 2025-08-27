package dto

type CreateTaskRequest struct {
	Task string `json:"task" binding:"required,string,max=20" example:"Buy milk"`
}

type GetTaskResponse struct {
	ID     int    `json:"id" example:"1"`
	Task   string `json:"task" example:"Buy milk"`
	Status string `json:"status" example:"pending"`
}

type ListTasksResponse struct {
	Tasks []GetTaskResponse `json:"tasks"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending in-progress completed" example:"completed"`
}

type UpdateStatusResponse struct {
	Message string `json:"message" example:"status updated"`
}

type DeleteTaskResponse struct {
	Message string `json:"message" example:"Task deleted successfully"`
}
