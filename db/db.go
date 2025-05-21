package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
)

type Sql struct {
	Db       *sqlx.DB
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

func (s *Sql) Connect() {
	dataSoure := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable ",
		s.Host, s.Port, s.Username, s.Password, s.Dbname)
	s.Db = sqlx.MustConnect("postgres", dataSoure)
	if err := s.Db.Ping(); err != nil {
		log.Error(err.Error())
		return
	}
	fmt.Println("Successfully connected to database")
}

func (s *Sql) Close() {
	s.Db.Close()
}
