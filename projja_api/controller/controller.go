package controller

import "database/sql"

type Controller struct {
	DB *sql.DB
}
