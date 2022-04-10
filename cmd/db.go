package cmd

import (
	"errors"
	"fmt"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/internal"
	"github.com/spf13/cobra"
)

var adminCmd = &cobra.Command{
	Use:   "db",
	Short: "whispering db control",
	Long: `TODO: A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Example: `db migrate
			  db create
              db drop
			  db status
			  db import
`,
	ValidArgs: dbSubCmds,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a db argument")
		}
		if isDbCmd(args[0]) {
			return nil
		}
		return fmt.Errorf("invalid db cmd specified: %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		tables := []interface{}{domain.User{}, domain.SocialAccount{}, domain.Cocktail{}, domain.CocktailIngredient{},
			domain.CocktailStep{}, domain.FavoriteCocktail{}}

		subcmd := args[0]
		internal.Migrate(cfgFile, subcmd, tables)
	},
}

func init() {
	rootCmd.AddCommand(adminCmd)

	// Here you will define your flags and configuration settings.
	//dbCmd.PersistentFlags().BoolVar(&isTestEnv, "test", false, "is test env or not")
	//adminCmd.PersistentFlags().StringVar(&sqlInsert, "sql", "serviceConfig.json", "insert sql statement")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	adminCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	adminCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
