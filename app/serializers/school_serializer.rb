class SchoolSerializer
  include JSONAPI::Serializer
  attributes :url, :name, :hours
  attribute :gps_coord do |object|
    {
      "X" => object.gpscoord["x"],
      "Y" => object.gpscoord["y"]
    }
  end
  cache_options enabled: true, cache_length: 12.hours
end
