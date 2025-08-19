class QuickwitIndexingService
  def self.index_all_restaurants
    new.index_all_restaurants
  end

  def self.index_restaurant(restaurant)
    new.index_restaurant(restaurant)
  end

  def index_all_restaurants
    Rails.logger.info "Starting Quickwit indexing for all restaurants"
    
    # Check if Quickwit is healthy
    unless QuickwitClient.health_check
      Rails.logger.error "Quickwit is not available, skipping indexing"
      return 0
    end
    
    # Create the index if it doesn't exist
    unless QuickwitClient.index_exists?
      Rails.logger.info "Creating Quickwit index..."
      unless QuickwitClient.create_index
        Rails.logger.error "Failed to create Quickwit index"
        return 0
      end
    end
    
    batch_size = 100
    indexed_count = 0
    
    Restaurant.includes(:suggestions_restaurant).find_in_batches(batch_size: batch_size) do |restaurants|
      documents = restaurants.map { |restaurant| prepare_document(restaurant) }
      
      if QuickwitClient.bulk_index_documents(documents)
        indexed_count += documents.size
        Rails.logger.info "Indexed batch of #{documents.size} restaurants (total: #{indexed_count})"
      else
        Rails.logger.error "Failed to index batch of restaurants"
      end
    end
    
    Rails.logger.info "Completed Quickwit indexing: #{indexed_count} restaurants indexed"
    indexed_count
  end

  def index_restaurant(restaurant)
    document = prepare_document(restaurant)
    QuickwitClient.index_document(document)
  end

  private

  def prepare_document(restaurant)
    keywords = restaurant.suggestions_restaurant.pluck(:keyword).join(' ')
    
    {
      id: restaurant.id,
      name: restaurant.name,
      keywords: keywords,
      timestamp: (restaurant.updated_at || restaurant.created_at || Time.current).to_i
    }
  end
end