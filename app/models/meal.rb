class Meal < ApplicationRecord
  self.table_name = "meal"
  belongs_to :restaurant

  validates :typemeal, presence: true
  validates :day, presence: true
end
