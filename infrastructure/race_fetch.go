package infrastructure

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/machinebox/graphql"
)

type AllRacesResponse struct {
	Races []Race `json:"races"`
}

func FetchRaces() {
	client := graphql.NewClient("https://www.dnd5eapi.co/graphql")

	req := graphql.NewRequest(`
query Races {
  races {
    name
    index
    ability_bonuses {
      bonus
      ability_score {
        index
        name
      }
    }
    subraces {
      name
      index
      ability_bonuses {
        ability_score {
          index
          name
        }
        bonus
      }
    }
  }
}
    `)

	var respData AllRacesResponse

	ctx := context.Background()
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatalf("GraphQL query error: %v", err)
	}

	outBytes, err := json.MarshalIndent(respData, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshal error: %v", err)
	}

	err = os.WriteFile("races.json", outBytes, 0644)
	if err != nil {
		log.Fatalf("Write file error: %v", err)
	}
}
