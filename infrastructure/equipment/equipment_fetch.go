package equipment

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/machinebox/graphql"
)

type AllEquipmentResponse struct {
	Equipment []Equipment `json:"equipments,omitempty"`
}

func FetchEquipment() {
	client := graphql.NewClient("https://www.dnd5eapi.co/graphql")

	req := graphql.NewRequest(`
    query Equipments {
  equipments {
    ... on Armor {
      name
      index
      equipment_category {
        name
        index
      }
      properties {
        index
        name
      }
      armor_class {
        base
        dex_bonus
      }
      armor_category
    }
    ... on Weapon {
      index
      name
      equipment_category {
        name
        index
      }
      properties {
        name
        index
      }
    }
  }
}
    `)

	var respData AllEquipmentResponse

	ctx := context.Background()
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatalf("GraphQL query error: %v", err)
	}

	outBytes, err := json.MarshalIndent(respData, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshal error: %v", err)
	}

	err = os.WriteFile("equipment.json", outBytes, 0644)
	if err != nil {
		log.Fatalf("Write file error: %v", err)
	}
}
