package commands

import (
	"fmt"
	"loopit/cli/utils"
	"loopit/internal/config"
	"loopit/internal/models"
	"os"

	"github.com/olekukonko/tablewriter"
)

func GetAllSocieties() {
	societies, err := SocietyService.GetAllSocieties()
	if err != nil {
		fmt.Println(config.Red + "Error fetching societies: " + err.Error() + config.Reset)
		return
	}

	if len(societies) == 0 {
		fmt.Println(config.Yellow + "No societies found." + config.Reset)
		return
	}

	fmt.Println(config.Green + "List of Societies:" + config.Reset)
	table := tablewriter.NewWriter(os.Stdout)
	table.Header("ID", "Name", "Location", "Pincode")

	for _, society := range societies {
		table.Append([]string{
			fmt.Sprintf("%d", society.ID),
			society.Name,
			society.Location,
			society.Pincode,
		})
	}

	table.Bulk(true)
	table.Render()
}

func CreateSociety(userCtx *models.UserContext) {
	name := utils.Input("Enter Society Name: ")
	location := utils.Input("Enter Society Location: ")
	pincode := utils.Input("Enter Society Pincode: ")

	if err := SocietyService.CreateSociety(name, location, pincode); err != nil {
		fmt.Println(config.Red + "Error creating society: " + err.Error() + config.Reset)
		return
	}

	fmt.Println("Society created successfully !")
}
