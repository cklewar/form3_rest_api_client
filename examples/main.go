package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/cklewar/form3_rest_api_client/api/client"
	"github.com/cklewar/form3_rest_api_client/api/response"
)

func main() {
	// Initialize api parameters
	parameters := client.Parameters{
		BaseURI:  "/v1/organisation/",
		Resource: "accounts",
	}

	// Construct client with default values
	c, err := client.NewClient("192.168.2.50", "", "", parameters)

	if err != nil {
		log.Fatal(err)
	}

	// We use input file for feeding JSON data into create operation
	cwd, err := os.Getwd()
	path := filepath.Join(cwd, "examples/json/org_acc_create.json")
	createInputData, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	createResp, err := c.Create(createInputData, 0)
	fmt.Println("Error: ", err)
	fmt.Println("ResponseCode: ", createResp.Code)
	data, err := client.JSONPrettyPrint(createResp.Body)
	fmt.Println(data)

	id, _ := client.GetObjID(createResp.Body)
	fetchResp, err := c.Fetch(id, 0)
	fmt.Println("Error: ", err)
	fmt.Println("ResponseCode: ", fetchResp.Code)
	data, err = client.JSONPrettyPrint(fetchResp.Body)
	fmt.Println(data)
	var account response.OrganisationAccountData
	json.Unmarshal([]byte(fetchResp.Body), &account)
	fmt.Println("Id:", account.Data.ID)
	fmt.Println("Type:", account.Data.Type)

	version, _ := client.GetObjVersion(createResp.Body)
	deleteResp, err := c.Delete(id, version, 0)
	fmt.Println("Error: ", err)
	fmt.Println("ResponseCode: ", deleteResp.Code)

	parameters = client.Parameters{
		BaseURI:  "/v2/organisation/",
		Resource: "accounts",
	}

	c.UpdateParameters(parameters)
}
