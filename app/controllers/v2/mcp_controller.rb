module V2
    class McpController < ApplicationController
      def index
        server = MCP::Server.new(
          name: "HackTheCrous.mcp",
          title: "MCP Server for HackTheCrous",
          version: "2.1.0",
          instructions: "Use the tools of this server to interact with the data of HackTheCrous",
          tools: [ Tools::CompassTool, Tools::FetchRestaurantTool, Tools::FetchMenuTool ],
          server_context: {}
        )
        render(json: server.handle_json(request.body.read))
      end
    end
end
