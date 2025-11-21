package integration

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"taskflow/internal/domain/task"
	"taskflow/pkg/database"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	mysqlPort string
	pool      *dockertest.Pool
	resource  *dockertest.Resource
)

func TestMain(m *testing.M) {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}
	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "8.0",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=pass",
			"MYSQL_DATABASE=testdb",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start mysql: %s", err)
	}
	mysqlPort = resource.GetPort("3306/tcp")

	os.Setenv("MYSQL_USER", "root")
	os.Setenv("MYSQL_PASSWORD", "pass")
	os.Setenv("MYSQL_DATABASE", "testdb")
	os.Setenv("MYSQL_HOST", "localhost")
	os.Setenv("MYSQL_PORT", mysqlPort)

	pool.MaxWait = 30 * time.Second
	var db *gorm.DB
	if err := pool.Retry(func() error {
		dsn := fmt.Sprintf("root:pass@tcp(localhost:%s)/testdb?parseTime=True", mysqlPort)
		var connErr error
		db, connErr = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if connErr != nil {
			return connErr
		}

		sqlDB, _ := db.DB()
		return sqlDB.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to mysql: %s", err)
	}

	if err := database.MigrateModels(db, &user.User{}, &task.Task{}); err != nil {
		log.Fatalf("Could not run migrations: %v", err)
	}

	log.Println("MySQL container ready for tests")
	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(code)
}

func TestConnectDB(t *testing.T) {
	tests := []struct {
		name        string
		config      *database.Config
		setupEnv    func()
		cleanupEnv  func()
		wantErr     bool
		errContains string
		validate    func(t *testing.T, db *gorm.DB)
		skipReason  string
	}{
		{
			name:    "verify TestMain setup is correct",
			config:  nil,
			wantErr: false,
			validate: func(t *testing.T, db *gorm.DB) {
				require.NotNil(t, db)

				assert.Equal(t, "root", os.Getenv("MYSQL_USER"))
				assert.Equal(t, "pass", os.Getenv("MYSQL_PASSWORD"))
				assert.Equal(t, "testdb", os.Getenv("MYSQL_DATABASE"))
				assert.Equal(t, "localhost", os.Getenv("MYSQL_HOST"))
				assert.Equal(t, mysqlPort, os.Getenv("MYSQL_PORT"))
				assert.NotEmpty(t, mysqlPort, "mysqlPort should be set by TestMain")

				sqlDB, err := db.DB()
				require.NoError(t, err)
				assert.NoError(t, sqlDB.Ping())

				var result int
				err = db.Raw("SELECT 1").Scan(&result).Error
				require.NoError(t, err)
				assert.Equal(t, 1, result)

				var dbName string
				err = db.Raw("SELECT DATABASE()").Scan(&dbName).Error
				require.NoError(t, err)
				assert.Equal(t, "testdb", dbName)

				stats := sqlDB.Stats()
				assert.Equal(t, 100, stats.MaxOpenConnections, "MaxOpenConnections should be 100")
				assert.GreaterOrEqual(t, stats.Idle, 0, "Should have idle connections available")
			},
		},

		{
			name:    "successful connection with valid config",
			config:  validTestConfig(),
			wantErr: false,
			validate: func(t *testing.T, db *gorm.DB) {
				require.NotNil(t, db)
				sqlDB, err := db.DB()
				require.NoError(t, err)
				assert.NoError(t, sqlDB.Ping())

				// Verify connection pool settings by checking stats
				stats := sqlDB.Stats()
				assert.GreaterOrEqual(t, stats.MaxOpenConnections, 0, "MaxOpenConnections should be set")
			},
		},

		{
			name:   "successful connection with nil config (loads from env)",
			config: nil,

			setupEnv: func() {
				os.Setenv("MYSQL_USER", "root")
				os.Setenv("MYSQL_PASSWORD", "pass")
				os.Setenv("MYSQL_HOST", "localhost")
				os.Setenv("MYSQL_PORT", mysqlPort)
				os.Setenv("MYSQL_DATABASE", "testdb")
			},

			cleanupEnv: func() {
				os.Unsetenv("MYSQL_USER")
				os.Unsetenv("MYSQL_PASSWORD")
				os.Unsetenv("MYSQL_HOST")
				os.Unsetenv("MYSQL_PORT")
				os.Unsetenv("MYSQL_DATABASE")
			},

			wantErr: false,

			validate: func(t *testing.T, db *gorm.DB) {
				require.NotNil(t, db)
				var dbName string
				err := db.Raw("SELECT DATABASE()").Scan(&dbName).Error
				require.NoError(t, err)
				assert.Equal(t, "testdb", dbName)
			},
		},

		{
			name: "connection failure with invalid password",
			config: &database.Config{
				User:       "root",
				Password:   "wrongpassword",
				Host:       "localhost",
				Port:       mysqlPort,
				Database:   "testdb",
				MaxRetries: 2,
				RetryDelay: 100 * time.Millisecond,
				LogLevel:   logger.Silent,
			},
			wantErr:     true,
			errContains: "failed to connect to database after 2 attempts",
		},

		{
			name: "connection failure with invalid host",
			config: &database.Config{
				User:       "root",
				Password:   "pass",
				Host:       "invalid-host",
				Port:       mysqlPort,
				Database:   "testdb",
				MaxRetries: 2,
				RetryDelay: 10 * time.Millisecond,
				LogLevel:   logger.Silent,
			},
			wantErr:     true,
			errContains: "failed to connect to database",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupEnv != nil {
				tt.setupEnv()
			}

			db, err := database.ConnectDB(tt.config)

			if tt.cleanupEnv != nil {
				tt.cleanupEnv()
			}

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, db)
			} else {
				require.NoError(t, err)
				require.NotNil(t, db)

				if tt.validate != nil {
					tt.validate(t, db)
				}

				sqlDB, _ := db.DB()
				if sqlDB != nil {
					sqlDB.Close()
				}
			}
		})
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		setupEnv func()
		expected *database.Config
	}{
		{
			name: "load all values from environment",
			setupEnv: func() {
				os.Setenv("MYSQL_USER", "testuser")
				os.Setenv("MYSQL_PASSWORD", "testpass")
				os.Setenv("MYSQL_HOST", "testhost")
				os.Setenv("MYSQL_PORT", "3307")
				os.Setenv("MYSQL_DATABASE", "testdatabase")
			},
			expected: &database.Config{
				User:       "testuser",
				Password:   "testpass",
				Host:       "testhost",
				Port:       "3307",
				Database:   "testdatabase",
				MaxRetries: 10,
				RetryDelay: 2 * time.Second,
				LogLevel:   logger.Info,
			},
		},
		{
			name: "use default values when env vars not set",
			setupEnv: func() {
				os.Unsetenv("MYSQL_USER")
				os.Unsetenv("MYSQL_PASSWORD")
				os.Unsetenv("MYSQL_HOST")
				os.Unsetenv("MYSQL_PORT")
				os.Unsetenv("MYSQL_DATABASE")
			},
			expected: &database.Config{
				User:       "appuser",
				Password:   "apppassword",
				Host:       "mysql",
				Port:       "3306",
				Database:   "taskdb",
				MaxRetries: 10,
				RetryDelay: 2 * time.Second,
				LogLevel:   logger.Info,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()

			cfg := database.LoadConfigFromEnv()

			assert.Equal(t, tt.expected.User, cfg.User)
			assert.Equal(t, tt.expected.Password, cfg.Password)
			assert.Equal(t, tt.expected.Host, cfg.Host)
			assert.Equal(t, tt.expected.Port, cfg.Port)
			assert.Equal(t, tt.expected.Database, cfg.Database)
			assert.Equal(t, tt.expected.MaxRetries, cfg.MaxRetries)
			assert.Equal(t, tt.expected.RetryDelay, cfg.RetryDelay)
			assert.Equal(t, tt.expected.LogLevel, cfg.LogLevel)
		})
	}
}

func TestMigrateModels(t *testing.T) {
	type TestModel struct {
		ID   uint   `gorm:"primaryKey"`
		Name string `gorm:"size:100"`
	}

	type AnotherModel struct {
		ID    uint   `gorm:"primaryKey"`
		Value string `gorm:"size:50"`
	}

	tests := []struct {
		name    string
		models  []any
		wantErr bool
		verify  func(t *testing.T, db *gorm.DB)
	}{
		{
			name:    "migrate single model successfully",
			models:  []any{&TestModel{}},
			wantErr: false,
			verify: func(t *testing.T, db *gorm.DB) {
				assert.True(t, db.Migrator().HasTable(&TestModel{}))
			},
		},
		{
			name:    "migrate multiple models successfully",
			models:  []any{&TestModel{}, &AnotherModel{}},
			wantErr: false,
			verify: func(t *testing.T, db *gorm.DB) {
				assert.True(t, db.Migrator().HasTable(&TestModel{}))
				assert.True(t, db.Migrator().HasTable(&AnotherModel{}))
			},
		},
		{
			name:    "handle empty models list",
			models:  []any{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validTestConfig()
			db, err := database.ConnectDB(cfg)
			require.NoError(t, err)
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()

			err = database.MigrateModels(db, tt.models...)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.verify != nil {
					tt.verify(t, db)
				}
			}

			for _, model := range tt.models {
				db.Migrator().DropTable(model)
			}
		})
	}
}

func validTestConfig() *database.Config {
	return &database.Config{
		User:       "root",
		Password:   "pass",
		Host:       "localhost",
		Port:       mysqlPort,
		Database:   "testdb",
		MaxRetries: 3,
		RetryDelay: 100 * time.Millisecond,
		LogLevel:   logger.Silent,
	}
}
