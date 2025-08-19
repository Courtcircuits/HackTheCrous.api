require 'faraday'
require 'json'

class QuickwitClient
  class << self
    def config
      @config ||= Rails.application.config_for(:quickwit)
    end

    def base_url
      "http://#{config['host']}:#{config['port']}"
    end

    def connection
      @connection ||= Faraday.new(
        url: base_url,
        headers: { 'Content-Type' => 'application/json' }
      ) do |faraday|
        faraday.request :json
        faraday.response :json
        faraday.adapter Faraday.default_adapter
      end
    end

    def search(query, options = {})
      search_params = {
        query: query,
        max_hits: options[:limit] || 20,
        start_offset: options[:offset] || 0
      }

      response = connection.get("/api/v1/#{config['index_name']}/search") do |req|
        req.params = search_params
        req.options.timeout = config['timeout'] || 30
      end
      
      if response.success?
        response.body
      else
        Rails.logger.error "Quickwit search failed: #{response.status} - #{response.body}"
        { "hits" => [], "num_hits" => 0 }
      end
    rescue Faraday::Error => e
      Rails.logger.error "Quickwit connection error: #{e.message}"
      { "hits" => [], "num_hits" => 0 }
    rescue => e
      Rails.logger.error "Quickwit search error: #{e.message}"
      { "hits" => [], "num_hits" => 0 }
    end

    def index_document(document)
      response = connection.post("/api/v1/#{config['index_name']}/docs") do |req|
        req.body = [document].to_json
      end

      response.success?
    rescue => e
      Rails.logger.error "Quickwit indexing error: #{e.message}"
      false
    end

    def bulk_index_documents(documents)
      return true if documents.empty?

      response = connection.post("/api/v1/#{config['index_name']}/docs") do |req|
        req.body = documents.to_json
      end

      response.success?
    rescue => e
      Rails.logger.error "Quickwit bulk indexing error: #{e.message}"
      false
    end

    def create_index
      index_config = {
        version: "0.7",
        index_id: config['index_name'],
        doc_mapping: {
          field_mappings: [
            {
              field_name: "id",
              type: "u64",
              indexed: false,
              stored: true
            },
            {
              field_name: "name",
              type: "text",
              indexed: true,
              stored: true,
              tokenizer: "default"
            },
            {
              field_name: "keywords",
              type: "text",
              indexed: true,
              stored: true,
              tokenizer: "default"
            },
            {
              field_name: "timestamp",
              type: "datetime",
              indexed: true,
              stored: true
            }
          ],
          timestamp_field: "timestamp"
        },
        search_settings: {
          default_search_fields: ["name", "keywords"]
        }
      }

      response = connection.post("/api/v1/indexes") do |req|
        req.body = index_config.to_json
      end

      response.success?
    rescue => e
      Rails.logger.error "Quickwit index creation error: #{e.message}"
      false
    end

    def health_check
      response = connection.get("/health")
      response.success? && response.body.to_s.include?("ready")
    rescue => e
      Rails.logger.error "Quickwit health check failed: #{e.message}"
      false
    end

    def index_exists?
      response = connection.get("/api/v1/indexes/#{config['index_name']}")
      response.success?
    rescue => e
      Rails.logger.debug "Index check failed: #{e.message}"
      false
    end
  end
end