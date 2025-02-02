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
		BuildJWTString(NewRandomString(16), NewRandomString(8), time.Duration(10))
	}
}

func BenchmarkCreateAndGetUserFromJWT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		key := NewRandomString(8)
		token, _ := BuildJWTString(NewRandomString(16), key, time.Duration(10))
		GetUsernameFromJWT(token, key)
	}

}
