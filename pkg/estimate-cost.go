package pkg

import (
	"fmt"
	"github.com/mariusmatioc/infractl/pkg/global"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func EstimateCost(craftPath string) (err error) {
	if craftPath != "" {
		err = global.SetOsEnvsFromCraft(craftPath)
		if err != nil {
			return
		}
	}
	// Create a new AWS session
	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		return fmt.Errorf("AWS_DEFAULT_REGION not found")
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		return
	}

	// Create a CostExplorer client
	svc := costexplorer.New(sess)
	// Specify the time range for the cost data
	now := time.Now()
	year := now.Year()
	month := now.Month()
	dateFormat := "%04d-%02d-01"
	start := fmt.Sprintf(dateFormat, year, month)
	if month == 12 {
		month = 1
		year++
	} else {
		month++
	}
	end := fmt.Sprintf(dateFormat, year, month)

	// Specify the granularity (DAILY, MONTHLY, etc.)
	granularity := "MONTHLY"

	// Specify the metrics you want to retrieve
	metrics := []string{"BlendedCost", "UsageQuantity"}

	// Create the input parameters
	params := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: &start,
			End:   &end,
		},
		Granularity: &granularity,
		Metrics:     aws.StringSlice(metrics),
	}

	// Call the GetCostAndUsage API
	result, err := svc.GetCostAndUsage(params)
	if err != nil {
		return fmt.Errorf("error getting cost and usage: %v", err)
	}

	// Print the estimated costs
	fmt.Println("Estimated Costs:")
	for _, resultByTime := range result.ResultsByTime {
		fmt.Println(resultByTime.String())
		//fmt.Printf("Period: %s - %s\n", *resultByTime.TimePeriod.Start, *resultByTime.TimePeriod.End)
		//fmt.Printf("Blended Cost: $%s\n", *resultByTime.Total.BlendedCost.Amount)
		//fmt.Printf("Usage Quantity: %s\n", *resultByTime.Total.UsageQuantity)
		//fmt.Println("------")
	}
	return
}
