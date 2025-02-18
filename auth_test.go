package ub_auth

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoginSuccess(t *testing.T) {
	err := godotenv.Load()
	assert.NoError(t, err, "Expected a nil but got error")

	username := os.Getenv("EMAIL_OR_NIM")
	password := os.Getenv("PASSWORD")

	studentDetails, err := Auth(username, password)
	assert.NoError(t, err, "Expected a nil but got error")

	assert.NotNil(t, studentDetails, "Expected student details but got nil")
	assert.NotEmpty(t, studentDetails.NIM, "Expected NIM but got empty")
	assert.NotEmpty(t, studentDetails.FullName, "Expected Full Name but got empty")
	assert.NotEmpty(t, studentDetails.Email, "Expected Email but got empty")
}

func TestLoginInvalidPassword(t *testing.T) {
	studentDetails, err := Auth("TEST_NIM_NGUAWUR_CIK", "PASSWORD")
	assert.Error(t, err, "Expected an error but got nil")

	assert.Contains(t, err.Error(), "invalid username or password", "Expected error message to contain 'invalid username or password'")
	assert.Empty(t, studentDetails.NIM, "Expected empty but got NIM")
	assert.Empty(t, studentDetails.FullName, "Expected empty but got Full Name")
	assert.Empty(t, studentDetails.Email, "Expected empty but got Email")
}
