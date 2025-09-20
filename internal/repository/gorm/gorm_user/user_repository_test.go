package gorm_user

import (
	"taskflow/internal/domain/user"
	"testing"

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

	err = db.AutoMigrate(&user.User{})
	require.NoError(t, err)

	return db
}

func TestUserRepository_Create(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(db *gorm.DB)
		inputUser user.User
		wantErr   bool
	}{
		{
			name: "success",
			inputUser: user.User{
				Email:    "abc@example.com",
				Password: "secret",
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			setup: func(db *gorm.DB) {
				db.Create(&user.User{Email: "dup@example.com", Password: "first"})
			},
			inputUser: user.User{
				Email:    "dup@example.com",
				Password: "second",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			if tt.setup != nil {
				tt.setup(db)
			}
			r := NewUserRepository(db)

			err := r.Create(&tt.inputUser)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotZero(t, tt.inputUser.ID)
				require.NotZero(t, tt.inputUser.CreatedAt)
			}
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(db *gorm.DB) user.User
		update    func(u *user.User)
		wantErr   bool
		wantCheck func(t *testing.T, db *gorm.DB, u user.User)
	}{
		{
			name: "success - update password",
			setup: func(db *gorm.DB) user.User {
				u := user.User{Email: "abc@example.com", Password: "oldpass"}
				db.Create(&u)
				return u
			},
			update: func(u *user.User) {
				u.Password = "newpass"
			},
			wantErr: false,
			wantCheck: func(t *testing.T, db *gorm.DB, u user.User) {
				var updated user.User
				err := db.First(&updated, u.ID).Error
				require.NoError(t, err)
				require.Equal(t, "newpass", updated.Password)
			},
		},
		{
			name: "error - update non-existent user",
			setup: func(db *gorm.DB) user.User {
				return user.User{ID: 9999, Email: "ghost@example.com", Password: "none"}
			},
			update: func(u *user.User) {
				u.Password = "newpass"
			},
			wantErr: true,
			wantCheck: func(t *testing.T, db *gorm.DB, u user.User) {
				var found user.User
				err := db.First(&found, u.ID).Error
				require.Error(t, err) // should not exist
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			r := NewUserRepository(db)

			u := tt.setup(db)
			tt.update(&u)

			err := r.Update(&u)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				tt.wantCheck(t, db, u)
			}
		})
	}
}

func TestUserRepository_Delete(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(db *gorm.DB) user.User
		deleteID int
		wantErr  bool
		checkDB  func(t *testing.T, db *gorm.DB, id int)
	}{
		{
			name: "success soft delete existing user",
			setup: func(db *gorm.DB) user.User {
				u := user.User{Email: "test@example.com", Password: "pass"}
				require.NoError(t, db.Create(&u).Error)
				return u
			},
			wantErr: false,
			checkDB: func(t *testing.T, db *gorm.DB, id int) {
				var u user.User
				err := db.First(&u, id).Error
				require.ErrorIs(t, err, gorm.ErrRecordNotFound)

				var u2 user.User
				err = db.Unscoped().First(&u2, id).Error
				require.NoError(t, err)
				require.NotNil(t, u2.DeletedAt.Time)
			},
		},
		{
			name: "delete non-existent user",
			setup: func(db *gorm.DB) user.User {
				return user.User{}
			},
			deleteID: 9999,
			wantErr:  false, // Delete does not error if record doesn't exist
			checkDB: func(t *testing.T, db *gorm.DB, id int) {
				var u user.User
				err := db.First(&u, id).Error
				require.ErrorIs(t, err, gorm.ErrRecordNotFound)
			},
		},
		{
			name:     "missing ID",
			deleteID: 0,
			wantErr:  true,
			checkDB:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			r := NewUserRepository(db)

			var id int
			if tt.deleteID != 0 {
				id = tt.deleteID
			} else if tt.setup != nil {
				u := tt.setup(db)
				id = u.ID
			}

			err := r.Delete(id)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.checkDB != nil {
					tt.checkDB(t, db, id)
				}
			}
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(db *gorm.DB) user.User
		id       int
		wantErr  bool
		wantUser *user.User
	}{
		{
			name: "success - existing user",
			setup: func(db *gorm.DB) user.User {
				u := user.User{Email: "idtest@example.com", Password: "pass"}
				require.NoError(t, db.Create(&u).Error)
				return u
			},
			wantErr:  false,
			wantUser: &user.User{Email: "idtest@example.com", Password: "pass"},
		},
		{
			name:     "non-existent user",
			setup:    nil,
			id:       9999,
			wantErr:  false, // no error for not found
			wantUser: nil,
		},
		{
			name:     "missing ID",
			setup:    nil,
			id:       0,
			wantErr:  true, // invalid ID
			wantUser: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			r := NewUserRepository(db)

			var id int
			if tt.setup != nil {
				u := tt.setup(db)
				id = u.ID
				tt.wantUser.ID = u.ID // set expected ID dynamically
			} else {
				id = tt.id
			}

			got, err := r.GetByID(id)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				if tt.wantUser == nil {
					require.Nil(t, got)
				} else {
					require.Equal(t, tt.wantUser.ID, got.ID)
					require.Equal(t, tt.wantUser.Email, got.Email)
					require.Equal(t, tt.wantUser.Password, got.Password)
				}
			}
		})
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(db *gorm.DB) user.User
		email    string
		wantErr  bool
		wantUser *user.User
	}{
		{
			name: "success - existing email",
			setup: func(db *gorm.DB) user.User {
				u := user.User{Email: "emailtest@example.com", Password: "pass"}
				require.NoError(t, db.Create(&u).Error)
				return u
			},
			wantErr:  false,
			wantUser: &user.User{Email: "emailtest@example.com", Password: "pass"},
		},
		{
			name:     "non-existent email",
			setup:    nil,
			email:    "ghost@example.com",
			wantErr:  false, // no error for not found
			wantUser: nil,
		},
		{
			name:     "missing email",
			setup:    nil,
			email:    "",
			wantErr:  true, // empty email is invalid
			wantUser: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			r := NewUserRepository(db)

			if tt.setup != nil {
				u := tt.setup(db)
				tt.email = u.Email
				tt.wantUser.ID = u.ID // set expected ID dynamically
			}

			got, err := r.GetByEmail(tt.email)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				if tt.wantUser == nil {
					require.Nil(t, got)
				} else {
					require.Equal(t, tt.wantUser.ID, got.ID)
					require.Equal(t, tt.wantUser.Email, got.Email)
					require.Equal(t, tt.wantUser.Password, got.Password)
				}
			}
		})
	}
}

func TestUserRepository_Exists(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(db *gorm.DB) user.User
		id      int
		want    bool
		wantErr bool
	}{
		{
			name: "success - existing user",
			setup: func(db *gorm.DB) user.User {
				u := user.User{Email: "exists@example.com", Password: "pass"}
				require.NoError(t, db.Create(&u).Error)
				return u
			},
			id:      1, // we'll override in the test run after setup
			want:    true,
			wantErr: false,
		},
		{
			name:    "non-existent user",
			id:      9999,
			want:    false,
			wantErr: false,
		},
		{
			name:    "invalid user ID",
			id:      0,
			want:    false,
			wantErr: true,
		},
		{
			name:    "negative user ID",
			id:      -5,
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			r := NewUserRepository(db)

			var userID int
			if tt.setup != nil {
				u := tt.setup(db)
				userID = u.ID
			} else {
				userID = tt.id
			}

			got, err := r.Exists(userID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
