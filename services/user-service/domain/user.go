package domain

type User struct {
    UserID   string `bson:"user_id"`
    Username string `bson:"username"`
    Password string `bson:"password"`
}


type UserRepository interface {
    Create(user *User) (string, error)
    Authenticate(username, password string) (string, error)
    GetUserByID(userID string) (*User, error)
}
