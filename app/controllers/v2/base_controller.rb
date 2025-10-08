module V2
  class BaseController < ApplicationController
    rescue_from ActiveRecord::RecordNotFound, with: :record_not_found

    def record_not_found
      render json: { error: "Not found" }, status: :not_found
    end
  end
end
