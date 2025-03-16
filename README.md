# FizzBuzz REST Server

## Overview
This is a production-ready FizzBuzz REST API server implemented in Go using the Fiber web framework.

## Features
- FizzBuzz generation endpoint
- Request statistics tracking
- Production-ready with middleware for logging, recovery, and CORS

## Endpoints

### FizzBuzz Generation
- **URL**: `/fizzbuzz`
- **Method**: POST
- **Request Body**:
```json
{
  "int1": 3,
  "int2": 5,
  "limit": 100,
  "str1": "fizz",
  "str2": "buzz"
}
```
- **Response**: Array of strings with FizzBuzz sequence

### Statistics
- **URL**: `/stats`
- **Method**: GET
- **Response**: 
```json
{
  "most_frequent_request": {
    "int1": 3,
    "int2": 5,
    "limit": 100,
    "str1": "fizz",
    "str2": "buzz"
  },
  "hits": 10
}
```

## Running the Server
1. Ensure you have Go 1.21+ installed
2. Clone the repository
3. Run `go mod tidy` to download dependencies
4. Run `go run main.go`

## Testing
Use tools like Postman or curl to test the endpoints:

```bash
# FizzBuzz endpoint
curl -X POST http://localhost:8080/fizzbuzz \
     -H "Content-Type: application/json" \
     -d '{"int1":3,"int2":5,"limit":15,"str1":"fizz","str2":"buzz"}'

# Statistics endpoint
curl http://localhost:8080/stats
```

## Design Considerations
- Concurrent-safe request tracking using mutex
- Input validation
- Comprehensive error handling
- Middleware for logging, recovery, and CORS
- Flexible FizzBuzz generation logic

## Performance
- Uses Fiber, a high-performance Go web framework
- Minimal memory allocation
- Efficient request handling