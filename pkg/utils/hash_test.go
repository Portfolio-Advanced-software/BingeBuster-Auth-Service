package utils

import "testing"

func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	// Test data
	password := "myPassword123"

	// Hash the password
	hashedPassword := HashPassword(password)

	// Ensure hashed password is not equal to the original password
	if hashedPassword == password {
		t.Errorf("Hashed password is equal to the original password")
	}

	// Check if the password matches the hash
	match := CheckPasswordHash(password, hashedPassword)
	if !match {
		t.Errorf("Password and hash do not match")
	}
}
