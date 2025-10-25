package class

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/machinebox/graphql"
)

type AllClassesResponse struct {
	Classes []Class `json:"classes"`
}

func FetchClasses() {
	client := graphql.NewClient("https://www.dnd5eapi.co/graphql")

	req := graphql.NewRequest(`
query Classes {
  classes {
    index
    name
    proficiencies {
      index
      name
    }
    subclasses {
      name
      index
    }
    proficiency_choices {
      from {
        options {
          item {
            ... on Proficiency {
              name
              index
            }
          }
        }
      }
      choose
    }
    spells {
      name
      level
    }
    class_levels {
      spellcasting {
        cantrips_known
        spell_slots_level_1
        spell_slots_level_2
        spell_slots_level_3
        spell_slots_level_4
        spell_slots_level_5
        spell_slots_level_6
        spell_slots_level_7
        spell_slots_level_8
        spell_slots_level_9
      }
    }
  }
}
    `)

	var respData AllClassesResponse

	ctx := context.Background()
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatalf("GraphQL query error: %v", err)
	}

	outBytes, err := json.MarshalIndent(respData, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshal error: %v", err)
	}

	err = os.WriteFile("classes.json", outBytes, 0644)
	if err != nil {
		log.Fatalf("Write file error: %v", err)
	}
}
