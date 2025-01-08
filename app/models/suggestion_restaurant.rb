class SuggestionsRestaurant < ApplicationRecord
  belongs_to :restaurant
  belongs_to :cat_suggestion
  validates :keyword, presence: true
end

