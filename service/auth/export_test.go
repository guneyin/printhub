package auth

func GenerateToken(uid string) (string, error) {
	return generateToken(uid)
}

func VerifyToken(hash string) (string, error) {
	return verifyToken(hash)
}
