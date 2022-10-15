package main

import (
	"encoding/json"
	"fmt"
	"io"
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


type DATA struct{
	Key string 
	Value string 
}
var db []DATA

func writeToDB(key, value string) bool  {	
    result, error := json.Marshal(DATA{Key:key, Value: value})
    if error != nil {
        fmt.Println(error)
    }

    f, erro := os.OpenFile("db.json", os.O_APPEND|os.O_WRONLY, 0666)
    if erro != nil {
        fmt.Println(erro)
    }

    n, err := io.WriteString(f, string(result))
    if err != nil {
        fmt.Println(n, err)
    }
    return true
}


func readFromDB(key string) string {
	file, err := ioutil.ReadFile("db.json")
	var obj DATA
	err = json.Unmarshal(file, &obj)
	if err != nil {
		fmt.Println(err)
	}
	if obj.Key == key {
		return "Key found: " + obj.Value
	}

	return "Key not Found!";
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