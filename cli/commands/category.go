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

	log.Info(fmt.Sprintf("CLI: Attempting to create category '%s' with price %.2f and security %.2f", name, price, security))
	err := initializer.CategoryService.CreateCategory(name, price, security)
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Failed to create category '%s': %v", name, err))
		fmt.Println(config.Red+"Failed to create category:"+config.Reset, err)
		return
	}

	log.Info(fmt.Sprintf("CLI: Successfully created category '%s'", name))
	fmt.Println(config.Green + "Category created successfully!" + config.Reset)
}

func GetAllCategories() {
	log.Info("CLI: Fetching all categories")
	categories, err := initializer.CategoryService.GetAllCategories()
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Error fetching categories: %v", err))
		fmt.Println(config.Red+"Failed to fetch categories:"+config.Reset, err)
		return
	}
	if len(categories) == 0 {
		log.Info("CLI: No categories available")
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
