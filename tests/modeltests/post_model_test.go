package modeltests

import (
	"log"
	"testing"

	"github.com/abdulhamidnugroho/go-full/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllPosts(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("error refreshing user and post table %v\n", err)
	}

	_, _, err = seedUsersAndPosts()
	if err != nil {
		log.Fatalf("error seeding user and post table %v\n", err)
	}

	posts, err := postInstance.FindAllPosts(server.DB)
	if err != nil {
		t.Errorf("error getting the posts: %v\n", err)
		return
	}

	assert.Equal(t, len(*posts), 2)
}

func TestSavePost(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("error refreshing user and post table %v\n", err)
	}

	user, err := seedOneUSer()
	if err != nil {
		log.Fatalf("cannot seed user %v\n", err)
	}

	newPost := models.Post{
		ID:       1,
		Title:    "This Title",
		Content:  "This Content",
		AuthorID: user.ID,
	}

	savedPost, err := newPost.SavePost(server.DB)
	if err != nil {
		t.Errorf("error getting the post: %v\n", err)
		return
	}

	assert.Equal(t, newPost.ID, savedPost.ID)
	assert.Equal(t, newPost.Title, savedPost.Title)
	assert.Equal(t, newPost.Content, savedPost.Content)
	assert.Equal(t, newPost.AuthorID, savedPost.AuthorID)
}

func TestGetPostByID(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("error refreshing user and post table: %v\n", err)
	}

	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("errong seeding table")
	}

	foundPost, err := postInstance.FindPostByID(server.DB, post.ID)
	if err != nil {
		t.Errorf("error getting one user: %v\n", err)
		return
	}

	assert.Equal(t, foundPost.ID, post.ID)
	assert.Equal(t, foundPost.Title, post.Title)
	assert.Equal(t, foundPost.Content, post.Content)
}

func TestUpdateAPost(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("error refreshing user and post table: %v\n", err)
	}

	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("errong seeding table")
	}

	postUpdate := models.Post{
		ID:       1,
		Title:    "mUpdate",
		Content:  "mupdate@gmail.com",
		AuthorID: post.AuthorID,
	}

	updatedPost, err := postUpdate.UpdateAPost(server.DB)
	if err != nil {
		t.Errorf("error updating the post: %v\n", err)
		return
	}

	assert.Equal(t, updatedPost.ID, postUpdate.ID)
	assert.Equal(t, updatedPost.Title, postUpdate.Title)
	assert.Equal(t, updatedPost.Content, postUpdate.Content)
	assert.Equal(t, updatedPost.AuthorID, postUpdate.AuthorID)
}

func TestDeleteAPost(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("error refreshing user and post table: %v\n", err)
	}

	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("errong seeding table")
	}

	isDeleted, err := postInstance.DeleteAPost(server.DB, post.ID, post.AuthorID)
	if err != nil {
		t.Errorf("error deleting the post: %v\n", err)
	}

	// assert.Equal(t, int(isDeleted), 1)
	assert.Equal(t, isDeleted, int64(1))
}
