package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	//client library
	"github.com/go-redis/redis/v8"
)

const redisAddress = "localhost:6379"

var validateInputRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

//this is read as type getHandler equals this next thing
type getHandler struct {
	store *redis.Client
}

func main() {

	fmt.Println("Program Starting")
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})

	// How we know it works / testing
	key := "Inuyasha"
	value := "Kagome"
	set(context.Background(), rdb, key, value)
	output, err := get(rdb, key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(output)

	getH := getHandler{

		store: rdb,
	}
	//takes in a string and a handler
	http.Handle("/cache/", getH)

	//listen to incoming requests and serve them using previous functions
	//if exits with error, log the error and end program
	//it won't run until you exit the server.
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//give a key and a store, try to look up the value
//can write a precondition: assumption regarding the key (on the inputs)
//a post condition: the idea is providing that someone meets the precondition, then the post condition makes some guarantees.
//suppose you had a function that squared a number but suppose you are implementing it yourself and you can only do so on positive numbers
//so precondition would be number is positive, postcondition is squared number
func get(ctx context.Context, store *redis.Client, key string) (string, error) {

	value, err := store.Get(ctx, key).Result()

	if err != nil {
		return "", err
	}

	return value, nil
}

//function object params, return type
func set(ctx context.Context, store *redis.Client, key string, value string) {
	err := store.Set(ctx, key, value, 0).Err()
	return err
}

//io.writestring is a function that takes a generic object and writes to it (ie could take a file or responsewriter or web socket
//if you were to put in the browser ?get=hello that is a query )
//ServeHTTP is a function on a getHandler
func (h getHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	//giving it a name so we don't have to worry about getting it from the URL anymore
	//when we do things like hit redis we don't want a huge query so a straightaway name assists with the URL and it will be a nice string
	key := strings.Split(req.URL.Path, "/")[2] //req.URL.Query().Get("key")
	//log.Println(key)
	//make sure key is useful now.
	if key == "" {
		io.WriteString(w, "You need to give me STUFF\n")
		return
	}
	//now we want to check it actually is a valid key (input)
	//! means "if not"
	isValid := validateInputRegex.MatchString(key)
	if !isValid {
		io.WriteString(w, "Input needs to be anime.\n")
		return
	}

	if req.Method == http.MethodGet {
		h.handleGetRequest(ctx, key, w, req)
	} else if req.Method == http.MethodPost {
		h.handlePostRequest(ctx, key, w, req)
	}
}

func (h getHandler) handlePostRequest(ctx context.Context, key string, w http.ResponseWriter, req *http.Request) {
	value := req.URL.Query().Get("value")
	err := set(ctx, h.store, key, value)
	if err != nil {
		io.WriteString(w, "Error writing value to store.\n")
		//engineer sees this
		log.Println("Error while writing to store.", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h getHandler) handleGetRequest(ctx context.Context, key string, w http.ResponseWriter, req *http.Request) {
	//now we have a valid key so now we can do get store
	value, err := get(ctx, h.store, key)

	type name struct {
		key   string
		value string
	}
	myStruct := name{
		key:   "hello",
		value: "hello",
	}
	//I'm gonna use a structure which has a key and a value.
	//anonymous struct - a struct that you only use in one place
	jsonbytes, err := json.Marshal(struct {
		Key   string
		Value string
	}{
		Key:   key,
		Value: value,
	})
	if err != nil {
		//user sees this message
		io.WriteString(w, "Error processing.\n")
		//engineer/person maintaining sees this message
		log.Println("Error marshalling json.", err.Error())
		return
	}

	//will get the query parameter we called "get" and grab the value
	io.WriteString(w, string(jsonbytes))
	io.WriteString(w, "\n")

}

//to test this one you can use a curl or use Postman, an HTTP request builder.
//cant post request from browser
//so go into ubuntu (or any)terminal
//$curl -X POST localhost:8080/cache/?value=newValue
//when you refresh you want to see :newValue in the json
//in the terminal you could run
//curl -iv -X POST localhost:8080/cache/?value=newValue
//http request should reflect
//developer tools in googlechrome > network tab
//you can see the error header thingy
