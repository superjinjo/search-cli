package search

import (
	"github.com/pkg/errors"
)

//FYI this is why I use float64: https://golang.org/pkg/encoding/json/#Unmarshal
type OrgJSONRepository struct {
	orgsIndex    map[float64]map[string]interface{} //map of json data indexed by org ID
	valueMatcher ValueMatcher
}

func NewOrgJSONRepository(orgs []map[string]interface{}) (*OrgJSONRepository, error) {

	repository := &OrgJSONRepository{
		orgsIndex:    make(map[float64]map[string]interface{}),
		valueMatcher: SearchValueMatches,
	}

	for i, org := range orgs {
		if err := repository.addOrg(org); err != nil {
			return nil, errors.WithMessagef(err, "Error with org at index %d", i)
		}
	}

	return repository, nil
}

//SetValueMatcher lets you set a different matcher which is useful for testing
func (repo *OrgJSONRepository) SetValueMatcher(matcherFn ValueMatcher) {
	repo.valueMatcher = matcherFn
}

func (repo *OrgJSONRepository) addOrg(org map[string]interface{}) error {
	orgID, isFloat := org["_id"].(float64) //FYI: in go, if a map doesn't have a key, it simply returns nil
	if !isFloat {
		return errors.New("Org is missing \"_id\" field or \"_id\" is not float64")
	}

	if _, exists := repo.orgsIndex[orgID]; exists {
		return errors.Errorf("Org with ID of %v already exists", orgID)
	}

	repo.orgsIndex[orgID] = org
	return nil
}

func (repo *OrgJSONRepository) FindByID(orgID float64) map[string]interface{} {
	return repo.orgsIndex[orgID]
}

func (repo *OrgJSONRepository) FindByField(fieldName string, searchVal interface{}) []map[string]interface{} {
	switch fieldName {
	case "_id":
		var orgList []map[string]interface{}

		orgID, isFloat := floatVal(searchVal)
		if org := repo.FindByID(orgID); isFloat && org != nil {
			orgList = append(orgList, org)
		}

		return orgList

	default:
		var orgList []map[string]interface{}

		for _, org := range repo.orgsIndex {
			if repo.valueMatcher(org[fieldName], searchVal) {
				orgList = append(orgList, org)
			}
		}
		return orgList
	}
}
