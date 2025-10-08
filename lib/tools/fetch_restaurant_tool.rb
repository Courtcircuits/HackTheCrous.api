module Tools
  class FetchRestaurantTool < MCP::Tool
    title "Find a restaurant"
    description "This tool is a search engine for restaurants, based on a query it returns all fitting restaurants with their meals"
    input_schema(
      properties: {
        query: { type: "string" }
      },
      required: %w[query],
    )
    annotations(
      read_only_hint: true,
      destructive_hint: false,
      idempotent_hint: true,
      open_world_hint: false,
      title: "Find a restaurant"
    )

    def self.call(query:)
      restaurants = RestaurantSearchService.new({
        q: query,
        page: 1,
        per_page: 100
      }).perform
      MCP::Tool::Response.new([ { type: "text", text: restaurants.to_json } ])
    end
  end
end
