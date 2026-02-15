package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

/*defere use case -> resouceClenup*/
func ReadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	// if file exist, before this func returns file.Close() will be called
	fmt.Println(file.Name())
	return nil
}

func Sample() error {
	return ReadFile("sample.txt")
}

// before this func returns tx.Rollback will be called if db.Begin() dont throw err
func PerformDBOperation(db *sql.DB) error {
	tx, err := db.Begin()
	//get a hold of tx type value first
	//if err happens, return err
	if err != nil {
		return err
	}
	//now since tx is defined, defer a rollback
	defer tx.Rollback()
	//execute queries
	//at each step return err if it occurs
	_, err = tx.Exec("INSERT into..")
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	//return nil err
	return nil
}

// Timers and Benchmarking
func MeasureTimeElapsed() {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Println("elapsed time: ", elapsed)
	}()
	time.Sleep(3 * time.Second)
}

//shutting down server with logs

func StartSever() {
	server := &http.Server{Addr: ":8080"}
	//cleanup->afterw this func ends
	defer func() {
		log.Println("shutting down server and cleaning up...")
	}()
	//ListenAndServe() will return when err is encountered, or if server shuts down
	log.Fatal(server.ListenAndServe())
}

type ParsedStruct struct{}

func SafeParsing(data []byte) (ps ParsedStruct, err error){
	defer func(){
		if r := recover(); r != nil {
			err = fmt.Errorf("parsed err: %v", r)
		}
	}()
	return unsafeParsing(data)
}

func unsafeParsing(data []byte) (ps ParsedStruct, err error){
	if len(data) == 0 {
		panic("lenght of data is 0, cant parse")
	}
	return ParsedStruct{}, nil
}
