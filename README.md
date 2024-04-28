## Batch Api Server

API for combining multiple requests into a single batch request.

### Example Request
```http request
POST /batch HTTP/1.1
Content-Type: application/json
Host: localhost:1323

{
  "requests": [
    {
      "id": "1",
      "path": "/api/user",
      "method": "POST",
      "body": {
        "name": "exampleName"
      },
      "headers": {
        "Content-type": "application/json; charset=UTF-8"
      }
    },
    {
      "id": "2",
      "path": "/api/user",
      "method": "GET"
    }
  ]
}
```

### Example Response
```json
{
  "responses": [
    {
      "id": "1",
      "statusCode": 201,
      "body": {
        "id": "2",
        "name": "exampleName"
      },
      "headers": {
        "Content-Type": [
          "application/json"
        ]
      }
    },
    {
      "id": "2",
      "statusCode": 200,
      "body": [
        {
          "id": "1",
          "name": "exampleName"
        },
        {
          "id": "2",
          "name": "exampleName"
        }
      ],
      "headers": {
        "Content-Type": [
          "application/json"
        ]
      }
    }
  ]
}
```

### Example Configuration
```yaml
port: 1323
base_url: http://localhost:8080
max_concurrent_requests: 10
allowed_paths:
    - ^/api/user.*
```
