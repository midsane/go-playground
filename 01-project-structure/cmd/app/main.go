package main
import (
	"log"
	"github.com/midsane/go-playground/01-project-structure/internal/user"
)
func main(){
	myRepo := user.NewRepository()
	myService, err := user.NewService(myRepo)
	if err != nil {
		log.Fatal("error in creating new service")
	}	
	myService.CreateUser(2)
}