package utils

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	random "math/rand"
	"net"
	"os"
	"time"

	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/golang-jwt/jwt/v4"
)

// Создает случайную строку заданной длины
func NewRandomString(size int) string {

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	b := make([]rune, size)
	for i := range b {
		b[i] = chars[random.Intn(len(chars))]
	}

	return string(b)
}

// Производит загрузку данных из файла
func LoadStorage(filePath string) ([]model.ShortURL, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(file)

	var aliases []model.ShortURL
	for {
		if data, err := buf.ReadBytes('\n'); err == nil {
			alias := &model.ShortURL{}

			err = json.Unmarshal(data, alias)
			if err != nil {
				return nil, err
			}
			aliases = append(aliases, *alias)
		} else {
			break
		}
	}
	return aliases, nil
}

// Создает файл с данными
func SaveStorage(aliases []model.ShortURL, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	for _, alias := range aliases {
		data, err := json.Marshal(alias)
		if err != nil {
			return err
		}
		_, err = file.Write(data)
		if err != nil {
			return err
		}

		_, err = file.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}
	return nil
}

// JWT структура с пользовательскими данными
type UserJWTClaims struct {
	jwt.RegisteredClaims
	Username string
}

// Производит парсинг JWT и возвращает имя пользователя
func GetUsernameFromJWT(tokenString string, key string) (string, error) {
	claims := &UserJWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("token is invalid")
	}
	return claims.Username, nil
}

// Производит создание JWT
func BuildJWTString(username string, key string, expires time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &UserJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expires)),
		},
		Username: username,
	})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Создает необходимые сущности для HTTPS
func CreateCert() (*bytes.Buffer, *bytes.Buffer, error) {
	cert := &x509.Certificate{
		// указываем уникальный номер сертификата
		SerialNumber: big.NewInt(1658),
		// заполняем базовую информацию о владельце сертификата
		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"RU"},
		},
		// разрешаем использование сертификата для 127.0.0.1 и ::1
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		// сертификат верен, начиная со времени создания
		NotBefore: time.Now(),
		// время жизни сертификата — 10 лет
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		// устанавливаем использование ключа для цифровой подписи,
		// а также клиентской и серверной авторизации
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}
	// создаём новый приватный RSA-ключ длиной 4096 бит
	// обратите внимание, что для генерации ключа и сертификата
	// используется rand.Reader в качестве источника случайных данных
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	// создаём сертификат x.509
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// кодируем сертификат и ключ в формате PEM, который
	// используется для хранения и обмена криптографическими ключами
	var certPEM *bytes.Buffer
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return nil, nil, err
	}
	var privateKeyPEM *bytes.Buffer
	err = pem.Encode(privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return nil, nil, err
	}
	return certPEM, privateKeyPEM, nil
}
