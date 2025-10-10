package infrastructure

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
