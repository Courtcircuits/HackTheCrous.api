class RestaurantSerializer
  include JSONAPI::Serializer
  attributes :url, :name, :hours
  attribute :gps_coord do |object|
    {
      "X" => object.gpscoord["x"],
      "Y" => object.gpscoord["y"]
    }
  end


  # has_many :suggestions_restaurant
  cache_options enabled: true, cache_length: 12.hours

  def self.serialize(resources, options = {}) # override to comply with v1 specification
    serialized = super(resources, options)
    if serialized.is_a?(Hash) && serialized.has_key?(:data)
      serialized[:data].map do |item|
        {
          id: item[:id].to_i,
          url: item[:attributes][:url],
          name: item[:attributes][:name],
          gps_coord: item[:attributes][:gps_coord],
          hours: item[:attributes][:hours]
        }
      end
    else
      serialized
    end
  end
end
