package user

type Service interface {
	CreateUser(id int) (User, error)
}

type service struct {
	repo Repository
}

func (s service) CreateUser(id int) (User, error){
	return s.repo.createUserUsingThisID(id)
}

func NewService(repo Repository) (Service, error) {
	return service{repo: repo}, nil
}