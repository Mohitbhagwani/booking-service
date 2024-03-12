package repository

import (
    "database/sql"
	"github.com/google/uuid"
    "booking-service/models"
	"log"
	"time"
	"strings" // Adjust the import path based on your project structure
)

type UserRepository struct {
    db *sql.DB // or *sql.Tx if you want to support transactions
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (ur *UserRepository) InsertUser(user models.User) (models.User, error) {
    // Generate a new UUID for the user
    userID := uuid.New()

    // Define the SQL query for inserting a user with a manually generated UUID
    query := `
        INSERT INTO "user" (id, first_name, last_name, password, role, username, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
    `
    // Execute the SQL query within the repository's database connection
    _, err := ur.db.Exec(query, userID, user.FirstName, user.LastName, user.Password, user.Role, user.Username)
    if err != nil {
        return models.User{}, err
    }

	log.Printf("Inserted user with ID: %s", userID)

    // Set the generated user ID to the user struct
    user.ID = userID

    return user, nil
}

func (ur *UserRepository) UpdateUser(user models.User, userID uuid.UUID) (models.User, error) {
    // Generate a new UUID for the user
	currentTime := time.Now()
    formattedTime := currentTime.Format("2006-01-02T15:04:05.999999Z")
    // Define the SQL query for inserting a user with a manually generated UUID
    query := `
	UPDATE public."user" SET first_name = $1, last_name = $2, role = $3, username =$4, updated_at=$6  WHERE id = $5
    `
    _, err := ur.db.Exec(query, user.FirstName, user.LastName, user.Role, strings.ToLower(user.Username), userID, formattedTime)
    if err != nil {
        return models.User{}, err
    }

	log.Printf("updated user by ID: %s", userID)

    // Set the generated user ID to the user struct
    user.ID = userID
	user.Username = strings.ToLower(user.Username)
	user.UpdatedAt, err = time.Parse("2006-01-02T15:04:05.999999Z", formattedTime)
	log.Printf("time ->> %s",formattedTime)
    return user, nil
}

func (ur *UserRepository) GetUserByID(userID uuid.UUID) (models.User, error) {
    // Define the SQL query for retrieving a user by ID
    query := `
        SELECT id, first_name, last_name, role, lower(username), created_at, updated_at, deleted_at
        FROM public."user"
        WHERE id = $1 and deleted_at is null
    `

    var user models.User
    err := ur.db.QueryRow(query, userID).Scan(
        &user.ID,
        &user.FirstName,
        &user.LastName,
        &user.Role,
        &user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
    )
    if err != nil {
        return models.User{}, err
    }

    return user, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (models.User, error) {
    // Define the SQL query for retrieving a user by ID
    query := `
        SELECT id, first_name, last_name, role, lower(username), password, created_at, updated_at, deleted_at
        FROM public."user"
        WHERE lower(username) = $1 and deleted_at is null
    `

    var user models.User
    err := ur.db.QueryRow(query, strings.ToLower(email)).Scan(
        &user.ID,
        &user.FirstName,
        &user.LastName,
        &user.Role,
        &user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
    )
    if err != nil {
        return models.User{}, err
    }

    return user, nil
}

// func (ur *UserRepository) HardDeleteUserById(id uuid.UUID) error {
//     // Define the SQL query for retrieving a user by ID
//     query := `
//         Delete
//         FROM public."user"
//         WHERE id = $1
//     `

// 	_, err := ur.db.Exec(query, id)
//     if err != nil {
//         return err
//     }
// 	return nil
// }

func (ur *UserRepository) SoftDeleteUserById(id uuid.UUID) error {
    // Define the SQL query for retrieving a user by ID
    query := `
	UPDATE public."user" SET deleted_at= Now() WHERE id = $1
    `

	_, err := ur.db.Exec(query, id)
    if err != nil {
        return err
    }
	return nil
}

func (ur *UserRepository) GetAllUsers() ([]models.User, error) {
    // Define the SQL query for retrieving all users
    query := `
        SELECT id, first_name, last_name, role, lower(username), created_at, updated_at, deleted_at
        FROM public."user"
    `

    // Execute the SQL query within the repository's database connection
    rows, err := ur.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Initialize a slice to store the retrieved users
    var users []models.User

    // Iterate through the query results and append users to the slice
    for rows.Next() {
        var user models.User
        err := rows.Scan(
            &user.ID,
            &user.FirstName,
            &user.LastName,
            &user.Role,
            &user.Username,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
        )
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    return users, nil
}