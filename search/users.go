package search

import (
	"github.com/pkg/errors"
)

//FYI this is why I use float64: https://golang.org/pkg/encoding/json/#Unmarshal
type UserJSONRepository struct {
	usersIndex   map[float64]map[string]interface{} //map of json data indexed by user ID
	orgsIndex    map[float64][]float64              //map of user IDs indext by org ID
	valueMatcher ValueMatcher
}

func NewUserJSONRepository(users []map[string]interface{}) (*UserJSONRepository, error) {

	repository := &UserJSONRepository{
		usersIndex:   make(map[float64]map[string]interface{}),
		orgsIndex:    make(map[float64][]float64),
		valueMatcher: SearchValueMatches,
	}

	for i, user := range users {
		if err := repository.addUser(user); err != nil {
			return nil, errors.WithMessagef(err, "Error with user at index %d", i)
		}
	}

	return repository, nil
}

//SetValueMatcher lets you set a different matcher which is useful for testing
func (repo *UserJSONRepository) SetValueMatcher(matcherFn ValueMatcher) {
	repo.valueMatcher = matcherFn
}

func (repo *UserJSONRepository) addUser(user map[string]interface{}) error {
	userID, isFloat := user["_id"].(float64) //FYI: in go, if a map doesn't have a key, it simply returns nil
	if !isFloat {
		return errors.New("User is missing \"_id\" field or \"_id\" is not float64")
	}

	if _, exists := repo.usersIndex[userID]; exists {
		return errors.Errorf("User with ID of %v already exists", userID)
	}

	repo.usersIndex[userID] = user

	if orgID, isFloat := user["organization_id"].(float64); isFloat {
		repo.orgsIndex[orgID] = append(repo.orgsIndex[orgID], userID)
	} else {
		repo.orgsIndex[0] = append(repo.orgsIndex[0], userID)
	}

	return nil
}

func (repo *UserJSONRepository) FindByID(userID float64) map[string]interface{} {
	return repo.usersIndex[userID]
}

func (repo *UserJSONRepository) FindByOrg(orgID float64) []map[string]interface{} {
	userIDs := repo.orgsIndex[orgID]

	userList := make([]map[string]interface{}, len(userIDs))

	for i := 0; i < len(userIDs); i++ {
		nextID := userIDs[i]

		if user := repo.FindByID(nextID); user != nil {
			userList[i] = user
		}
	}

	return userList
}

func (repo *UserJSONRepository) FindByField(fieldName string, searchVal interface{}) []map[string]interface{} {
	switch fieldName {
	case "_id":
		userList := []map[string]interface{}{}

		userID, isFloat := floatVal(searchVal)
		if user := repo.FindByID(userID); isFloat && user != nil {
			userList = append(userList, user)
		}

		return userList

	case "organization_id":

		if orgID, isFloat := floatVal(searchVal); isFloat {
			return repo.FindByOrg(orgID)
		}

		return []map[string]interface{}{}

	default:
		userList := []map[string]interface{}{}

		for _, user := range repo.usersIndex {
			if repo.valueMatcher(user[fieldName], searchVal) {
				userList = append(userList, user)
			}
		}
		return userList
	}
}
