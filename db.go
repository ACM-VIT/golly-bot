package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func isExistDB() bool{
	if _, err := os.Stat("db.json"); err == nil {
		return true;
	} else {
		return false
	}
}

func createDB() bool {
	os.Create("db.json");
	return true
}

func writeToDB(key, value string) bool  {
	if isExistDB() {

		jsonFile, err := os.Open("db.json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		// Declared an empty map interface
		var result map[string]interface{}
		result = make(map[string]interface{})
		
		// Unmarshal or Decode the JSON to the interface.
		json.Unmarshal(byteValue, &result)

		result[key] = value

		// Marshal or Encode the interface data
		jsonResult, _ := json.Marshal(result)
		// Write the JSON data to the file
		ioutil.WriteFile("db.json", jsonResult, 0644)

		return true
	} else {
		return false
	}
}


func readFromDB(key string) (string, bool) {
	if isExistDB() {
		jsonFile, err := os.Open("db.json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		// Declared an empty map interface
		var result map[string]interface{}

		// Unmarshal or Decode the JSON to the interface.
		json.Unmarshal(byteValue, &result)

		// Reading each value by its key
		if result[key] != nil {
			return result[key].(string), true
		} else {
			return "", false
		}
	} else {
		return "", false
	}
}


func main() {
	fmt.Println("Enter the number:\n1. Check the DataBase\n2. Make a new DataBase\n3. Add data to the database\n4. Check wheather the data is exists or not in the DataBase")
	var mainCase int
	fmt.Scan(&mainCase);
	switch mainCase {
	case 1:
		fmt.Println(isExistDB())
		break
	case 2:
		fmt.Println(createDB())
		break
	case 3:
		fmt.Println("Enter Key and Value")
		var key string
		var value string
		fmt.Scanln(&key)
		fmt.Scanln(&value)
		fmt.Println(writeToDB(key, value))
		break
	case 4:
		var key string
		fmt.Scanln(&key)
		fmt.Println(readFromDB(key))
		break
	default:
		fmt.Errorf("No Number Pressed")
		
	}
		
}
