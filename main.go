package main

import (
	"fmt"
	"net/http"
	"os"
	"io/ioutil"
	"log"
	"encoding/json"
	"encoding/csv"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	fmt.Println("Starting server...")

	// Create router
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// Load data from endpoints (buyers, products, transactions)
	router.Get("/getData", func(w http.ResponseWriter, r *http.Request) {
		// Get buyers data
		buyerResponse, err := http.Get("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/buyers")

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
	
		buyerResponseBody, err := ioutil.ReadAll(buyerResponse.Body)
		if err != nil {
			log.Fatal(err)
		}

		var buyerResponseJSON []Buyer
		json.Unmarshal(buyerResponseBody, &buyerResponseJSON)

		//fmt.Println(buyerResponseJSON) // Print buyers

		// Get products data
		productResponse, err := http.Get("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/products")

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		productResponseBody, err := ioutil.ReadAll(productResponse.Body)
		if err != nil {
			log.Fatal(err)
		}

		reader := csv.NewReader(strings.NewReader(string(productResponseBody)))
		reader.Comma = '\'' // Use different delimter

		readerData, err := reader.ReadAll()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var productResponseJSON []Product

		for _, line := range readerData {
			priceValue, err := strconv.Atoi(line[2])
			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}

			tempProduct := Product {
				Id: line[0],
				Name: line[1],
				Price: priceValue,
			}

			productResponseJSON = append(productResponseJSON, tempProduct)
		}

		//fmt.Println(productResponseJSON) // Print products

		// Get transactions data
		transactionResponse, err := http.Get("https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/transactions")

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		transactionResponseBody, err := ioutil.ReadAll(transactionResponse.Body)
		if err != nil {
			log.Fatal(err)
		}
		tempBodyString := string(transactionResponseBody)

		var transactionResponseBodyString string
		doubleZero := 0

		fmt.Println("Began execution")

		for _, letter := range tempBodyString {
			r, _ := utf8.DecodeRuneInString(string(letter)) // _ = size
			//fmt.Printf("%d %v\n", r, size)
			
			if r == 0 && doubleZero == 1 {
				transactionResponseBodyString = transactionResponseBodyString[:len(transactionResponseBodyString)-1]
				transactionResponseBodyString += string('\n')
				doubleZero++
			}else if r == 0 {
				//fmt.Println("Found a 0!")
				transactionResponseBodyString += string('\'')
				doubleZero++
			}else if string(r) == "(" || string(r) == ")" {
				// Do nothing
			}else {
				doubleZero = 0

				transactionResponseBodyString += string(letter)
			}
		}

		fmt.Println("Finished execution!")
		//fmt.Println(transactionResponseBodyString) // Print transactions

		reader = csv.NewReader(strings.NewReader(transactionResponseBodyString))
		reader.Comma = '\'' // Use different delimter

		readerData, err = reader.ReadAll()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var transactionResponseJSON []Transaction

		for _, line := range readerData {
			productsIdsArray := strings.Fields(line[4])

			tempTransaction := Transaction {
				Id: line[0],
				BuyerId: line[1],
				IP: line[2],
				Device: line[3],
				ProductsIds: productsIdsArray,
			}

			transactionResponseJSON = append(transactionResponseJSON, tempTransaction)
		}

		fmt.Println(transactionResponseJSON) // Print transactions

		w.Write([]byte("Testing..."))
	})

	// Inform where the application is running
	http.ListenAndServe("127.0.0.1:3000", router)
}

// Buyers struct
type Buyer struct {
    Id string `json:"id"`
    Name string `json:"name"`
	Age int `json:"age"`
}

// Products struct
type Product struct {
    Id string `json:"id"`
    Name string `json:"name"`
	Price int `json:"price"`
}

// Transactions struct
type Transaction struct {
    Id string `json:"id"`
    BuyerId string `json:"buyerId"`
	IP string `json:"ip"`
	Device string `json:"device"`
	ProductsIds []string `json:"productsIds"`
}