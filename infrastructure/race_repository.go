package infrastructure

import (
    "encoding/json"
    "os"
)

type Race struct {
    Index    string    `json:"index"`
    Name     string    `json:"name"`
    Subraces []Subrace `json:"subraces"`
}

type Subrace struct {
    Index string `json:"index"`
    Name  string `json:"name"`
}

func OpenRaceFile() ([]byte, error) {
	return os.ReadFile("races.json")
}

func LoadRacesAndSubraces() ([]string, error) {
    data, err := OpenRaceFile()
    if err != nil {
        return nil, err
    }

    var races []Race
    if err := json.Unmarshal(data, &races); err != nil {
        return nil, err
    }

    var all []string
    for _, r := range races {
        all = append(all, r.Name)
        for _, sub := range r.Subraces {
            all = append(all, sub.Name)
        }
    }
	
    return all, nil
}
