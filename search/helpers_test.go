package search_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/superjinjo/zendesk-search/search"
)

func Test_SearchValueMatches(t *testing.T) {
	//This pattern is called table testing and it's helpful when there are a lot of scenarios to test
	tests := []struct {
		name           string
		subject        interface{}
		term           interface{}
		expectedResult bool
	}{
		//float subject----------------------------------------
		{
			name:           "Match float to a float",
			subject:        float64(2),
			term:           float64(2),
			expectedResult: true,
		},
		{
			name:           "match string to a float",
			subject:        float64(3),
			term:           "3",
			expectedResult: true,
		},
		{
			name:           "match empty string to zero float",
			subject:        float64(0),
			term:           "",
			expectedResult: true,
		},
		{
			name:           "match nil to zero float",
			subject:        float64(0),
			term:           nil,
			expectedResult: true,
		},
		{
			name:           "non matching float",
			subject:        float64(123),
			term:           float64(456),
			expectedResult: false,
		},
		{
			name:           "non float value",
			subject:        float64(12),
			term:           "twelve",
			expectedResult: false,
		},

		//bool subject-------------------------------------
		{
			name:           "match bool to bool",
			subject:        false,
			term:           false,
			expectedResult: true,
		},
		{
			name:           "match \"true\" string to true bool",
			subject:        true,
			term:           "true",
			expectedResult: true,
		},
		{
			name:           "match string 1 to true bool",
			subject:        true,
			term:           "1",
			expectedResult: true,
		},
		{
			name:           "match float 1 to true bool",
			subject:        true,
			term:           float64(1),
			expectedResult: true,
		},
		{
			name:           "match \"false\" string to false bool",
			subject:        false,
			term:           "false",
			expectedResult: true,
		},
		{
			name:           "match string 0 to false bool",
			subject:        false,
			term:           "0",
			expectedResult: true,
		},
		{
			name:           "match float 0 to false bool",
			subject:        false,
			term:           float64(0),
			expectedResult: true,
		},
		{
			name:           "match empty string false bool",
			subject:        false,
			term:           "",
			expectedResult: true,
		},
		{
			name:           "match nil false bool",
			subject:        false,
			term:           nil,
			expectedResult: true,
		},
		{
			name:           "non matching bool",
			subject:        false,
			term:           true,
			expectedResult: false,
		},

		//string subject---------------------------------------
		{
			name:           "match string to string",
			subject:        "Mr Caitlin",
			term:           "Mr Caitlin",
			expectedResult: true,
		},
		{
			name:           "match float to string",
			subject:        "1",
			term:           float64(1),
			expectedResult: true,
		},
		{
			name:           "match bool to string",
			subject:        "true",
			term:           true,
			expectedResult: true,
		},
		{
			name:           "match nil to string",
			subject:        "",
			term:           nil,
			expectedResult: true,
		},

		//slice subject--------------------------------------------------------
		{
			name:           "match one slice item",
			subject:        []interface{}{"Appleton", "Boston", "Chicago"},
			term:           []interface{}{"Appleton"},
			expectedResult: true,
		},
		{
			name:           "match two slice items not in order",
			subject:        []interface{}{"Appleton", "Boston", "Chicago"},
			term:           []interface{}{"Chicago", "Appleton"},
			expectedResult: true,
		},
		{
			name:           "match normal string to slice",
			subject:        []interface{}{"Appleton", "Boston", "Chicago"},
			term:           "Appleton",
			expectedResult: true,
		},
		{
			name:           "match comma deliniated string to slice",
			subject:        []interface{}{"Appleton", "Boston", "Chicago"},
			term:           "Boston,Appleton",
			expectedResult: true,
		},
		{
			name:           "match empty string to empty slice",
			subject:        []interface{}{},
			term:           "",
			expectedResult: true,
		},
		{
			name:           "match nil to empty slice",
			subject:        []interface{}{},
			term:           nil,
			expectedResult: true,
		},
		{
			name:           "one slice item matches and one doesn't",
			subject:        []interface{}{"Appleton", "Boston", "Chicago"},
			term:           []interface{}{"Appleton", "Westeros"},
			expectedResult: false,
		},
		{
			name:           "non matching with an empty slice",
			subject:        []interface{}{"Appleton", "Boston", "Chicago"},
			term:           []interface{}{},
			expectedResult: false,
		},
		{
			name:           "match empty slice to empty slice",
			subject:        []interface{}{},
			term:           []interface{}{},
			expectedResult: true,
		},
	}

	//This is where the actual testing happens
	for _, tt := range tests {
		tt := tt
		// this is not being run concurrently but good to keep in mind
		// https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedResult, search.SearchValueMatches(tt.subject, tt.term))
		})
	}
}
