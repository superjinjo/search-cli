package search

type UserRepository interface {
	FindByID(userID float64) map[string]interface{}
	FindByOrg(orgID float64) []map[string]interface{}
	FindByField(fieldName string, searchVal interface{}) []map[string]interface{}
}

type OrgRepository interface {
	FindByID(orgID float64) map[string]interface{}
	FindByField(fieldName string, searchVal interface{}) []map[string]interface{}
}

type TicketRepository interface {
	FindByID(ticketID string) map[string]interface{}
	FindByAssignee(userID float64) []map[string]interface{}
	FindBySubmitter(userID float64) []map[string]interface{}
	FindByOrg(orgID float64) []map[string]interface{}
	FindByField(fieldName string, searchVal interface{}) []map[string]interface{}
}

type SearchRepository struct {
	userRepository   UserRepository
	orgRepository    OrgRepository
	ticketRepository TicketRepository
}

func NewSearchRepository(users UserRepository, orgs OrgRepository, tickets TicketRepository) *SearchRepository {
	return &SearchRepository{
		userRepository:   users,
		orgRepository:    orgs,
		ticketRepository: tickets,
	}
}

func (repo *SearchRepository) findOrgRelation(item map[string]interface{}) map[string]interface{} {
	if orgID, isFloat := item["organization_id"].(float64); isFloat {
		return repo.orgRepository.FindByID(orgID)
	}

	return nil
}

func (repo *SearchRepository) FindUsers(fieldName string, searchValue interface{}) []map[string]interface{} {
	users := repo.userRepository.FindByField(fieldName, searchValue)

	for i, user := range users {
		org := repo.findOrgRelation(user)
		if org != nil {
			user["organization"] = org
		}

		if userID, isFloat := user["_id"].(float64); isFloat {
			user["submitted_tickets"] = repo.ticketRepository.FindBySubmitter(userID)
			user["assigned_tickets"] = repo.ticketRepository.FindByAssignee(userID)
		}

		users[i] = user
	}

	return users
}

func (repo *SearchRepository) FindOrgs(fieldName string, searchValue interface{}) []map[string]interface{} {
	orgs := repo.orgRepository.FindByField(fieldName, searchValue)

	for i, org := range orgs {
		if orgID, isFloat := org["_id"].(float64); isFloat {
			org["users"] = repo.userRepository.FindByOrg(orgID)
			org["tickets"] = repo.ticketRepository.FindByOrg(orgID)
		}

		orgs[i] = org
	}

	return orgs
}

func (repo *SearchRepository) FindTickets(fieldName string, searchValue interface{}) []map[string]interface{} {
	tickets := repo.orgRepository.FindByField(fieldName, searchValue)

	for i, ticket := range tickets {
		org := repo.findOrgRelation(ticket)
		if org != nil {
			ticket["organization"] = org
		}

		if userID, isFloat := ticket["submitter_id"].(float64); isFloat {
			ticket["submitted_user"] = repo.userRepository.FindByID(userID)
		}

		if userID, isFloat := ticket["assignee_id"].(float64); isFloat {
			ticket["assigned_user"] = repo.userRepository.FindByID(userID)
		}

		tickets[i] = ticket
	}

	return tickets
}
