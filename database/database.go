package database

import (
	"fmt"

	//"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
	//"database/sql"
)

//Config is ...
type Config struct {
	User       string
	Password   string
	Host       string
	Port       int
	Name       string
	DisableTLS bool
}

//PsqlInfo is use to export psql information
func PsqlInfo() string {
	const (
		host     = "172.18.0.2" //172.18.0.2
		port     = 5432
		user     = "postgres"
		password = "0135"
		dbname   = "omega"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	return psqlInfo
}

//ReservePsqlInfo is use to export psql information
func ReservePsqlInfo() string {
	const (
		host     = "142.93.177.152" 
		port     = 5432
		user     = "postgres"
		password = "0135"
		dbname   = "omega"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	return psqlInfo
}
