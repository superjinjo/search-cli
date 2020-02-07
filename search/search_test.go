package search_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/superjinjo/zendesk-search/search"
)

func Test_SearchRepository_FindUsers(t *testing.T) {
	usersRepo := new(OrgUserMockRepo)
	orgsRepo := new(OrgUserMockRepo)
	ticketsRepo := new(TicketMockRepo)

	repo := search.NewSearchRepository(usersRepo, orgsRepo, ticketsRepo)

	users := []map[string]interface{}{
		{
			"_id":             float64(1),
			"name":            "same_name",
			"organization_id": float64(22),
		},
		{
			"_id":             float64(2),
			"name":            "user 2",
			"organization_id": float64(33),
		},
		{
			"_id":  float64(3),
			"name": "same_name",
		},
	}

	orgs := []map[string]interface{}{
		{
			"_id":  float64(22),
			"name": "org 1",
		},
		{
			"_id":  float64(33),
			"name": "org 2",
		},
	}

	tickets := []map[string]interface{}{
		{
			"_id":             "abcd",
			"subject":         "ticket 1",
			"assignee_id":     float64(1),
			"submitter_id":    float64(2),
			"organization_id": float64(22),
		},
		{
			"_id":             "efgh",
			"subject":         "ticket 2",
			"assignee_id":     float64(2),
			"submitter_id":    float64(1),
			"organization_id": float64(22),
		},
	}

	usersRepo.On("FindByField", "_id", float64(2)).Return([]map[string]interface{}{users[1]})
	usersRepo.On("FindByField", "name", "same_name").Return([]map[string]interface{}{users[0], users[2]})
	orgsRepo.On("FindByID", float64(22)).Return(orgs[0])
	orgsRepo.On("FindByID", float64(33)).Return(orgs[1])

	ticketsRepo.On("FindBySubmitter", float64(2)).Return([]map[string]interface{}{tickets[0]})
	ticketsRepo.On("FindByAssignee", float64(2)).Return([]map[string]interface{}{tickets[1]})

	ticketsRepo.On("FindBySubmitter", float64(1)).Return([]map[string]interface{}{tickets[1]})
	ticketsRepo.On("FindByAssignee", float64(1)).Return([]map[string]interface{}{tickets[0]})

	ticketsRepo.On("FindBySubmitter", float64(3)).Return([]map[string]interface{}{})
	ticketsRepo.On("FindByAssignee", float64(3)).Return([]map[string]interface{}{})

	results1 := repo.FindUsers("_id", float64(2))
	expected1 := []map[string]interface{}{
		{
			"_id":               float64(2),
			"name":              "user 2",
			"organization_id":   float64(33),
			"organization":      orgs[1],
			"submitted_tickets": []map[string]interface{}{tickets[0]},
			"assigned_tickets":  []map[string]interface{}{tickets[1]},
		},
	}

	require.Equal(t, expected1, results1)

	results2 := repo.FindUsers("name", "same_name")
	expected2 := []map[string]interface{}{
		{
			"_id":               float64(1),
			"name":              "same_name",
			"organization_id":   float64(22),
			"organization":      orgs[0],
			"submitted_tickets": []map[string]interface{}{tickets[1]},
			"assigned_tickets":  []map[string]interface{}{tickets[0]},
		},
		{
			"_id":               float64(3),
			"name":              "same_name",
			"submitted_tickets": []map[string]interface{}{},
			"assigned_tickets":  []map[string]interface{}{},
		},
	}

	require.Equal(t, expected2, results2)

}

func Test_SearchRepository_FindOrgs(t *testing.T) {
	usersRepo := new(OrgUserMockRepo)
	orgsRepo := new(OrgUserMockRepo)
	ticketsRepo := new(TicketMockRepo)

	repo := search.NewSearchRepository(usersRepo, orgsRepo, ticketsRepo)

	users := []map[string]interface{}{
		{
			"_id":             float64(1),
			"name":            "same_name",
			"organization_id": float64(22),
		},
		{
			"_id":             float64(2),
			"name":            "user 2",
			"organization_id": float64(33),
		},
		{
			"_id":  float64(3),
			"name": "same_name",
		},
	}

	orgs := []map[string]interface{}{
		{
			"_id":  float64(22),
			"name": "org 1",
		},
		{
			"_id":  float64(33),
			"name": "org 2",
		},
	}

	tickets := []map[string]interface{}{
		{
			"_id":             "abcd",
			"subject":         "ticket 1",
			"assignee_id":     float64(1),
			"submitter_id":    float64(2),
			"organization_id": float64(22),
		},
		{
			"_id":             "efgh",
			"subject":         "ticket 2",
			"assignee_id":     float64(2),
			"submitter_id":    float64(1),
			"organization_id": float64(22),
		},
	}

	orgsRepo.On("FindByField", "_id", float64(22)).Return([]map[string]interface{}{orgs[0]})
	orgsRepo.On("FindByField", "name", "org 2").Return([]map[string]interface{}{orgs[1]})

	ticketsRepo.On("FindByOrg", float64(22)).Return([]map[string]interface{}{tickets[0], tickets[1]})
	ticketsRepo.On("FindByOrg", float64(33)).Return([]map[string]interface{}{})

	usersRepo.On("FindByOrg", float64(22)).Return([]map[string]interface{}{users[0]})
	usersRepo.On("FindByOrg", float64(33)).Return([]map[string]interface{}{users[1]})

	results1 := repo.FindOrgs("_id", float64(22))
	expected1 := []map[string]interface{}{
		{
			"_id":     float64(22),
			"name":    "org 1",
			"users":   []map[string]interface{}{users[0]},
			"tickets": []map[string]interface{}{tickets[0], tickets[1]},
		},
	}

	require.Equal(t, expected1, results1)

	results2 := repo.FindOrgs("name", "org 2")
	expected2 := []map[string]interface{}{
		{
			"_id":     float64(33),
			"name":    "org 2",
			"users":   []map[string]interface{}{users[1]},
			"tickets": []map[string]interface{}{},
		},
	}

	require.Equal(t, expected2, results2)
}

func Test_SearchRepository_FindTickets(t *testing.T) {

	usersRepo := new(OrgUserMockRepo)
	orgsRepo := new(OrgUserMockRepo)
	ticketsRepo := new(TicketMockRepo)

	repo := search.NewSearchRepository(usersRepo, orgsRepo, ticketsRepo)

	users := []map[string]interface{}{
		{
			"_id":             float64(1),
			"name":            "same_name",
			"organization_id": float64(22),
		},
		{
			"_id":             float64(2),
			"name":            "user 2",
			"organization_id": float64(33),
		},
		{
			"_id":  float64(3),
			"name": "same_name",
		},
	}

	orgs := []map[string]interface{}{
		{
			"_id":  float64(22),
			"name": "org 1",
		},
		{
			"_id":  float64(33),
			"name": "org 2",
		},
	}

	tickets := []map[string]interface{}{
		{
			"_id":             "abcd",
			"subject":         "ticket 1",
			"assignee_id":     float64(1),
			"submitter_id":    float64(2),
			"organization_id": float64(22),
		},
		{
			"_id":             "efgh",
			"subject":         "ticket 2",
			"assignee_id":     float64(2),
			"submitter_id":    float64(1),
			"organization_id": float64(22),
		},
	}

	ticketsRepo.On("FindByField", "_id", "abcd").Return([]map[string]interface{}{tickets[0]})
	ticketsRepo.On("FindByField", "subject", "ticket 2").Return([]map[string]interface{}{tickets[1]})
	orgsRepo.On("FindByID", float64(22)).Return(orgs[0])
	usersRepo.On("FindByID", float64(1)).Return(users[0])
	usersRepo.On("FindByID", float64(2)).Return(users[1])

	results1 := repo.FindTickets("_id", "abcd")
	expected1 := []map[string]interface{}{
		{
			"_id":             "abcd",
			"subject":         "ticket 1",
			"assignee_id":     float64(1),
			"submitter_id":    float64(2),
			"organization_id": float64(22),
			"organization":    orgs[0],
			"submitted_user":  users[1],
			"assigned_user":   users[0],
		},
	}

	require.Equal(t, expected1, results1)

	results2 := repo.FindTickets("subject", "ticket 2")
	expected2 := []map[string]interface{}{
		{
			"_id":             "efgh",
			"subject":         "ticket 2",
			"assignee_id":     float64(2),
			"submitter_id":    float64(1),
			"organization_id": float64(22),
			"organization":    orgs[0],
			"submitted_user":  users[0],
			"assigned_user":   users[1],
		},
	}

	require.Equal(t, expected2, results2)

}

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) FindByOrg(orgID float64) []map[string]interface{} {
	args := m.Called(orgID)

	return args.Get(0).([]map[string]interface{})
}

func (m *MockRepository) FindByField(fieldName string, searchVal interface{}) []map[string]interface{} {
	args := m.Called(fieldName, searchVal)

	return args.Get(0).([]map[string]interface{})
}

type OrgUserMockRepo struct {
	MockRepository
}

func (m *OrgUserMockRepo) FindByID(mockID float64) map[string]interface{} {
	args := m.Called(mockID)

	return args.Get(0).(map[string]interface{})
}

type TicketMockRepo struct {
	MockRepository
}

func (m *TicketMockRepo) FindByID(mockID string) map[string]interface{} {
	args := m.Called(mockID)

	return args.Get(0).(map[string]interface{})
}

func (m *TicketMockRepo) FindBySubmitter(userID float64) []map[string]interface{} {
	args := m.Called(userID)

	return args.Get(0).([]map[string]interface{})
}

func (m *TicketMockRepo) FindByAssignee(userID float64) []map[string]interface{} {
	args := m.Called(userID)

	return args.Get(0).([]map[string]interface{})
}
