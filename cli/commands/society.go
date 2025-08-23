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

func GetAllSocieties() {
	log.Info("CLI: Fetching all societies")
	societies, err := initializer.SocietyService.GetAllSocieties()
	if err != nil {
		log.Error(fmt.Sprintf("CLI: Error fetching societies: %v", err))
		fmt.Println(config.Red + "Error fetching societies: " + err.Error() + config.Reset)
		return
	}

	if len(societies) == 0 {
		log.Info("CLI: No societies found")
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
	log.Info(fmt.Sprintf("CLI: User %d creating society", userCtx.ID))
	name := utils.Input("Enter Society Name: ")
	location := utils.Input("Enter Society Location: ")
	pincode := utils.Input("Enter Society Pincode: ")

	if err := initializer.SocietyService.CreateSociety(name, location, pincode); err != nil {
		log.Error(fmt.Sprintf("CLI: Error creating society: %v", err))
		fmt.Println(config.Red + "Error creating society: " + err.Error() + config.Reset)
		return
	}

	log.Info(fmt.Sprintf("CLI: Society '%s' created successfully", name))
	fmt.Println("Society created successfully !")
}
