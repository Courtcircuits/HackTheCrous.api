class RestaurantSearchService
  def initialize(params)
    @query = params[:q]
    @page = params[:page] || 1
    @per_page = params[:per_page] || 10
  end

  def perform
    return Restaurant.none if @query.blank?

    # Calculate pagination parameters for Quickwit
    limit = @per_page.to_i
    offset = (@page.to_i - 1) * limit

    # Search using Quickwit
    search_results = QuickwitClient.search(@query, limit: limit, offset: offset)
    
    # Extract restaurant IDs from Quickwit results
    restaurant_ids = search_results["hits"]&.map { |hit| hit["id"] } || []
    
    # Return empty relation if no results
    return Restaurant.none if restaurant_ids.empty?

    # Fetch restaurants maintaining Quickwit order and apply pagination
    restaurants = Restaurant.where(id: restaurant_ids)
    
    # Apply manual ordering based on Quickwit results order
    ordered_restaurants = restaurant_ids.map { |id| restaurants.find { |r| r.id == id } }.compact

    # Convert to paginated collection to maintain interface compatibility
    paginated_collection = Kaminari.paginate_array(
      ordered_restaurants, 
      total_count: search_results["num_hits"] || 0
    ).page(@page).per(@per_page)

    paginated_collection
  end

  private

  # Fallback to original SQL search if Quickwit is unavailable
  def fallback_search
    Restaurant.joins(:suggestions_restaurant)
             .where("UPPER(suggestions_restaurant.keyword) LIKE :query", query: "%#{@query.upcase}%")
             .distinct
             .page(@page)
             .per(@per_page)
  end
end
