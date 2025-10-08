class Meal < ApplicationRecord
  self.table_name = "meal"
  belongs_to :restaurant, foreign_key: "idrestaurant"

  validates :typemeal, presence: true
  validates :day, presence: true
end
