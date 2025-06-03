package repository

import (
	"context"
	"diplom/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryInterface interface {
	GetByID(id uuid.UUID) (*domain.User, error)
	Login(login string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uuid.UUID) error
	Create(user *domain.User) (uuid.UUID, error)
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) (uuid.UUID, error) {
	id := uuid.New()
	query := `INSERT INTO diplom.users (id, name, login, password, role) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(context.Background(), query, id, user.Name, user.Login, user.Password, user.Role)
	return id, err
}

func (r *UserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, name, login, password, role FROM diplom.users WHERE id = $1`
	user := &domain.User{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.Role)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	query := `UPDATE diplom.users SET name = $2, login = $3, password = $4, role = $5 WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, user.ID, user.Name, user.Login, user.Password, user.Role)
	return err
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM diplom.users WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

//func (r *UserRepository) GetAll() ([]domain.User, error) {
//	query := `SELECT id, name, login, password, role FROM users`
//	rows, err := r.db.Query(context.Background(), query)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var users []domain.User
//	for rows.Next() {
//		user := domain.User{}
//		if err := rows.Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.Role); err != nil {
//			return nil, err
//		}
//		users = append(users, user)
//	}
//	return users, nil
//}

func (r *UserRepository) Login(login string) (*domain.User, error) {
	query := `SELECT id, name, login, password, role FROM diplom.users WHERE login = $1`
	user := &domain.User{}
	err := r.db.QueryRow(context.Background(), query, login).Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.Role)

	if err != nil {
		return nil, err
	}
	return user, nil
}
