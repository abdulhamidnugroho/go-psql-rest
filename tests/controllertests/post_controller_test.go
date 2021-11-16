package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/abdulhamidnugroho/go-full/api/models"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreatePost(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("cannot seed user %v\n", err)
	}

	token, err := server.SignIn(user.Email, "password")
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}

	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON    string
		statusCode   int
		title        string
		content      string
		author_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"title":"The title", "content": "the content", "author_id": 1}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			title:        "The title",
			content:      "the content",
			author_id:    user.ID,
			errorMessage: "",
		},
		{
			inputJSON:    `{"title":"The title", "content": "the content", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "title already taken",
		},
		{
			// When no token is passed
			inputJSON:    `{"title":"When no token is passed", "content": "the content", "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"title":"When incorrect token is passed", "content": "the content", "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "unauthorized",
		},
		{
			inputJSON:    `{"title": "", "content": "The content", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "required title",
		},
		{
			inputJSON:    `{"title": "This is a title", "content": "", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "required content",
		},
		{
			inputJSON:    `{"title": "This is an awesome title", "content": "the content"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "required author",
		},
		{
			// When user 2 uses user 1 token
			inputJSON:    `{"title": "This is an awesome title", "content": "the content", "author_id": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("POST", "/posts", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("the error: %v\n", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreatePost)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["content"], v.content)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id)) //just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetPosts(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	_, _, err = seedUsersAndPost()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/posts", nil)
	if err != nil {
		t.Errorf("the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetPosts)
	handler.ServeHTTP(rr, req)

	var posts []models.Post
	err = json.Unmarshal([]byte(rr.Body.String()), &posts)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(posts), 2)
}

func TestGetPostByID(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatal(err)
	}

	postSample := []struct {
		id           string
		statusCode   int
		title        string
		content      string
		author_id    uint32
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(post.ID)),
			statusCode: 200,
			title:      post.Title,
			content:    post.Content,
			author_id:  post.AuthorID,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}

	for _, v := range postSample {
		req, err := http.NewRequest("GET", "/posts", nil)
		if err != nil {
			t.Errorf("the error: %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetPost)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, post.Title, responseMap["title"])
			assert.Equal(t, post.Content, responseMap["content"])
			assert.Equal(t, float64(post.AuthorID), responseMap["author_id"])
		}
	}
}

func TestUpdatePost(t *testing.T) {
	var PostUserEmail, PostUserPassword string
	var AuthPostAuthorID uint32
	var AuthPostID uint64

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	users, posts, err := seedUsersAndPost()
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		PostUserEmail = user.Email
		PostUserPassword = "password"
	}

	token, err := server.SignIn(PostUserEmail, PostUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}

	tokenString := fmt.Sprintf("Bearer %v", token)

	for _, post := range posts {
		if post.ID == 2 {
			continue
		}
		AuthPostID = post.ID
		AuthPostAuthorID = post.AuthorID
	}

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		title        string
		content      string
		author_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"title":"The updated post", "content": "This is the updated content", "author_id": 1}`,
			statusCode:   200,
			title:        "The updated post",
			content:      "This is the updated content",
			author_id:    AuthPostAuthorID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// When no token is provided
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"title":"This is still another title", "content": "This is the updated content", "author_id": 1}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "unauthorized",
		},
		{
			// When incorrect token is provided
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"title":"This is still another title", "content": "This is the updated content", "author_id": 1}`,
			tokenGiven:   "this is an incorrect token",
			statusCode:   401,
			errorMessage: "unauthorized",
		},
		{
			//Note: "Title 2" belongs to post 2, and title must be unique
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"title":"Title 2", "content": "This is the updated content", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "title already taken",
		},
		{
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"title":"", "content": "This is the updated content", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "required title",
		},
		{
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"title":"Awesome title", "content": "", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "required content",
		},
		{
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"title":"This is another title", "content": "This is the updated content"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "unauthorized",
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"title":"This is still another title x", "content": "This is the updated content", "author_id": 2}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "unauthorized",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("POST", "/posts", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("the error: %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdatePost)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["content"], v.content)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id)) //just to match the type of the json we receive thats why we used float64
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeletePost(t *testing.T) {
	var PostUserEmail, PostUserPassword string
	var PostUserID uint32
	var AuthPostID uint64

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	users, posts, err := seedUsersAndPost()
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		if user.ID == 1 {
			continue
		}
		PostUserEmail = user.Email
		PostUserPassword = "password"
	}

	token, err := server.SignIn(PostUserEmail, PostUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the second post
	for _, post := range posts {
		if post.ID == 1 {
			continue
		}
		AuthPostID = post.ID
		PostUserID = post.AuthorID
	}

	postSample := []struct {
		id           string
		author_id    uint32
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthPostID)),
			author_id:    PostUserID,
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When empty token is passed
			id:           strconv.Itoa(int(AuthPostID)),
			author_id:    PostUserID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "unauthorized",
		},
		{
			// When incorrect token is passed
			id:           strconv.Itoa(int(AuthPostID)),
			author_id:    PostUserID,
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "unauthorized",
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(1)),
			author_id:    1,
			statusCode:   401,
			errorMessage: "unauthorized",
		},
	}

	for _, v := range postSample {
		req, _ := http.NewRequest("GET", "/posts", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeletePost)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {

			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
