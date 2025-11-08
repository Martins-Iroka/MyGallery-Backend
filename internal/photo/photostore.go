package photopost

import "database/sql"

type PhotoPost struct {
	ID           int64
	Photographer string
	Original     string
	Large2x      string
	Large        string
	Medium       string
	Small        string
	Portrait     string
	Landscape    string
	Tiny         string
}

type PhotoStore struct {
	Db *sql.DB
}
