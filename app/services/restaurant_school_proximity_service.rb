class RestaurantSchoolProximityService
  def initialize(params)
    @restaurant = params[:id]
  end

  def self.find_for_restaurant(restaurant_id)
    Restaurant.find_by_sql(<<-SQL, [ restaurant_id ])
      SELECT#{' '}
        s.name AS school_name,
        s.long_name AS school_full_name,
        s.idschool,
        (point(r.gpscoord[0], r.gpscoord[1]) <-> point(s.coords[0], s.coords[1])) * 111.111 AS distance_km
      FROM restaurant r
      CROSS JOIN LATERAL (
        SELECT#{' '}
          s.idschool,
          s.name,
          s.long_name,
          s.coords
        FROM school s
        ORDER BY point(r.gpscoord[0], r.gpscoord[1]) <-> point(s.coords[0], s.coords[1])
        LIMIT 5
      ) s
      WHERE r.idrestaurant = $1;
    SQL
  end
end
