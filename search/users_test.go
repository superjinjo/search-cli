package search_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/superjinjo/zendesk-search/search"
)

func Test_NewUserRepository(t *testing.T) {
	//empty list is valid
	emptyList := []map[string]interface{}{}
	_, err1 := search.NewUserRepository(emptyList)
	require.Nil(t, err1)

	//single item in list
	goodList1 := []map[string]interface{}{
		{
			"_id":             float64(1),
			"name":            "Nicole Martinez",
			"organization_id": float64(123),
		},
	}
	_, err2 := search.NewUserRepository(goodList1)
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
	_, err3 := search.NewUserRepository(goodList2)
	require.Nil(t, err3)

	//"_id" field is required
	badList := []map[string]interface{}{
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
	_, err4 := search.NewUserRepository(badList)
	require.NotNil(t, err4)

}
