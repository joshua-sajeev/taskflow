package user_handler

import (
	"net/http"

	"taskflow/internal/auth"
	"taskflow/internal/common"
	"taskflow/internal/dto"
	user_service "taskflow/internal/service/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	service  user_service.UserServiceInterface
	userAuth auth.UserAuthInterface
}

func NewUserHandler(s user_service.UserServiceInterface, ua auth.UserAuthInterface) *UserHandler {
	return &UserHandler{service: s, userAuth: ua}
}

var _ UserHandlerInterface = (*UserHandler)(nil)

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User registration data"
// @Success 201 {object} dto.CreateUserResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 409 {object} common.ErrorResponse "Email already exists"
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: err.Error()})
		return
	}

	resp, err := h.service.CreateUser(&req)
	if err != nil {
		if err.Error() == "email already exists" {
			c.JSON(http.StatusConflict, common.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.AuthRequest true "User credentials"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse "Invalid credentials"
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: err.Error()})
		return
	}

	resp, err := h.service.AuthenticateUser(&req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, common.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdatePassword godoc
// @Summary Update user password
// @Description Update user's password (requires authentication)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdatePasswordRequest true "Password update data"
// @Success 200 {object} dto.UpdatePasswordResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /users/password [patch]
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var req dto.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, common.ErrorResponse{Message: "unauthorized"})
		return
	}

	req.ID = userID.(int)

	resp, err := h.service.UpdatePassword(&req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, common.ErrorResponse{Message: "user not found"})
			return
		}
		if err.Error() == "invalid old password" {
			c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteUser godoc
// @Summary Delete user account
// @Description Delete user account (requires authentication)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.DeleteUserResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /users/account [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, common.ErrorResponse{Message: "unauthorized"})
		return
	}

	req := dto.DeleteUserRequest{ID: userID.(int)}

	resp, err := h.service.DeleteUser(&req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, common.ErrorResponse{Message: "user not found"})
			return
		}
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
