## Batch Api Server

API for combining multiple requests into a single batch request.

### Usage

`POST /batch`

| Property Name |  Type  | Description                     |
|:-------------:|:------:|---------------------------------|
|   requests    | `list` | List of requests to be batched. |

| Request Property Name |   Type   | Description                                                                              |
|:---------------------:|:--------:|------------------------------------------------------------------------------------------|
|          id           | `string` | Identifier for the request.                                                              |
|         path          | `string` | Path of the request.                                                                     |
|        method         | `string` | Method (`GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `HEAD`, `OPTIONS`, `TRACE`, `CONNECT`). |
|         body          | `object` | Body of the request (not required).                                                      |
|        headers        | `object` | Headers of the request (not required).                                                   |

#### Example

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

### Response

| Property Name |  Type  | Description        |
|:-------------:|:------:|--------------------|
|   responses   | `list` | List of responses. |

| Response Property Name |   Type   | Description                                                            |
|:----------------------:|:--------:|------------------------------------------------------------------------|
|           id           | `string` | Identifier for the response.                                           |
|       statusCode       |  `int`   | Status code of the response (`500` if requests fails for some reason). |
|          body          | `object` | Body of the response.                                                  |
|        headers         | `object` | Headers of the response.                                               |

#### Example

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

### Configuration

Configuration file `config.yaml` is automatically created in the current working directory.

|      Property Name      |   Type   | Description                                        |
|:-----------------------:|:--------:|----------------------------------------------------|
|          port           |  `int`   | Port number for the server.                        |
|        base_url         | `string` | Base URL for the requests.                         |
| max_concurrent_requests |  `int`   | Maximum number of concurrent requests.             |
|      allowed_paths      |  `list`  | List of regular expressions for the allowed paths. |

#### Example

```yaml
port: 1323
base_url: http://localhost:8080
max_concurrent_requests: 10
allowed_paths:
  - ^/api/user.*
```
