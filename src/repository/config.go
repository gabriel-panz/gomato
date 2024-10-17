package repository

import (
	"database/sql"
	"errors"

	"github.com/gabriel-panz/gomato/db"
	localTypes "github.com/gabriel-panz/gomato/types"
)

type TimerConfigRepository struct {
	db *sql.DB
}

var TimerConfigRepo *TimerConfigRepository = nil

func GetTimerRepo() *TimerConfigRepository {
	if TimerConfigRepo == nil {
		InitConfigRepo()
	}
	return TimerConfigRepo
}

func InitConfigRepo() {
	db := db.GetDb()
	TimerConfigRepo = &TimerConfigRepository{
		db: db}
}

func (repo *TimerConfigRepository) GetAllConfigs() ([]localTypes.TimerConfig, error) {
	res := make([]localTypes.TimerConfig, 0)
	q := `
	SELECT id, name, work_time, pause_time, notification_level FROM configuration
	ORDER BY name
	`
	rs, err := repo.db.Query(q)
	if err != nil {
		return nil, err
	}
	for rs.Next() {
		c := &localTypes.TimerConfig{}
		err := rs.Scan(&c.Id, &c.Name, &c.WorkTime, &c.PauseTime, &c.NotificationLevel)
		if err != nil {
			return nil, err
		}
		res = append(res, *c)
	}

	return res, nil
}

func (repo *TimerConfigRepository) GetConfig(id int) (*localTypes.TimerConfig, error) {
	c := &localTypes.TimerConfig{}
	q := `
SELECT id, name, work_time, pause_time, notification_level FROM configuration
WHERE id = ? LIMIT 1
	`
	r := repo.db.QueryRow(q, id)
	err := r.Scan(&c.Id, &c.Name, &c.WorkTime, &c.PauseTime, &c.NotificationLevel)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (repo *TimerConfigRepository) GetConfigByName(name string) (*localTypes.TimerConfig, error) {
	c := &localTypes.TimerConfig{}
	q := `
SELECT id, name, work_time, pause_time, notification_level FROM configuration
WHERE name = ? LIMIT 1
	`
	r := repo.db.QueryRow(q, name)
	err := r.Scan(&c.Id, &c.Name, &c.WorkTime, &c.PauseTime, &c.NotificationLevel)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrConfigNotFound
		}
		return nil, err
	}
	return c, nil
}

func (repo *TimerConfigRepository) GetDefaultConfig() (*localTypes.TimerConfig, error) {
	c := &localTypes.TimerConfig{}
	q := `
SELECT id, name, work_time, pause_time, notification_level FROM configuration
WHERE id IN (SELECT default_config FROM defaults LIMIT 1)
	`
	r := repo.db.QueryRow(q)
	err := r.Scan(&c.Id, &c.Name, &c.WorkTime, &c.PauseTime, &c.NotificationLevel)
	if err != nil {
		if err == sql.ErrNoRows {
			return repo.recoverDefaultConfig(), nil
		}
		return nil, err
	}
	return c, nil
}

func (repo *TimerConfigRepository) recoverDefaultConfig() *localTypes.TimerConfig {
	q := `
INSERT INTO configuration (
    id, name, work_time, pause_time, notification_level
) VALUES (
  	1, 'default', 25000000000, 5000000000, 2
) ON CONFLICT DO UPDATE SET 
	name='default',
	work_time=25000000000,
	pause_time=5000000000,
	notification_level=2
RETURNING *`

	r, err := repo.db.Query(q)
	if err != nil {
		panic(err)
	}

	c := &localTypes.TimerConfig{}
	err = r.Scan(&c.Id, &c.Name, &c.WorkTime, &c.PauseTime, &c.NotificationLevel)
	if err != nil {
		panic(err)
	}

	return nil
}

func (repo *TimerConfigRepository) InsertConfig(c *localTypes.TimerConfig) error {
	q := `
INSERT INTO configuration (
    name, work_time, pause_time, notification_level
) VALUES (
  ?, ?, ?, ?
)`

	r, err := repo.db.Exec(q, &c.Name, &c.WorkTime, &c.PauseTime, &c.NotificationLevel)
	if err != nil {
		return err
	}

	c.Id, err = r.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (repo *TimerConfigRepository) UpdateConfig(c *localTypes.TimerConfig) error {
	q := `
UPDATE configuration SET
    name=?,
	work_time=?,
	pause_time=?,
	notification_level=?
WHERE
	id=?
`

	_, err := repo.db.Exec(q, &c.Name, &c.WorkTime, &c.PauseTime, &c.NotificationLevel, &c.Id)
	if err != nil {
		return err
	}

	return nil
}

func (repo *TimerConfigRepository) DeleteConfig(c *localTypes.TimerConfig) error {
	q := `DELETE configuration WHERE id=?`

	_, err := repo.db.Exec(q, &c.Id)
	if err != nil {
		return err
	}

	return nil
}

func (repo *TimerConfigRepository) SetDefaultConfig(id int) (*localTypes.TimerConfig, error) {
	c := &localTypes.TimerConfig{}
	q := `UPDATE defaults SET default_config = ?`
	_, err := repo.db.Exec(q, id)
	if err != nil {
		return nil, err
	}
	return c, nil
}

var ErrConfigNotFound error = errors.New("configuration not found")
