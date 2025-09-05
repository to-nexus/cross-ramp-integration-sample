package test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sample-game-backend/internal/database"
	"sample-game-backend/internal/handlers"
	"sample-game-backend/internal/middleware"
	"sample-game-backend/internal/models"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SimpleResultRequest 간단한 결과 요청 구조체
type SimpleResultRequest struct {
	UUID    string `json:"uuid"`
	TxHash  string `json:"tx_hash"`
	Receipt struct {
		Status            string       `json:"status"`
		CumulativeGasUsed string       `json:"cumulativeGasUsed"`
		LogsBloom         string       `json:"logsBloom"`
		TransactionHash   common.Hash  `json:"transactionHash"`
		GasUsed           *hexutil.Big `json:"gasUsed"`
		Logs              []struct {
			Address string   `json:"address"`
			Topics  []string `json:"topics"`
			Data    string   `json:"data"`
		} `json:"logs"`
	} `json:"receipt"`
	Intent models.ExchangeIntent `json:"intent"`
}

// generateHMACSignature 가이드에 따라 HMAC 서명을 생성하는 함수
func generateHMACSignature(data []byte, salt string) (string, error) {
	// Base64 URL 디코딩
	saltBytes, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(salt)
	if err != nil {
		return "", err
	}

	// HMAC-SHA256 생성
	h := hmac.New(sha256.New, saltBytes)
	h.Write(data)
	hashBytes := h.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}

// setupTestRouter 테스트용 라우터 설정
func setupTestRouter() *gin.Engine {
	// 데이터베이스 초기화
	err := database.InitDB()
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 미들웨어 설정
	r.Use(middleware.AuthMiddleware())

	// 라우트 설정
	api := r.Group("/api")
	{
		validate := api.Group("/validate")
		validate.Use(middleware.AuthMiddleware())
		{
			validate.POST("", handlers.ValidateUserActionHandler)
		}

		result := api.Group("/result")
		{
			result.POST("", handlers.ExchangeResultHandler)
		}
	}

	return r
}

// TestValidateResultWorkflow validate -> result 워크플로우 테스트
func TestValidateResultWorkflow(t *testing.T) {
	// 테스트 라우터 설정
	router := setupTestRouter()
	defer database.CloseDB()

	// 테스트 데이터
	testUUID := "test-workflow-uuid-123"
	testSessionID := "test-session-workflow"
	testUserSig := "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	testUserAddress := "0xB777C937fa1afC99606aFa85c5b83cFe7f82BabD"
	testProjectID := "test-project-id"
	testDigest := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	testBloom := "0x561234561234561234561234561234561234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	// HMAC 키 설정 (가이드에 따라)
	testHMACKey := "my_secret_salt_value_!@#$%^&*" // 가이드의 예시 키 사용

	// 1단계: Validate API 호출
	validateReq := models.ValidateRequest{
		UUID:        testUUID,
		UserSig:     testUserSig,
		UserAddress: testUserAddress,
		ProjectID:   testProjectID,
		Digest:      testDigest,
		Intent: models.ExchangeIntent{
			Type:   "assemble",
			Method: "mint",
			From: []models.PairAsset{
				{Type: "asset", AssetID: "asset_money", Amount: 1000},
				{Type: "asset", AssetID: "asset_gold", Amount: 500},
			},
			To: []models.PairAsset{
				{Type: "erc20", AssetID: "0x1234", Amount: 1000},
			},
		},
	}

	validateReqBytes, err := json.Marshal(validateReq)
	require.NoError(t, err, "Failed to marshal validate request")

	validateHMACSignature, err := generateHMACSignature(validateReqBytes, testHMACKey)
	require.NoError(t, err, "Failed to generate HMAC signature for validate request")
	fmt.Printf("🔐 Validate API HMAC Signature: %s\n", validateHMACSignature)

	// Validate API 요청
	validateReqHTTP := httptest.NewRequest("POST", "/api/validate", bytes.NewBuffer(validateReqBytes))
	validateReqHTTP.Header.Set("Content-Type", "application/json")
	validateReqHTTP.Header.Set("Authorization", "Bearer test_cross_auth_jwt_token")
	validateReqHTTP.Header.Set("X-Dapp-Authorization", "Bearer test_dapp_access_token")
	validateReqHTTP.Header.Set("X-Dapp-SessionID", testSessionID)
	validateReqHTTP.Header.Set("X-HMAC-SIGNATURE", validateHMACSignature)

	validateRecorder := httptest.NewRecorder()
	router.ServeHTTP(validateRecorder, validateReqHTTP)

	// Validate API 응답 확인
	assert.Equal(t, http.StatusOK, validateRecorder.Code, "Validate API should return 200")

	var validateResp models.ValidateResponse
	err = json.Unmarshal(validateRecorder.Body.Bytes(), &validateResp)
	require.NoError(t, err, "Failed to unmarshal validate response")

	assert.True(t, validateResp.Success, "Validate API should return success")
	assert.NotEmpty(t, validateResp.Data.ValidatorSig, "Validator signature should not be empty")

	fmt.Printf("✅ Validate API 성공: UUID=%s, SessionID=%s\n", testUUID, testSessionID)

	// 2단계: UUID 매핑 확인
	retrievedSessionID, err := database.GetSessionIDByUUID(testUUID)
	assert.NoError(t, err, "Should be able to retrieve session ID by UUID")
	assert.Equal(t, testSessionID, retrievedSessionID, "Retrieved session ID should match")

	fmt.Printf("✅ UUID 매핑 확인: UUID=%s -> SessionID=%s\n", testUUID, retrievedSessionID)

	// 3단계: Result API 호출
	resultReq := SimpleResultRequest{
		UUID:   testUUID,
		TxHash: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		Receipt: struct {
			Status            string       `json:"status"`
			CumulativeGasUsed string       `json:"cumulativeGasUsed"`
			LogsBloom         string       `json:"logsBloom"`
			TransactionHash   common.Hash  `json:"transactionHash"`
			GasUsed           *hexutil.Big `json:"gasUsed"`
			Logs              []struct {
				Address string   `json:"address"`
				Topics  []string `json:"topics"`
				Data    string   `json:"data"`
			} `json:"logs"`
		}{
			Status:            "0x1",
			CumulativeGasUsed: "0x1000000",
			LogsBloom:         testBloom,
			TransactionHash:   common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
			GasUsed:           &hexutil.Big{},
			Logs: []struct {
				Address string   `json:"address"`
				Topics  []string `json:"topics"`
				Data    string   `json:"data"`
			}{},
		},
		Intent: models.ExchangeIntent{
			Type:   "assemble",
			Method: "mint",
			From: []models.PairAsset{
				{Type: "asset", AssetID: "asset_money", Amount: 1000},
				{Type: "asset", AssetID: "asset_gold", Amount: 500},
			},
			To: []models.PairAsset{
				{Type: "erc20", AssetID: "0x1234", Amount: 1000},
			},
		},
	}

	resultReqBytes, err := json.Marshal(resultReq)
	require.NoError(t, err, "Failed to marshal result request")

	resultHMACSignature, err := generateHMACSignature(resultReqBytes, testHMACKey)
	require.NoError(t, err, "Failed to generate HMAC signature for result request")
	fmt.Printf("🔐 Result API HMAC Signature: %s\n", resultHMACSignature)

	// Result API 요청
	resultReqHTTP := httptest.NewRequest("POST", "/api/result", bytes.NewBuffer(resultReqBytes))
	resultReqHTTP.Header.Set("Content-Type", "application/json")
	resultReqHTTP.Header.Set("X-HMAC-SIGNATURE", resultHMACSignature)

	resultRecorder := httptest.NewRecorder()
	router.ServeHTTP(resultRecorder, resultReqHTTP)

	// Result API 응답 확인
	assert.Equal(t, http.StatusOK, resultRecorder.Code, "Result API should return 200")

	var resultResp map[string]interface{}
	err = json.Unmarshal(resultRecorder.Body.Bytes(), &resultResp)
	require.NoError(t, err, "Failed to unmarshal result response")

	assert.Equal(t, true, resultResp["success"], "Result API should return success")

	fmt.Printf("✅ Result API 성공: UUID=%s\n", testUUID)

	// 4단계: 자산 변경 확인
	sessionAssets, err := database.GetOrCreateSessionAssets(testSessionID)
	require.NoError(t, err, "Should be able to get session assets")

	// item_gem 자산이 추가되었는지 확인
	itemGemBalance, exists := sessionAssets.Assets["item_gem"]
	assert.True(t, exists, "item_gem should exist in assets")
	assert.Equal(t, "1000", itemGemBalance, "item_gem balance should be 1000")

	// asset_money가 추가되었는지 확인
	moneyBalance, exists := sessionAssets.Assets["asset_money"]
	assert.True(t, exists, "asset_money should exist in assets")
	assert.NotEmpty(t, moneyBalance, "asset_money balance should not be empty")

	fmt.Printf("✅ 자산 변경 확인: item_gem=%s, asset_money=%s\n", itemGemBalance, moneyBalance)
}

// TestValidateResultWorkflowWithInsufficientBalance 잔액 부족 시나리오 테스트
func TestValidateResultWorkflowWithInsufficientBalance(t *testing.T) {
	// 테스트 라우터 설정
	router := setupTestRouter()
	defer database.CloseDB()

	// 테스트 데이터
	testUUID := "test-insufficient-uuid"
	testSessionID := "test-session-insufficient"
	// HMAC 키 설정 (가이드에 따라)
	testHMACKey := "my_secret_salt_value_!@#$%^&*" // 가이드의 예시 키 사용

	// 1단계: Validate API 호출 (잔액 부족 시나리오)
	validateReq := models.ValidateRequest{
		UUID:        testUUID,
		UserSig:     "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		UserAddress: "0xB777C937fa1afC99606aFa85c5b83cFe7f82BabD",
		ProjectID:   "test-project-id",
		Digest:      "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		Intent: models.ExchangeIntent{
			Type:   "assemble",
			Method: "mint",
			From: []models.PairAsset{
				{Type: "asset", AssetID: "asset_money", Amount: 999999999}, // 매우 큰 금액
			},
			To: []models.PairAsset{
				{Type: "erc20", AssetID: "0x1234", Amount: 1000},
			},
		},
	}

	validateReqBytes, err := json.Marshal(validateReq)
	require.NoError(t, err, "Failed to marshal validate request")

	validateHMACSignature, err := generateHMACSignature(validateReqBytes, testHMACKey)
	require.NoError(t, err, "Failed to generate HMAC signature for validate request")
	fmt.Printf("🔐 Insufficient Balance Test - Validate API HMAC Signature: %s\n", validateHMACSignature)

	// Validate API 요청
	validateReqHTTP := httptest.NewRequest("POST", "/api/validate", bytes.NewBuffer(validateReqBytes))
	validateReqHTTP.Header.Set("Content-Type", "application/json")
	validateReqHTTP.Header.Set("Authorization", "Bearer test_cross_auth_jwt_token")
	validateReqHTTP.Header.Set("X-Dapp-Authorization", "Bearer test_dapp_access_token")
	validateReqHTTP.Header.Set("X-Dapp-SessionID", testSessionID)
	validateReqHTTP.Header.Set("X-HMAC-SIGNATURE", validateHMACSignature)

	validateRecorder := httptest.NewRecorder()
	router.ServeHTTP(validateRecorder, validateReqHTTP)

	// Validate API 응답 확인 (잔액 부족으로 실패해야 함)
	assert.Equal(t, http.StatusBadRequest, validateRecorder.Code, "Validate API should return 400 for insufficient balance")

	var validateResp models.ValidateResponse
	err = json.Unmarshal(validateRecorder.Body.Bytes(), &validateResp)
	require.NoError(t, err, "Failed to unmarshal validate response")

	assert.False(t, validateResp.Success, "Validate API should return failure")
	assert.Equal(t, "INSUFFICIENT_BALANCE", *validateResp.ErrorCode, "Error code should be INSUFFICIENT_BALANCE")

	fmt.Printf("✅ 잔액 부족 시나리오 테스트 성공: UUID=%s\n", testUUID)
}

// TestValidateResultWorkflowWithInvalidUUID 잘못된 UUID 시나리오 테스트
func TestValidateResultWorkflowWithInvalidUUID(t *testing.T) {
	// 테스트 라우터 설정
	router := setupTestRouter()
	defer database.CloseDB()

	// 테스트 데이터
	invalidUUID := "invalid-uuid-not-stored"
	// HMAC 키 설정 (가이드에 따라)
	testHMACKey := "my_secret_salt_value_!@#$%^&*" // 가이드의 예시 키 사용

	// Result API 호출 (UUID가 저장되지 않은 경우)
	resultReq := SimpleResultRequest{
		UUID:   invalidUUID,
		TxHash: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		Receipt: struct {
			Status            string       `json:"status"`
			CumulativeGasUsed string       `json:"cumulativeGasUsed"`
			LogsBloom         string       `json:"logsBloom"`
			TransactionHash   common.Hash  `json:"transactionHash"`
			GasUsed           *hexutil.Big `json:"gasUsed"`
			Logs              []struct {
				Address string   `json:"address"`
				Topics  []string `json:"topics"`
				Data    string   `json:"data"`
			} `json:"logs"`
		}{
			Status:            "0x1",
			CumulativeGasUsed: "1000000",
			TransactionHash:   common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
			GasUsed:           &hexutil.Big{},
		},
		Intent: models.ExchangeIntent{
			Type:   "assemble",
			Method: "mint",
			From: []models.PairAsset{
				{Type: "asset", AssetID: "asset_money", Amount: 1000},
				{Type: "asset", AssetID: "asset_gold", Amount: 500},
			},
			To: []models.PairAsset{
				{Type: "erc20", AssetID: "0x1234", Amount: 1000},
			},
		},
	}

	resultReqBytes, err := json.Marshal(resultReq)
	require.NoError(t, err, "Failed to marshal result request")

	// HMAC 서명 생성 (가이드에 따라)
	resultHMACSignature, err := generateHMACSignature(resultReqBytes, testHMACKey)
	require.NoError(t, err, "Failed to generate HMAC signature for result request")
	fmt.Printf("🔐 Invalid UUID Test - Result API HMAC Signature: %s\n", resultHMACSignature)

	// Result API 요청
	resultReqHTTP := httptest.NewRequest("POST", "/api/result", bytes.NewBuffer(resultReqBytes))
	resultReqHTTP.Header.Set("Content-Type", "application/json")
	resultReqHTTP.Header.Set("X-HMAC-SIGNATURE", resultHMACSignature)

	resultRecorder := httptest.NewRecorder()
	router.ServeHTTP(resultRecorder, resultReqHTTP)

	// Result API 응답 확인 (잘못된 UUID로 실패해야 함)
	assert.Equal(t, http.StatusBadRequest, resultRecorder.Code, "Result API should return 400 for invalid UUID")

	var resultResp map[string]interface{}
	err = json.Unmarshal(resultRecorder.Body.Bytes(), &resultResp)
	require.NoError(t, err, "Failed to unmarshal result response")

	assert.Equal(t, "Invalid UUID or session not found", resultResp["error"], "Error message should indicate invalid UUID")

	fmt.Printf("✅ 잘못된 UUID 시나리오 테스트 성공: UUID=%s\n", invalidUUID)
}

// TestValidateResultWorkflowConcurrent 동시 요청 테스트
func TestValidateResultWorkflowConcurrent(t *testing.T) {
	// 테스트 라우터 설정
	router := setupTestRouter()
	defer database.CloseDB()

	// 동시 요청 테스트
	done := make(chan bool, 5)
	// HMAC 키 설정 (가이드에 따라)
	testHMACKey := "my_secret_salt_value_!@#$%^&*" // 가이드의 예시 키 사용

	for i := 0; i < 5; i++ {
		go func(id int) {
			testUUID := fmt.Sprintf("concurrent-uuid-%d", id)
			testSessionID := fmt.Sprintf("concurrent-session-%d", id)

			// Validate API 호출
			validateReq := models.ValidateRequest{
				UUID:        testUUID,
				UserSig:     "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				UserAddress: "0xB777C937fa1afC99606aFa85c5b83cFe7f82BabD",
				ProjectID:   "test-project-id",
				Digest:      "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				Intent: models.ExchangeIntent{
					Type:   "assemble",
					Method: "mint",
					From: []models.PairAsset{
						{Type: "asset", AssetID: "asset_money", Amount: 100},
					},
					To: []models.PairAsset{
						{Type: "erc20", AssetID: "0x1234", Amount: 1000},
					},
				},
			}

			validateReqBytes, _ := json.Marshal(validateReq)

			validateHMACSignature, _ := generateHMACSignature(validateReqBytes, testHMACKey)

			validateReqHTTP := httptest.NewRequest("POST", "/api/validate", bytes.NewBuffer(validateReqBytes))
			validateReqHTTP.Header.Set("Content-Type", "application/json")
			validateReqHTTP.Header.Set("Authorization", "Bearer test_cross_auth_jwt_token")
			validateReqHTTP.Header.Set("X-Dapp-Authorization", "Bearer test_dapp_access_token")
			validateReqHTTP.Header.Set("X-Dapp-SessionID", testSessionID)
			validateReqHTTP.Header.Set("X-HMAC-SIGNATURE", validateHMACSignature)

			validateRecorder := httptest.NewRecorder()
			router.ServeHTTP(validateRecorder, validateReqHTTP)

			// Validate 성공 확인
			assert.Equal(t, http.StatusOK, validateRecorder.Code)

			// Result API 호출
			resultReq := SimpleResultRequest{
				UUID:   testUUID,
				TxHash: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				Receipt: struct {
					Status            string       `json:"status"`
					CumulativeGasUsed string       `json:"cumulativeGasUsed"`
					LogsBloom         string       `json:"logsBloom"`
					TransactionHash   common.Hash  `json:"transactionHash"`
					GasUsed           *hexutil.Big `json:"gasUsed"`
					Logs              []struct {
						Address string   `json:"address"`
						Topics  []string `json:"topics"`
						Data    string   `json:"data"`
					} `json:"logs"`
				}{
					Status:            "0x1",
					CumulativeGasUsed: "0x1000000",
				},
				Intent: models.ExchangeIntent{
					Type:   "disassemble",
					Method: "mint",
					From: []models.PairAsset{
						{Type: "erc20", AssetID: "0x1234", Amount: 1000},
					},
					To: []models.PairAsset{
						{AssetID: "item_gem", Amount: 100},
					},
				},
			}

			resultReqBytes, _ := json.Marshal(resultReq)
			resultHMACSignature, _ := generateHMACSignature(resultReqBytes, testHMACKey)

			resultReqHTTP := httptest.NewRequest("POST", "/api/result", bytes.NewBuffer(resultReqBytes))
			resultReqHTTP.Header.Set("Content-Type", "application/json")
			resultReqHTTP.Header.Set("X-HMAC-SIGNATURE", resultHMACSignature)

			resultRecorder := httptest.NewRecorder()
			router.ServeHTTP(resultRecorder, resultReqHTTP)

			// Result 성공 확인
			assert.Equal(t, http.StatusOK, resultRecorder.Code)

			done <- true
		}(i)
	}

	// 모든 고루틴 완료 대기
	for i := 0; i < 5; i++ {
		<-done
	}

	fmt.Printf("✅ 동시 요청 테스트 성공: 5개의 동시 워크플로우 완료\n")
}
