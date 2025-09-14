package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"test/character"

	_ "modernc.org/sqlite" // pure Go SQLite driver
)

type Character struct {
	Name  string
	Race  string
	Class string
	Level int
	Strength int
	Dexterity int
	Constitution int
	Intelligence int
	Wisdom int 
	Charisma int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: create | serve")
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
		name TEXT,
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

		char, err := getCharacter(*name)
		if err != nil {
			log.Fatal(err)
		}

		printCharacter(char)
		// TODO check of character deleted is/bestaat
		
	case "delete":
		// TODO delete by name
	case "serve":
    // Serve static files from ./static
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Debug
		fmt.Println(getCharacter("niets"))

		chars, err := getCharacter("Alice")
		if err != nil {
    		log.Fatal(err) // or handle gracefully
		}
        tmpl := template.Must(template.ParseFiles("templates/list.html"))
        tmpl.Execute(w, chars)
    })

    fmt.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))

	}
}
func getCharacter(name string) (Character, error) {
    db, err := sql.Open("sqlite", "./characters.db")
    if err != nil {
        return Character{}, err
    }
    defer db.Close()

    var c Character
    err = db.QueryRow(`SELECT name, race, class, level, strength, dexterity, constitution, intelligence, wisdom, charisma 
                       FROM characters WHERE name = ?`, name).
        Scan(&c.Name, &c.Race, &c.Class, &c.Level, 
             &c.Strength, &c.Dexterity, &c.Constitution, 
             &c.Intelligence, &c.Wisdom, &c.Charisma)

    if err != nil {
        return Character{}, err
    }

    return c, nil
}
func modifier(score int) int {
    if (score-10)%2 != 0 && score < 10 {
        return (score-10)/2 - 1
    }
    return (score - 10) / 2
}


func printCharacter(char Character) {
    strMod := modifier(char.Strength)
    dexMod := modifier(char.Dexterity)
    conMod := modifier(char.Constitution)
    intMod := modifier(char.Intelligence)
    wisMod := modifier(char.Wisdom)
    chaMod := modifier(char.Charisma)

    fmt.Printf(
`Name: %s
Class: %s
Race: %s
Background: Acolyte
Level: %d
Ability scores:
  STR: %d (%+d)
  DEX: %d (%+d)
  CON: %d (%+d)
  INT: %d (%+d)
  WIS: %d (%+d)
  CHA: %d (%+d)
Proficiency bonus: +2
Skill proficiencies: animal handling, athletics, insight, religion
`, char.Name, char.Class, char.Race, char.Level,
        char.Strength, strMod,
        char.Dexterity, dexMod,
        char.Constitution, conMod,
        char.Intelligence, intMod,
        char.Wisdom, wisMod,
        char.Charisma, chaMod)
}