package search

import (
	"github.com/pkg/errors"
)

//FYI this is why I use float64: https://golang.org/pkg/encoding/json/#Unmarshal
type TicketJSONRepository struct {
	ticketsIndex   map[string]map[string]interface{} //map of json data indexed by ticket ID
	orgsIndex      map[float64][]string              //map of ticket IDs indexed by org ID
	submitterIndex map[float64][]string              //map of ticket IDs indexed by submitter ticket ID
	assigneeIndex  map[float64][]string              //map of ticket IDs indexed by assignee ticket ID
	valueMatcher   ValueMatcher
}

func NewTicketJSONRepository(tickets []map[string]interface{}) (*TicketJSONRepository, error) {

	repository := &TicketJSONRepository{
		ticketsIndex:   make(map[string]map[string]interface{}),
		orgsIndex:      make(map[float64][]string),
		submitterIndex: make(map[float64][]string),
		assigneeIndex:  make(map[float64][]string),
		valueMatcher:   SearchValueMatches,
	}

	for i, ticket := range tickets {
		if err := repository.addTicket(ticket); err != nil {
			return nil, errors.WithMessagef(err, "Error with ticket at index %d", i)
		}
	}

	return repository, nil
}

//SetValueMatcher lets you set a different matcher which is useful for testing
func (repo *TicketJSONRepository) SetValueMatcher(matcherFn ValueMatcher) {
	repo.valueMatcher = matcherFn
}

func (repo *TicketJSONRepository) addTicket(ticket map[string]interface{}) error {
	ticketID, isString := ticket["_id"].(string) //FYI: in go, if a map doesn't have a key, it simply returns nil
	if !isString {
		return errors.New("User is missing \"_id\" field or \"_id\" is not string")
	}

	if _, exists := repo.ticketsIndex[ticketID]; exists {
		return errors.Errorf("User with ID of %v already exists", ticketID)
	}

	repo.ticketsIndex[ticketID] = ticket

	if orgID, isFloat := ticket["organization_id"].(float64); isFloat {
		repo.orgsIndex[orgID] = append(repo.orgsIndex[orgID], ticketID)
	} else {
		repo.orgsIndex[0] = append(repo.orgsIndex[0], ticketID)
	}

	if submitterID, isFloat := ticket["submitter_id"].(float64); isFloat {
		repo.submitterIndex[submitterID] = append(repo.submitterIndex[submitterID], ticketID)
	} else {
		repo.submitterIndex[0] = append(repo.submitterIndex[0], ticketID)
	}

	if assigneeID, isFloat := ticket["assignee_id"].(float64); isFloat {
		repo.assigneeIndex[assigneeID] = append(repo.assigneeIndex[assigneeID], ticketID)
	} else {
		repo.assigneeIndex[0] = append(repo.assigneeIndex[0], ticketID)
	}

	return nil
}

func (repo *TicketJSONRepository) FindByID(ticketID string) map[string]interface{} {
	return repo.ticketsIndex[ticketID]
}

func (repo *TicketJSONRepository) FindByOrg(orgID float64) []map[string]interface{} {
	ticketIDs := repo.orgsIndex[orgID]

	ticketList := make([]map[string]interface{}, len(ticketIDs))

	for i := 0; i < len(ticketIDs); i++ {
		nextID := ticketIDs[i]

		if ticket := repo.FindByID(nextID); ticket != nil {
			ticketList[i] = ticket
		}
	}

	return ticketList
}

func (repo *TicketJSONRepository) FindBySubmitter(userID float64) []map[string]interface{} {
	ticketIDs := repo.submitterIndex[userID]

	ticketList := make([]map[string]interface{}, len(ticketIDs))

	for i := 0; i < len(ticketIDs); i++ {
		nextID := ticketIDs[i]

		if ticket := repo.FindByID(nextID); ticket != nil {
			ticketList[i] = ticket
		}
	}

	return ticketList
}

func (repo *TicketJSONRepository) FindByAssignee(userID float64) []map[string]interface{} {
	ticketIDs := repo.assigneeIndex[userID]

	ticketList := make([]map[string]interface{}, len(ticketIDs))

	for i := 0; i < len(ticketIDs); i++ {
		nextID := ticketIDs[i]

		if ticket := repo.FindByID(nextID); ticket != nil {
			ticketList[i] = ticket
		}
	}

	return ticketList
}

func (repo *TicketJSONRepository) FindByField(fieldName string, searchVal interface{}) []map[string]interface{} {
	switch fieldName {
	case "_id":
		var ticketList []map[string]interface{}

		ticketID := stringVal(searchVal)
		if ticket := repo.FindByID(ticketID); ticket != nil {
			ticketList = append(ticketList, ticket)
		}

		return ticketList

	case "organization_id":

		if orgID, isFloat := floatVal(searchVal); isFloat {
			return repo.FindByOrg(orgID)
		}

		return []map[string]interface{}{}

	case "submitter_id":

		if userID, isFloat := floatVal(searchVal); isFloat {
			return repo.FindBySubmitter(userID)
		}

		return []map[string]interface{}{}

	case "assignee_id":

		if userID, isFloat := floatVal(searchVal); isFloat {
			return repo.FindByAssignee(userID)
		}

		return []map[string]interface{}{}

	default:
		var ticketList []map[string]interface{}

		for _, ticket := range repo.ticketsIndex {
			if repo.valueMatcher(ticket[fieldName], searchVal) {
				ticketList = append(ticketList, ticket)
			}
		}
		return ticketList
	}
}
