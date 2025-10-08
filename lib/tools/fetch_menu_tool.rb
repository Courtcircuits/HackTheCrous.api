module Tools
  class FetchMenuTool < MCP::Tool
    title "Fetch a menu"
    description "This tool returns a menu for a given restaurant, either by name or by id but can't do both at the same time"
    input_schema(
      properties: {
        name: { type: "string" },
        id: { type: "integer" }
      }
    )
    annotations(
      read_only_hint: true,
      destructive_hint: false,
      idempotent_hint: true,
      open_world_hint: false,
      title: "Find a meal"
    )

    def self.call(name: nil, id: nil)
      if name.present?
        query = name
      elsif id.present?
        query = id
      end
      meals = MealFetchService.new({
        q: query,
        page: 1,
        per_page: 100,
        by_name: name.present?
      }).perform
      MCP::Tool::Response.new([ { type: "text", text: meals.to_json } ])
    end
  end
end
