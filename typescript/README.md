# Sample Game Backend API - TypeScript

This is the TypeScript implementation of the Sample Game Backend API for Ramp integration.

## Features

- **Express.js** with TypeScript
- **HMAC Authentication** for secure API requests
- **In-memory Database** for session management
- **JWT-like Authentication** middleware
- **CORS** enabled
- **Helmet** for security headers
- **Morgan** for logging

## Prerequisites

- Node.js 18+ 
- npm or yarn

## Installation

```bash
cd typescript
npm install
```

## Development

```bash
# Start development server
npm run dev

# Build for production
npm run build

# Start production server
npm start
```

## Testing

```bash
# Run tests
npm test

# Run tests in watch mode
npm run test:watch
```

## API Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/assets` | Get user assets | Yes |
| POST | `/api/validate` | Validate user action | Yes + HMAC |
| POST | `/api/result` | Process exchange result | HMAC only |
| GET | `/health` | Health check | No |

## Authentication

### Required Headers

- `Authorization: Bearer <CROSS_AUTH_JWT>`
- `X-Dapp-Authorization: Bearer <DAPP_ACCESS_TOKEN>`
- `X-Dapp-SessionID: <DAPP_SESSION_ID>`

### HMAC Authentication

For POST requests to `/api/validate` and `/api/result`, include:

- `X-HMAC-Signature: <HMAC_SIGNATURE>`

The HMAC signature is calculated using SHA256 with the request body and a shared salt.

## Project Structure

```
typescript/
├── src/
│   ├── controllers/     # API controllers
│   ├── database/        # In-memory database
│   ├── middleware/      # Authentication & HMAC middleware
│   ├── routes/          # Route definitions
│   ├── services/        # Business logic
│   ├── types/           # TypeScript type definitions
│   └── index.ts         # Main server file
├── test/                # Test files
├── package.json         # Dependencies
├── tsconfig.json        # TypeScript configuration
└── README.md           # This file
```

## Environment Variables

Create a `.env` file in the typescript directory:

```env
PORT=8080
NODE_ENV=development
```

## TODO Items

1. **HMAC Salt**: Load from environment variables instead of hardcoded value
2. **JWT Validation**: Implement proper JWT token validation
3. **Database**: Replace in-memory database with persistent storage
4. **Duplicate UUID Prevention**: Implement defense logic to prevent duplicate UUIDs
5. **Validator Private Key**: Use real validator private key instead of mock
6. **Error Handling**: Improve error handling and logging
7. **Rate Limiting**: Add rate limiting middleware
8. **Input Validation**: Add comprehensive input validation
9. **Testing**: Add more comprehensive test coverage
10. **Documentation**: Add API documentation with Swagger

## Security Considerations

- HMAC salt should be loaded from environment variables
- JWT tokens should be properly validated
- Consider implementing rate limiting
- Use HTTPS in production
- Validate all input data
- Implement proper error handling without exposing sensitive information 