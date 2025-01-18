class School < ApplicationRecord
  self.table_name = "school"
  validates :name, presence: true
  validates :coords, presence: true
  validates :long_name, presence: true
end
