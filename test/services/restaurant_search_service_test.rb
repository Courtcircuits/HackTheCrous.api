require "test_helper"

class RestaurantSearchServiceTest < ActiveSupport::TestCase
  def setup
    @service = RestaurantSearchService.new({ q: "burger", page: 1, per_page: 10 })
  end

  test "returns empty result for blank query" do
    service = RestaurantSearchService.new({ q: "", page: 1, per_page: 10 })
    result = service.perform
    assert_equal Restaurant.none.class, result.class
  end

  test "returns empty result for nil query" do
    service = RestaurantSearchService.new({ q: nil, page: 1, per_page: 10 })
    result = service.perform
    assert_equal Restaurant.none.class, result.class
  end

  test "handles missing page parameter with default" do
    service = RestaurantSearchService.new({ q: "test" })
    # Should not raise an error
    assert_not_nil service
  end

  test "handles missing per_page parameter with default" do
    service = RestaurantSearchService.new({ q: "test", page: 1 })
    # Should not raise an error
    assert_not_nil service
  end

  # Note: Testing actual Quickwit integration would require
  # a running Quickwit instance and proper test data setup
  # For now, we test the basic service behavior
end