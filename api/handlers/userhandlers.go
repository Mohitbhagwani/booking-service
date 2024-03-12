package handlers

import (
	"booking-service/auth"
	"booking-service/db"
	"booking-service/models"
	"booking-service/repository"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"time"
)

type UserResponse struct {
	ID        uuid.UUID  `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Role      string     `json:"role"`
	Username  string     `json:"username"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to connect to the database")
		return
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	users, err := userRepo.GetAllUsers()
	if err != nil {
		log.Printf("Failed to fetch users: %s", err)
		errorMessage := "Something went wrong while fetching users"
		respondWithError(w, http.StatusInternalServerError, errorMessage)
		return
	}

	for _, existingUser := range users {
		if existingUser.Username == user.Username {
			errorMessage := "User already exists with username: " + existingUser.Username
			respondWithError(w, http.StatusConflict, errorMessage)
			return
		}
	}
	insertedUser, err := userRepo.InsertUser(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	userResponse := UserResponse{
		ID:        insertedUser.ID,
		FirstName: insertedUser.FirstName,
		LastName:  insertedUser.LastName,
		Role:      insertedUser.Role,
		Username:  insertedUser.Username,
	}

	log.Printf("%s  --> checking user created", userResponse)

	respondWithJSON(w, http.StatusCreated, userResponse)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idStr := vars["id"]

	userID, err1 := uuid.Parse(idStr)
	if err1 != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to connect to the database")
		return
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)

	existingUser, err := userRepo.GetUserByID(userID)
	if err != nil {
		log.Printf("checking this failed GetUserByID %s", err)
		errorMessage := "User not found with ID: " + userID.String()
		respondWithError(w, http.StatusConflict, errorMessage)
		return
	}

	insertedUser, err := userRepo.UpdateUser(user, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	log.Printf("%s  --> checking user updated", insertedUser)

	userResponse := UserResponse{
		ID:        insertedUser.ID,
		FirstName: insertedUser.FirstName,
		LastName:  insertedUser.LastName,
		Role:      insertedUser.Role,
		Username:  strings.ToLower(insertedUser.Username),
		CreatedAt: existingUser.CreatedAt,
		UpdatedAt: insertedUser.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
	respondWithJSON(w, http.StatusOK, userResponse)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	userID, err1 := uuid.Parse(idStr)
	if err1 != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to connect to the database")
		return
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)

	user1, err := userRepo.GetUserByID(userID)
	if err != nil {
		errorMessage := "User not exists with ID: " + userID.String()
		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}
	err = userRepo.SoftDeleteUserById(user1.ID)

	if err != nil {
		errorMessage := "Failed to delete user with ID: " + userID.String()
		respondWithError(w, http.StatusInternalServerError, errorMessage)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return

}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Get the "id" path variable from the request URL using Gorilla Mux
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Parse the "id" string as a UUID
	userID, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	// Create a database connection
	db, err := db.ConnectDB()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to connect to the database")
		return
	}
	defer db.Close()

	// Create a UserRepository instance
	userRepo := repository.NewUserRepository(db)

	// Retrieve the user by ID from the repository
	user, err := userRepo.GetUserByID(userID)
	log.Printf("%s  --> checking user ", user)
	if err != nil {
		errorMessage := "User not found with ID: " + userID.String()
		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	// Create a UserResponse object without the password
	userResponse := UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}

	// Respond with the retrieved user (excluding the password)
	respondWithJSON(w, http.StatusOK, userResponse)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Retrieve all users from the repository
	db, err := db.ConnectDB()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to connect to the database")
		return
	}
	defer db.Close()
	userRepo := repository.NewUserRepository(db)
	users, err := userRepo.GetAllUsers()
	if err != nil {

		respondWithError(w, http.StatusInternalServerError, "Failed to get users")
		return
	}

	if len(users) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Create a slice of UserResponse objects without the password for all users
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			DeletedAt: user.DeletedAt,
		}
	}

	// Respond with the list of users (excluding passwords) as JSON
	respondWithJSON(w, http.StatusOK, userResponses)
}

type User struct {
	ID        uuid.UUID
	Username  string
	Password  string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Implement a login handler
func Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to connect to the database")
		return
	}
	defer db.Close()
	userRepo := repository.NewUserRepository(db)
	// Replace this with your database query to fetch the user by username and password
	user, err := userRepo.GetUserByEmail(loginRequest.Username)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Username name not exists")
		return
	}
	
	log.Printf("checking condition %s and %s", user, loginRequest)
	if user.Password != loginRequest.Password {
		log.Printf("checking condition %s and %s", user.Password, loginRequest.Password)
		respondWithError(w, http.StatusUnauthorized, "Incorrect password")
		return
	}
	log.Printf("role type check->  %s", user.Role)
	if user.Role == "admin" {
		// Generate a JWT token with admin role
		tokenString, err := auth.GenerateJWT(user.ID, []string{"admin"}, 600) // 3600 seconds = 1 hour
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// Respond with the JWT token
		respondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
	} else {
		// Non-admin users are not allowed to log in
		http.Error(w, "Only admin users are allowed", http.StatusForbidden)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
