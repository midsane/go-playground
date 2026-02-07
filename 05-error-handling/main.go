package main

import (
	"errors"
	"fmt"
	"log/slog"

	"go.uber.org/zap"
)

type User struct {
	name  string
	email string
}

var ErrUserNotFound = errors.New("user not found")

func createuser() (User, error) {
	//lets consider we got error while creating user
	return User{}, ErrUserNotFound
}

/*
custom error type -> errors.New instansiate a error type interface having
Error() string functin
*/
type permissionDenied struct {
	userID string
	action string
}

func (e *permissionDenied) Error() string {
	return "user with userID:" + e.userID + " don't have permission for '" + e.action + "'"
}

// lets also a .New to instansita an object of this type
func New(userID string, action string) permissionDenied {
	return permissionDenied{userID, action}
}

func authenticateUser(userID string) permissionDenied {
	return permissionDenied{
		userID: userID,
		action: "loggin in",
	}
}

func main() {
	user, err := createuser()
	if errors.Is(err, ErrUserNotFound) {
		fmt.Println(err)
	}
	fmt.Println(user)

	pd := New("2", "loggin in")
	fmt.Println(pd.Error())

	//now lets consider returning htis error type
	err2 := authenticateUser("2")
	fmt.Println("err2:", err2)
	var perr *permissionDenied
	if errors.As(&err2, &perr) {
		slog.Info("permission denied")
		zap.String("user_id", err2.userID)
		zap.String("action", err2.action)
	}
}


/*
to do - tomorrow
zap not working,
difference between errors.as and errors.is
reference and ponter brush up

*/