class CatSuggestion < ApplicationRecord
  has_many :suggestions_restaurants, dependent: :destroy
  has_many :restaurants, through: :suggestions_restaurants
  validates :namecat, presence: true
end
