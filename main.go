package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GoLang_Dungeon/db"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

type User struct {
	Fname      string `json:"fname"`
	Lname      string `json:"lname"`
	Location   string `json:"location"`
	Bio        string `json:"bio"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	ProfilePic string `json:"profilePic"`
}

type BookResponse struct {
	BookID string `json:"bookID"`
	Title  string `json:"title"`
}

type PostId struct {
	PostId string `json:"postId"`
}

type Post struct {
	Picture                    string `json:"picture"`
	Description                string `json:"description"`
	Duration                   string `json:"duration"`
	RequirementsAndRestriction string `json:"requirementsAndRestriction"`
	Location                   string `json:"location"`
	Cost                       string `json:"cost"`
	Username                   string `json:"username"`
}

type Dude struct {
	Username string `json:"username"`
}

type City struct {
	Location string `json:"location"`
}

type Book struct {
	Username string `json:"username"`
	PostID   string `json:"postID"`
	Title    string `json:"title"`
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	createdUser, err := client.User.CreateOne(
		db.User.LastName.Set(user.Lname),
		db.User.FirstName.Set(user.Fname),
		db.User.ProfilePicture.Set(user.ProfilePic),
		db.User.Location.Set(user.Location),
		db.User.Bio.Set(user.Bio),
		db.User.Email.Set(user.Email),
		db.User.Username.Set(user.Username),
	).Exec(ctx)

	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	result, _ := json.MarshalIndent(createdUser, "", "  ")
	fmt.Printf("Created user: %s\n", result)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	posts, err := client.Post.FindMany().Exec(ctx)

	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func GetPostByCity(w http.ResponseWriter, r *http.Request) {
	var city City
	err := json.NewDecoder(r.Body).Decode(&city)

	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	posts, err := client.Post.FindMany(
		db.Post.Location.Equals(city.Location),
	).Exec(ctx)

	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func GetPostByUsername(w http.ResponseWriter, r *http.Request) {
	var peep Dude
	err := json.NewDecoder(r.Body).Decode(&peep)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()
	foundUser, err := client.User.FindUnique(
		db.User.Username.Equals(peep.Username),
	).Exec(ctx)
	if err != nil {
		http.Error(w, "Failed to find user", http.StatusInternalServerError)
		return
	}
	if foundUser == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	posts, err := client.Post.FindMany(
		db.Post.UserID.Equals(foundUser.UserID),
	).Exec(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	foundUser, err := client.User.FindUnique(
		db.User.Username.Equals(post.Username),
	).Exec(ctx)
	if err != nil {
		http.Error(w, "Failed to find user", http.StatusInternalServerError)
		return
	}

	if foundUser == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	createdPost, err := client.Post.CreateOne(
		db.Post.Picture.Set(post.Picture),
		db.Post.Description.Set(post.Description),
		db.Post.Duration.Set(post.Duration),
		db.Post.RequirementsAndRestriction.Set(post.RequirementsAndRestriction),
		db.Post.Location.Set(post.Location),
		db.Post.Cost.Set(post.Cost),
		db.Post.User.Link(db.User.UserID.Equals(foundUser.UserID)),
	).Exec(ctx)

	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	result, _ := json.MarshalIndent(createdPost, "", "  ")
	fmt.Printf("Created post: %s\n", result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Post created successfully"))
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()
	foundUser, err := client.User.FindUnique(
		db.User.Username.Equals(book.Username),
	).Exec(ctx)
	if err != nil {
		http.Error(w, "Failed to find user", http.StatusInternalServerError)
		return
	}
	if foundUser == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	createdBook, err := client.Book.CreateOne(
		db.Book.Title.Set(book.Title),
		db.Book.Post.Link(db.Post.PostID.Equals(book.PostID)),
		db.Book.User.Link(db.User.UserID.Equals(foundUser.UserID)),
	).Exec(ctx)
	if err != nil {
		http.Error(w, "Failed to create book", http.StatusInternalServerError)
		return
	}
	result, _ := json.MarshalIndent(createdBook, "", "  ")
	fmt.Printf("Created book: %s\n", result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Book created successfully"))
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	var peep Dude
	err := json.NewDecoder(r.Body).Decode(&peep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	foundUser, err := client.User.FindUnique(
		db.User.Username.Equals(peep.Username),
	).With(
		db.User.Posts.Fetch(),
		db.User.Comments.Fetch(),
		db.User.Books.Fetch(),
	).Exec(ctx)
	if err != nil {
		http.Error(w, "Failed to find user", http.StatusInternalServerError)
		return
	}

	if foundUser == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	userJSON, err := json.Marshal(foundUser)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJSON)
}

func GetBookingUsername(w http.ResponseWriter, r *http.Request) {
	var peep Dude
	err := json.NewDecoder(r.Body).Decode(&peep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			log.Println("Error disconnecting from database:", err)
		}
	}()
	ctx := context.Background()
	foundUser, err := client.User.FindUnique(
		db.User.Username.Equals(peep.Username),
	).Exec(ctx)
	if err != nil {
		http.Error(w, "Failed to find user", http.StatusInternalServerError)
		return
	}
	if foundUser == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	books, err := client.Book.FindMany(
		db.Book.UserID.Equals(foundUser.UserID),
	).Exec(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
		return
	}
	var booksResponse []BookResponse
	for _, book := range books {
		booksResponse = append(booksResponse, BookResponse{
			BookID: book.BookID,
			Title:  book.Title,
		})
	}

	booksJSON, err := json.Marshal(booksResponse)
	if err != nil {
		log.Println("Failed to marshal books JSON:", err)
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(booksJSON)
}

func GetBookingPost(w http.ResponseWriter, r *http.Request) {
	var peep PostId
	err := json.NewDecoder(r.Body).Decode(&peep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()
	ctx := context.Background()
	foundPost, err := client.Post.FindUnique(
		db.Post.PostID.Equals(peep.PostId),
	).Exec(ctx)
	if err != nil {
		http.Error(w, "Failed to find post", http.StatusInternalServerError)
		return
	}

	if foundPost == nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	books, err := client.Book.FindMany(
		db.Book.PostID.Equals(foundPost.PostID),
	).Exec(ctx)
	if err != nil {
		http.Error(w, "Failed to retrieve books", http.StatusInternalServerError)
		return
	}
	booksJSON, err := json.Marshal(books)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(booksJSON)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/postuser", PostUser).Methods("POST")
	r.HandleFunc("/getPosts", GetAllPosts).Methods("GET")
	r.HandleFunc("/getPostsByCity", GetPostByCity).Methods("GET")
	r.HandleFunc("/getPostsByUsername", GetPostByUsername).Methods("GET")
	r.HandleFunc("/createPost", CreatePost).Methods("POST")
	r.HandleFunc("/createBooking", CreateBook).Methods("POST")
	r.HandleFunc("/getUserInfo", GetUserInfo).Methods("GET")
	r.HandleFunc("/getBookingUsername", GetBookingUsername).Methods("GET")
	r.HandleFunc("/getBookingByPost", GetBookingPost).Methods("GET")

	corsHandler := cors.Default().Handler(r)

	http.ListenAndServe(":4000", corsHandler)
}
