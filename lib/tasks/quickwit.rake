namespace :quickwit do
  desc "Check Quickwit health"
  task health: :environment do
    puts "Checking Quickwit health..."
    
    if QuickwitClient.health_check
      puts "✓ Quickwit is healthy and responding"
      
      if QuickwitClient.index_exists?
        puts "✓ Search index exists"
      else
        puts "⚠ Search index does not exist"
      end
    else
      puts "✗ Quickwit is not available"
      exit 1
    end
  end

  desc "Create Quickwit index"
  task create_index: :environment do
    puts "Creating Quickwit index..."
    if QuickwitClient.create_index
      puts "✓ Quickwit index created successfully"
    else
      puts "✗ Failed to create Quickwit index"
      exit 1
    end
  end

  desc "Index all restaurants in Quickwit"
  task index_restaurants: :environment do
    puts "Starting restaurant indexing..."
    count = QuickwitIndexingService.index_all_restaurants
    puts "✓ Indexed #{count} restaurants in Quickwit"
  end

  desc "Setup Quickwit (health check, create index, and index data)"
  task setup: [:health, :create_index, :index_restaurants]

  desc "Reindex all restaurants in Quickwit"
  task reindex: [:health, :create_index, :index_restaurants]

  desc "Test Quickwit search"
  task :test_search, [:query] => :environment do |t, args|
    query = args[:query] || "test"
    puts "Testing Quickwit search with query: '#{query}'"
    
    results = QuickwitClient.search(query)
    puts "Results: #{results['num_hits']} hits found"
    
    results['hits']&.each_with_index do |hit, index|
      puts "#{index + 1}. #{hit['name']} (ID: #{hit['id']})"
    end
  end
end