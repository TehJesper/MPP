package spells

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/machinebox/graphql"
)

type AllSpellsResponse struct {
	Spells []Spell `json:"spells"`
}

func FetchSpells() {
	client := graphql.NewClient("https://www.dnd5eapi.co/graphql")

	req := graphql.NewRequest(`
query Spells {
  spells {
    name
    index
    classes {
      name
      index
      spellcasting {
        info {
          name
        }
        spellcasting_ability {
          name
        }
      }
      spells {
        name
        level
      }
    }
  }
}
    `)

	var respData AllSpellsResponse

	ctx := context.Background()
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatalf("GraphQL query error: %v", err)
	}

	outBytes, err := json.MarshalIndent(respData, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshal error: %v", err)
	}

	err = os.WriteFile("spells.json", outBytes, 0644)
	if err != nil {
		log.Fatalf("Write file error: %v", err)
	}
}
