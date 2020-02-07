package search

import (
	"github.com/pkg/errors"
)

//FYI this is why I use float64: https://golang.org/pkg/encoding/json/#Unmarshal
type UserRepository struct {
	usersIndex map[float64]map[string]interface{} //map of json data indexed by user ID
	orgsIndex  map[float64][]float64              //map of user IDs indext by org ID
}

func NewUserRepository(users []map[string]interface{}) (*UserRepository, error) {

	repository := &UserRepository{
		usersIndex: make(map[float64]map[string]interface{}),
		orgsIndex:  make(map[float64][]float64),
	}

	for i, user := range users {
		if err := repository.AddUser(user); err != nil {
			return nil, errors.WithMessagef(err, "Error with user at index %d", i)
		}
	}

	return repository, nil
}

func (repo *UserRepository) AddUser(user map[string]interface{}) error {
	userID, isInt := user["_id"].(float64) //FYI: in go, if a map doesn't have a key, it simply returns nil
	if !isInt {
		return errors.New("User is missing \"_id\" field")
	}

	if _, exists := repo.usersIndex[userID]; exists {
		return errors.Errorf("User with ID of %v already exists", userID)
	}

	repo.usersIndex[userID] = user

	if orgID, isInt := user["organization_id"].(float64); isInt {
		repo.orgsIndex[orgID] = append(repo.orgsIndex[orgID], userID)
	}

	return nil
}

func (repo *UserRepository) FindByID(userID float64) (map[string]interface{}, bool) {
	user, exists := repo.usersIndex[userID]
	return user, exists
}

func (repo *UserRepository) FindByOrg(orgID float64) []map[string]interface{} {
	userIDs := repo.orgsIndex[orgID]

	userList := make([]map[string]interface{}, len(userIDs))

	for i := 0; i < len(userIDs); i++ {
		nextID := userIDs[i]

		if user, exists := repo.FindByID(nextID); exists {
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
		if user, exists := repo.FindByID(userID); isInt && exists {
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
			if valuesMatch(user[fieldName], searchVal) {
				userList = append(userList, user)
			}
		}
		return userList
	}
}
