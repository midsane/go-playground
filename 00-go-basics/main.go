package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
)

//lets understand types in general in go
//type keyword created a new named type
//types basically define how the values are gonna look like in memory
//what methods are associated with them
//what operation are allowed on them
//ex int -> (can )
//GO HAVE no classes, only type+methods

//types are of two kinds -> builtin like int32
//and composite -> user defined -> built using type keyword

//struct -> when u need multiple fields , methods attached to it
//interface -> a set of method signatures

// lets define reader interface-> Read function will accept stream of byte and return a int value from it
type Reader interface {
	Read(p []byte) (n int, err error)
}

// now lets define a struct that will implement it
type File struct{}

func (f File) Read(p []byte) (n int, err error) {
	if len(p) > 0 {
		return int(p[0]), nil
	}
	return 0, errors.New("could not read from file stream")
}

// func types
// a named function signature
// lets give named type to a function signature func(a int32)
type func1 func(a int32)

// basically if we directly want to give a type to function that is not part of a struct
// we can also add more methods to it like this
func (f func1) additionalFunc(val bool) bool {
	return !val
}

// lets create a function that accepts another funciton of type func1
// this is how callbacks works in go
func type2(f1 func1) {
	f1(2)
	fmt.Println(f1.additionalFunc(true))
}

// now lets give name to a single type -> why name it again
// name it again after changing it a bit like a making fixes sized array (not a slice)
type A [3]int

/*
alias vs new type
*/
// type newInt int
//here newInt is a new type

type MyInt = int

//here MyInt is just a new alias

/*
Types define meaning
Interfaces define capability
Functions define behavior
Structs define state
*/

/*
Comparable constraint in golang -
it is a constraint that means the given type is comparable with == and !=.
so it's important to use in generic function signature if implementation uses equalities.
otherwise any type too can be uses.
any is actually type interface{}
*/

// generics
func findInSlice[t comparable](arr []t, val t) int {
	//this fucntion finds if chutmal exist in sexual
	for i, v := range arr {
		if v == val {
			return i
		}
	}
	return -1
}

// if function dont need comparsion , just use any
func printEl[t any](arr []t) {
	for _, v := range arr {
		fmt.Println(v)
	}
}

// goroutines and channels
func goRoutines() {
	fmt.Println("hello")
	for i := range 4 {
		go fmt.Println(i)
	}
	fmt.Println("world")
}

func populate(c chan int) {
	val := 0
	for i := range 10 {
		val += i + rand.Intn(i+1)
	}
	c <- val
}

func DummyFunc() (e error) {
	e = errors.New("yoy")
	panic("panic with err")
}

func Yep(f func() (e error)) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("err: %v", r)
		}
	}()
	return f()
}

/*
golang reflect -
*/
func ReflectPrimer(){
	var lavru any = 12
	tl := reflect.TypeOf(lavru)
	//at complile time(static type) of lavru is any -> but runtime type will be int, can be checked by reflect.TypeOf(variable0)
	vl := reflect.ValueOf(lavru)
	//even runtime value can be inspected via reflect.ValueOf
	fmt.Println(tl,vl)
}


func main() {
	file := File{}
	dta, err := file.Read([]byte{2, 3})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Println(dta)
	type2(func(a int32) {
		fmt.Println("custom type1 function")
	})
	fmt.Println(findInSlice([]int{1, 2, 3}, 3))
	printEl([]int{1, 2, 3, 4})

	goRoutines()
	fmt.Println("------------")

	c := make(chan int)
	for range 4 {
		go populate(c)
	}

	// a, b, c2, d := <-c, <-c, <-c, <-c
	// fmt.Println(a, b, c2, d)

	fmt.Println("------------")

	//if you already max capacit of the array just add that to avoid again and again
	//realocation of memory for array after everytime lenght exceeds capacity
	out := make([]int, 0, 2)
	for i := range 2 {
		out = append(out, i)
	}
	// fmt.Println(out)

	// fmt.Println("------------")

	// fmt.Println(yep(dummyFunc))

	// StartSever()

	
	fmt.Println("------------")
	ReflectPrimer()
}
