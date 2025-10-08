module V2
  module Schools
    class SchoolController < V2::BaseController
      def index
        @schools = School.all
        render json: @schools
      end

      def show
        @school = School.find(params[:id])
        render json: @school
      end
    end
  end
end
