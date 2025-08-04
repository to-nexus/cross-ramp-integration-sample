import { Response } from 'express';
import { AuthenticatedRequest } from '../middleware/authMiddleware';
import { AssetsResponse, Asset } from '../types';
import { MemoryDatabase } from '../database/memoryDb';

export class AssetsController {
  private db: MemoryDatabase;

  constructor(db: MemoryDatabase) {
    this.db = db;
  }

  /**
   * Get assets for a session
   * @param req - Request object
   * @param res - Response object
   */
  getAssets(req: AuthenticatedRequest, res: Response): void {
    try {
      const sessionId = req.sessionId;
      const language = (req.query.language as string) || 'ko';

      if (!sessionId) {
        res.status(400).json({
          success: false,
          errorCode: 'INVALID_SESSION_ID',
          message: 'Session ID is required'
        });
        return;
      }

      // Get session assets from database
      const sessionAssets = this.db.getSessionAssets(sessionId);
      
      // TODO: In actual implementation, fetch from real database
      // For demo purposes, create mock assets if not found
      const assets: Asset[] = sessionAssets?.assets || [
        { id: 'asset_money', balance: '2000' },
        { id: 'item_gem', balance: '1500' }
      ];

      // Create guide information
      const guide = {
        Authorization: `Bearer ${req.headers.authorization}`,
        'X-Dapp-Authorization': req.headers['x-dapp-authorization'] as string,
        'X-Dapp-SessionID': sessionId,
        message: 'guide 필드는 요청 시 헤더 정보를 표기합니다. 올바르게 게임사와 프로토콜을 맞췄는지 확인하는 용도이고 게임사에는 제공되지 않습니다. ramp frontend 개발자 참고용입니다.',
        session_info: {
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString()
        }
      };

      const response: AssetsResponse = {
        success: true,
        errorCode: null,
        data: {
          v1: {
            player_id: 'C1',
            name: 'playerName_C1',
            wallet_address: '0xaaaa',
            server: 'test',
            assets
          },
          guide
        }
      };

      res.json(response);
    } catch (error) {
      console.error('Error in getAssets:', error);
      res.status(500).json({
        success: false,
        errorCode: 'DB_ERROR',
        message: 'Internal server error'
      });
    }
  }
} 