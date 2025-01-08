class Restaurant < ApplicationRecord
  self.table_name = "restaurant"
  has_many :meals, dependent: :destroy
  has_many :suggestions_restaurant, dependent: :destroy, foreign_key: "idrestaurant"
  has_many :cat_suggestions, through: :suggestions_restaurants, source: :suggetions_restaurant

  validates :name, presence: true
end
