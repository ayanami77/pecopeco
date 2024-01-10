package cmd

import (
	"errors"
	"fmt"
	"os"
	"unicode/utf8"

	foodFactory "github.com/Seiya-Tagami/pecopeco-cli/api/factory/food"
	"github.com/Seiya-Tagami/pecopeco-cli/api/repository/food"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run pecopeco CLI",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		selectOption()
	},
}

func selectOption() {
	prompt := promptui.Select{
		Label: "What would you like to do?",
		Items: []string{"Search food", "Show favorites", "Exit"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed%v\n", err)
		return
	}

	switch result {
	case "Search food":
		factory := foodFactory.CreateFactory()

		searchFoodInput := getSearchFoodInput()
		request := food.GetFoodListRequest{
			City: searchFoodInput.city,
			Food: searchFoodInput.food,
		}

		foodList, err := factory.GetFoodList(request)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, v := range foodList {
			fmt.Printf("-----------------\n店名: %s\n住所: %s\n最寄り駅: %s\nジャンル: %s\nURL: %s\n", v.Name, v.Address, v.StationName, v.GenreName, v.URL)
		}

		selectOption()
	case "Show favorites":
		fmt.Printf("%s called\n", result)
	case "Exit":
		fmt.Print("Bye!\n")
		os.Exit(1)
	}
}

type searchFoodInput struct {
	city string
	food string
}

func getSearchFoodInput() searchFoodInput {
	promptForCity := promptui.Prompt{
		Label: "Which city?",
		Validate: func(input string) error {
			if utf8.RuneCountInString(input) == 0 {
				return errors.New("Please enter a city.")
			}
			return nil
		},
	}
	city, err := promptForCity.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return searchFoodInput{}
	}

	promptForFood := promptui.Prompt{
		Label: "What food?",
		Validate: func(input string) error {
			if utf8.RuneCountInString(input) == 0 {
				return errors.New("Please enter food.")
			}
			return nil
		},
	}
	food, err := promptForFood.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return searchFoodInput{}
	}

	return searchFoodInput{city, food}
}

func init() {
	rootCmd.AddCommand(runCmd)
}
