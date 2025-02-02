package utils

import (
	"testing"
	"time"
)

func BenchmarkGenerateString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewRandomString(16)
	}
}

func BenchmarkBuildJWT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := BuildJWTString(NewRandomString(16), NewRandomString(8), time.Duration(10))
		if err != nil {
			return
		}
	}
}

func BenchmarkCreateAndGetUserFromJWT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		key := NewRandomString(8)
		token, err := BuildJWTString(NewRandomString(16), key, time.Duration(10))
		if err != nil {
			return
		}
		_, err = GetUsernameFromJWT(token, key)
		if err != nil {
			return
		}
	}

}
