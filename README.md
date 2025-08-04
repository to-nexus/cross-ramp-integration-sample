# Sample Game Backend API

This project is a sample backend for Ramp integration guides by language.

It contains post-processing content for asset information queries, validation, and on-chain result reception.

Other processing is simplified. (For example, when querying asset information, the code is written to temporarily generate it in memory and expose the results. In actual implementation, assets on the pre-implemented DB should be exposed.)

Please focus on API requests and responses, and refer to TODO comments for considerations.

## ⚠️ Important Notice

**This code is provided only as a reference implementation for Ramp integration.**

- This sample code is designed to show integration patterns with the Ramp protocol
- **Modifications may be required according to each game company's requirements**
- **Users are responsible for any problems or damages caused by using this code**
- This implementation is not ready for production use and requires thorough testing and customization before deployment
- The author is not responsible for any results caused by using this sample code

## Supported Languages

This project is implemented in the following languages:

- **Go** (Go 1.23.0) - Golang implementation
- **C#** (.NET 8.0) - ASP.NET Core based
- **C++** (C++17) - CMake based

## Main Features
- Virtual game server for CROSSRAMP integration
- In-game asset query API
- User order validation API
- On-chain result post-processing API

## API Specification

### 1. Asset Information Query API

#### Request
```bash
curl -X GET "http://localhost:8080/api/assets?language=ko" \
  -H "Authorization: Bearer <CROSS_AUTH_JWT>" \
  -H "X-Dapp-Authorization: Bearer <DAPP_ACCESS_TOKEN>" \
  -H "X-Dapp-SessionID: <DAPP_SESSION_ID>"
```

#### Response (status=200)
```json
{
  "success": true,
  "errorCode": null,
  "data": {
    "v1": {
      "player_id": "C1",            // unique key
      "name": "playerName_C1",
      "wallet_address": "0xaaaa",   // mapped wallet address
      "server": "test",
      "assets": [
        {
          "id": "asset_money",      // asset id registered in ramp console
          "balance": "2000"
        },
        {
          "id": "item_gem",
          "balance": "1500"
        }
      ]
    },
    "guide": {
      "Authorization": "Bearer <token>",
      "X-Dapp-Authorization": "Bearer <token>",
      "X-Dapp-SessionID": "<session_id>",
      "message": "The guide field displays header information at request time. It is for checking if the game company and protocol are correctly matched and is not provided to game companies. For ramp frontend developer reference.",
      "session_info": {
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    }
  }
}
```

### 2. Order Validation API

#### Request
```bash
curl -X POST "http://localhost:8080/api/validate" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <CROSS_AUTH_JWT>" \
  -H "X-Dapp-Authorization: Bearer <DAPP_ACCESS_TOKEN>" \
  -H "X-Dapp-SessionID: <DAPP_SESSION_ID>" \
  -d '{
    "user_sig": "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    "user_address": "0xB777C937fa1afC99606aFa85c5b83cFe7f82BabD",
    "project_id": "acjviwejsi",
    "intent": {
      "method": "mint",
      "from": [
        { "type": "asset", "id": "asset_money", "amount": 1000 }
      ],
      "to": [
        { "type": "asset", "id": "item_gem", "amount": 1000 }
      ]
    }
  }'
```

#### Response (status=200)
```json
{
  "success": true,
  "errorCode": null,
  "data": {
    "userSig": "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    "validatorSig": "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
  }
}
```

### 3. Exchange Result Processing API

#### Request
```bash
curl -X POST "http://localhost:8080/api/result" \
  -H "Content-Type: application/json" \
  -d '{
    "uuid": "exchange-uuid-123",
    "intent": {
      "outputs": [
        { "asset_id": "item_gem", "amount": 1000 }
      ]
    },
    "receipt": {
      "status": 1
    }
  }'
```

#### Response (status=200)
```json
{
  "success": true
}
```

## Installation and Execution

### Go Version (Recommended)

#### Prerequisites
- Go 1.23.0 or higher
- Git

#### 1. Clone Repository
```bash
git clone <repository-url>
cd sample-game-backend/golang
```

#### 2. Install Dependencies
```bash
go mod tidy
```

#### 3. Run Server
```bash
go run main.go
```

The server starts on port 8080. Session-specific asset information is stored in **go-memdb** in-memory database.

### C# Version

#### Prerequisites
- .NET 8.0 SDK
- Git

#### 1. Clone Repository
```bash
git clone <repository-url>
cd sample-game-backend/csharp/SampleGameBackend
```

#### 2. Restore Dependencies
```bash
dotnet restore
```

#### 3. Run Server
```bash
dotnet run
```

### C++ Version

#### Prerequisites
- CMake 3.16 or higher
- C++17 compatible compiler
- nlohmann/json library
- OpenSSL
- Git

#### 1. Clone Repository
```bash
git clone <repository-url>
cd sample-game-backend/cpp
```

#### 2. Build
```bash
mkdir build && cd build
cmake ..
make
```

#### 3. Run Server
```bash
./sample_game_backend_cpp
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/assets` | Asset information query |
| POST | `/api/validate` | Order validation |
| POST | `/api/result` | Exchange result processing |
| GET | `/health` | Health check |

## Error Codes

| Error Code | Description |
|------------|-------------|
| INVALID_REQUEST | Invalid request format or missing required fields |
| INVALID_SESSION_ID | Missing or invalid session ID |
| DB_ERROR | Database operation failure |
| INVALID_INTENT | Invalid intent structure or method |
| UUID_MAPPING_FAILED | UUID mapping save failure |
| INSUFFICIENT_BALANCE | Insufficient asset balance for operation |
| SIGNATURE_GENERATION_FAILED | Validator signature generation failure |

## Testing

### Manual Test Examples

#### Asset Information Query
```bash
curl -X GET "http://localhost:8080/api/assets?language=ko" \
  -H "Authorization: Bearer test_cross_auth_jwt_token" \
  -H "X-Dapp-Authorization: Bearer test_dapp_access_token" \
  -H "X-Dapp-SessionID: test_session_id"
```

#### Order Validation
```bash
curl -X POST "http://localhost:8080/api/validate" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test_cross_auth_jwt_token" \
  -H "X-Dapp-Authorization: Bearer test_dapp_access_token" \
  -H "X-Dapp-SessionID: test_session_id" \
  -d '{
    "user_sig": "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    "user_address": "0xB777C937fa1afC99606aFa85c5b83cFe7f82BabD",
    "project_id": "acjviwejsi",
    "intent": {
      "method": "mint",
      "from": [
        { "type": "asset", "id": "asset_money", "amount": 1000 }
      ],
      "to": [
        { "type": "asset", "id": "item_gem", "amount": 1000 }
      ]
    }
  }'
```

#### Health Check
```bash
curl -X GET "http://localhost:8080/health"
```

## Project Structure

```
sample-game-backend/
├── golang/                 # Go implementation (main)
├── csharp/                # C# implementation
├── cpp/                   # C++ implementation
└── README_KR.md           # Korean documentation
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes
4. Add tests for new features
5. Submit a pull request

## Support

If you have questions or issues, please create an issue in the repository. 