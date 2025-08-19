class Restaurant < ApplicationRecord
  self.table_name = "restaurant"
  has_many :meals, dependent: :destroy
  has_many :suggestions_restaurant, dependent: :destroy, foreign_key: "idrestaurant"
  has_many :cat_suggestions, through: :suggestions_restaurants, source: :suggetions_restaurant

  validates :name, presence: true

  # Auto-index in Quickwit after create/update
  after_commit :index_in_quickwit, on: [:create, :update]

  private

  def index_in_quickwit
    # Skip indexing in test environment unless explicitly enabled
    return if Rails.env.test? && ENV['QUICKWIT_INDEX_IN_TESTS'] != 'true'
    
    # Use background job in production, but for now do it synchronously
    QuickwitIndexingService.index_restaurant(self)
  rescue => e
    Rails.logger.error "Failed to index restaurant #{id} in Quickwit: #{e.message}"
    # Don't raise the error to avoid breaking the save operation
  end
end
