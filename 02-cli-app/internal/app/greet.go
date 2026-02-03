package app

import "math/rand"

func Greet() string {
	formattedStrings := []string{
		"Hey There, How you doin %v",
		"Welcome %v bro",
		"Nice Weather, ain't it %v",
	}

	return formattedStrings[rand.Intn(len(formattedStrings))]
}
