package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/template"

	"test/domain"
	"test/infrastructure"
	"test/infrastructure/class"
	"test/infrastructure/race"
	"test/infrastructure/spells"

	// "test/infrastructure/equipment"
	"test/services"

	_ "modernc.org/sqlite"
)

func initAPIData() {
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		if _, err := os.Stat("classes.json"); errors.Is(err, os.ErrNotExist) {
			class.FetchClasses()
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := os.Stat("races.json"); errors.Is(err, os.ErrNotExist) {
			race.FetchRaces()
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := os.Stat("spells.json"); errors.Is(err, os.ErrNotExist) {
			spells.FetchSpells()
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := os.Stat("equipment.json"); errors.Is(err, os.ErrNotExist) {
			// equipment.FetchEquipment()
		}
	}()

	wg.Wait()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: create | view | delete | equip | serve | prepare-spell | learn-spell")
		os.Exit(1)
	}

	// Init API data
	initAPIData()

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
	proficiency_bonus INTEGER,
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
	skills TEXT,
	mainhand TEXT,
	offhand TEXT,
	shield TEXT,
	armor TEXT,
	armor_class INTEGER,
	initiative INTEGER,
	passive_perception INTEGER
)`)
	if err != nil {
		log.Fatal(err)
	}
	// Init service
	repo := infrastructure.NewSQLRepository(db)
	
	// eqService := &services.EquipmentService{}
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
		name := viewCmd.String("name", "", "character name to delete")
		viewCmd.Parse(os.Args[2:])

		service.DeleteCharacter(*name)
		fmt.Printf("deleted %s", *name)
	case "equip":
		viemCmd := flag.NewFlagSet("equip", flag.ExitOnError)
		name := viemCmd.String("name", "", "character name to equip")
		weapon := viemCmd.String("weapon", "", "equipment to equip")
		shield := viemCmd.String("shield", "", "shield to equip")
		armor := viemCmd.String("armor", "", "armor to equip")
		slot := viemCmd.String("slot", "", "slot to equip")
		viemCmd.Parse(os.Args[2:])

		err := service.EquipCharacter(*name, *weapon, *slot, *armor, *shield)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if *weapon != "" {
			fmt.Print("Equipped weapon ", *weapon, " to ", *slot)
		}
		if *shield != "" {
			fmt.Print("Equipped shield ", *shield)
		}
		if *armor != "" {
			fmt.Print("Equipped armor ", *armor)
		}

	case "serve":
		viemCmd := flag.NewFlagSet("serve", flag.ExitOnError)
		name := viemCmd.String("name", "", "character name to view")
		viemCmd.Parse(os.Args[2:])

		fs := http.FileServer(http.Dir("./static"))
		http.Handle("/static/", http.StripPrefix("/static/", fs))

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			chars, err := repo.View(*name)
			if err != nil {
				log.Fatal(err)
			}
			   // Convert comma-separated skills to a map
			skillMap := make(map[string]bool)
			for _, skill := range strings.Split(chars.Skills, ",") {
				skillMap[strings.TrimSpace(skill)] = true
			}

			// Pass both the character and the map to the template
			tmplData := struct {
				Character *domain.Character
				SkillMap map[string]bool
			}{
				Character: &chars,
				SkillMap: skillMap,
			}
			fmt.Print(tmplData)
			tmpl := template.Must(template.ParseFiles("templates/list.html"))
			tmpl.Execute(w, tmplData)
		})

		fmt.Println("Server running on http://localhost:8080")
		log.Fatal(http.ListenAndServe(":8080", nil))

	// prepare-spell
	case "prepare-spell":
	viemCmd := flag.NewFlagSet("prepare-spell", flag.ExitOnError)
	spell := viemCmd.String("spell", "", "spell to prepare")
	name := viemCmd.String("name", "", "character name to view")

	viemCmd.Parse(os.Args[2:])

	service.PrepareSpell(*name, *spell)
	// fmt.Print("Prepared spell ", *spell)

	// learn-spell
	case "learn-spell":
	viemCmd := flag.NewFlagSet("learn-spell", flag.ExitOnError)
	spell := viemCmd.String("spell", "", "spell to learn")
	name := viemCmd.String("name", "", "character name to view")

	viemCmd.Parse(os.Args[2:])

	service.LearnSpell(*name, *spell)

	default:
		fmt.Println("Usage: create | view | delete | equip | serve | prepare-spell | learn-spell")
	}

}
