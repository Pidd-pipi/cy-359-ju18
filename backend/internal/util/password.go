package util

import "golang.org/x/crypto/bcrypt"

// HashPassword 使用 bcrypt 生成密码哈希。
func HashPassword(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword 校验明文密码与哈希是否匹配。
func CheckPassword(hashed, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
