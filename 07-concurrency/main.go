package concurrency

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)


func main(){
	log.Fatal(http.ListenAndServe(":8080", nil))
	gin.Default()
}