package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/datafactory/mgmt/datafactory"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/date"

	cfg "github.com/chelnak/godf/internal/config"
)

type PipelineRun struct {
	Name       string
	Status     string
	Start      *date.Time
	End        *date.Time
	DurationMs *int32
}

func getPipelineRunsClient() datafactory.PipelineRunsClient {

	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		panic("Could not authenticate with Azure CLI credentials!")
	}

	subscriptionId := cfg.GetString("SubscriptionId")
	client := datafactory.NewPipelineRunsClient(subscriptionId)
	client.Authorizer = authorizer

	return client
}

func getParsedTime() time.Time {
	timeString := time.Now().Format(time.RFC3339)
	time, err := time.Parse(time.RFC3339, timeString)

	if err != nil {
		panic(err)
	}

	return time
}

func GetPipelineRuns() []PipelineRun {

	var continuationToken string
	var runList []PipelineRun

	resourceGroup := cfg.GetString("ResourceGroupName")
	factory := cfg.GetString("DataFactoryName")

	client := getPipelineRunsClient()

	orderBy := []datafactory.RunQueryOrderBy{{
		OrderBy: "RunStart",
		Order:   "DESC",
	}}

	now := getParsedTime()
	LastUpdatedAfter := now.Add(-3 * time.Hour)
	lastUpdateBefore := now

	filterParameters := datafactory.RunFilterParameters{
		ContinuationToken: &continuationToken,
		LastUpdatedAfter:  &date.Time{LastUpdatedAfter},
		LastUpdatedBefore: &date.Time{lastUpdateBefore},
		Filters:           &[]datafactory.RunQueryFilter{},
		OrderBy:           &orderBy,
	}

	for {

		pipelineRuns, err := client.QueryByFactory(context.Background(), resourceGroup, factory, filterParameters)
		if err != nil {
			error := fmt.Sprintf("We encountered and error while requesting pipeline runs: %s", err.Error())
			panic(error)
		}

		for _, p := range *pipelineRuns.Value {
			runList = append(runList, PipelineRun{
				Name:       *p.PipelineName,
				Status:     *p.Status,
				Start:      p.RunStart,
				End:        p.RunEnd,
				DurationMs: p.DurationInMs,
			})
		}

		if pipelineRuns.ContinuationToken != nil {
			continuationToken = *pipelineRuns.ContinuationToken
		} else {
			break
		}
	}

	return runList

}
