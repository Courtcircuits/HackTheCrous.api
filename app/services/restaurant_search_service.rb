class RestaurantSearchService
  def initialize(params)
    @query = params[:q]
    @page = params[:page]
    @per_page = params[:per_page]
  end

  def perform
    return Restaurant.none if @query.blank?

    Restaurant.joins(:suggestions_restaurant)
             .where("UPPER(suggestions_restaurant.keyword) LIKE :query", query: "%#{@query.upcase}%")
             .distinct
             .page(@page)
             .per(@per_page)
  end
end
