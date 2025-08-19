# Quickwit Search Integration

This application uses Quickwit as the search engine for restaurant data, replacing the previous SQL LIKE-based fuzzy search.

## Configuration

Quickwit settings are configured in `config/quickwit.yml`:

```yaml
development:
  host: localhost
  port: 7280
  index_name: restaurants
  timeout: 30
```

Environment variables can override these settings:
- `QUICKWIT_HOST`: Quickwit server hostname
- `QUICKWIT_PORT`: Quickwit server port

## Setup

1. **Start Quickwit server:**
   ```bash
   # Using Docker
   docker run -p 7280:7280 quickwit/quickwit:latest
   ```

2. **Create the search index:**
   ```bash
   bundle exec rake quickwit:create_index
   ```

3. **Index existing restaurant data:**
   ```bash
   bundle exec rake quickwit:index_restaurants
   ```

## Usage

### Search API
The search endpoint remains unchanged:
```
GET /v2/restaurants/search?q=burger&page=1&per_page=10
```

### Indexing
- New restaurants are automatically indexed when created
- Updated restaurants are automatically re-indexed when modified
- Manual re-indexing: `bundle exec rake quickwit:reindex`

### Testing Search
```bash
# Test search functionality
bundle exec rake quickwit:test_search["burger"]
```

## Architecture

- **`QuickwitClient`**: HTTP client for Quickwit API communication
- **`QuickwitIndexingService`**: Handles indexing restaurant data
- **`RestaurantSearchService`**: Updated to use Quickwit instead of SQL queries
- **Rake Tasks**: Management commands for indexing and testing

## Index Schema

The Quickwit index stores:
- `id`: Restaurant unique identifier
- `name`: Restaurant name (searchable)
- `keywords`: Combined keywords from suggestions (searchable)
- `timestamp`: Last update time

## Migration from SQL Search

The old SQL-based search using LIKE queries has been replaced with Quickwit's indexed search:

**Before:**
```sql
SELECT * FROM restaurant 
JOIN suggestions_restaurant 
WHERE UPPER(keyword) LIKE '%QUERY%'
```

**After:**
```
POST /api/v1/restaurants/search
{ "query": "QUERY", "max_hits": 20 }
```

This provides:
- Better search relevance
- Faster query performance
- Full-text search capabilities
- Scalable architecture