package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/cklewar/form3_clib/api/client"
)

func main() {

	// Initialize client parameters
	parameters := client.Parameters{
		Host:     "192.168.2.50",
		BaseURI:  "/v1/organisation/",
		Resource: "accounts",
	}

	// Construct client
	cStatus, c := client.NewClient(parameters)

	if cStatus {
		//status, body := c.Fetch("ad27e265-9604-4b4b-a0e5-3003ea9cc4db")
		//fmt.Println("Status:", status)
		// fmt.Println(client.JSONPrettyPrint(body))
		//var account response.OrganisationAccountData
		//json.Unmarshal([]byte(body), &account)
		//fmt.Println("ID:", account.Data.ID)
		//fmt.Println("TYPE:", account.Data.Type)

		// We use input file for feeding JSON data into create operation
		cwd, err := os.Getwd()
		path := filepath.Join(cwd, "examples/json/org_acc_create.json")
		createInputData, err := ioutil.ReadFile(path)

		if err != nil {
			log.Fatal(err)
		}

		createStatus, createResp := c.Create(createInputData)
		fmt.Println("Status: ", createStatus)
		fmt.Println(client.JSONPrettyPrint(createResp))

		fetchStatus, fetchResp := c.Fetch(c.GetObjID(createResp))
		fmt.Println("Status: ", fetchStatus)
		fmt.Println(client.JSONPrettyPrint(fetchResp))

		deleteStatus := c.Delete(c.GetObjID(createResp), c.GetObjVersion(createResp))
		fmt.Println("Status: ", deleteStatus)
	}
}
