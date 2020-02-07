package search

import (
	"github.com/pkg/errors"
)

//FYI this is why I use float64: https://golang.org/pkg/encoding/json/#Unmarshal
type UserRepository struct {
	usersIndex   map[float64]map[string]interface{} //map of json data indexed by user ID
	orgsIndex    map[float64][]float64              //map of user IDs indext by org ID
	valueMatcher ValueMatcher
}

func NewUserRepository(users []map[string]interface{}) (*UserRepository, error) {

	repository := &UserRepository{
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
func (repo *UserRepository) SetValueMatcher(matcherFn ValueMatcher) {
	repo.valueMatcher = matcherFn
}

func (repo *UserRepository) addUser(user map[string]interface{}) error {
	userID, isInt := user["_id"].(float64) //FYI: in go, if a map doesn't have a key, it simply returns nil
	if !isInt {
		return errors.New("User is missing \"_id\" field or \"_id\" is not float64")
	}

	if _, exists := repo.usersIndex[userID]; exists {
		return errors.Errorf("User with ID of %v already exists", userID)
	}

	repo.usersIndex[userID] = user

	if orgID, isInt := user["organization_id"].(float64); isInt {
		repo.orgsIndex[orgID] = append(repo.orgsIndex[orgID], userID)
	} else {
		repo.orgsIndex[0] = append(repo.orgsIndex[0], userID)
	}

	return nil
}

func (repo *UserRepository) FindByID(userID float64) map[string]interface{} {
	return repo.usersIndex[userID]
}

func (repo *UserRepository) FindByOrg(orgID float64) []map[string]interface{} {
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

func (repo *UserRepository) FindByField(fieldName string, searchVal interface{}) []map[string]interface{} {
	switch fieldName {
	case "_id":
		var userList []map[string]interface{}

		userID, isInt := floatVal(searchVal)
		if user := repo.FindByID(userID); isInt && user != nil {
			userList = append(userList, user)
		}

		return userList

	case "organization_id":

		if orgID, isInt := floatVal(searchVal); isInt {
			return repo.FindByOrg(orgID)
		}

		return []map[string]interface{}{}

	default:
		var userList []map[string]interface{}

		for _, user := range repo.usersIndex {
			if repo.valueMatcher(user[fieldName], searchVal) {
				userList = append(userList, user)
			}
		}
		return userList
	}
}
