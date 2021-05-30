package tokenstorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/storage/memory"
	"go.uber.org/zap"
	"sync"
	"testing"
	"time"
)

func createTSWithTestData(t *testing.T) (*TokenStorage, map[string]interface{}) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	ts, err := NewTokenStorage([]byte("testsecretkey"), logger, memory.New(), time.Second, time.Second*2)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	data := map[string]interface{}{
		"test": "12345678",
		"tmp":  "87654321",
	}
	return ts, data
}

func checkToken(ts *TokenStorage, tokenString string, claims map[string]interface{}) error {
	tokenData, err := ts.ParseToken(tokenString)
	if err == nil {
		for key, value := range claims {
			if tokenData[key] != value {
				return fmt.Errorf("wrong value in claims: got %v, expected %v", tokenData[key], value)
			}
		}
		return nil
	}
	return err
}

func Test_TokenStorage_NewTokenPair(t *testing.T) {
	t.Parallel()
	ts, _ := createTSWithTestData(t)
	_, _, err := ts.NewTokenPair(map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_TokenStorage_ParseTokenCorrect(t *testing.T) {
	t.Parallel()
	ts, data := createTSWithTestData(t)
	access, refresh, err := ts.NewTokenPair(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err = checkToken(ts, access, data); err != nil {
		t.Fatalf("invalid access token: %v", err)
	}
	if err = checkToken(ts, refresh, data); err != nil {
		t.Fatalf("invalid refresh token: %v", err)
	}
}

func Test_TokenStorage_ExpireTokenAccess(t *testing.T) {
	t.Parallel()
	ts, data := createTSWithTestData(t)
	access, _, err := ts.NewTokenPair(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := ts.ExpireToken(access); !errors.Is(err, ErrExpireNonRefresh) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_TokenStorage_ExpireTokenRefresh(t *testing.T) {
	t.Parallel()
	ts, data := createTSWithTestData(t)
	access, refresh, err := ts.NewTokenPair(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := ts.ExpireToken(refresh); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := checkToken(ts, access, data); !errors.Is(err, ErrTokenExpired) {
		t.Fatalf("access must be expired")
	}
	if err := checkToken(ts, refresh, data); !errors.Is(err, ErrTokenExpired) {
		t.Fatalf("refresh must be expired")
	}
}

func Test_TokenStorage_ParseWrongKey(t *testing.T) {
	t.Parallel()
	ts, _ := createTSWithTestData(t)
	if !errors.Is(checkToken(ts, "erfermfjiermfi", map[string]interface{}{}), ErrInvalidToken) {
		t.Fatalf("treated wrong token as correct")
	}
}

func Test_TokenStorage_ParseWrongSignature(t *testing.T) {
	t.Parallel()
	ts, data := createTSWithTestData(t)
	ts1, err := NewTokenStorage([]byte("testsecretkey1"), nil, memory.New(), time.Second, time.Second*2)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	access, refresh, err := ts1.NewTokenPair(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := checkToken(ts, access, data); !errors.Is(err, ErrInvalidSignature) {
		t.Fatalf("signature not checked, err: %v", err)
	}
	if err := checkToken(ts, refresh, data); !errors.Is(err, ErrInvalidSignature) {
		t.Fatalf("signature not checked, err: %v", err)
	}
}

func Test_TokenStorage_TokenExpirationCheck(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	t.Parallel()
	ts, data := createTSWithTestData(t)
	access, refresh, err := ts.NewTokenPair(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	token, err := ts.ParseToken(access)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	expAt, err := token["exp"].(json.Number).Int64()
	if err != nil {
		t.Fatalf("failed to parse exp from token")
	}
	time.Sleep(time.Until(time.Unix(expAt, 0).Add(time.Second)))
	if err = checkToken(ts, access, data); !errors.Is(err, ErrTokenExpired) {
		t.Fatalf("access must be expired")
	}
	token, err = ts.ParseToken(refresh)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	expAt, err = token["exp"].(json.Number).Int64()
	if err != nil {
		t.Fatalf("failed to parse exp from token")
	}
	time.Sleep(time.Until(time.Unix(expAt, 0).Add(time.Second)))
	if err := checkToken(ts, refresh, data); !errors.Is(err, ErrTokenExpired) {
		t.Fatalf("refresh must be expired")
	}
}

const iterCount = 1000
const coroutinesCount = 100

func TestTokenStorage_ParallelAccess(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	t.Parallel()
	ts, data := createTSWithTestData(t)
	var wg sync.WaitGroup
	var globalErrMutex sync.Mutex
	var globalErr error
	for coroutine := 0; coroutine < coroutinesCount; coroutine++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for iter := 0; iter < iterCount; iter++ {
				access, refresh, err := ts.NewTokenPair(data)
				if err != nil {
					globalErrMutex.Lock()
					globalErr = fmt.Errorf("unexpected error: %w", err)
					globalErrMutex.Unlock()
				}
				if err = checkToken(ts, access, data); err != nil {
					globalErrMutex.Lock()
					globalErr = fmt.Errorf("invalid access token: %w", err)
					globalErrMutex.Unlock()
				}
				if err = checkToken(ts, refresh, data); err != nil {
					globalErrMutex.Lock()
					globalErr = fmt.Errorf("invalid refresh token: %w", err)
					globalErrMutex.Unlock()
				}
				if err := ts.ExpireToken(refresh); err != nil {
					globalErrMutex.Lock()
					globalErr = fmt.Errorf("unexpected error: %w", err)
					globalErrMutex.Unlock()
				}
				if err := checkToken(ts, access, data); !errors.Is(err, ErrTokenExpired) {
					globalErrMutex.Lock()
					globalErr = fmt.Errorf("access must be expired")
					globalErrMutex.Unlock()
				}
				if err := checkToken(ts, refresh, data); !errors.Is(err, ErrTokenExpired) {
					globalErrMutex.Lock()
					globalErr = fmt.Errorf("refresh must be expired")
					globalErrMutex.Unlock()
				}
			}
		}()
	}
	wg.Wait()
	if globalErr != nil {
		t.Fatalf("%v", globalErr)
	}
}
