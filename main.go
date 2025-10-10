package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"test/infrastructure"
	"test/services"

	_ "modernc.org/sqlite"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: create | view | delete")
		os.Exit(1)
	}
	// Init DB
	db, err := sql.Open("sqlite", "./characters.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS characters (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT UNIQUE,
	race TEXT,
	class TEXT,
	level INTEGER,
	strength INTEGER,
	dexterity INTEGER,
	constitution INTEGER,
	intelligence INTEGER,
	wisdom INTEGER,
	charisma INTEGER,
	strength_mod INTEGER,
	dexterity_mod INTEGER,
	constitution_mod INTEGER,
	intelligence_mod INTEGER,
	wisdom_mod INTEGER,
	charisma_mod INTEGER,
	skills TEXT
)`)
	if err != nil {
		log.Fatal(err)
	}
	// Init service
	repo := infrastructure.NewSQLRepository(db)
	service := services.NewService(repo)

	// Init API data
	if _, err := os.Stat("classes.json"); errors.Is(err, os.ErrNotExist) {
		infrastructure.FetchClasses()
	}

	if _, err := os.Stat("races.json"); errors.Is(err, os.ErrNotExist) {
		infrastructure.FetchRaces()
	}

	switch os.Args[1] {
	case "create":
		createCmd := flag.NewFlagSet("create", flag.ExitOnError)
		name := createCmd.String("name", "", "character name (required)")
		class := createCmd.String("class", "", "character class (required)")
		race := createCmd.String("race", "", "character race (required)")
		level := createCmd.Int("level", 1, "character level")
		strength := createCmd.Int("str", 10, "strength")
		dexterity := createCmd.Int("dex", 10, "dexterity")
		constitution := createCmd.Int("con", 10, "constitution")
		intelligence := createCmd.Int("int", 10, "intelligence")
		wisdom := createCmd.Int("wis", 10, "wisdom")
		charisma := createCmd.Int("cha", 10, "charisma")

		createCmd.Parse(os.Args[2:])

		if *name == "" {
			fmt.Println("name is required")
			os.Exit(2)
		}

		char, err := service.CreateNewCharacter(
			*name, *class, *race, *level,
			*strength, *dexterity, *constitution,
			*intelligence, *wisdom, *charisma,
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("saved character", char.Name)

	case "view":
		viewCmd := flag.NewFlagSet("view", flag.ExitOnError)
		name := viewCmd.String("name", "", "character name to view (optional)")
		viewCmd.Parse(os.Args[2:])

		if *name != "" {
			c := service.ViewCharacter(*name)
			fmt.Println(c)
		}

	case "delete":
		viewCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		name := viewCmd.String("name", "", "character name to view (optional)")
		viewCmd.Parse(os.Args[2:])

		service.DeleteCharacter(*name)
		fmt.Printf("deleted %s", *name)
	case "equip":
		viemCmd := flag.NewFlagSet("equip", flag.ExitOnError)
		name := viemCmd.String("name", "", "character name to equip")
		equipment := viemCmd.String("weapon", "", "equipment to equip")
		viemCmd.Parse(os.Args[2:])

		fmt.Print(*name, *equipment)
	case "serve":
		fs := http.FileServer(http.Dir("./static"))
		http.Handle("/static/", http.StripPrefix("/static/", fs))

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			chars, err := repo.View("hobbit2")
			if err != nil {
				log.Fatal(err)
			}
			tmpl := template.Must(template.ParseFiles("templates/list.html"))
			tmpl.Execute(w, chars)
		})

		fmt.Println("Server running on http://localhost:8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	default:
		fmt.Println("Usage: create | view | delete")
	}

}
