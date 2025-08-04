import express from 'express';
import cors from 'cors';
import helmet from 'helmet';
import morgan from 'morgan';
import dotenv from 'dotenv';
import { setupRoutes } from './routes';

// Load environment variables
dotenv.config();

const app = express();
const PORT = process.env.PORT || 8080;

// Middleware
app.use(helmet());
app.use(morgan('combined'));
app.use(cors({
  origin: '*',
  methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  allowedHeaders: [
    'Authorization', 
    'X-Dapp-Authorization', 
    'X-Dapp-SessionID', 
    'Content-Type',
    'X-HMAC-Signature'
  ]
}));

// Body parsing middleware
app.use(express.json({ limit: '10mb' }));
app.use(express.urlencoded({ extended: true }));

// Setup routes
app.use(setupRoutes());

// Error handling middleware
app.use((err: Error, req: express.Request, res: express.Response, next: express.NextFunction) => {
  console.error('Unhandled error:', err);
  res.status(500).json({
    success: false,
    errorCode: 'INTERNAL_ERROR',
    message: 'Internal server error'
  });
});

// 404 handler
app.use('*', (req, res) => {
  res.status(404).json({
    success: false,
    errorCode: 'NOT_FOUND',
    message: 'Endpoint not found'
  });
});

// Start server
app.listen(PORT, () => {
  console.log('🚀 Server started on port', PORT);
  console.log('📡 API endpoint: http://localhost:8080/api/assets?language=ko');
  console.log('🔐 Order validation API: http://localhost:8080/api/validate');
  console.log('💚 Health check: http://localhost:8080/health');
  console.log('💾 Session-specific asset information is stored in memory');
});

export default app; 