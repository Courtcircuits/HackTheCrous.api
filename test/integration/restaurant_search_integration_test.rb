require "test_helper"

class RestaurantSearchIntegrationTest < ActionDispatch::IntegrationTest
  test "search endpoint responds successfully" do
    get "/v2/restaurants/search", params: { q: "burger" }
    
    assert_response :success
    response_data = JSON.parse(response.body)
    
    # Should return properly formatted JSON response
    assert response_data.key?("data")
  end

  test "search endpoint handles empty query" do
    get "/v2/restaurants/search", params: { q: "" }
    
    assert_response :success
    response_data = JSON.parse(response.body)
    
    # Should return empty results for blank query
    assert_equal [], response_data["data"]
  end

  test "search endpoint handles pagination parameters" do
    get "/v2/restaurants/search", params: { q: "test", page: 2, per_page: 5 }
    
    assert_response :success
    # Should not fail even if no results exist
  end

  test "search endpoint handles missing query parameter" do
    get "/v2/restaurants/search"
    
    assert_response :success
    response_data = JSON.parse(response.body)
    
    # Should return empty results for missing query
    assert_equal [], response_data["data"]
  end
end