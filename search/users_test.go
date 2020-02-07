package search_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/superjinjo/zendesk-search/search"
)

func Test_NewUserJSONRepository(t *testing.T) {
	//empty list is valid
	emptyList := []map[string]interface{}{}
	_, err1 := search.NewUserJSONRepository(emptyList)
	require.Nil(t, err1)

	//single item in list
	goodList1 := []map[string]interface{}{
		{
			"_id":             float64(1),
			"name":            "Nicole Martinez",
			"organization_id": float64(123),
		},
	}
	_, err2 := search.NewUserJSONRepository(goodList1)
	require.Nil(t, err2)

	//no organization_id is okay
	goodList2 := []map[string]interface{}{
		{
			"_id":             float64(1),
			"name":            "Nicole Martinez",
			"organization_id": float64(123),
		},
		{
			"_id":  float64(2),
			"name": "Johnny Jarvis",
		},
	}
	_, err3 := search.NewUserJSONRepository(goodList2)
	require.Nil(t, err3)

	//"_id" field is required
	badList1 := []map[string]interface{}{
		{
			"_id":             float64(1),
			"name":            "Nicole Martinez",
			"organization_id": float64(123),
		},
		{
			"name":            "Johnny Jarvis",
			"organization_id": float64(123),
		},
	}
	_, err4 := search.NewUserJSONRepository(badList1)
	require.NotNil(t, err4)

	//_id must be a float64
	badList2 := []map[string]interface{}{
		{
			"_id":             float64(1),
			"name":            "Nicole Martinez",
			"organization_id": float64(123),
		},
		{
			"id":              2,
			"name":            "Johnny Jarvis",
			"organization_id": float64(123),
		},
	}
	_, err5 := search.NewUserJSONRepository(badList2)
	require.NotNil(t, err5)

	//duplicate _id fields
	badList3 := []map[string]interface{}{
		{
			"_id":             float64(1),
			"name":            "Nicole Martinez",
			"organization_id": float64(123),
		},
		{
			"_id":             float64(1),
			"name":            "Johnny Jarvis",
			"organization_id": float64(123),
		},
	}
	_, err6 := search.NewUserJSONRepository(badList3)
	require.NotNil(t, err6)

}

func Test_UserJSONRepository_FindByID(t *testing.T) {
	userList := []map[string]interface{}{
		{
			"_id":             float64(1),
			"name":            "Nicole Martinez",
			"organization_id": float64(123),
		},
		{
			"_id":             float64(2),
			"name":            "Johnny Jarvis",
			"organization_id": float64(123),
		},
	}
	repository, err := search.NewUserJSONRepository(userList)
	require.Nil(t, err)

	result1 := repository.FindByID(float64(1))
	require.Equal(t, userList[0], result1)

	result2 := repository.FindByID(float64(2))
	require.Equal(t, userList[1], result2)

	result3 := repository.FindByID(float64(404))
	require.Nil(t, result3)
}

func Test_UserJSONRepository_FindByOrg(t *testing.T) {
	userList := []map[string]interface{}{
		{
			"_id":             float64(1),
			"name":            "Nicole Martinez",
			"organization_id": float64(123),
		},
		{
			"_id":             float64(2),
			"name":            "Johnny Jarvis",
			"organization_id": float64(123),
		},
		{
			"_id":             float64(3),
			"name":            "Tamara Tamarind",
			"organization_id": float64(456),
		},
		{
			"_id":  float64(4),
			"name": "Mikey NoOrgs",
		},
	}
	repository, err := search.NewUserJSONRepository(userList)
	require.Nil(t, err)

	result1 := repository.FindByOrg(123)
	require.Len(t, result1, 2)
	require.Contains(t, result1, userList[0])
	require.Contains(t, result1, userList[1])

	result2 := repository.FindByOrg(456)
	require.Len(t, result2, 1)
	require.Contains(t, result2, userList[2])

	result3 := repository.FindByOrg(789)
	require.Empty(t, result3)

	result4 := repository.FindByOrg(0)
	require.Len(t, result4, 1)
	require.Contains(t, result4, userList[3])
}

func Test_UserJSONRepository_FindByField(t *testing.T) {
	userList := []map[string]interface{}{
		{
			"_id":             float64(1),
			"name":            "Nicole Martinez",
			"organization_id": float64(123),
		},
		{
			"_id":             float64(2),
			"name":            "Johnny Jarvis",
			"alias":           "Mr Caitlin",
			"active":          true,
			"organization_id": float64(123),
		},
		{
			"_id":             float64(3),
			"name":            "Tamara Tamarind",
			"organization_id": float64(456),
		},
		{
			"_id":  float64(4),
			"name": "Mickey NoOrgs",
		},
	}
	repository, err := search.NewUserJSONRepository(userList)
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
			name:           "Match float to user ID (no value matcher)",
			field:          "_id",
			searchVal:      float64(3),
			expectedResult: []map[string]interface{}{userList[2]},
		},
		{
			name:           "Match string to user ID (no value matcher)",
			field:          "_id",
			searchVal:      "3",
			expectedResult: []map[string]interface{}{userList[2]},
		},
		{
			name:           "Match float to organization_id (no value matcher)",
			field:          "organization_id",
			searchVal:      float64(123),
			expectedResult: []map[string]interface{}{userList[0], userList[1]},
		},
		{
			name:           "Match string to organization_id (no value matcher)",
			field:          "organization_id",
			searchVal:      "123",
			expectedResult: []map[string]interface{}{userList[0], userList[1]},
		},
		{
			name:           "Match empty string to users with no org (no value matcher)",
			field:          "organization_id",
			searchVal:      "",
			expectedResult: []map[string]interface{}{userList[3]},
		},
		{
			name:           "Match nil to users with no org (no value matcher)",
			field:          "organization_id",
			searchVal:      nil,
			expectedResult: []map[string]interface{}{userList[3]},
		},
		{
			name:           "Match any other attribute using value matcher",
			field:          "name",
			searchVal:      "test search value",
			expectedResult: []map[string]interface{}{userList[1], userList[3]},
		},
	}

	//This is where the actual testing happens
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.field == "organization_id" || tt.field == "_id" {
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
