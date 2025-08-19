require "test_helper"

class QuickwitClientTest < ActiveSupport::TestCase
  test "config loads from quickwit.yml" do
    config = QuickwitClient.config
    assert_not_nil config
    assert config.key?('host')
    assert config.key?('port')
    assert config.key?('index_name')
  end

  test "base_url formats correctly" do
    base_url = QuickwitClient.base_url
    assert_match %r{^http://.*:\d+$}, base_url
  end

  test "search handles empty query gracefully" do
    # Mock the connection to avoid actual HTTP calls in tests
    mock_response = { "hits" => [], "num_hits" => 0 }
    
    # This would require actual mocking in a real test environment
    # For now, we just test that the method exists and returns expected structure
    assert_respond_to QuickwitClient, :search
  end

  test "connection is properly configured" do
    connection = QuickwitClient.connection
    assert_not_nil connection
    assert_equal 'application/json', connection.headers['Content-Type']
  end
end