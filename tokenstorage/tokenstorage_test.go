package tokenstorage

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/storage/memory"
	"sync"
	"testing"
	"time"
)

var ts = TokenStorage{
	secret:            []byte("testsecretkey"),
	storage:           memory.New(),
	accessExpiration:  time.Second,
	refreshExpiration: time.Second * 2,
}

var data = map[string]interface{}{
	"test": "12345678",
	"tmp":  "87654321",
}

func Test_TokenStorage_NewTokenPair(t *testing.T) {
	_, _, err := ts.NewTokenPair(map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func checkToken(tokenString string, claims map[string]interface{}) error {
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

func Test_TokenStorage_ParseTokenCorrect(t *testing.T) {
	access, refresh, err := ts.NewTokenPair(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err = checkToken(access, data); err != nil {
		t.Fatalf("invalid access token: %v", err)
	}
	if err = checkToken(refresh, data); err != nil {
		t.Fatalf("invalid refresh token: %v", err)
	}
}

func Test_TokenStorage_ExpireTokenAccess(t *testing.T) {
	access, refresh, err := ts.NewTokenPair(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := ts.ExpireToken(access); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := checkToken(access, data); err != ErrTokenExpired {
		t.Fatalf("access must be expired")
	}
	if err := checkToken(refresh, data); err != ErrTokenExpired {
		t.Fatalf("refresh must be expired")
	}
}

func Test_TokenStorage_ExpireTokenRefresh(t *testing.T) {
	access, refresh, err := ts.NewTokenPair(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := ts.ExpireToken(refresh); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := checkToken(access, data); err != ErrTokenExpired {
		t.Fatalf("access must be expired")
	}
	if err := checkToken(refresh, data); err != ErrTokenExpired {
		t.Fatalf("refresh must be expired")
	}
}

func Test_TokenStorage_ParseWrongKey(t *testing.T) {
	if checkToken("erfermfjiermfi", map[string]interface{}{}) == ErrInvalidToken {
		t.Fatalf("treated wrong token as correct")
	}
}

func Test_TokenStorage_ParseWrongSignature(t *testing.T) {
	var ts1 = TokenStorage{
		secret:            []byte("testsecretkey123"),
		storage:           memory.New(),
		accessExpiration:  time.Second,
		refreshExpiration: time.Second * 2,
	}
	access, refresh, err := ts1.NewTokenPair(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := checkToken(access, data); err != ErrInvalidSignature {
		t.Fatalf("signature not checked, err: %v", err)
	}
	if err := checkToken(refresh, data); err != ErrInvalidSignature {
		t.Fatalf("signature not checked, err: %v", err)
	}
}

func Test_TokenStorage_TokenExpirationCheck(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
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
	time.Sleep(time.Unix(expAt, 0).Add(time.Second).Sub(time.Now()))
	if err := checkToken(access, data); err != ErrTokenExpired {
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
	time.Sleep(time.Unix(expAt, 0).Add(time.Second).Sub(time.Now()))
	if err := checkToken(refresh, data); err != ErrTokenExpired {
		t.Fatalf("refresh must be expired")
	}
}

const iterCount = 1000
const coroutinesCount = 100

func TestTokenStorage_ParallelAccess(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	var wg sync.WaitGroup
	for coroutine := 0; coroutine < coroutinesCount; coroutine++ {
		wg.Add(1)
		go func() {
			for iter := 0; iter < iterCount; iter++ {
				Test_TokenStorage_ParseTokenCorrect(t)
				Test_TokenStorage_ExpireTokenAccess(t)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
