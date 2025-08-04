package database

import (
	"fmt"
	"sample-game-backend/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreAndGetUUIDMapping(t *testing.T) {
	// DB 초기화
	err := InitDB()
	require.NoError(t, err, "Failed to initialize test database")
	defer CloseDB()

	// 테스트 데이터
	testUUID := "test-uuid-123"
	testSessionID := "session-test-456"

	// UUID 매핑 저장 테스트
	err = StoreUUIDMapping(testUUID, testSessionID)
	assert.NoError(t, err, "Failed to store UUID mapping")

	// UUID로 SessionID 조회 테스트
	retrievedSessionID, err := GetSessionIDByUUID(testUUID)
	assert.NoError(t, err, "Failed to get session ID by UUID")
	assert.Equal(t, testSessionID, retrievedSessionID, "Retrieved session ID should match stored session ID")

	// 존재하지 않는 UUID 조회 테스트
	_, err = GetSessionIDByUUID("non-existent-uuid")
	assert.Error(t, err, "Should return error for non-existent UUID")
	assert.Contains(t, err.Error(), "uuid mapping not found", "Error message should indicate mapping not found")
}

func TestGetOrCreateSessionAssets(t *testing.T) {
	// DB 초기화
	err := InitDB()
	require.NoError(t, err, "Failed to initialize test database")
	defer CloseDB()

	// 테스트 세션 ID
	testSessionID := "test-session-789"

	// 새로운 세션 자산 생성 테스트
	sessionAssets, err := GetOrCreateSessionAssets(testSessionID)
	assert.NoError(t, err, "Failed to get or create session assets")
	assert.NotNil(t, sessionAssets, "Session assets should not be nil")
	assert.Equal(t, testSessionID, sessionAssets.SessionID, "Session ID should match")
	assert.NotEmpty(t, sessionAssets.Assets, "Assets should not be empty")

	// 기존 세션 자산 조회 테스트
	sessionAssets2, err := GetOrCreateSessionAssets(testSessionID)
	assert.NoError(t, err, "Failed to get existing session assets")
	assert.Equal(t, sessionAssets.SessionID, sessionAssets2.SessionID, "Session IDs should match")
	assert.Equal(t, sessionAssets.Assets, sessionAssets2.Assets, "Assets should be the same")

	// 자산 종류 확인
	expectedAssets := []string{"asset_money", "asset_gold", "item_gem", "item_banana", "asset_silver", "item_apple", "item_fish", "item_branch", "item_horn", "item_maple"}
	for _, expectedAsset := range expectedAssets {
		_, exists := sessionAssets.Assets[expectedAsset]
		assert.True(t, exists, "Asset %s should exist", expectedAsset)
	}
}

func TestCheckAndDeductAssets(t *testing.T) {
	// DB 초기화
	err := InitDB()
	require.NoError(t, err, "Failed to initialize test database")
	defer CloseDB()

	// 테스트 세션 ID
	testSessionID := "test-session-deduct"

	// 세션 자산 생성
	sessionAssets, err := GetOrCreateSessionAssets(testSessionID)
	require.NoError(t, err, "Failed to create session assets")

	// 초기 자산 잔액 확인
	initialMoneyBalance := sessionAssets.Assets["asset_money"]
	initialGoldBalance := sessionAssets.Assets["asset_gold"]

	// 자산 차감 테스트
	deductAssets := []struct {
		Type   string `json:"type" binding:"required"`
		ID     string `json:"id" binding:"required"`
		Amount int    `json:"amount" binding:"required"`
	}{
		{Type: "asset", ID: "asset_money", Amount: 1000},
		{Type: "asset", ID: "asset_gold", Amount: 500},
	}

	err = CheckAndDeductAssets(testSessionID, deductAssets)
	assert.NoError(t, err, "Failed to deduct assets")

	// 차감 후 자산 잔액 확인
	updatedSessionAssets, err := GetOrCreateSessionAssets(testSessionID)
	require.NoError(t, err, "Failed to get updated session assets")

	// 잔액이 차감되었는지 확인
	assert.NotEqual(t, initialMoneyBalance, updatedSessionAssets.Assets["asset_money"], "Money balance should be deducted")
	assert.NotEqual(t, initialGoldBalance, updatedSessionAssets.Assets["asset_gold"], "Gold balance should be deducted")
}

func TestAddAssets(t *testing.T) {
	// DB 초기화
	err := InitDB()
	require.NoError(t, err, "Failed to initialize test database")
	defer CloseDB()

	// 테스트 세션 ID
	testSessionID := "test-session-add"

	// 세션 자산 생성
	sessionAssets, err := GetOrCreateSessionAssets(testSessionID)
	require.NoError(t, err, "Failed to create session assets")

	// 초기 자산 잔액 확인
	initialMoneyBalance := sessionAssets.Assets["asset_money"]

	// 자산 증가 테스트
	addAssets := []models.PairAsset{
		{AssetID: "asset_money", Amount: 1000},
		{AssetID: "asset_gold", Amount: 500},
		{AssetID: "new_asset", Amount: 200}, // 새로운 자산
	}

	err = AddAssets(testSessionID, addAssets)
	assert.NoError(t, err, "Failed to add assets")

	// 증가 후 자산 잔액 확인
	updatedSessionAssets, err := GetOrCreateSessionAssets(testSessionID)
	require.NoError(t, err, "Failed to get updated session assets")

	// 잔액이 증가되었는지 확인
	assert.NotEqual(t, initialMoneyBalance, updatedSessionAssets.Assets["asset_money"], "Money balance should be increased")

	// 새로운 자산이 생성되었는지 확인
	newAssetBalance, exists := updatedSessionAssets.Assets["new_asset"]
	assert.True(t, exists, "New asset should exist")
	assert.Equal(t, "200", newAssetBalance, "New asset balance should be 200")
}

func TestUUIDMappingWorkflow(t *testing.T) {
	// DB 초기화
	err := InitDB()
	require.NoError(t, err, "Failed to initialize test database")
	defer CloseDB()

	// 시나리오: validate -> result 워크플로우 테스트
	testUUID := "workflow-test-uuid"
	testSessionID := "workflow-session"

	// 1. UUID 매핑 저장 (validate 단계)
	err = StoreUUIDMapping(testUUID, testSessionID)
	assert.NoError(t, err, "Failed to store UUID mapping in validate step")

	// 2. 세션 자산 생성
	_, err = GetOrCreateSessionAssets(testSessionID)
	require.NoError(t, err, "Failed to create session assets")

	// 3. UUID로 SessionID 조회 (result 단계)
	retrievedSessionID, err := GetSessionIDByUUID(testUUID)
	assert.NoError(t, err, "Failed to get session ID by UUID in result step")
	assert.Equal(t, testSessionID, retrievedSessionID, "Retrieved session ID should match")

	// 4. 자산 증가 처리 (result 단계)
	addAssets := []models.PairAsset{
		{AssetID: "asset_money", Amount: 1000},
		{AssetID: "asset_gold", Amount: 500},
	}

	err = AddAssets(retrievedSessionID, addAssets)
	assert.NoError(t, err, "Failed to add assets in result step")

	// 5. 최종 자산 확인
	finalSessionAssets, err := GetOrCreateSessionAssets(retrievedSessionID)
	require.NoError(t, err, "Failed to get final session assets")

	// 자산이 증가되었는지 확인
	moneyBalance := finalSessionAssets.Assets["asset_money"]
	goldBalance := finalSessionAssets.Assets["asset_gold"]
	assert.NotEmpty(t, moneyBalance, "Money balance should exist")
	assert.NotEmpty(t, goldBalance, "Gold balance should exist")
}

func TestConcurrentAccess(t *testing.T) {
	// DB 초기화
	err := InitDB()
	require.NoError(t, err, "Failed to initialize test database")
	defer CloseDB()

	// 동시 접근 테스트
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			sessionID := fmt.Sprintf("concurrent-session-%d", id)

			// 자산 생성
			_, err := GetOrCreateSessionAssets(sessionID)
			assert.NoError(t, err)

			// 자산 추가
			addAssets := []models.PairAsset{
				{AssetID: "asset_money", Amount: uint(id * 100)},
			}
			err = AddAssets(sessionID, addAssets)
			assert.NoError(t, err)

			done <- true
		}(i)
	}

	// 모든 고루틴 완료 대기
	for i := 0; i < 10; i++ {
		<-done
	}
}
