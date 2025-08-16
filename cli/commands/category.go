package commands

import (
	"fmt"
	"loopit/cli/initializer"
	"loopit/cli/utils"
	"loopit/internal/config"
	"os"

	"github.com/olekukonko/tablewriter"
)

func CreateCategory() {
	fmt.Println("\nCreate Category")
	name := utils.Input("Enter Category Name")
	price := utils.FloatConversion(utils.Input("Enter Price"))
	security := utils.FloatConversion(utils.Input("Enter Security Deposit"))

	err := initializer.CategoryService.CreateCategory(name, price, security)
	if err != nil {
		fmt.Println(config.Red+"Failed to create category:"+config.Reset, err)
		return
	}

	fmt.Println(config.Green + "Category created successfully!" + config.Reset)
}

func GetAllCategories() {
	categories, err := initializer.CategoryService.GetAllCategories()
	if err != nil {
		fmt.Println(config.Red+"Failed to fetch categories:"+config.Reset, err)
		return
	}
	if len(categories) == 0 {
		fmt.Println(config.Yellow + "No categories available." + config.Reset)
		return
	}
	fmt.Println(config.Green + "Available Categories:" + config.Reset)

	table := tablewriter.NewWriter(os.Stdout)
	table.Header("ID", "Name", "Price", "Security")

	for _, category := range categories {
		table.Append([]string{
			fmt.Sprintf("%d", category.ID),
			category.Name,
			fmt.Sprintf("%.2f", category.Price),
			fmt.Sprintf("%.2f", category.Security),
		})
	}
	table.Bulk(true)
	table.Render()
}
