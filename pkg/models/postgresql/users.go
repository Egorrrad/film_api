package postgresql

import (
	"database/sql"
	"errors"
	"films-api/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Get(api_key string) (*models.User, error) {
	stmt := `SELECT role, api_key FROM users WHERE api_key = $1`

	row := m.DB.QueryRow(stmt, api_key)

	usr := &models.User{}

	err := row.Scan(&usr.Role, &usr.Api_key)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return usr, nil
}
