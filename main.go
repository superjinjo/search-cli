package main

import (
	"encoding/json"
	"log"

	"github.com/markbates/pkger"
	"github.com/superjinjo/zendesk-search/cmd"
	"github.com/superjinjo/zendesk-search/search"
)

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
	//pkger requires that you use hardcoded strings with their functions
	//in order to properly pack the files into the binary
	pkger.Include("/data/users.json")
	userData, err := openJSONFile("/data/users.json")
	if err != nil {
		return nil, err
	}

	userRepo, err := search.NewUserJSONRepository(userData)
	if err != nil {
		return nil, err
	}

	pkger.Include("/data/organizations.json")
	orgData, err := openJSONFile("/data/organizations.json")
	if err != nil {
		return nil, err
	}

	orgRepo, err := search.NewOrgJSONRepository(orgData)
	if err != nil {
		return nil, err
	}

	pkger.Include("/data/tickets.json")
	ticketData, err := openJSONFile("/data/tickets.json")
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
