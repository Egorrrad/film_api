package postgresql

import (
	"database/sql"
	"errors"
	"films-api/pkg/models"
)

type ActorModel struct {
	DB *sql.DB
}

func (m *ActorModel) Insert(name1, gender, birthday string) (int, error) {

	stmt := `INSERT INTO actors (name, gender, birhday) VALUES($1,$2,$3) RETURNING id`

	lastInsertId := 0
	err := m.DB.QueryRow(stmt, name1, gender, birthday).Scan(&lastInsertId)

	id := lastInsertId
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *ActorModel) Get_By_Id(id int) (*models.Actor, error) {
	stmt := `SELECT id, name, gender, birhday FROM actors WHERE id = $1`

	row := m.DB.QueryRow(stmt, id)

	act := &models.Actor{}

	err := row.Scan(&act.Id, &act.Name, &act.Gender, &act.Birthday)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	stmt = `SELECT t1.name FROM films t1 JOIN actors_films t3
		    ON t1.id = t3.film_id where t3.actor_id = $1`

	rows, err := m.DB.Query(stmt, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var films []string

	for rows.Next() {
		name := ""
		err = rows.Scan(&name)

		if err != nil {
			return nil, err
		}

		films = append(films, name)
	}

	act.Films = films

	return act, nil
}

func (m *ActorModel) Get_Actors() (*models.Actors, error) {
	stmt := `SELECT t1.id, t1.name, t1.gender, t1.birhday ,t2.name FROM actors t1 JOIN actors_films t3
                                  ON t1.id = t3.actor_id JOIN films t2 ON t3.film_id = t2.id
                                               ORDER BY t1.id`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	actors := &models.Actors{}

	var actorsf []models.Actor

	for rows.Next() {
		act := models.Actor{Films: []string{}}
		film := ""
		err = rows.Scan(&act.Id, &act.Name, &act.Gender, &act.Birthday, &film)
		if err != nil {
			return nil, err
		}
		l := len(actorsf)

		if l > 0 && actorsf[l-1].Id == act.Id {
			actorsf[l-1].Films = append(actorsf[l-1].Films, film)
		} else {
			act.Films = append(act.Films, film)
			actorsf = append(actorsf, act)
		}

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	actors.Sort = "by names"
	actors.Count = len(actorsf)
	actors.Actors = actorsf

	return actors, nil
}

func (m *ActorModel) Change(name1, gender, birthday string, id int) error {

	if len(name1) > 0 && len(gender) > 0 && len(birthday) > 0 {
		stmt := `UPDATE  actors SET name = $1, gender = $2, birhday = $3
                WHERE id = $4;`
		err := m.DB.QueryRow(stmt, name1, gender, birthday, id).Err()
		if err != nil {
			return err
		}
	} else {
		if len(name1) > 0 {
			err := m.Change_name(name1, id)
			if err != nil {
				return err
			}
		}
		if len(gender) > 0 {
			err := m.Change_gender(gender, id)
			if err != nil {
				return err
			}
		}
		if len(birthday) > 0 {
			err := m.Change_birthday(birthday, id)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *ActorModel) Change_name(name1 string, id int) error {
	stmt := `UPDATE  actors SET name = $1
                WHERE id = $2;`

	err := m.DB.QueryRow(stmt, name1, id).Err()

	if err != nil {
		return err
	}

	return nil
}

func (m *ActorModel) Change_gender(gender string, id int) error {
	stmt := `UPDATE  actors SET gender = $1
                WHERE id = $2;`

	err := m.DB.QueryRow(stmt, gender, id).Err()

	if err != nil {
		return err
	}

	return nil
}

func (m *ActorModel) Change_birthday(birhday string, id int) error {
	stmt := `UPDATE  actors SET birhday = $1
                WHERE id = $2;`

	err := m.DB.QueryRow(stmt, birhday, id).Err()

	if err != nil {
		return err
	}

	return nil
}

func (m *ActorModel) Delete(id int) error {
	stmt := `DELETE FROM actors WHERE id=$1`

	err := m.DB.QueryRow(stmt, id).Err()

	if err != nil {
		return err
	}

	return nil
}
