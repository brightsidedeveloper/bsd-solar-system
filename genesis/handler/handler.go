package handler

import (
	"database/sql"
	"solar-system/genesis/util"
)

type Handler struct {
	DB   *sql.DB
	JSON *util.JSON
}
