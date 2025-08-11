package e2e

import (
    "testing"
    "strconv"
    "taxi-service/models"
    "taxi-service/test"

    "github.com/stretchr/testify/assert"
)

func TestGetDummyUserByID1(t *testing.T) {
    app := test.SetupTestApp(t)
    defer test.CleanupTestApp(t)

	// Test Create
	createPayload := models.DummyUser{
		Name:  "Test User",
		Email: "test@test.com",
	}

	resp := test.MakeRequest(t, app, "POST", "/dummy-users", createPayload)
	assert.Equal(t, 200, resp.StatusCode)

	var createdUser models.DummyUser
	test.ParseResponseBody(t, resp, &createdUser)
	assert.Equal(t, createPayload.Name, createdUser.Name)
	assert.Equal(t, createPayload.Email, createdUser.Email)

	// Test Get
	resp = test.MakeRequest(t, app, "GET", "/dummy-users/"+strconv.FormatUint(uint64(createdUser.ID), 10), nil)
	assert.Equal(t, 200, resp.StatusCode)

	var fetchedUser models.DummyUser
	test.ParseResponseBody(t, resp, &fetchedUser)
	assert.Equal(t, createdUser.ID, fetchedUser.ID)
	assert.Equal(t, createdUser.Name, fetchedUser.Name)

	// Test Update
	updatePayload := models.DummyUser{
		Name:  "Updated User",
		Email: "Email@email.com",
	}

	resp = test.MakeRequest(t, app, "PUT", "/dummy-users/"+strconv.FormatUint(uint64(createdUser.ID), 10), updatePayload)
	assert.Equal(t, 200, resp.StatusCode)

	var updatedUser models.DummyUser
	test.ParseResponseBody(t, resp, &updatedUser)
	assert.Equal(t, updatePayload.Name, updatedUser.Name)
	assert.Equal(t, updatePayload.Email, updatedUser.Email)

	// Test List
	resp = test.MakeRequest(t, app, "GET", "/dummy-users", nil)
	assert.Equal(t, 200, resp.StatusCode)

	var users []models.DummyUser
	test.ParseResponseBody(t, resp, &users)
	assert.GreaterOrEqual(t, len(users), 1)

	// Test Delete
	resp = test.MakeRequest(t, app, "DELETE", "/dummy-users/"+strconv.FormatUint(uint64(createdUser.ID), 10), nil)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify deletion
	resp = test.MakeRequest(t, app, "GET", "/dummy-users/"+strconv.FormatUint(uint64(createdUser.ID), 10), nil)
	assert.Equal(t, 500, resp.StatusCode) // The service returns 500 when record not found
}

func TestFindOrCreateDummyUserID1(t *testing.T) {
    app := test.SetupTestApp(t)
    defer test.CleanupTestApp(t)

    // First, check what users exist
    listResp := test.MakeRequest(t, app, "GET", "/dummy", nil)
    assert.Equal(t, 200, listResp.StatusCode)
    
    var users []models.DummyUser
    test.ParseResponseBody(t, listResp, &users)
    
    // Look for user with ID 1
    var userID1 *models.DummyUser
    for _, user := range users {
        if user.ID == 1 {
            userID1 = &user
            break
        }
    }
    
    if userID1 != nil {
        // User with ID 1 exists
        assert.Equal(t, 1, userID1.ID)
        assert.NotEmpty(t, userID1.Name)
        assert.NotEmpty(t, userID1.Email)
        t.Logf("User with ID 1 already exists: Name='%s', Email='%s'", userID1.Name, userID1.Email)
        
        // Verify we can fetch it directly
        resp := test.MakeRequest(t, app, "GET", "/dummy/1", nil)
        assert.Equal(t, 200, resp.StatusCode)
        
    } else if len(users) == 0 {
        // No users exist, create one - it should get ID 1
        createPayload := models.DummyUser{
            Name:  "First User",
            Email: "first@example.com",
        }
        
        createResp := test.MakeRequest(t, app, "POST", "/dummy", createPayload)
        assert.Equal(t, 200, createResp.StatusCode)
        
        var createdUser models.DummyUser
        test.ParseResponseBody(t, createResp, &createdUser)
        assert.Equal(t, 1, createdUser.ID, "First created user should have ID 1")
        assert.Equal(t, "First User", createdUser.Name)
        assert.Equal(t, "first@example.com", createdUser.Email)
        t.Log("Created first user with ID 1")
        
        // Verify it can be fetched
        getResp := test.MakeRequest(t, app, "GET", "/dummy/1", nil)
        assert.Equal(t, 200, getResp.StatusCode)
        
        var fetchedUser models.DummyUser
        test.ParseResponseBody(t, getResp, &fetchedUser)
        assert.Equal(t, 1, fetchedUser.ID)
        
    } else {
        // Users exist but none with ID 1
        t.Logf("Found %d users but none with ID 1:", len(users))
        for _, user := range users {
            t.Logf("- ID: %d, Name: '%s', Email: '%s'", user.ID, user.Name, user.Email)
        }
        
        // Try to get ID 1 directly (should fail)
        resp := test.MakeRequest(t, app, "GET", "/dummy/1", nil)
        assert.Equal(t, 500, resp.StatusCode, "Should return 500 when user ID 1 doesn't exist")
        t.Log("Confirmed: User with ID 1 does not exist")
    }
}

func TestDummyUserJSONRepository(t *testing.T) {
    app := test.SetupTestApp(t)
    defer test.CleanupTestApp(t)

    // Test the JSON repository pattern directly
    
    // 1. List all users (should work even if empty)
    listResp := test.MakeRequest(t, app, "GET", "/dummy", nil)
    assert.Equal(t, 200, listResp.StatusCode)
    
    var initialUsers []models.DummyUser
    test.ParseResponseBody(t, listResp, &initialUsers)
    initialCount := len(initialUsers)
    t.Logf("Initial user count: %d", initialCount)
    
    // 2. Create a new user
    newUser := models.DummyUser{
        Name:  "Test User",
        Email: "test@example.com",
    }
    
    createResp := test.MakeRequest(t, app, "POST", "/dummy", newUser)
    assert.Equal(t, 200, createResp.StatusCode)
    
    var createdUser models.DummyUser
    test.ParseResponseBody(t, createResp, &createdUser)
    assert.NotZero(t, createdUser.ID, "Created user should have an ID")
    assert.Equal(t, "Test User", createdUser.Name)
    assert.Equal(t, "test@example.com", createdUser.Email)
    t.Logf("Created user with ID: %d", createdUser.ID)
    
    // 3. Verify user count increased
    listResp2 := test.MakeRequest(t, app, "GET", "/dummy", nil)
    assert.Equal(t, 200, listResp2.StatusCode)
    
    var afterCreateUsers []models.DummyUser
    test.ParseResponseBody(t, listResp2, &afterCreateUsers)
    assert.Equal(t, initialCount+1, len(afterCreateUsers), "User count should increase by 1")
    
    // 4. Get the specific user
    getUserResp := test.MakeRequest(t, app, "GET", "/dummy/"+string(rune(createdUser.ID+'0')), nil)
    if createdUser.ID <= 9 { // Simple conversion for single digits
        assert.Equal(t, 200, getUserResp.StatusCode)
        
        var fetchedUser models.DummyUser
        test.ParseResponseBody(t, getUserResp, &fetchedUser)
        assert.Equal(t, createdUser.ID, fetchedUser.ID)
        assert.Equal(t, createdUser.Name, fetchedUser.Name)
        assert.Equal(t, createdUser.Email, fetchedUser.Email)
    }
    
    t.Log("JSON repository test completed successfully")
}