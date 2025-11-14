package race_test

import (
    "os"
    "path/filepath"
    "reflect"
    "testing"
	rce "test/infrastructure/race"
)

func TestGetRaceSkillsByName(t *testing.T) {
    t.Parallel()

    sample := `{
  "races": [
    {
      "index": "half-elf",
      "name": "Half-Elf",
      "traits": [
        {
          "index": "menacing",
          "name": "Menacing",
          "desc": [
            "You gain proficiency in the Intimidation skill."
          ],
          "proficiencies": [
            {
              "index": "skill-intimidation",
              "name": "Skill: Intimidation"
            }
          ]
        }
      ]
    }
  ]
}`

    tmp := t.TempDir()
    filePath := filepath.Join(tmp, "races.json")
    if err := os.WriteFile(filePath, []byte(sample), 0o644); err != nil {
        t.Fatalf("write sample races.json: %v", err)
    }

    // switch working directory so OpenRaceFile() will read our temp races.json
    origWd, err := os.Getwd()
    if err != nil {
        t.Fatalf("getwd: %v", err)
    }
    defer func() { _ = os.Chdir(origWd) }()

    if err := os.Chdir(tmp); err != nil {
        t.Fatalf("chdir: %v", err)
    }

    skills, err := rce.GetRaceSkillsByName("Half-Elf")
    if err != nil {
        t.Fatalf("GetRaceSkillsByName returned error: %v", err)
    }

    want := []string{"Intimidation"}
    if !reflect.DeepEqual(skills, want) {
        t.Fatalf("skills = %v, want %v", skills, want)
    }
}