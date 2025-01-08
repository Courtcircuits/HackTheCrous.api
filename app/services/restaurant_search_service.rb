class RestaurantSearchService
  def initialize(params)
    @query = params[:q]
    @page = params[:page]
    @per_page = params[:per_page]
  end

  def perform
    return Restaurant.none if @query.blank?


    Restaurant.find_by_sql([
      "SELECT * FROM restaurant
WHERE idrestaurant IN (SELECT r.idrestaurant FROM restaurant r JOIN suggestions_restaurant sr ON sr.idrestaurant=r.idrestaurant WHERE UPPER(sr.keyword) LIKE :query)",
{
  query: "%#{@query.upcase}%"
}
    ])
  end
end
