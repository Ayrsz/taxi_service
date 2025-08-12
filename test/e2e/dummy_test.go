package e2e

import (
	"strconv"
	"testing"

	"github.com/gjcms/taxi_service/models"
	"github.com/gjcms/taxi_service/test"

	"github.com/stretchr/testify/assert"
)

func TestDummyCRUD(t *testing.T) {
	app := test.SetupTestApp(t)
	defer test.CleanupTestApp(t)

	// Test Create
	createPayload := models.DummyUser{
		Name:  "Test User",
		Email: "test@test.com",
	}

	resp := test.MakeRequest(t, app, "POST", "/dummy", createPayload)
	assert.Equal(t, 200, resp.StatusCode)

	var createdUser models.DummyUser
	test.ParseResponseBody(t, resp, &createdUser)
	assert.Equal(t, createPayload.Name, createdUser.Name)
	assert.Equal(t, createPayload.Email, createdUser.Email)

	// Test Get
	resp = test.MakeRequest(t, app, "GET", "/dummy/"+strconv.FormatUint(uint64(createdUser.ID), 10), nil)
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

	resp = test.MakeRequest(t, app, "PUT", "/dummy/"+strconv.FormatUint(uint64(createdUser.ID), 10), updatePayload)
	assert.Equal(t, 200, resp.StatusCode)

	var updatedUser models.DummyUser
	test.ParseResponseBody(t, resp, &updatedUser)
	assert.Equal(t, updatePayload.Name, updatedUser.Name)
	assert.Equal(t, updatePayload.Email, updatedUser.Email)

	// Test List
	resp = test.MakeRequest(t, app, "GET", "/dummy", nil)
	assert.Equal(t, 200, resp.StatusCode)

	var users []models.DummyUser
	test.ParseResponseBody(t, resp, &users)
	assert.GreaterOrEqual(t, len(users), 1)

	// Test Delete
	resp = test.MakeRequest(t, app, "DELETE", "/dummy/"+strconv.FormatUint(uint64(createdUser.ID), 10), nil)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify deletion
	resp = test.MakeRequest(t, app, "GET", "/dummy/"+strconv.FormatUint(uint64(createdUser.ID), 10), nil)
	assert.Equal(t, 500, resp.StatusCode) // The service returns 500 when record not found
}
