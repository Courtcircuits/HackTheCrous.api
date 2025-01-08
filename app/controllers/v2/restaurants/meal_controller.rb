class V2::Restaurants::MealController < ApplicationController
  def index
    meal = Meal.all
    render json: meal
  end

  def find
    meal = Meal.where(idrestaurant: params[:id])
    render json: meal
  end
end
