class MealFetchService
  def initialize(params)
    @query = params[:q]
    @page = params[:page]
    @per_page = params[:per_page]
    @by_name = params[:by_name] # this is a boolean
  end

  def perform
    return Meal.none if @query.blank?

    if @by_name
      Meal.joins(:restaurant)
          .where("restaurant.name ILIKE :query", query: "%#{@query.upcase}%")
          .distinct
          .page(@page)
          .per(@per_page)
    else
      Meal.joins(:restaurant)
          .where("restaurant.idrestaurant = :query", query: @query)
          .distinct
    end

  end
end
