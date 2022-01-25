package cmath

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func AesDecrypt(encryptedString string, keyString string) (string, error) {
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", err
	}

	enc, err := hex.DecodeString(encryptedString)
	if err != nil {
		return "", err
	}

	// Создание нового блок шифра из ключа
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Создание нового GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Получаю размер одноразового номера
	nonceSize := aesGCM.NonceSize()

	// Извлекаю одноразовый номер из зашифрованных данных
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	// Расшифровываю данные
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", plaintext), nil
}

func AesEncrypt(stringToEncrypt string, keyString string) (string, error) {
	// Преобразование строкового ключа в байты
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", err
	}

	plaintext := []byte(stringToEncrypt)

	// Создание нового блок шифра из ключа
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Создание нового GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Создание одноразового номера. Nonce должен быть из GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Шифрование данных с использованием aesGCM.Seal
	// Не сохраняю одноразовый номер, поэтому добавляю его в качестве префикса к зашифрованным данным.
	// Первым одноразовым аргументом в Seal является префикс.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext), nil
}
