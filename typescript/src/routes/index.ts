import { Router } from 'express';
import { AssetsController } from '../controllers/assetsController';
import { ValidationController } from '../controllers/validationController';
import { ExchangeController } from '../controllers/exchangeController';
import { authMiddleware } from '../middleware/authMiddleware';
import { hmacMiddleware } from '../middleware/hmacMiddleware';
import { MemoryDatabase } from '../database/memoryDb';
import { ValidationService } from '../services/validationService';

export function setupRoutes(): Router {
  const router = Router();
  
  // Initialize services
  const db = new MemoryDatabase();
  const validationService = new ValidationService(db);
  
  // Initialize controllers
  const assetsController = new AssetsController(db);
  const validationController = new ValidationController(validationService, db);
  const exchangeController = new ExchangeController(validationService, db);

  // API routes
  const apiRouter = Router();

  // Assets endpoint (requires authentication)
  apiRouter.get('/assets', authMiddleware, (req, res) => {
    assetsController.getAssets(req, res);
  });

  // Validation endpoint (requires authentication and HMAC)
  apiRouter.post('/validate', authMiddleware, hmacMiddleware, (req, res) => {
    validationController.validateUserAction(req, res);
  });

  // Exchange result endpoint (requires HMAC)
  apiRouter.post('/result', hmacMiddleware, (req, res) => {
    exchangeController.exchangeResult(req, res);
  });

  // Mount API routes
  router.use('/api', apiRouter);

  // Health check endpoint
  router.get('/health', (req, res) => {
    res.json({
      status: 'healthy',
      message: 'Server is running normally'
    });
  });

  return router;
} 