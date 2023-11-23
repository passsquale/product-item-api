package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

// NewPostgres returns DB
func NewPostgres(dsn, driver string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create database connection")

		return nil, err
	}
	return db, nil
}
