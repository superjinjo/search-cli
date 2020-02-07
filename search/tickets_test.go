package search_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/superjinjo/zendesk-search/search"
)

func Test_NewTicketRepository(t *testing.T) {
	//empty list is valid
	emptyList := []map[string]interface{}{}
	_, err1 := search.NewTicketRepository(emptyList)
	require.Nil(t, err1)

	//single item in list
	goodList1 := []map[string]interface{}{
		{
			"_id":             "436bf9b0-1147-4c0a-8439-6f79833bff5b",
			"subject":         "Nicole Martinez ticket",
			"organization_id": float64(123),
			"assignee_id":     float64(1),
			"submitter_id":    float64(2),
		},
	}
	_, err2 := search.NewTicketRepository(goodList1)
	require.Nil(t, err2)

	//no organization_id, submitter_id, and assignee_id  is okay
	goodList2 := []map[string]interface{}{
		{
			"_id":             "436bf9b0-1147-4c0a-8439-6f79833bff5b",
			"subject":         "Nicole Martinez ticket",
			"organization_id": float64(123),
			"assignee_id":     float64(1),
			"submitter_id":    float64(2),
		},
		{
			"_id":     "1a227508-9f39-427c-8f57-1b72f3fab87c",
			"subject": "Johnny Jarvis ticket",
		},
	}
	_, err3 := search.NewTicketRepository(goodList2)
	require.Nil(t, err3)

	//"_id" field is required
	badList1 := []map[string]interface{}{
		{
			"_id":             "436bf9b0-1147-4c0a-8439-6f79833bff5b",
			"subject":         "Nicole Martinez ticket",
			"organization_id": float64(123),
			"assignee_id":     float64(1),
			"submitter_id":    float64(2),
		},
		{
			"subject":         "Johnny Jarvis ticket",
			"organization_id": float64(123),
		},
	}
	_, err4 := search.NewTicketRepository(badList1)
	require.NotNil(t, err4)

	//_id must be a string
	badList2 := []map[string]interface{}{
		{
			"_id":             float64(199),
			"subject":         "Nicole Martinez ticket",
			"organization_id": float64(123),
			"assignee_id":     float64(1),
			"submitter_id":    float64(2),
		},
		{
			"_id":     "1a227508-9f39-427c-8f57-1b72f3fab87c",
			"subject": "Johnny Jarvis ticket",
		},
	}
	_, err5 := search.NewTicketRepository(badList2)
	require.NotNil(t, err5)

	//duplicate _id fields
	badList3 := []map[string]interface{}{
		{
			"_id":             "1a227508-9f39-427c-8f57-1b72f3fab87c",
			"subject":         "Nicole Martinez ticket",
			"organization_id": float64(123),
			"assignee_id":     float64(1),
			"submitter_id":    float64(2),
		},
		{
			"_id":     "1a227508-9f39-427c-8f57-1b72f3fab87c",
			"subject": "Johnny Jarvis ticket",
		},
	}
	_, err6 := search.NewTicketRepository(badList3)
	require.NotNil(t, err6)

}

func Test_TicketRepository_FindByID(t *testing.T) {
	ticketList := []map[string]interface{}{
		{
			"_id":             "436bf9b0-1147-4c0a-8439-6f79833bff5b",
			"subject":         "Nicole Martinez ticket",
			"organization_id": float64(123),
			"assignee_id":     float64(1),
			"submitter_id":    float64(2),
		},
		{
			"_id":     "1a227508-9f39-427c-8f57-1b72f3fab87c",
			"subject": "Johnny Jarvis ticket",
		},
	}
	repository, err := search.NewTicketRepository(ticketList)
	require.Nil(t, err)

	result1 := repository.FindByID("436bf9b0-1147-4c0a-8439-6f79833bff5b")
	require.Equal(t, ticketList[0], result1)

	result2 := repository.FindByID("1a227508-9f39-427c-8f57-1b72f3fab87c")
	require.Equal(t, ticketList[1], result2)

	result3 := repository.FindByID("non existing")
	require.Nil(t, result3)
}

func Test_TicketRepository_FindByOrg(t *testing.T) {
	ticketList := []map[string]interface{}{
		{
			"_id":             "436bf9b0-1147-4c0a-8439-6f79833bff5b",
			"subject":         "Nicole Martinez ticket",
			"organization_id": float64(123),
		},
		{
			"_id":     "1a227508-9f39-427c-8f57-1b72f3fab87c",
			"subject": "Johnny Jarvis ticket",
		},
		{
			"_id":             "2217c7dc-7371-4401-8738-0a8a8aedc08d",
			"subject":         "Nicole Martinez ticket again",
			"organization_id": float64(123),
		},
	}
	repository, err := search.NewTicketRepository(ticketList)
	require.Nil(t, err)

	result1 := repository.FindByOrg(123)
	require.Len(t, result1, 2)
	require.Contains(t, result1, ticketList[0])
	require.Contains(t, result1, ticketList[2])

	result2 := repository.FindByOrg(0)
	require.Len(t, result2, 1)
	require.Contains(t, result2, ticketList[1])
}

func Test_TicketRepository_FindBySubmitter(t *testing.T) {
	ticketList := []map[string]interface{}{
		{
			"_id":          "436bf9b0-1147-4c0a-8439-6f79833bff5b",
			"subject":      "Nicole Martinez ticket",
			"submitter_id": float64(123),
		},
		{
			"_id":     "1a227508-9f39-427c-8f57-1b72f3fab87c",
			"subject": "Johnny Jarvis ticket",
		},
		{
			"_id":          "2217c7dc-7371-4401-8738-0a8a8aedc08d",
			"subject":      "Nicole Martinez ticket again",
			"submitter_id": float64(123),
		},
	}
	repository, err := search.NewTicketRepository(ticketList)
	require.Nil(t, err)

	result1 := repository.FindBySubmitter(123)
	require.Len(t, result1, 2)
	require.Contains(t, result1, ticketList[0])
	require.Contains(t, result1, ticketList[2])

	result2 := repository.FindBySubmitter(0)
	require.Len(t, result2, 1)
	require.Contains(t, result2, ticketList[1])
}

func Test_TicketRepository_FindByAssignee(t *testing.T) {
	ticketList := []map[string]interface{}{
		{
			"_id":         "436bf9b0-1147-4c0a-8439-6f79833bff5b",
			"subject":     "Nicole Martinez ticket",
			"assignee_id": float64(123),
		},
		{
			"_id":     "1a227508-9f39-427c-8f57-1b72f3fab87c",
			"subject": "Johnny Jarvis ticket",
		},
		{
			"_id":         "2217c7dc-7371-4401-8738-0a8a8aedc08d",
			"subject":     "Nicole Martinez ticket again",
			"assignee_id": float64(123),
		},
	}
	repository, err := search.NewTicketRepository(ticketList)
	require.Nil(t, err)

	result1 := repository.FindByAssignee(123)
	require.Len(t, result1, 2)
	require.Contains(t, result1, ticketList[0])
	require.Contains(t, result1, ticketList[2])

	result2 := repository.FindByAssignee(0)
	require.Len(t, result2, 1)
	require.Contains(t, result2, ticketList[1])
}

func Test_TicketRepository_FindByField(t *testing.T) {
	ticketList := []map[string]interface{}{
		{
			"_id":             "436bf9b0-1147-4c0a-8439-6f79833bff5b",
			"subject":         "Nicole Martinez ticket",
			"assignee_id":     float64(123),
			"submitter_id":    float64(2),
			"organization_id": float64(222),
		},
		{
			"_id":             "1a227508-9f39-427c-8f57-1b72f3fab87c",
			"subject":         "Johnny Jarvis ticket",
			"organization_id": float64(222),
		},
		{
			"_id":          "2217c7dc-7371-4401-8738-0a8a8aedc08d",
			"subject":      "Nicole Martinez ticket again",
			"assignee_id":  float64(123),
			"submitter_id": float64(2),
		},
	}
	repository, err := search.NewTicketRepository(ticketList)
	require.Nil(t, err)

	defaultMatcher := func(subject interface{}, term interface{}) bool {
		return false
	}

	repository.SetValueMatcher(defaultMatcher)

	tests := []struct {
		name           string
		field          string
		searchVal      interface{}
		expectedResult []map[string]interface{}
	}{
		{
			name:           "Match string to ticket ID (no value matcher)",
			field:          "_id",
			searchVal:      "436bf9b0-1147-4c0a-8439-6f79833bff5b",
			expectedResult: []map[string]interface{}{ticketList[0]},
		},
		{
			name:           "Match float to organization_id (no value matcher)",
			field:          "organization_id",
			searchVal:      float64(222),
			expectedResult: []map[string]interface{}{ticketList[0], ticketList[1]},
		},
		{
			name:           "Match string to organization_id (no value matcher)",
			field:          "organization_id",
			searchVal:      "222",
			expectedResult: []map[string]interface{}{ticketList[0], ticketList[1]},
		},
		{
			name:           "Match empty string to tickets with no org (no value matcher)",
			field:          "organization_id",
			searchVal:      "",
			expectedResult: []map[string]interface{}{ticketList[2]},
		},
		{
			name:           "Match nil to tickets with no org (no value matcher)",
			field:          "organization_id",
			searchVal:      nil,
			expectedResult: []map[string]interface{}{ticketList[2]},
		},
		{
			name:           "Match float to submitter_id (no value matcher)",
			field:          "submitter_id",
			searchVal:      float64(2),
			expectedResult: []map[string]interface{}{ticketList[0], ticketList[2]},
		},
		{
			name:           "Match string to submitter_id (no value matcher)",
			field:          "submitter_id",
			searchVal:      "2",
			expectedResult: []map[string]interface{}{ticketList[0], ticketList[2]},
		},
		{
			name:           "Match empty string to tickets with no submitter (no value matcher)",
			field:          "submitter_id",
			searchVal:      "",
			expectedResult: []map[string]interface{}{ticketList[1]},
		},
		{
			name:           "Match nil to tickets with no submitter (no value matcher)",
			field:          "submitter_id",
			searchVal:      nil,
			expectedResult: []map[string]interface{}{ticketList[1]},
		},

		{
			name:           "Match float to assignee_id (no value matcher)",
			field:          "assignee_id",
			searchVal:      float64(123),
			expectedResult: []map[string]interface{}{ticketList[0], ticketList[2]},
		},
		{
			name:           "Match string to assignee_id (no value matcher)",
			field:          "assignee_id",
			searchVal:      "123",
			expectedResult: []map[string]interface{}{ticketList[0], ticketList[2]},
		},
		{
			name:           "Match empty string to tickets with no assignee (no value matcher)",
			field:          "assignee_id",
			searchVal:      "",
			expectedResult: []map[string]interface{}{ticketList[1]},
		},
		{
			name:           "Match nil to tickets with no assignee (no value matcher)",
			field:          "assignee_id",
			searchVal:      nil,
			expectedResult: []map[string]interface{}{ticketList[1]},
		},
		{
			name:           "Match any other attribute using value matcher",
			field:          "subject",
			searchVal:      "test search value",
			expectedResult: []map[string]interface{}{ticketList[1], ticketList[2]},
		},
	}

	//This is where the actual testing happens
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.field == "organization_id" || tt.field == "_id" || tt.field == "assignee_id" || tt.field == "submitter_id" {
				repository.SetValueMatcher(defaultMatcher)
			} else {
				matcher := func(subject interface{}, term interface{}) bool {
					require.Equal(t, tt.searchVal, term)

					for _, expUser := range tt.expectedResult {
						if subject == expUser[tt.field] {
							return true
						}
					}
					return false
				}
				repository.SetValueMatcher(matcher)
			}
			require.ElementsMatch(t, tt.expectedResult, repository.FindByField(tt.field, tt.searchVal))
		})
	}

}
