namespace :quickwit do
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

  desc "Reindex all restaurants in Quickwit"
  task reindex: [:create_index, :index_restaurants]

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