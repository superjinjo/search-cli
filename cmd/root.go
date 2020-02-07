package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/superjinjo/zendesk-search/search"
)

type BaseCommand struct {
	cobra.Command
	repository *search.SearchRepository

	Formatter func([]map[string]interface{}) string
}

func newBaseCommand(command cobra.Command, repo *search.SearchRepository) BaseCommand {
	return BaseCommand{
		Command:    command,
		repository: repo,
	}
}

func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "zensearch",
		Short: "zensearch allows you to search users, organizations, and tickets",
		Long:  `zensearch allows you to search users, organizations, and tickets by any attribute and will provide any related data`,
	}
}

func Execute(repository *search.SearchRepository) {

	rootCmd := NewRootCmd()

	usersCmd := NewUsersCommand(repository)
	orgsCmd := NewOrganizationsCommand(repository)
	ticketsCmd := NewTicketsCommand(repository)

	rootCmd.AddCommand(usersCmd, orgsCmd, ticketsCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
