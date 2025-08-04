import { SessionAssets, UuidMapping, Asset } from '../types';

export class MemoryDatabase {
  private sessionAssets: Map<string, SessionAssets> = new Map();
  private uuidMappings: Map<string, UuidMapping> = new Map();

  /**
   * Store session assets
   * @param sessionId - Session ID
   * @param assets - Assets array
   */
  storeSessionAssets(sessionId: string, assets: Asset[]): void {
    const now = new Date();
    this.sessionAssets.set(sessionId, {
      sessionId,
      assets,
      createdAt: now,
      updatedAt: now
    });
  }

  /**
   * Get session assets
   * @param sessionId - Session ID
   * @returns Session assets or null if not found
   */
  getSessionAssets(sessionId: string): SessionAssets | null {
    return this.sessionAssets.get(sessionId) || null;
  }

  /**
   * Store UUID mapping
   * @param uuid - UUID
   * @param sessionId - Session ID
   */
  storeUuidMapping(uuid: string, sessionId: string): void {
    const now = new Date();
    this.uuidMappings.set(uuid, {
      uuid,
      sessionId,
      createdAt: now
    });
  }

  /**
   * Get session ID by UUID
   * @param uuid - UUID
   * @returns Session ID or null if not found
   */
  getSessionIdByUuid(uuid: string): string | null {
    const mapping = this.uuidMappings.get(uuid);
    return mapping ? mapping.sessionId : null;
  }

  /**
   * Update session assets
   * @param sessionId - Session ID
   * @param assets - New assets array
   */
  updateSessionAssets(sessionId: string, assets: Asset[]): void {
    const existing = this.sessionAssets.get(sessionId);
    if (existing) {
      existing.assets = assets;
      existing.updatedAt = new Date();
    }
  }

  /**
   * Remove UUID mapping
   * @param uuid - UUID to remove
   */
  removeUuidMapping(uuid: string): void {
    this.uuidMappings.delete(uuid);
  }

  /**
   * Get all session assets (for debugging)
   * @returns All session assets
   */
  getAllSessionAssets(): SessionAssets[] {
    return Array.from(this.sessionAssets.values());
  }

  /**
   * Get all UUID mappings (for debugging)
   * @returns All UUID mappings
   */
  getAllUuidMappings(): UuidMapping[] {
    return Array.from(this.uuidMappings.values());
  }

  /**
   * Clear all data (for testing)
   */
  clear(): void {
    this.sessionAssets.clear();
    this.uuidMappings.clear();
  }
} 