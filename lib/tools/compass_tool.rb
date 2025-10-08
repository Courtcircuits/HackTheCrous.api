module Tools
  class CompassTool < MCP::Tool
    title "Compass"
    description "This tool returns the closests restaurants or schools from a given restaurant or school"
    input_schema(
      properties: {
        id: { type: "integer" },
        type: { type: "string" },
        target_type: { type: "string" }
      },
      required: %w[id type target_type],
    )
    annotations(
      read_only_hint: true,
      destructive_hint: false,
      idempotent_hint: true,
      open_world_hint: false,
      title: "Compass"
    )

    def self.call(id:, type:, target_type:)
      begin
        entity = EntityProximityService.new({
          entity_id: id,
          entity_type: type,
          searched_entity_type: target_type
        }).find_closest_entities
        puts entity
        MCP::Tool::Response.new([ { type: "text", text: entity.to_json } ])
      rescue => e
        puts e.message
        MCP::Tool::Response.new([ { type: "text", text: e.message } ])
      end
    end
  end
end
