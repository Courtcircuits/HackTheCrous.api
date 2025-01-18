module V2
  module Restaurants
    class RestaurantController < ApplicationController
      rescue_from ActiveRecord::RecordNotFound, with: :record_not_found

      def index
        restaurants = Restaurant.page(params[:page])
        render json: RestaurantSerializer.new(restaurants).serializable_hash.to_json
      end

      def schools
        render json: RestaurantSchoolProximityService.find_for_restaurant(params[:id]).to_json
      end

      def show
        restaurant = Restaurant.find(params[:id])
        render json: RestaurantSerializer.new(restaurant).serializable_hash.to_json
      end

      def search
        restaurants = RestaurantSearchService.new(search_params).perform
        render json: RestaurantSerializer.new(restaurants).serializable_hash.to_json
      end

      private

      def search_params
        params.permit(:q, :page, :per_page)
      end

      def record_not_found
        render json: { error: "Restaurant not found" }, status: :not_found
      end
    end
  end
end
