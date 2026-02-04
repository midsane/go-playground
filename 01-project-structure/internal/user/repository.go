/*
define an interface and struct that implements it
*/
package user

type Repository interface {
	createUserUsingThisID(id int) (User, error)
}

type memoRepo struct {}

func (mr memoRepo) createUserUsingThisID(id int) (User, error){
	return User{ID: id, Name: "satmak"}, nil
}

func NewRepository() Repository {
	return memoRepo{}
}


