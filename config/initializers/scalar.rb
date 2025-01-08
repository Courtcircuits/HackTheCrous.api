# config/initializers/scalar.rb
Scalar.setup do |config|
  config.specification = File.read(Rails.root.join('docs/openapi.yaml'))
end
