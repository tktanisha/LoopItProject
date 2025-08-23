package commands

import (
	"fmt"
	"loopit/cli/initializer"
	"loopit/cli/utils"
	"loopit/internal/config"
	"loopit/internal/models"
	"os"

	"github.com/olekukonko/tablewriter"
)

func GetAllProducts() {
	log.Info("CLI: Fetching all products")
	products, err := initializer.ProductService.GetAllProducts()
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Error fetching products: %v", err))
		fmt.Println(config.Red+"Error fetching products:"+config.Reset, err)
		return
	}
	if len(products) == 0 {
		log.Info("CLI: No products found")
		fmt.Println(config.Yellow + "No products found." + config.Reset)
		return
	}
	fmt.Println("\nAvailable Products:")
	table := tablewriter.NewWriter(os.Stdout)
	table.Header(
		"ID", "Name", "Category", "Price (₹)", "Security (₹)",
		"Lender Name", "Duration (days)", "Available",
	)
	for _, p := range products {
		table.Append([]string{
			fmt.Sprintf("%d", p.Product.ID),
			p.Product.Name,
			p.Category.Name,
			fmt.Sprintf("%.2f", p.Category.Price),
			fmt.Sprintf("%.2f", p.Category.Security),
			p.User.FullName,
			fmt.Sprintf("%d", p.Product.Duration),
			fmt.Sprintf("%t", p.Product.IsAvailable),
		})
	}
	table.Bulk(true)
	table.Render()
}

func GetProductByID() {
	id := utils.IntConversion(utils.Input("Enter Product ID: "))
	log.Info(fmt.Sprintf("CLI: Fetching product by ID %d", id))
	product, err := initializer.ProductService.GetProductByID(id)
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Error fetching product by ID %d: %v", id, err))
		fmt.Println(config.Red+"Error fetching product:"+config.Reset, err)
		return
	}
	fmt.Println("\nProduct Details:")
	table := tablewriter.NewWriter(os.Stdout)
	table.Header(
		"ID", "Name", "Category", "Price (₹)", "Security (₹)",
		"Lender Name", "Duration (days)", "Available",
	)
	table.Append([]string{
		fmt.Sprintf("%d", product.Product.ID),
		product.Product.Name,
		product.Category.Name,
		fmt.Sprintf("%.2f", product.Category.Price),
		fmt.Sprintf("%.2f", product.Category.Security),
		product.User.FullName,
		fmt.Sprintf("%d", product.Product.Duration),
		fmt.Sprintf("%t", product.Product.IsAvailable),
	})
	table.Bulk(true)
	table.Render()
}

func CreateProduct(userCtx *models.UserContext) {
	name := utils.Input("Enter Product Name: ")
	description := utils.Input("Enter Product Description: ")
	duration := utils.IntConversion(utils.Input("Enter Product Duration (e.g., 30 for 30 days): "))

	log.Info(fmt.Sprintf("CLI: User %d creating product with name: %s", userCtx.ID, name))
	categories, err := initializer.CategoryService.GetAllCategories()
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Error fetching categories: %v", err))
		fmt.Println(config.Red+"Error fetching categories:"+config.Reset, err)
		return
	}
	categoryOptions := []string{}
	for _, cat := range categories {
		categoryOptions = append(categoryOptions, fmt.Sprintf("%s (₹%.2f)", cat.Name, cat.Price))
	}
	index, _ := utils.SelectFromList("Select Product Category", categoryOptions)
	if index == -1 {
		log.Warning("CLI: Category selection failed during product creation")
		fmt.Println("Category selection failed.")
		return
	}
	selectedCategory := categories[index]
	product := &models.Product{
		Name:        name,
		Description: description,
		CategoryID:  selectedCategory.ID,
		Duration:    duration,
		IsAvailable: true,
	}
	err = initializer.ProductService.CreateProduct(product, userCtx)
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Error creating product '%s': %v", name, err))
		fmt.Println(config.Red+"Error creating product:"+config.Reset, err)
		return
	}
	log.Info(fmt.Sprintf("CLI: Product '%s' created successfully by user %d", name, userCtx.ID))
	fmt.Println(config.Green + "Product created successfully: " + product.Name + config.Reset)
}
