package modeltests

import (
	"log"
	"testing"

	"github.com/abdulhamidnugroho/go-full/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllUsers(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}

	users, err := userInstance.FindAllUser(server.DB)
	if err != nil {
		t.Errorf("error getting the users: %v\n", err)
		return
	}

	assert.Equal(t, len(*users), 2)
}

func TestSaveUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	newUser := models.User{
		ID:       1,
		Email:    "test@gmail.com",
		Nickname: "test",
		Password: "password",
	}

	savedUser, err := newUser.SaveUser(server.DB)
	if err != nil {
		t.Errorf("error getting the users: %v\n", err)
		return
	}

	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Email, savedUser.Email)
	assert.Equal(t, newUser.Nickname, savedUser.Nickname)
}

func TestGetUserByID(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUSer()
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}

	foundUser, err := userInstance.FindUserByID(server.DB, user.ID)
	if err != nil {
		t.Errorf("error getting the users: %v\n", err)
		return
	}

	assert.Equal(t, foundUser.ID, user.ID)
	assert.Equal(t, foundUser.Email, user.Email)
	assert.Equal(t, foundUser.Nickname, user.Nickname)
}

func TestUpdateAUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUSer()
	if err != nil {
		log.Fatalf("cannot seed users: %v", err)
	}

	userUpdate := models.User{
		ID:       1,
		Nickname: "mUpdate",
		Email:    "mupdate@gmail.com",
		Password: "password",
	}

	updatedUser, err := userUpdate.UpdateAUser(server.DB, user.ID)
	if err != nil {
		t.Errorf("error updating the users: %v\n", err)
		return
	}

	assert.Equal(t, updatedUser.ID, userUpdate.ID)
	assert.Equal(t, updatedUser.Email, userUpdate.Email)
	assert.Equal(t, updatedUser.Nickname, userUpdate.Nickname)
}

func TestDeleteAUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUSer()
	if err != nil {
		log.Fatalf("cannot seed users: %v", err)
	}

	isDeleted, err := userInstance.DeleteAUser(server.DB, user.ID)
	if err != nil {
		t.Errorf("error deleting the users: %v\n", err)
		return
	}

	// assert.Equal(t, int(isDeleted), 1)
	assert.Equal(t, isDeleted, int64(1))
}
