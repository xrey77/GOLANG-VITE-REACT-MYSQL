package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"src/golang_mysql8/dto"
	authmiddleware "src/golang_mysql8/middleware"
	auth "src/golang_mysql8/middleware/auth"
	middleware "src/golang_mysql8/middleware/products"
	users "src/golang_mysql8/middleware/users"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/pquerna/otp/totp"
)

// ===TEST CASE 1 - USER LOGIN==============
func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 2. Mock Data
	loginReq := dto.UserLogin{
		Username: "Rey",
		Password: "rey",
	}
	body, _ := json.Marshal(loginReq)

	// 3. Create Request and Recorder
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Create a real http request to satisfy BindJSON
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/auth/signin", bytes.NewBuffer(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	// 4. Call the Handler
	// NOTE: This will currently FAIL if your DB/Kafka aren't running.
	auth.Login(ctx)

	// 5. Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Login Successfull.", response["message"])
}

// ====TEST CASE 2 - USER REGISTRATION======
func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 2. Mock Data
	registrationReq := dto.UserRegister{
		Firstname: "Roger",
		Lastname:  "Galam",
		Email:     "roger@yahoo.com",
		Mobile:    "324234243",
		Username:  "Roger",
		Password:  "rey",
	}
	body, _ := json.Marshal(registrationReq)

	// 3. Create Request and Recorder
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Create a real http request to satisfy BindJSON
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	// 4. Call the Handler
	// NOTE: This will currently FAIL if your DB/Kafka aren't running.
	auth.Register(ctx)

	// 5. Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Registration Successful, please login now.", response["message"])
}

// ====TEST CASE 3 - GET USER ID==============
func TestGetUserid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	authGuard := router.Group("/api")
	// Use your actual middleware here
	authGuard.Use(authmiddleware.AuthMiddleware())
	{
		authGuard.GET("/getuserid/:id", users.GetUserid)
	}

	// 2. Create the Request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/getuserid/1", nil)

	bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJleUB5YWhvby5jb20iLCJpc3MiOiJCQVJDTEFZUyBCQU5LIiwic3ViIjoidXNlcl9hdXRoZW50aWNhdGlvbiIsImV4cCI6MTc3Njg2NTczOSwibmJmIjoxNzc2ODM2OTM5LCJpYXQiOjE3NzY4MzY5MzksImp0aSI6IjEifQ.6YLLoQYHGZId94UE1i1VS3gLZJgnMDEW4L-AlKHq0Ok"
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	router.ServeHTTP(w, req)

	// 5. Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.Users
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response)
	assert.Equal(t, "1", response[0].Id)
}

// ====TEST CASE 4 - GET ALL USERS=====
func TestGetAllusers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	authGuard := router.Group("/api")
	// Use your actual middleware here
	authGuard.Use(authmiddleware.AuthMiddleware())
	{
		authGuard.GET("/getallusers", users.GetAllusers)
	}

	// 2. Create the Request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/getallusers", nil)

	bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJleUB5YWhvby5jb20iLCJpc3MiOiJCQVJDTEFZUyBCQU5LIiwic3ViIjoidXNlcl9hdXRoZW50aWNhdGlvbiIsImV4cCI6MTc3Njg2NTczOSwibmJmIjoxNzc2ODM2OTM5LCJpYXQiOjE3NzY4MzY5MzksImp0aSI6IjEifQ.6YLLoQYHGZId94UE1i1VS3gLZJgnMDEW4L-AlKHq0Ok"
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	router.ServeHTTP(w, req)

	// 5. Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.Users
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response)
	assert.Equal(t, "1", response[0].Id)
}

// ====TEST CASE 5 - UPDATE USER PROFILE============
func TestUpdateProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 1. Setup Router (Matches your actual routing logic)
	router := gin.New()
	authGuard := router.Group("/api")
	// Note: You might need to mock AuthMiddleware or use a real token
	authGuard.Use(authmiddleware.AuthMiddleware())
	{
		authGuard.PATCH("/updateprofile/:id", users.UpdateProfile)
	}

	profileReq := dto.ProfileData{
		Firstname: "Pator Joey",
		Lastname:  "Galam",
		Mobile:    "324234243",
	}
	body, _ := json.Marshal(profileReq)

	// 2. Create Request (Include the body here!)
	req, _ := http.NewRequest(http.MethodPatch, "/api/updateprofile/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJleUB5YWhvby5jb20iLCJpc3MiOiJCQVJDTEFZUyBCQU5LIiwic3ViIjoidXNlcl9hdXRoZW50aWNhdGlvbiIsImV4cCI6MTc3Njg2NTczOSwibmJmIjoxNzc2ODM2OTM5LCJpYXQiOjE3NzY4MzY5MzksImp0aSI6IjEifQ.6YLLoQYHGZId94UE1i1VS3gLZJgnMDEW4L-AlKHq0Ok"
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	// 3. Use the Router to serve the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req) // This triggers the full lifecycle (Middleware -> Handler)

	// 4. Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Your Profile has been successfully changed.", response["message"])
}

// ===TEST CASE 6 - CHANGE USER PASSWORD=====
func TestChangePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 1. Setup Router (Matches your actual routing logic)
	router := gin.New()
	authGuard := router.Group("/api")
	// Note: You might need to mock AuthMiddleware or use a real token
	authGuard.Use(authmiddleware.AuthMiddleware())
	{
		authGuard.PATCH("/changepassword/:id", users.ChangePassword)
	}

	passwordReq := dto.ChangePassword{
		Password: "nald",
	}
	body, _ := json.Marshal(passwordReq)

	// 2. Create Request (Include the body here!)
	req, _ := http.NewRequest(http.MethodPatch, "/api/changepassword/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJleUB5YWhvby5jb20iLCJpc3MiOiJCQVJDTEFZUyBCQU5LIiwic3ViIjoidXNlcl9hdXRoZW50aWNhdGlvbiIsImV4cCI6MTc3Njg2NTczOSwibmJmIjoxNzc2ODM2OTM5LCJpYXQiOjE3NzY4MzY5MzksImp0aSI6IjEifQ.6YLLoQYHGZId94UE1i1VS3gLZJgnMDEW4L-AlKHq0Ok"
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	// 3. Use the Router to serve the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req) // This triggers the full lifecycle (Middleware -> Handler)

	// 4. Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "You have changed your password successfully.", response["message"])
}

// ===TEST CASE 7 - ACTIVATE / DEACTIVATE MULTI-FACTO AUTHENTICATOR================
func TestMfaActivate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 1. Setup Router (Matches your actual routing logic)
	router := gin.New()
	authGuard := router.Group("/api")
	// Note: You might need to mock AuthMiddleware or use a real token
	authGuard.Use(authmiddleware.AuthMiddleware())
	{
		authGuard.PATCH("/mfa/activate/:id", auth.MfaActivate)
	}

	mfaReq := dto.MfaActivation{
		TwoFactorEnabled: true,
	}
	body, _ := json.Marshal(mfaReq)

	// 2. Create Request (Include the body here!)
	req, _ := http.NewRequest(http.MethodPatch, "/api/mfa/activate/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJleUB5YWhvby5jb20iLCJpc3MiOiJCQVJDTEFZUyBCQU5LIiwic3ViIjoidXNlcl9hdXRoZW50aWNhdGlvbiIsImV4cCI6MTc3Njg2NTczOSwibmJmIjoxNzc2ODM2OTM5LCJpYXQiOjE3NzY4MzY5MzksImp0aSI6IjEifQ.6YLLoQYHGZId94UE1i1VS3gLZJgnMDEW4L-AlKHq0Ok"
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	// 3. Use the Router to serve the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req) // This triggers the full lifecycle (Middleware -> Handler)

	// 4. Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Multi-Factor Authenticator has been enabled.", response["message"])
}

// ====TEST CASE 8 - VERIFY OTP Code=============
func TestMfaVerifyotp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 1. Setup Router (Matches your actual routing logic)
	router := gin.New()
	authGuard := router.Group("/api")
	// Note: You might need to mock AuthMiddleware or use a real token
	authGuard.Use(authmiddleware.AuthMiddleware())
	{
		authGuard.PATCH("/mfa/verifytotp/:id", auth.MfaVerifyotp)
	}

	secret := "E2G257XPK6YITCYIPII44KD3NWEYWSYD"

	validOtp, _ := totp.GenerateCode(secret, time.Now())

	mfaverifyReq := dto.MfaKeys{
		Otp: validOtp, // Use the dynamic code here
	}

	mfaverifyReq = dto.MfaKeys{
		Otp: validOtp,
	}

	body, _ := json.Marshal(mfaverifyReq)

	req, _ := http.NewRequest(http.MethodPatch, "/api/mfa/verifytotp/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJleUB5YWhvby5jb20iLCJpc3MiOiJCQVJDTEFZUyBCQU5LIiwic3ViIjoidXNlcl9hdXRoZW50aWNhdGlvbiIsImV4cCI6MTc3Njg3MzQ3MiwibmJmIjoxNzc2ODQ0NjcyLCJpYXQiOjE3NzY4NDQ2NzIsImp0aSI6IjEifQ.gIMpGEBLAW-GkjqWQzfO5GiGB6tS_n6q4rOnYM8tMeQ"
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 4. Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "OTP code is successfully validated.", response["message"])
}

// ====TEST CASE 9 - PRODUCT LIST PAGINATION======
func ProductList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/products/list/:page", middleware.ProductList)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/products/list/1", nil)

	router.ServeHTTP(w, req)

	// 5. Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.Products
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response)
	assert.Equal(t, "1", response[0].Id)
}

// === TEST CASE 10 - SEARCH PRODUCT AND PAGINATE======
func TestProductSearch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/products/search/:page/:key", middleware.ProductList)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/products/search/1/cineo", nil)

	router.ServeHTTP(w, req)

	// 5. Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.Products
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response)
	assert.Equal(t, "1", response[0].Id)
}

// ===TEST CASE 11 - GET ALL PRODUCT TO GENERATE PDF REPORT ====
func TestGetProductdata(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/productreport", middleware.GetProductdata)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/productreport", nil)

	router.ServeHTTP(w, req)

	// 5. Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.Products
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response)
	assert.Equal(t, "1", response[0].Id)
}

// ===TEST CASE 12 - GET ALL PRODUCTS FOR MASTER AND DETAILS FOR PDF REPORT====
func TestGetAllMasterDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/categories", middleware.GetProductdata)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/categories", nil)

	router.ServeHTTP(w, req)

	// 5. Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.Products
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response)
	assert.Equal(t, "1", response[0].Id)
}
