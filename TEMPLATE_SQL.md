Мы в тестконтейнерах поднимаем субд, создаем одну бд с миграцичми, а потом на каждый тест создаём свою бд, используя TEMPLATE в постгре и все тесты параллельно гоняем на изолированных бд. Считай юнит тесты

Еще можно юзать темпфс, чтобы не тратить время для записи на диск, и отключить у постгреса fsync, что мы нашли допустимым для тестов, и тогда время выполнения тестов, почти как с моками, а качество и надежность — в разы выше

```go
package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	tc "github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

const (
	baseDBName = "base_testdb"
	testDBUser = "testuser"
	testDBPass = "testpass"
)

func runContainer(ctx context.Context) (*tcpg.PostgresContainer, error) {
	return tcpg.RunContainer(ctx,
		tcpg.WithDatabase(baseDBName),
		tcpg.WithUsername(testDBUser),
		tcpg.WithPassword(testDBPass),
		tc.WithWaitStrategy(
			tcwait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
}

func setupBaseDatabase(ctx context.Context, connString string) error {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer conn.Close(ctx)

	// Выполняем миграции или создаем необходимую структуру в базовой БД
	_, err = conn.Exec(ctx, `
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	return nil
}

func createTestDatabase(ctx context.Context, baseConnString, testDBName string) (string, error) {
	conn, err := pgx.Connect(ctx, baseConnString)
	if err != nil {
		return "", fmt.Errorf("failed to connect to base database: %w", err)
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s", testDBName, baseDBName))
	if err != nil {
		return "", fmt.Errorf("failed to create test database: %w", err)
	}

	// Формируем новую строку подключения для тестовой БД
	testConnString := fmt.Sprintf("%s dbname=%s", baseConnString[:len(baseConnString)-len(baseDBName)], testDBName)
	return testConnString, nil
}

func TestExample(t *testing.T) {
	ctx := context.Background()

	postgres, err := runContainer(ctx)
	if err != nil {
		t.Fatalf("failed to start postgres container: %s", err)
	}
	defer func() {
		if err := postgres.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	connString, err := postgres.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	if err := setupBaseDatabase(ctx, connString); err != nil {
		t.Fatalf("failed to setup base database: %s", err)
	}

	// Создаем тестовую базу данных
	testDBName := fmt.Sprintf("%s_%s", baseDBName, t.Name())
	testConnString, err := createTestDatabase(ctx, connString, testDBName)
	if err != nil {
		t.Fatalf("failed to create test database: %s", err)
	}

	// Подключаемся к тестовой базе данных
	conn, err := pgx.Connect(ctx, testConnString)
	if err != nil {
		t.Fatalf("failed to connect to test database: %s", err)
	}
	defer conn.Close(ctx)

	// Выполняем тестовые операции
	_, err = conn.Exec(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)", "John Doe", "john@example.com")
	if err != nil {
		t.Fatalf("failed to insert user: %s", err)
	}

	var count int
	err = conn.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("failed to count users: %s", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 user, got %d", count)
	}
}

func main() {
	log.Println("This is a test file. Run 'go test' to execute the tests.")
}
```
но как нам переиспользовать контейнер?

```go
genericContainerReq := testcontainers.GenericContainerRequest{
  ContainerRequest: req,
  Started:          true,
  Logger:           logger,
  Reuse:            true, // правильный ответ
}
```

только флага `Reuse` нет в модуле Postgres, приходится изобретать синглтон:

```go
const (
	baseDBName     = "base_testdb"
	testDBUser     = "testuser"
	testDBPass     = "testpass"
	containerName  = "reusable-postgres-container"
	containerImage = "postgres:13-alpine"
)

var (
	postgresContainer *tcpg.PostgresContainer
	containerLock     sync.Mutex
)

func getOrCreatePostgresContainer(ctx context.Context) (*tcpg.PostgresContainer, error) {
	containerLock.Lock()
	defer containerLock.Unlock()

	if postgresContainer != nil {
		return postgresContainer, nil
	}

	req := tc.ContainerRequest{
		Image:        containerImage,
		ExposedPorts: []string{"5432/tcp"},
		Name:         containerName,
		Env: map[string]string{
			"POSTGRES_DB":       baseDBName,
			"POSTGRES_USER":     testDBUser,
			"POSTGRES_PASSWORD": testDBPass,
		},
		WaitingFor: tcwait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(5 * time.Second),
	}

	container, err := tcpg.RunContainer(ctx,
		tcpg.WithDatabase(baseDBName),
		tcpg.WithUsername(testDBUser),
		tcpg.WithPassword(testDBPass),
    tcpg.WithSSLMode("disable"),
		tc.WithContainerRequest(req),
	)

	if err != nil {
		return nil, err
	}

	postgresContainer = container
	return container, nil
}
```

Ожидаю: https://github.com/testcontainers/testcontainers-go/issues/2726

Пока так:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

const (
	baseDBName     = "base_testdb"
	testDBUser     = "testuser"
	testDBPass     = "testpass"
	containerName  = "reusable-postgres-container"
	containerImage = "postgres:13-alpine"
)

var (
	postgresContainer *tcpg.PostgresContainer
	containerLock     sync.Mutex
)

func getOrCreatePostgresContainer(ctx context.Context) (*tcpg.PostgresContainer, error) {
	containerLock.Lock()
	defer containerLock.Unlock()

	if postgresContainer != nil {
		return postgresContainer, nil
	}

	req := tc.ContainerRequest{
		Image:        containerImage,
		ExposedPorts: []string{"5432/tcp"},
		Name:         containerName,
		Env: map[string]string{
			"POSTGRES_DB":       baseDBName,
			"POSTGRES_USER":     testDBUser,
			"POSTGRES_PASSWORD": testDBPass,
		},
		WaitingFor: tcwait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(5 * time.Second),
	}

	container, err := tcpg.RunContainer(ctx,
		tcpg.WithDatabase(baseDBName),
		tcpg.WithUsername(testDBUser),
		tcpg.WithPassword(testDBPass),
    tcpg.WithSSLMode("disable"),
    tc.WithContainerRequest(req),
	)

	if err != nil {
		return nil, err
	}

	postgresContainer = container
	return container, nil
}

func setupBaseDatabase(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	return nil
}

func createTestDatabase(ctx context.Context, basePool *pgxpool.Pool, testDBName string) (*pgxpool.Pool, error) {
	_, err := basePool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
	if err != nil {
		return nil, fmt.Errorf("failed to drop existing test database: %w", err)
	}

	_, err = basePool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s", testDBName, baseDBName))
	if err != nil {
		return nil, fmt.Errorf("failed to create test database: %w", err)
	}

	connConfig, err := pgxpool.ParseConfig(basePool.Config().ConnString())
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	connConfig.ConnConfig.Database = testDBName
	testPool, err := pgxpool.ConnectConfig(ctx, connConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	return testPool, nil
}

func TestExample(t *testing.T) {
	ctx := context.Background()

	postgres, err := getOrCreatePostgresContainer(ctx)
	require.NoError(t, err, "failed to start or get postgres container")

	connString, err := postgres.ConnectionString(ctx)
	require.NoError(t, err, "failed to get connection string")

	basePool, err := pgxpool.New(ctx, connString+" sslmode=disable")
	require.NoError(t, err, "failed to create connection pool")
	defer basePool.Close()

	err = setupBaseDatabase(ctx, basePool)
	require.NoError(t, err, "failed to setup base database")

	testDBName := fmt.Sprintf("%s_%s", baseDBName, t.Name())
	testPool, err := createTestDatabase(ctx, basePool, testDBName)
	require.NoError(t, err, "failed to create test database")
	defer testPool.Close()

	_, err = testPool.Exec(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)", "John Doe", "john@example.com")
	require.NoError(t, err, "failed to insert user")

	var count int
	err = testPool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	require.NoError(t, err, "failed to count users")

	require.Equal(t, 1, count, "Expected 1 user, got %d", count)
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	
	_, err := getOrCreatePostgresContainer(ctx)
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}

	code := m.Run()

	if postgresContainer != nil {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("Could not stop container: %s", err)
		}
	}

	os.Exit(code)
}

func main() {
	log.Println("This is a test file. Run 'go test' to execute the tests.")
}
```

у твоего tcpg есть такая настройка `tcpg.WithInitScripts(filepath.Join(".", "testdata", "dev-db.sql"))`, Чтобы руками в тестах не делать setupBaseDatabase