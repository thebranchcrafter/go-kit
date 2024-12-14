package module_generator_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thebranchcrafter/go-kit/pkg/command/module_generator"
	"os"
	"testing"
)

func TestFileCreatorForTestModule(t *testing.T) {
	fc := module_generator.NewFileCreator("./templates/")

	moduleName := "test_module"
	err := createTestFolder()
	assert.NoError(t, err)

	err = fc.Create(moduleName, "./test", "github.com/thebranchcrafter/go-kit")
	assert.NoError(t, err)

	t.Run("Check aggregate file", func(t *testing.T) {
		expectedAggregateFileContent := getExpectedAggregateFileContent()
		aggregateFile, err := os.ReadFile(fmt.Sprintf("./test/domain/%s.go", moduleName))
		assert.NoError(t, err)
		assert.Equal(t, expectedAggregateFileContent, aggregateFile)
	})

	t.Run("Check domain repository file", func(t *testing.T) {
		expectedDomainRepository := getExpectedDomainRepository()
		domainRepositoryFile, err := os.ReadFile(fmt.Sprintf("./test/domain/%s_repository.go", moduleName))
		assert.NoError(t, err)
		assert.Equal(t, expectedDomainRepository, domainRepositoryFile)
	})

	t.Run("Check infrastructure repository file", func(t *testing.T) {
		expectedInfrastructureRepository := getExpectedInfrastructureRepository()
		infrastructureRepositoryFile, err := os.ReadFile(fmt.Sprintf("./test/infrastructure/postgres_%s_repository.go", moduleName))
		assert.NoError(t, err)

		assert.Equal(t, expectedInfrastructureRepository, infrastructureRepositoryFile)
	})

	t.Run("Check module file", func(t *testing.T) {
		expectedModuleFile := getExpectedModuleFile()
		moduleFile, err := os.ReadFile("./test/module.go")
		assert.NoError(t, err)

		assert.Equal(t, expectedModuleFile, moduleFile)
	})

	err = deleteTestFolder()
	assert.NoError(t, err)
}

func getExpectedModuleFile() []byte {
	return []byte(`package test_module

import "github.com/thebranchcrafter/go-kit/pkg/infrastructure/app"

type TestModuleModule struct {
	*app.BaseModule
}

func (u *TestModuleModule) Routes() []app.Route {
	return []app.Route{

	}
}

func InitTestModuleModule(d app.CommonDependencies) *TestModuleModule {
	um := &TestModuleModule{
		&app.BaseModule{
			CommonDependencies: d,
		},
	}

	return um
}

func (u *TestModuleModule) Name() string {
	return "test_module_module"
}
`)
}

func getExpectedDomainRepository() []byte {
	return []byte(`package test_module_domain

import "context"

type TestModuleRepository interface {
	// Save a new TestModule
	Save(ctx context.Context, tm *TestModule) error

	// GetByID a TestModule by its ID
	GetByID(ctx context.Context, id TestModuleId) (*TestModule, error)

	// GetAll TestModule with optional filters (if necessary)
	GetAll(ctx context.Context, filters map[string]interface{}) ([]*TestModule, error)

	// Update an existing TestModule
	Update(ctx context.Context, tm *TestModule) error

	// Delete a TestModule by its ID
	Delete(ctx context.Context, id TestModuleId) error
}
`)
}

func getExpectedInfrastructureRepository() []byte {
	return []byte(`package test_module_infrastructure

import (
	"context"
	test_module_domain "github.com/thebranchcrafter/go-kit/internal/app/module/test_module/domain"
	"github.com/Masterminds/squirrel"
    "github.com/jackc/pgx/v4/pgxpool"
    "time"
)

var columns = []string{"id", "created_at", "updated_at"}

// PostgresTest_moduleRepository is a Postgres implementation of Test_moduleRepository using Squirrel and pgxpool.
type PostgresTest_moduleRepository struct {
	tableName string
	DB        *pgxpool.Pool
}

// NewPostgresTest_moduleRepository initializes a new Postgres Test_module repository with a connection pool.
func NewPostgresTest_moduleRepository(poolConfig string) (*PostgresTest_moduleRepository, error) {
	pool, err := pgxpool.Connect(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	return &PostgresTest_moduleRepository{
		DB: pool,
	}, nil
}

// Save creates or updates a Test_module record.
func (r *PostgresTest_moduleRepository) Save(ctx context.Context, tm *test_module_domain.TestModule) error {
	query := squirrel.Insert(r.tableName).
		Columns(columns...).
		Values(tm.Id(), tm.CreatedAt(), tm.UpdatedAt()).
		Suffix("ON CONFLICT (id) DO UPDATE SET updated_at = EXCLUDED.updated_at").
		PlaceholderFormat(squirrel.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.DB.Exec(ctx, sqlQuery, args...)
	return err
}

// GetByID retrieves a Test_module by its ID.
func (r *PostgresTest_moduleRepository) GetByID(ctx context.Context, id test_module_domain.TestModuleId) (*test_module_domain.TestModule, error) {
	queryBuilder := squirrel.Select(columns...).
		From(r.tableName).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sqlQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.DB.QueryRow(ctx, sqlQuery, args...)
	return scanRow(row)
}

// GetAll retrieves all Test_module records, with optional filters.
func (r *PostgresTest_moduleRepository) GetAll(ctx context.Context, filters map[string]interface{}) ([]*test_module_domain.TestModule, error) {
	query := squirrel.Select(columns...).
		From(r.tableName).
		PlaceholderFormat(squirrel.Dollar)

	for key, value := range filters {
		query = query.Where(squirrel.Eq{key: value})
	}

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*test_module_domain.TestModule
	for rows.Next() {
		test_module, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, test_module)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// Update modifies an existing Test_module record.
func (r *PostgresTest_moduleRepository) Update(ctx context.Context, tm *test_module_domain.TestModule) error {
	query := squirrel.Update(r.tableName).
		Set("updated_at", tm.UpdatedAt()).
		Where(squirrel.Eq{"id": tm.Id()}).
		PlaceholderFormat(squirrel.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.DB.Exec(ctx, sqlQuery, args...)
	return err
}

// Delete removes a Test_module by its ID.
func (r *PostgresTest_moduleRepository) Delete(ctx context.Context, id test_module_domain.TestModuleId) error {
	query := squirrel.Delete(r.tableName).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.DB.Exec(ctx, sqlQuery, args...)
	return err
}

// scanRow scans a row into a Test_module domain object.
func scanRow(row squirrel.RowScanner) (*test_module_domain.TestModule, error) {
	var (
		id        string
		createdAt time.Time
		updatedAt time.Time
	)

	err := row.Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	return test_module_domain.FromPrimitives(
		id,
		createdAt,
		updatedAt,
	), nil
}`)
}

func getExpectedAggregateFileContent() []byte {
	return []byte(`package test_module_domain

import (
	"time"
)

type TestModuleId string

type TestModule struct {
	id                TestModuleId
	createdAt         time.Time
	updatedAt         time.Time
}

func (t *TestModule) Id() TestModuleId {
	return t.id
}

func (t *TestModule) CreatedAt() time.Time {
	return t.createdAt
}

func (t *TestModule) UpdatedAt() time.Time {
	return t.updatedAt
}

func FromPrimitives(id string, createdAt, updatedAt time.Time) *TestModule {
	return &TestModule{
		id:        TestModuleId(id),
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}`)
}

func deleteTestFolder() error {
	return os.RemoveAll("./test")
}

func createTestFolder() error {
	deleteTestFolder()
	err := os.Mkdir("test", 0777)
	if err != nil {
		return err
	}

	return nil
}
