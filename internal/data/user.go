package data

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Age      int
	School   string
	Role     string // "teacher" or "coach"
}
