package main

import (
	"embed"
	"log"
	"os"

	"github.com/joho/godotenv"

	mdb "github.com/Whadislov/TTCompanion/internal/my_db"
	mu "github.com/Whadislov/TTCompanion/internal/my_frontend/my_ui"
)

//go:embed translation/*
var translations embed.FS

func main() {

	// Load translations
	mu.InitTranslations(translations)

	// Load env variables
	err := godotenv.Load("credentials.env")
	if err != nil {
		log.Fatal("Cannot load variables from .env")
	}

	mdb.SetPsqlInfo(os.Getenv("LOCAL_DB_LINK"))
	mdb.SetDBName(os.Getenv("DB_NAME"))
	mu.Display("local")
}
