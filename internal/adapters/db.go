package adapters

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/krajeswaran/gostartup/internal/models"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

//DBAdapter - Struct to logically bind all the database related functions
type DBAdapter struct{}

//DBInit initializes DB connection
func (d *DBAdapter) DBInit() (*sqlx.DB, error) {
	// create postgres connection
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s sslrootcert=%s connect_timeout=%s",
		viper.GetString("DB_HOST"), viper.GetString("DB_PORT"),
		viper.GetString("DB_USER"), viper.GetString("DB_NAME"),
		viper.GetString("DB_SSL_MODE"), viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_SSL_CERT"), viper.GetString("DB_CONNECT_TIMEOUT"))

	// throws error if application_name is empty
	appName := viper.GetString("DB_APPLICATION_NAME")
	if appName != "" {
		connectionString = connectionString + fmt.Sprintf(" application_name=%s", appName)
	}

	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// DB parameters
	db.SetConnMaxLifetime(viper.GetDuration("DB_CONNECTION_LIFETIME"))
	db.SetMaxOpenConns(viper.GetInt("DB_MAX_OPEN_CONNECTIONS"))
	if viper.GetInt("DB_MAX_IDLE_CONNECTIONS") > 0 {
		db.SetMaxIdleConns(viper.GetInt("DB_MAX_IDLE_CONNECTIONS"))
	}

	return db, nil
}

//DeepStatus checks for a DB connection
func (d *DBAdapter) DeepStatus() error {
	if err := db.Ping(); err != nil {
		return errors.New("SERVICE_DB_DOWN")
	}
	return nil
}

//FetchUser Fetches a user based on id
func (d *DBAdapter) FetchUser(id string) (*models.User, error) {
	var user models.User
	if err := db.Get(&user, "SELECT * FROM user WHERE id=?", id); err != nil {
		return nil, err
	}

	return &user, nil
}

//CreateUser Creates user given a user name
func (d *DBAdapter) CreateUser(name string) (*models.User, error) {
	row, err := db.Queryx("INSERT INTO user (name) VALUES (?)", name)
	if err != nil {
		return nil, err
	}
	var user models.User
	if err = row.StructScan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

//InitSchema Initializes schema for app start
func (d *DBAdapter) InitSchema() error {
	schema := `DROP TABLE IF EXISTS public.user; CREATE TABLE public.user (
    created_at timestamp with time zone DEFAULT timezone('UTC'::text, now()) NOT NULL,
    modified_at timestamp with time zone DEFAULT timezone('UTC'::text, now()) NOT NULL,
    id serial NOT NULL,
	name text NOT NULL
    );`

	if _, err := db.Exec(schema); err != nil {
		return err
	}

	return nil
}
