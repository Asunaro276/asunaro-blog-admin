# API Integration Details

## Current Data Provider Implementation

### Blog Posts API
- **Endpoint**: Uses `API_URL` environment variable
- **GET /**: Fetches all blog posts
- **Current Status**: Read-only implementation
- **Data Transformation**: 
  - Maps API response to Refine format
  - Transforms: `id`, `title`, `body -> content`, `category_id -> category.id`, `status`, `created_at -> createdAt`

### Categories
- **Current Status**: Mock data implementation
- **Data**: Static array with Technology, Lifestyle, Business categories
- **IDs**: String-based ("1", "2", "3")

### CRUD Operations Status
- **getList**: ✅ Implemented for blog_posts, mock for categories
- **getOne**: ✅ Implemented for blog_posts, mock for categories  
- **create**: ⚠️ Mock implementation (returns generated ID)
- **update**: ⚠️ Mock implementation (returns input data)
- **deleteOne**: ⚠️ Mock implementation (returns ID only)

### Authentication Provider
- File exists but implementation needs to be checked
- Handles login, register, forgot password flows

## Integration Notes
- API URL configured via environment variable
- Error handling implemented for missing articles
- Ready for REST API expansion when backend endpoints are available
- Mock implementations allow frontend development without full backend