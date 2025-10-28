module Tools
  class GetRestaurants < MCP::Tool
    title "Get restaurants"
    description "This tool returns all restaurants"
    annotations(
      read_only_hint: true,
      destructive_hint: false,
      idempotent_hint: true,
      open_world_hint: false,
      title: "Get restaurants"
    )

    def self.call
      begin
        restaurants = Restaurant.all
        MCP::Tool::Response.new([ { type: "text", text: restaurants.to_json } ])
      rescue => e
        puts e
        MCP::Tool::Response.new([ { type: "text", text: e.message } ])
      end
    end
  end
end
