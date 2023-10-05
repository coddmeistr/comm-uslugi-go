package initializers

import (
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb() {
	sqlDB, err := sql.Open("mysql", "root:12345@/main?parseTime=true")
	if err != nil {
		panic("Failed to connect to database.")
	}
	var err2 error
	DB, err2 = gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err2 != nil {
		panic("Failed to connect sql connection to gorm database.")
	}
}
