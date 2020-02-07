package search_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/superjinjo/zendesk-search/search"
)

func Test_NewOrgRepository(t *testing.T) {
	//empty list is valid
	emptyList := []map[string]interface{}{}
	_, err1 := search.NewOrgRepository(emptyList)
	require.Nil(t, err1)

	//happy path
	goodList1 := []map[string]interface{}{
		{
			"_id":  float64(1),
			"name": "Umbrella Corp.",
		},
		{
			"_id":  float64(2),
			"name": "Aperture Science",
		},
	}
	_, err3 := search.NewOrgRepository(goodList1)
	require.Nil(t, err3)

	//"_id" field is required
	badList1 := []map[string]interface{}{
		{
			"_id":  float64(1),
			"name": "Umbrella Corp.",
		},
		{
			"name": "Aperture Science",
		},
	}
	_, err4 := search.NewOrgRepository(badList1)
	require.NotNil(t, err4)

	//_id must be a float64
	badList2 := []map[string]interface{}{
		{
			"_id":  float64(1),
			"name": "Umbrella Corp.",
		},
		{
			"_id":  2,
			"name": "Aperture Science",
		},
	}
	_, err5 := search.NewOrgRepository(badList2)
	require.NotNil(t, err5)

	//duplicate _id fields
	badList3 := []map[string]interface{}{
		{
			"_id":  float64(1),
			"name": "Umbrella Corp.",
		},
		{
			"_id":  float64(1),
			"name": "Aperture Science",
		},
	}
	_, err6 := search.NewOrgRepository(badList3)
	require.NotNil(t, err6)

}

func Test_OrgRepository_FindByID(t *testing.T) {
	orgList := []map[string]interface{}{
		{
			"_id":  float64(1),
			"name": "Umbrella Corp.",
		},
		{
			"_id":  float64(2),
			"name": "Aperture Science",
		},
	}
	repository, err := search.NewOrgRepository(orgList)
	require.Nil(t, err)

	result1 := repository.FindByID(float64(1))
	require.Equal(t, orgList[0], result1)

	result2 := repository.FindByID(float64(2))
	require.Equal(t, orgList[1], result2)

	result3 := repository.FindByID(float64(404))
	require.Nil(t, result3)
}

func Test_OrgRepository_FindByField(t *testing.T) {
	orgList := []map[string]interface{}{
		{
			"_id":  float64(1),
			"name": "Umbrella Corp.",
		},
		{
			"_id":  float64(2),
			"name": "Aperture Science",
		},
		{
			"_id":  float64(4),
			"name": "Maliwan",
		},
		{
			"_id":  float64(5),
			"name": "Atlas Corperation",
		},
	}
	repository, err := search.NewOrgRepository(orgList)
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
			name:           "Match float to org ID (no value matcher)",
			field:          "_id",
			searchVal:      float64(2),
			expectedResult: []map[string]interface{}{orgList[1]},
		},
		{
			name:           "Match string to org ID (no value matcher)",
			field:          "_id",
			searchVal:      "2",
			expectedResult: []map[string]interface{}{orgList[1]},
		},
		{
			name:           "Match any other attribute using value matcher",
			field:          "name",
			searchVal:      "test search value",
			expectedResult: []map[string]interface{}{orgList[2], orgList[3]},
		},
	}

	//This is where the actual testing happens
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.field == "_id" {
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
