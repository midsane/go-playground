package app

import (
	"fmt"
	"math/rand"
)

func Greet() string {
	formattedStrings := []string{
		"Hey There, How you doin %v",
		"Welcome %v bro",
		"Nice Weather, ain't it %v",
	}

	return formattedStrings[rand.Intn(len(formattedStrings))]
}

func GreetToCli(name string, formal bool) string {
	if formal {
		return fmt.Sprintf("Good day, Sir %v", name)
	}
	formattedStrings := []string{
		"Hey There, How you doin %v",
		"Welcome %v bro",
		"Nice Weather, ain't it %v",
	}

	formatter := formattedStrings[rand.Intn(len(formattedStrings))]
	return fmt.Sprintf(formatter, name)
}
