Rails.application.routes.draw do
  # Define your application routes per the DSL in https://guides.rubyonrails.org/routing.html

  # Reveal health status on /up that returns 200 if the app boots with no exceptions, otherwise 500.
  # Can be used by load balancers and uptime monitors to verify that the app is live.
  get "up" => "rails/health#show", as: :rails_health_check

  mount Scalar::UI, at: "/docs"

  namespace :v2 do
    namespace :restaurants do
      root "restaurant#index"

      get "meals/:id", to: "meal#find"
      get "search", to: "restaurant#search"
      get ":id", to: "restaurant#show"
    end
  end

  # Defines the root path route ("/")
  # root "posts#index"
end

# SELECT idrestaurant, url, name, gpscoord, hours FROM restaurant
# WHERE idrestaurant IN (SELECT r.idrestaurant FROM restaurant r JOIN suggestions_restaurant sr ON sr.idrestaurant=r.idrestaurant WHERE UPPER(sr.keyword) LIKE $1)
