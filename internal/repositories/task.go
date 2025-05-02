package repositories

import (
	"context"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/cabrn3t/todo-list-server/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Repository struct {
	pool *pgxpool.Pool
}

func (t *Repository) CreateTask(ctx context.Context, task *models.Task) error {
	query := psql.Insert("tasks").Columns("title", "description", "status").Values(task.Title, task.Description.String, task.Status)
	sql, args, err := query.ToSql()
	if err != nil {
		slog.Error("error while forming request", slog.Any("err", err))
		return err
	}
	_, err = t.pool.Exec(ctx, sql, args...)
	if err != nil {
		slog.Error("error while sending request", slog.Any("err", err))
		return err
	}
	slog.Info("task successfully created")
	return nil
}

func (t *Repository) DeleteTask(ctx context.Context, id int64) error {
	query := psql.Delete("tasks").Where(squirrel.Eq{"id": id}).Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		slog.Error("error while forming request", slog.Any("err", err))
		return err
	}
	var deletedID int64
	err = t.pool.QueryRow(ctx, sql, args...).Scan(&deletedID)
	if err != nil {
		slog.Error("there are no tasks with selected id", slog.Any("err", err))
		return err
	}

	slog.Info("Task deleted successfully", "id", deletedID)
	return nil
}

func (t *Repository) GetTasks(ctx context.Context) ([]models.Task, error) {
	query := psql.Select("*").From("tasks")
	sql, args, err := query.ToSql()
	if err != nil {
		slog.Error("error while forming request", slog.Any("err", err))
		return []models.Task{}, err
	}
	rows, err := t.pool.Query(ctx, sql, args...)
	if err != nil {
		slog.Error("error while sending request", slog.Any("err", err))
		return []models.Task{}, err
	}
	defer rows.Close()
	var taskSlice []models.Task

	for rows.Next() {
		var t models.Task

		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			slog.Error("error while reading tasks", slog.Any("err", err))
			return []models.Task{}, err
		}
		taskSlice = append(taskSlice, t)
	}

	return taskSlice, nil
}

func (t *Repository) UpdateTask(ctx context.Context, task *models.Task) error {
	query := psql.Update("tasks").Where(squirrel.Eq{"id": task.ID})
	if task.Title != "" {
		query = query.Set("title", task.Title)
	}
	if task.Description.Valid {
		query = query.Set("description", task.Description.String)
	}
	if task.Status != "" {
		query = query.Set("status", task.Status)
	}
	query = query.Set("updated_at", time.Now()).Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		slog.Error("error while forming request", slog.Any("err", err))
		return err
	}

	var updatedID int64
	err = t.pool.QueryRow(ctx, sql, args...).Scan(&updatedID)
	if err != nil {
		slog.Error("there are no tasks with selected id", slog.Any("err", err))
		return err
	}
	slog.Info("task successfully updated")
	return nil
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}
