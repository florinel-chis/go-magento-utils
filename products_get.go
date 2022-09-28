package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

func get_base_url() string {
	url := os.Getenv("BASE_URL")
	final_url := ""
	final_url = url + "/rest/all/V1/products?searchCriteria[pageSize]=20&searchCriteria[currentPage]="
	return final_url
}

func get_token() string {
	token := os.Getenv("M24_TOKEN")
	return token
}

func get_products(page_no int, folder_name string) {
	token := get_token()
	base_url := get_base_url()
	url := base_url + strconv.Itoa(page_no)
	var bearer = "Bearer " + token

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	//log.Println(string([]byte(body)))

	//save body to file
	err = ioutil.WriteFile(folder_name+"/products_"+strconv.Itoa(page_no)+".json", body, 0644)
	if err != nil {
		log.Println("Error while saving the response bytes:", err)
	}
}

func get_products_from_json(filename string) map[string]interface{} {
	jsonFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened first page.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	return result
}

func main() {

	//create a folder with the current timestamp
	currentTime := time.Now()
	folder_name := currentTime.Format("2006-01-02-15_04_05")
	//check if the folder exists
	if _, err := os.Stat(folder_name); os.IsNotExist(err) {
		//create the folder
		os.Mkdir(folder_name, os.ModePerm)
	}
	get_products(1, folder_name)
	// Open our jsonFile
	fname := "./" + folder_name + "/products_1.json"
	result := get_products_from_json(fname)
	sc := result["search_criteria"].(map[string]interface{})
	total_count := result["total_count"].(float64)
	//current_page := sc["current_page"].(float64)
	page_size := sc["page_size"].(float64)
	number_of_pages := math.Ceil(total_count / page_size)
	if number_of_pages > 1 {
		for i := 2; i <= int(number_of_pages); i++ {
			get_products(i, folder_name)
		}
	}
	fmt.Println("Total count:", total_count)
	fmt.Println("Pages:", number_of_pages)
	fmt.Println(sc["current_page"])
}
