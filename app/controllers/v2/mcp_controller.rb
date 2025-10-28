module V2
    class McpController < ApplicationController
      def index
        server = MCP::Server.new(
          name: "HackTheCrous.mcp",
          title: "MCP Server for HackTheCrous",
          version: "2.1.0",
          instructions: "Use the tools of this server to interact with the data of HackTheCrous. The tools inputs must be written in French.",
          tools: [ Tools::CompassTool, Tools::FetchRestaurantTool, Tools::FetchMenuTool, Tools::GetRestaurants ],
          server_context: {}
        )
        render(json: server.handle_json(request.body.read))
      end
    end
end
