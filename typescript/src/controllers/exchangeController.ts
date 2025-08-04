import { Request, Response } from 'express';
import { ExchangeResultRequest, ExchangeResultResponse } from '../types';
import { ValidationService } from '../services/validationService';
import { MemoryDatabase } from '../database/memoryDb';

export class ExchangeController {
  private validationService: ValidationService;
  private db: MemoryDatabase;

  constructor(validationService: ValidationService, db: MemoryDatabase) {
    this.validationService = validationService;
    this.db = db;
  }

  /**
   * Process exchange result
   * @param req - Request object
   * @param res - Response object
   */
  exchangeResult(req: Request, res: Response): void {
    try {
      const requestBody = req.body as ExchangeResultRequest;

      // Request validation
      if (!requestBody) {
        console.error('Failed to bind request body');
        res.status(400).json({ error: 'Failed to bind request body' });
        return;
      }

      console.log('ResultHandler:', { requestBody });

      if (requestBody.intent.outputs.length > 0) {
        // Get SessionID by UUID
        let sessionId: string;
        try {
          sessionId = this.db.getSessionIdByUuid(requestBody.uuid) || '';
          if (!sessionId) {
            throw new Error('Session not found');
          }
        } catch (error) {
          console.error('Failed to get session ID by UUID:', requestBody.uuid, error);
          res.status(400).json({ error: 'Invalid UUID or session not found' });
          return;
        }

        // Process exchange result
        try {
          // Note: In TypeScript version, we'll use a simple receipt status
          // In real implementation, you would extract this from the receipt
          const receiptStatus = 1; // Assuming success
          this.validationService.processExchangeResult(sessionId, requestBody.intent.outputs, receiptStatus);
        } catch (error) {
          console.error('Failed to process exchange result:', error);
          res.status(500).json({ error: 'Failed to process exchange result' });
          return;
        }
      }

      const response: ExchangeResultResponse = { success: true };
      res.json(response);
    } catch (error) {
      console.error('ExchangeResult error:', error);
      res.status(500).json({ error: 'Internal server error' });
    }
  }
} 