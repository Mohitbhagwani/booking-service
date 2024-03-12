package models

import (
    "github.com/google/uuid"
    "time"
)

type User struct {
    ID        uuid.UUID    `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Password  string `json:"password"`
    Role      string `json:"role"`
    Username  string `json:"username"`
    CreatedAt time.Time `json:"created_at`
    UpdatedAt time.Time `json:"updated_at`
    DeletedAt *time.Time `json:"deleted_at`
}
