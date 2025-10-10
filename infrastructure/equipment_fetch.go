package infrastructure

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/machinebox/graphql"
)

type AllEquipmentResponse struct {
	Equipment []Equipment `json:"equipment"`
}

func FetchEquipment() {
	client := graphql.NewClient("https://www.dnd5eapi.co/graphql")

	req := graphql.NewRequest(`
    
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
