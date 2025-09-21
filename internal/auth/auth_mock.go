package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockUserAuth struct {
	mock.Mock
}

var _ UserAuthInterface = (*MockUserAuth)(nil)

func (m *MockUserAuth) AuthMiddleware() gin.HandlerFunc {
	args := m.Called()
	return args.Get(0).(gin.HandlerFunc)
}

func (m *MockUserAuth) OptionalAuthMiddleware() gin.HandlerFunc {
	args := m.Called()
	return args.Get(0).(gin.HandlerFunc)
}
