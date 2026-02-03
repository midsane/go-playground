package main

import (
	"fmt"
	"github.com/midsane/go-playground/02-cli-app/internal/app"
)
/*
we import packages not files, hence importing app package in this path,
using app.Greet function here. go.mod is defined at root so that we can track cmd/internal here.
so cmd programs can use function from internal packages.
*/

func main() {
	/*
		print greeting message depending, get a person's name
	*/
	name := "satmak"
	msg := fmt.Sprintf(app.Greet(), name)
	fmt.Println(msg)

}
