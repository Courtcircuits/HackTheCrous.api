class EntityProximityService
  def initialize(params)
    @entity_id = params[:entity_id]
    @entity_type = params[:entity_type]
    @searched_entity_type = params[:searched_entity_type]
  end

  def find_closest_entities
    coords_col_name = if @searched_entity_type == "restaurant"
                       "rout.gpscoord"
    else
                       "rout.coords"
    end

    inner_query = if @searched_entity_type == "restaurant"
      <<~SQL
        SELECT
          rin.idrestaurant,
          rin.url,
          rin.name,
          rin.gpscoord,
          rin.hours
        FROM restaurant rin
        WHERE rin.idrestaurant = $1
        ORDER BY point(#{coords_col_name}[0], #{coords_col_name}[1]) <-> point(rin.gpscoord[0], rin.gpscoord[1])
        LIMIT 10
      SQL
    else
      <<~SQL
        SELECT#{' '}
          rin.idschool,
          rin.name,
          rin.long_name,
          rin.coords
        FROM school rin
          WHERE rin.idschool = $1;
        ORDER BY point(#{coords_col_name}[0], #{coords_col_name}[1]) <-> point(rin.coords[0], rin.coords[1])
        LIMIT 10
      SQL
    end
    query = if @searched_entity_type == "restaurant"
              <<~SQL
                SELECT
                rout.idrestaurant AS idrestaurant,
                rout.url AS url,
                rout.name AS name,
                rout.gpscoord,
                rout.hours AS hours
                FROM restaurant rout
                CROSS JOIN LATERAL (
                  #{inner_query}
                ) rin;
              SQL
    else
              <<~SQL
                SELECT
                rout.idschool AS idschool,
                rout.name AS name,
                rout.long_name AS long_name,
                rout.coords
                FROM school rout
                CROSS JOIN LATERAL (
                  #{inner_query}
                ) rin;
                SQL
    end

    if @searched_entity_type == "restaurant"
      Restaurant.find_by_sql(query, [ @entity_id ])
    else
      School.find_by_sql(query, [ @entity_id ])
    end
  end
end
