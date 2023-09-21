package test

import (
	"testing"

	"github.com/Courtcircuits/HackTheCrous.api/types"
)

func TestParseMeal(t *testing.T) {
	foods_string := `[{"food": ["1 boisson chaude (café, thé, chocolat)", "+ 1 viennoiserie", "+ 1 jus d'orange ou 1 yaourt", "Tarif : 2.15 TTC (1.95 € HT)"], "type": "Formule petit déjeuner"}, {"food": ["1 boisson chaude (café, thé, chocolat)", "+ 1 viennoiserie", "+ 1 jus d'orange ", "+ 1 yaourt", "+ 1 compote", "Tarif : 3.19 TTC (2.90 € HT)", "HAPPY HOUR", "Tous les jours de 7h30 à 8h30, le petit dej' gourmand au prix du petit dej' classique ! (soit 2.15 TTC)"], "type": "Formule petit déjeuner gourmand"}]`
	foods, err := types.ParseFoods(foods_string)

	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(foods)
}
