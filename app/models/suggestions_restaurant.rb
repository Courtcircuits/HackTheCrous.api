class SuggestionsRestaurant < ApplicationRecord
  self.table_name = "suggestions_restaurant"
  belongs_to :restaurant, foreign_key: "idrestaurant"
  belongs_to :cat_suggestion
  validates :keyword, presence: true
end

