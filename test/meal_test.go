package test

import (
	"testing"

	"github.com/Courtcircuits/HackTheCrous.api/types"
)

func TestParseMeal(t *testing.T) {
	foods_string := `[
  {
    "type": "Formule goûter",
    "content": [
      "menu non communiqué"
    ]
  },
  {
    "type": "Pizza",
    "content": [
      "menu non communiqué"
    ]
  },
  {
    "type": "Plat cuisiné",
    "content": [
      "Pâtes bolognaise",
      "Carré de seitan"
    ]
  }
	]`

	foods, err := types.ParseFoods(foods_string)

	food_compare := []types.Foods{
		{
			Type: "Formule goûter",
			Food: []string{
				"menu non communiqué",
			},
		},
		{
			Type: "Pizza",
			Food: []string{
				"menu non communiqué",
			},
		},
		{
			Type: "Plat cuisiné",
			Food: []string{
				"Pâtes bolognaise",
				"Carré de seitan",
			},
		},
	}

	is_equal := true
	for i, food_cat := range food_compare {
		if food_cat.Type != foods[i].Type {
			is_equal = false
			break
		}
		for j, food := range food_cat.Food {
			if food != foods[i].Food[j] {
				is_equal = false
				break
			}
		}
	}

	if err != nil {
		t.Fatal(err.Error())
	}

	if !is_equal {
		t.Fatal("foods are not equal")
	}

	t.Log(foods)
}
