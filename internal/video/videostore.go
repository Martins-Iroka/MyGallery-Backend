package video

import "database/sql"

type VideoPost struct {
	ID        int64
	Video_Url string
	Duration  int
}

type VideoDownloadFile struct {
	ID            int64
	Video_Post_Id int64
	Video_Link    string
	Video_Size    int32
}

type VideoStore struct {
	Db *sql.DB
}
