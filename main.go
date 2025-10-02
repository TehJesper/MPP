package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"test/infrastructure"
	"test/services"
	"test/domain"

	_ "modernc.org/sqlite"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: create | view | delete")
		os.Exit(1)
	}

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
		skills TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}

	repo := infrastructure.NewSQLRepository(db)
	service := services.NewService(repo)

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

	default:
		fmt.Println("Usage: create | view | delete")
	}

}


func formatCharacter(c domain.Character) string {
	return fmt.Sprintf(
		`Name: %s
Class: %s
Race: %s
Background: acolyte
Level: %d
Ability scores:
  STR: %d
  DEX: %d
  CON: %d
  INT: %d
  WIS: %d
  CHA: %d
Proficiency bonus: +2
Skill proficiencies: %s`,
		c.Name, c.Class, c.Race, c.Level,
		c.Strength, c.Dexterity, c.Constitution,
		c.Intelligence, c.Wisdom, c.Charisma,
		c.Skills,
	)
}

	// case "serve":
    // // Serve static files from ./static
    // fs := http.FileServer(http.Dir("./static"))
    // http.Handle("/static/", http.StripPrefix("/static/", fs))

    // http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	// Debug
	// 	fmt.Println(getCharacter("niets"))

	// 	chars, err := getCharacter("Alice")
	// 	if err != nil {
    // 		log.Fatal(err) // or handle gracefully
	// 	}
    //     tmpl := template.Must(template.ParseFiles("templates/list.html"))
    //     tmpl.Execute(w, chars)
    // })

    // fmt.Println("Server running on http://localhost:8080")
    // log.Fatal(http.ListenAndServe(":8080", nil))

	// }
// }
// }