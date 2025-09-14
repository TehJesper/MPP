package main

import (
	"database/sql"
	"flag"
	"fmt"
	// "html/template"
	"log"
	// "net/http"
	"os"
	"test/character"

	_ "modernc.org/sqlite" // pure Go SQLite driver
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: create | view |serve")
		os.Exit(1)
	}

	db, err := sql.Open("sqlite", "./characters.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ensure table exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS characters (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE,
		race TEXT,
		class TEXT,
		level INTEGER,
		strength,
		dexterity,
		constitution,
		intelligence,
		wisdom,
		charisma
		)`)
	if err != nil {
		log.Fatal(err)
	}

	repo := character.NewSQLRepository(db)
	service := character.NewService(repo)

	switch os.Args[1] {
	case "create":
		createCmd := flag.NewFlagSet("create", flag.ExitOnError)
		name := createCmd.String("name", "", "character name (required)")
		class := createCmd.String("class", "", "character class (required)")
		race := createCmd.String("race", "", "charcter race (required)")
		level := createCmd.Int("level", 1, "character level")
		strength := createCmd.Int("str", 1, "character strength")
		dexterity := createCmd.Int("dex", 1, "character dexterity")
		constitution := createCmd.Int("con", 1, "character constitution")
		intelligence := createCmd.Int("int", 1, "character intelligence")
		wisdom := createCmd.Int("wis", 1, "character wisdom")
		charisma := createCmd.Int("cha", 1, "character charisma")
		createCmd.Parse(os.Args[2:])

		char, err := service.CreateNewCharacter(*name, *class, *race, *level, *strength, *dexterity, *constitution, *intelligence, *wisdom, *charisma)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Created character", char.Name)
	
	
	case "view":
		viewCmd := flag.NewFlagSet("view", flag.ExitOnError)
		name := viewCmd.String("name", "", "character name (required)")
		_ = viewCmd.Parse(os.Args[2:])

		if *name == "" {
			fmt.Println("name is required")
			os.Exit(2)
		}

		// Call the service
		output, err := service.ViewCharacterByName(*name)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println(output)
}
		
	case "delete":
		// TODO delete by name
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
}
}