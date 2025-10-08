module V2
  module Restaurants
    class MealController < V2::BaseController
      def index
        meal = Meal.all
        render json: meal
      end

      def find
        meal = Meal.where(idrestaurant: params[:id])
        render json: meal
      end
    end
  end
end
