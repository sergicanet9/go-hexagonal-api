package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/sergicanet9/go-hexagonal-api/core/domain"
)

// UserRepository struct of an user repository for postgres
type UserRepository struct {
	PostgresRepository
}

// NewUserRepository creates a user repository for postgres
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		PostgresRepository{
			db,
		},
	}
}

func (r *UserRepository) Create(ctx context.Context, entity interface{}) (string, error) {
	q := `
    INSERT INTO users (name, surnames, email, password_hash, claims, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id;
    `

	u := entity.(domain.User)
	row := r.db.QueryRowContext(
		ctx, q, u.Name, u.Surnames, u.Email, u.PasswordHash, pq.Array(u.Claims), u.CreatedAt, u.UpdatedAt,
	)

	err := row.Scan(&u.ID)
	if err != nil {
		return "", err
	}

	return u.ID, nil
}

func (r *UserRepository) Get(ctx context.Context, filter map[string]interface{}, skip, take *int) ([]interface{}, error) {
	var where string
	for k, v := range filter {
		if where == "" {
			where = "WHERE"
		} else {
			where = fmt.Sprintf("%s AND", where)
		}
		where = fmt.Sprintf("%s %s = '%v'", where, k, v)
	}
	if skip != nil {
		where = fmt.Sprintf("%s OFFSET %d", where, *skip)
	}
	if take != nil {
		where = fmt.Sprintf("%s LIMIT %d", where, *take)
	}

	q := fmt.Sprintf(`
	SELECT id, name, surnames, email, password_hash, claims, created_at, updated_at
	    FROM users %s;
	`, where)

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []interface{}
	for rows.Next() {
		var u domain.User
		err = rows.Scan(&u.ID, &u.Name, &u.Surnames, &u.Email, &u.PasswordHash, pq.Array(&u.Claims), &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

func (r *UserRepository) GetByID(ctx context.Context, ID string) (interface{}, error) {
	q := `
    SELECT id, name, surnames, email, password_hash, claims, created_at, updated_at
        FROM users WHERE id = $1;
    `

	row := r.db.QueryRowContext(ctx, q, ID)

	var u domain.User
	err := row.Scan(&u.ID, &u.Name, &u.Surnames, &u.Email, &u.PasswordHash, pq.Array(&u.Claims), &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, ID string, entity interface{}) error {
	q := `
	UPDATE users set name=$1, surnames=$2, email=$3, password_hash=$4, claims=$5, updated_at=$6
	    WHERE id=$7;
	`

	u := entity.(domain.User)
	result, err := r.db.ExecContext(
		ctx, q, u.Name, u.Surnames, u.Email, u.PasswordHash, pq.Array(u.Claims), u.UpdatedAt, ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Println("RowsAffected Error", err)
	}
	if rows < 1 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, ID string) error {
	q := `DELETE FROM users WHERE id=$1;`

	result, err := r.db.ExecContext(ctx, q, ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Println("RowsAffected Error", err)
	}
	if rows < 1 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserRepository) InsertMany(ctx context.Context, entities []interface{}) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// `tx` is an instance of `*sql.Tx` through which we can execute our queries

	for _, entity := range entities {
		u := entity.(domain.User)

		// Here, the query is executed on the transaction instance, and not applied to the database yet
		_, err = tx.ExecContext(ctx, "INSERT INTO users (name, surnames, email, password_hash, claims, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)", u.Name, u.Surnames, u.Email, u.PasswordHash, pq.Array(u.Claims), u.CreatedAt, u.UpdatedAt)
		if err != nil {
			// Incase we find any error in the query execution, rollback the transaction
			tx.Rollback()
			return err
		}
	}

	// Finally, if no errors are recieved from the queries, commit the transaction
	// this applies the above changes to our database
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
