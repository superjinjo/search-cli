package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/superjinjo/zendesk-search/search"
)

var userFields = []string{
	"_id",
	"url",
	"external_id",
	"name",
	"alias",
	"created_at",
	"active",
	"verified",
	"shared",
	"locale",
	"timezone",
	"last_login_at",
	"email",
	"phone",
	"signature",
	"organization_id",
	"tags",
	"suspended",
	"role",
}

type UserSearchCommand struct {
	cobra      *cobra.Command
	repository *search.SearchRepository

	Formatter func([]map[string]interface{}) (string, error)
}

func NewUserSearchCommand(repo *search.SearchRepository) *UserSearchCommand {
	userCmd := &UserSearchCommand{
		Formatter:  formatJSONOutput,
		repository: repo,
	}

	command := &cobra.Command{
		Use:   "search [field] [search term]",
		Short: "search zendesk users by field.",
		Long:  `search zendesk users by field. If the search term is omitted, it will return all users that have the chosen field empty.`,
		Args: func(command *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a field argument")
			}

			if len(args) > 2 {
				return errors.New("too many arguments")
			}

			if !stringInSlice(args[0], userFields) {
				return fmt.Errorf(`Invalid field "%v"`, args[0])
			}

			return nil
		},
		RunE: userCmd.RunCommand,
	}

	userCmd.cobra = command

	return userCmd
}

func (uc *UserSearchCommand) RunCommand(command *cobra.Command, args []string) error {
	var fieldName = args[0]
	var searchTerm string

	if len(args) > 1 {
		searchTerm = args[1]
	}

	searchResults := uc.repository.FindUsers(fieldName, searchTerm)

	formattedResults, err := uc.Formatter(searchResults)
	if err != nil {
		return err
	}

	fmt.Println(formattedResults)
	return nil
}

func NewUsersCommand(repo *search.SearchRepository) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "users",
		Short: "zendesk users operations",
		Long:  `zendesk users operations`,
	}

	fieldsCmd := &cobra.Command{
		Use:   "fields",
		Short: "list valid user fields to search by",
		Long:  `list valid user fields to search by`,
		Run: func(command *cobra.Command, args []string) {
			for _, fieldName := range userFields {
				fmt.Println(fieldName)
			}
		},
	}

	searchCmd := NewUserSearchCommand(repo)

	rootCmd.AddCommand(fieldsCmd, searchCmd.cobra)

	return rootCmd
}
