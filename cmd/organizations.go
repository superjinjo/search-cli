package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/superjinjo/zendesk-search/search"
)

var organizationFields = []string{
	"_id",
	"url",
	"external_id",
	"name",
	"domain_names",
	"created_at",
	"details",
	"shared_tickets",
	"tags",
}

type OrganizationSearchCommand struct {
	cobra      *cobra.Command
	repository *search.SearchRepository

	Formatter func([]map[string]interface{}) (string, error)
}

func NewOrganizationSearchCommand(repo *search.SearchRepository) *OrganizationSearchCommand {
	organizationCmd := &OrganizationSearchCommand{
		Formatter:  formatJSONOutput,
		repository: repo,
	}

	command := &cobra.Command{
		Use:   "search [field] [search term]",
		Short: "search zendesk organizations by field.",
		Long:  `search zendesk organizations by field. If the search term is omitted, it will return all organizations that have the chosen field empty.`,
		Args: func(command *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a field argument")
			}

			if len(args) > 2 {
				return errors.New("too many arguments")
			}

			if !stringInSlice(args[0], organizationFields) {
				return fmt.Errorf(`Invalid field "%v"`, args[0])
			}

			return nil
		},
		RunE: organizationCmd.RunCommand,
	}

	organizationCmd.cobra = command

	return organizationCmd
}

func (oc *OrganizationSearchCommand) RunCommand(command *cobra.Command, args []string) error {
	var fieldName = args[0]
	var searchTerm string

	if len(args) > 1 {
		searchTerm = args[1]
	}

	searchResults := oc.repository.FindOrgs(fieldName, searchTerm)

	formattedResults, err := oc.Formatter(searchResults)
	if err != nil {
		return err
	}

	fmt.Println(formattedResults)
	return nil
}

func NewOrganizationsCommand(repo *search.SearchRepository) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "organizations",
		Short: "zendesk organizations operations",
		Long:  `zendesk organizations operations`,
	}

	fieldsCmd := &cobra.Command{
		Use:   "fields",
		Short: "list valid organization fields to search by",
		Long:  `list valid organization fields to search by`,
		Run: func(command *cobra.Command, args []string) {
			for _, fieldName := range organizationFields {
				fmt.Println(fieldName)
			}
		},
	}

	searchCmd := NewOrganizationSearchCommand(repo)

	rootCmd.AddCommand(fieldsCmd, searchCmd.cobra)

	return rootCmd
}
