package test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/assert"
)

// 테스트용 구조체 정의
type User struct {
	ID    string
	Name  string
	Email string
	Age   int
}

// go-memdb 테스트 함수
func TestMemDBPutGet(t *testing.T) {
	// 스키마 정의
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"user": {
				Name: "user",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
					"email": {
						Name:    "email",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Email"},
					},
					"age": {
						Name:    "age",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "Age"},
					},
				},
			},
		},
	}

	// 데이터베이스 생성
	db, err := memdb.NewMemDB(schema)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// 쓰기 트랜잭션 시작
	txn := db.Txn(true)

	// 사용자 데이터 삽입
	users := []*User{
		{ID: "1", Name: "김철수", Email: "kim@example.com", Age: 25},
		{ID: "2", Name: "이영희", Email: "lee@example.com", Age: 30},
		{ID: "3", Name: "박민수", Email: "park@example.com", Age: 28},
		{ID: "4", Name: "정수진", Email: "jung@example.com", Age: 35},
	}

	for _, user := range users {
		err := txn.Insert("user", user)
		assert.NoError(t, err)
	}

	// 트랜잭션 커밋
	txn.Commit()

	// 읽기 전용 트랜잭션 시작
	txn = db.Txn(false)
	defer txn.Abort()

	// ID로 사용자 조회
	raw, err := txn.First("user", "id", "1")
	assert.NoError(t, err)
	assert.NotNil(t, raw)

	user := raw.(*User)
	assert.Equal(t, "김철수", user.Name)
	assert.Equal(t, "kim@example.com", user.Email)
	assert.Equal(t, 25, user.Age)

	fmt.Printf("조회된 사용자: ID=%s, 이름=%s, 이메일=%s, 나이=%d\n",
		user.ID, user.Name, user.Email, user.Age)

	// 이메일로 사용자 조회
	raw, err = txn.First("user", "email", "lee@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, raw)

	user = raw.(*User)
	assert.Equal(t, "이영희", user.Name)
	assert.Equal(t, "2", user.ID)

	fmt.Printf("이메일로 조회된 사용자: ID=%s, 이름=%s\n", user.ID, user.Name)

	// 모든 사용자 조회
	it, err := txn.Get("user", "id")
	assert.NoError(t, err)

	fmt.Println("모든 사용자 목록:")
	count := 0
	for obj := it.Next(); obj != nil; obj = it.Next() {
		user := obj.(*User)
		fmt.Printf("  ID=%s, 이름=%s, 이메일=%s, 나이=%d\n",
			user.ID, user.Name, user.Email, user.Age)
		count++
	}
	assert.Equal(t, 4, count)

	// 나이 범위로 조회 (25-30세)
	it, err = txn.LowerBound("user", "age", 25)
	assert.NoError(t, err)

	fmt.Println("25-30세 사용자:")
	ageCount := 0
	for obj := it.Next(); obj != nil; obj = it.Next() {
		user := obj.(*User)
		if user.Age > 30 {
			break
		}
		fmt.Printf("  %s (나이: %d)\n", user.Name, user.Age)
		ageCount++
	}
	assert.Equal(t, 3, ageCount) // 25, 28, 30세 사용자
}

// 업데이트 테스트
func TestMemDBUpdate(t *testing.T) {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"user": {
				Name: "user",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	assert.NoError(t, err)

	// 쓰기 트랜잭션
	txn := db.Txn(true)

	// 사용자 삽입
	user := &User{ID: "1", Name: "김철수", Email: "kim@example.com", Age: 25}
	err = txn.Insert("user", user)
	assert.NoError(t, err)

	txn.Commit()

	// 업데이트 트랜잭션
	txn = db.Txn(true)

	// 기존 사용자 조회
	raw, err := txn.First("user", "id", "1")
	assert.NoError(t, err)
	assert.NotNil(t, raw)

	// 사용자 정보 업데이트
	updatedUser := raw.(*User)
	updatedUser.Name = "김철수(수정됨)"
	updatedUser.Age = 26

	// 기존 데이터 삭제 후 새 데이터 삽입
	err = txn.Delete("user", updatedUser)
	assert.NoError(t, err)

	err = txn.Insert("user", updatedUser)
	assert.NoError(t, err)

	txn.Commit()

	// 읽기 트랜잭션으로 확인
	txn = db.Txn(false)
	defer txn.Abort()

	raw, err = txn.First("user", "id", "1")
	assert.NoError(t, err)
	assert.NotNil(t, raw)

	user = raw.(*User)
	assert.Equal(t, "김철수(수정됨)", user.Name)
	assert.Equal(t, 26, user.Age)

	fmt.Printf("업데이트된 사용자: ID=%s, 이름=%s, 나이=%d\n",
		user.ID, user.Name, user.Age)
}

// 삭제 테스트
func TestMemDBDelete(t *testing.T) {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"user": {
				Name: "user",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	assert.NoError(t, err)

	// 쓰기 트랜잭션
	txn := db.Txn(true)

	// 사용자 삽입
	user := &User{ID: "1", Name: "김철수", Email: "kim@example.com", Age: 25}
	err = txn.Insert("user", user)
	assert.NoError(t, err)

	txn.Commit()

	// 삭제 트랜잭션
	txn = db.Txn(true)

	// 사용자 조회 후 삭제
	raw, err := txn.First("user", "id", "1")
	assert.NoError(t, err)
	assert.NotNil(t, raw)

	err = txn.Delete("user", raw)
	assert.NoError(t, err)

	txn.Commit()

	// 읽기 트랜잭션으로 삭제 확인
	txn = db.Txn(false)
	defer txn.Abort()

	raw, err = txn.First("user", "id", "1")
	assert.NoError(t, err)
	assert.Nil(t, raw) // 삭제되었으므로 nil이어야 함

	fmt.Println("사용자가 성공적으로 삭제되었습니다.")
}
