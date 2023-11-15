package postgres

import (
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	config "github.com/ryanadiputraa/spendr-backend/configs"
)

const (
	maxOpenConns    = 60
	connMaxLifeTime = 120
	maxIdleConn     = 30
	connMaxIdleTime = 20
)

func NewDB(c *config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.DBName,
		c.Postgres.SSLMode,
		c.Postgres.Password,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifeTime * time.Second)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)

	err = pingWithRetry(db, 3*time.Second, 5)
	return db, err
}

func pingWithRetry(db *sqlx.DB, interval time.Duration, maxRetries int) (err error) {
	for i := 0; i < maxRetries; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		time.Sleep(interval)
	}
	return
}
