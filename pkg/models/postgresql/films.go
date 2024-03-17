package postgresql

import (
	"database/sql"
	"errors"
	"films-api/pkg/models"
	"github.com/lib/pq"
)

type FilmModel struct {
	DB *sql.DB
}

func (m *FilmModel) Insert(name, description, date string, rating int, actors []string) (int, error) {

	stmt := `INSERT INTO films (name, description, date, rating)
    VALUES($1,$2,$3,$4) RETURNING id`

	lastInsertId := 0
	err := m.DB.QueryRow(stmt, name, description, date, rating).Scan(&lastInsertId)

	id := lastInsertId
	if err != nil {
		return 0, err
	}

	stmt = `SELECT id FROM actors WHERE name = ANY ($1)`
	rows, err := m.DB.Query(stmt, pq.Array(actors))
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	if rows != nil {
		//var actors_id []int
		for rows.Next() {
			act := 0
			err = rows.Scan(&act)
			if err != nil {
				return 0, err
			}
			stmt = `INSERT INTO actors_films (actor_id, film_id) VALUES ($1,$2)`
			err = m.DB.QueryRow(stmt, act, id).Err()
			if err != nil {
				return 0, err
			}
			//actors_id = append(actors_id, act)

		}

		/*
			stmt = `INSERT INTO actors_films (actor_id, film_id) SELECT DISTINCT $1, unnest($2)`

			err = m.DB.QueryRow(stmt, id, pq.Array(actors_id)).Err()
			if err != nil {
				return 0, err
			}

		*/

	}

	return id, nil
}

func (m *FilmModel) Get_By_Id(id int) (*models.Film, error) {
	stmt := `SELECT id, name, description, date, rating FROM films
    WHERE id = $1`

	row := m.DB.QueryRow(stmt, id)

	fil := &models.Film{}

	err := row.Scan(&fil.Id, &fil.Name, &fil.Description, &fil.Date, &fil.Rating)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return fil, nil
}

func (m *FilmModel) Get_Films(orderBy string) (*models.Films, error) {
	stmt := `SELECT id, name, description, date, rating FROM films ORDER BY rating DESC`
	if orderBy == "name" {
		stmt = `SELECT id, name, description, date, rating FROM films ORDER BY name DESC`
	} else if orderBy == "date" {
		stmt = `SELECT id, name, description, date, rating FROM films ORDER BY date DESC`
	}

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	films := &models.Films{}
	films.Sort = orderBy

	var films_arr []models.Film

	for rows.Next() {
		fil := models.Film{}
		err = rows.Scan(&fil.Id, &fil.Name, &fil.Description, &fil.Date, &fil.Rating)

		if err != nil {
			return nil, err
		}

		films_arr = append(films_arr, fil)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	films.Films = films_arr
	films.Count = films.Count

	return films, nil
}

func (m *FilmModel) Delete(id int) error {
	stmt := `DELETE FROM films WHERE id=$1`

	err := m.DB.QueryRow(stmt, id).Err()

	if err != nil {
		return err
	}

	return nil
}

func (m *FilmModel) Change(name, description, date string, rating int, actors []string, id int) error {
	if len(name) > 0 && len(description) > 0 && len(date) > 0 && rating > 0 {
		stmt := `UPDATE films SET name = $1, description = $2, rating = $3
                WHERE id = $4;`
		err := m.DB.QueryRow(stmt, name, description, date, rating, id).Err()

		if err != nil {
			return err
		}
	} else {
		if len(name) > 0 {
			err := m.Change_name(id, name)

			if err != nil {
				return err
			}
		}
		if len(description) > 0 {
			err := m.Change_description(id, description)

			if err != nil {
				return err
			}
		}
		if len(date) > 0 {
			err := m.Change_date(id, date)
			if err != nil {
				return err
			}
		}
		if rating > 0 {
			err := m.Change_rating(id, rating)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *FilmModel) Change_name(id int, name string) error {
	stmt := `UPDATE films SET name = $1
                WHERE id = $2;`
	err := m.DB.QueryRow(stmt, name, id).Err()

	if err != nil {
		return err
	}

	return nil
}

func (m *FilmModel) Change_description(id int, description string) error {
	stmt := `UPDATE films SET description = $1
                WHERE id = $2;`
	err := m.DB.QueryRow(stmt, description, id).Err()

	if err != nil {
		return err
	}

	return nil
}

func (m *FilmModel) Change_date(id int, date string) error {
	stmt := `UPDATE films SET date = $1
                WHERE id = $2;`
	err := m.DB.QueryRow(stmt, date, id).Err()
	if err != nil {
		return err
	}

	return nil
}

func (m *FilmModel) Change_rating(id, rating int) error {
	stmt := `UPDATE films SET rating = $1
                WHERE id = $2;`
	err := m.DB.QueryRow(stmt, rating, id).Err()

	if err != nil {
		return err
	}

	return nil
}
