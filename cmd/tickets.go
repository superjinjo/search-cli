package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/superjinjo/zendesk-search/search"
)

var ticketFields = []string{
	"_id",
	"url",
	"external_id",
	"created_at",
	"type",
	"subject",
	"description",
	"priority",
	"status",
	"submitter_id",
	"assignee_id",
	"organization_id",
	"tags",
	"has_incidents",
	"due_at",
	"via",
}

type TicketSearchCommand struct {
	cobra      *cobra.Command
	repository *search.SearchRepository

	Formatter func([]map[string]interface{}) (string, error)
}

func NewTicketSearchCommand(repo *search.SearchRepository) *TicketSearchCommand {
	ticketCmd := &TicketSearchCommand{
		Formatter:  formatJSONOutput,
		repository: repo,
	}

	command := &cobra.Command{
		Use:   "search [field] [search term]",
		Short: "search zendesk tickets by field.",
		Long:  `search zendesk tickets by field. If the search term is omitted, it will return all tickets that have the chosen field empty.`,
		Args: func(command *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a field argument")
			}

			if len(args) > 2 {
				return errors.New("too many arguments")
			}

			if !stringInSlice(args[0], ticketFields) {
				return fmt.Errorf(`Invalid field "%v"`, args[0])
			}

			return nil
		},
		RunE: ticketCmd.RunCommand,
	}

	ticketCmd.cobra = command

	return ticketCmd
}

func (tc *TicketSearchCommand) RunCommand(command *cobra.Command, args []string) error {
	var fieldName = args[0]
	var searchTerm string

	if len(args) > 1 {
		searchTerm = args[1]
	}

	searchResults := tc.repository.FindTickets(fieldName, searchTerm)

	formattedResults, err := tc.Formatter(searchResults)
	if err != nil {
		return err
	}

	fmt.Println(formattedResults)
	return nil
}

func NewTicketsCommand(repo *search.SearchRepository) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tickets",
		Short: "zendesk tickets operations",
		Long:  `zendesk tickets operations`,
	}

	fieldsCmd := &cobra.Command{
		Use:   "fields",
		Short: "list valid ticket fields to search by",
		Long:  `list valid ticket fields to search by`,
		Run: func(command *cobra.Command, args []string) {
			for _, fieldName := range ticketFields {
				fmt.Println(fieldName)
			}
		},
	}

	searchCmd := NewTicketSearchCommand(repo)

	rootCmd.AddCommand(fieldsCmd, searchCmd.cobra)

	return rootCmd
}
