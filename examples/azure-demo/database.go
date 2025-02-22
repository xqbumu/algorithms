package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/phuslu/log"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func RunDatabase() {
	dsn := "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
	if s := os.Getenv("AZURE_DSN"); s != "" {
		dsn = s
	}
	log.Printf("using %s", dsn)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Printf("conntected to %s", dsn)

	{
		var rows []map[string]any
		if err := db.Table("information_schema.schemata").Select("schema_name").Scan(&rows).Error; err != nil {
			panic(err)
		}
		for _, row := range rows {
			log.Printf("%s", row["schema_name"])
		}
	}

	{
		var rows []map[string]any
		if err := db.Table("sys.tables").Select("name").Scan(&rows).Error; err != nil {
			panic(err)
		}
		for _, row := range rows {
			log.Printf("%s", row["name"])
		}
	}

	{

		if err := db.AutoMigrate(&User{}); err != nil {
			panic(err)
		}
		{
			user := &User{Name: "Bob"}
			db.Create(user)
			log.Printf("created user %s", user.Name)
		}
	}
}

type User struct {
	gorm.Model
	Name string
}

func (User) TableName() string { return "playground.user" }
