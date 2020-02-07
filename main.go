package main

import (
	"encoding/json"
	"log"

	"github.com/markbates/pkger"
	"github.com/superjinjo/zendesk-search/cmd"
	"github.com/superjinjo/zendesk-search/search"
)

//pkger requires that file paths be absolute to the module root
const userFilePath = "/data/users.json"
const orgFilePath = "/data/organizations.json"
const ticketFilePath = "/data/tickets.json"

func openJSONFile(pkgerPath string) ([]map[string]interface{}, error) {
	file, err := pkger.Open(pkgerPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	jsonDecoder := json.NewDecoder(file)

	var fileJSON []map[string]interface{}

	if err := jsonDecoder.Decode(&fileJSON); err != nil {
		return nil, err
	}

	return fileJSON, nil
}

func buildRepository() (*search.SearchRepository, error) {
	userData, err := openJSONFile(userFilePath)
	if err != nil {
		return nil, err
	}

	userRepo, err := search.NewUserJSONRepository(userData)
	if err != nil {
		return nil, err
	}

	orgData, err := openJSONFile(orgFilePath)
	if err != nil {
		return nil, err
	}

	orgRepo, err := search.NewOrgJSONRepository(orgData)
	if err != nil {
		return nil, err
	}

	ticketData, err := openJSONFile(ticketFilePath)
	if err != nil {
		return nil, err
	}

	ticketRepo, err := search.NewTicketJSONRepository(ticketData)
	if err != nil {
		return nil, err
	}

	return search.NewSearchRepository(userRepo, orgRepo, ticketRepo), nil
}

func main() {
	searchRepo, err := buildRepository()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Execute(searchRepo)
}
