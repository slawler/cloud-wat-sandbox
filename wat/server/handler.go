package server

import (
	"fmt"
	"net/http"
	"wat/jobs"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/labstack/echo/v4"
)

// Status checks if server is available
func Status() echo.HandlerFunc {
	return func(c echo.Context) error { return c.JSON(http.StatusOK, "available") }
}

// PushEvent sends message to appropriate queue
func PushEvent(svc *sqs.SQS) echo.HandlerFunc {
	return func(c echo.Context) error {

		pluginName := c.Param("plugin")
		queueURL := fmt.Sprintf("%v/queue/events", svc.Endpoint)

		var jsonMessage []byte

		event, err := registeredPlugins(pluginName)
		if err := c.Bind(event); err != nil {
			resp := fmt.Sprintf("binding error, verify manifest is correct (%v)", err)
			return c.JSON(http.StatusBadRequest, resp)
		}

		event.SetPlugin(pluginName)
		jsonMessage, err = event.Payload()
		if err != nil {
			resp := fmt.Sprintf("error: %v", err)
			return c.JSON(http.StatusBadRequest, resp)
		}

		sqsResponse, err := SendMsg(svc, &queueURL, string(jsonMessage))
		if err != nil {
			resp := fmt.Sprintf("error: %v", err)
			return c.JSON(http.StatusBadRequest, resp)
		}

		return c.JSON(http.StatusOK, sqsResponse)
	}
}

// SubmitJob provides path to job payload
func SubmitJob(svc *sqs.SQS) echo.HandlerFunc {
	return func(c echo.Context) error {

		queueURL := fmt.Sprintf("%v/queue/events", svc.Endpoint)

		jobConfig := c.QueryParam("config")

		job, err := jobs.ParseJob(jobConfig)
		if err != nil {
			return c.JSON(http.StatusNotFound, err)
		}

		dag, err := job.Plan()
		if err != nil {
			return c.JSON(http.StatusNotFound, err)
		}

		sqsResponse, err := SendMsg(svc, &queueURL, string(dag))
		if err != nil {
			resp := fmt.Sprintf("error: %v", err)
			return c.JSON(http.StatusBadRequest, resp)
		}

		// return c.JSON(http.StatusOK, strings.Join(responses, ";"))
		return c.JSON(http.StatusOK, sqsResponse)
	}

}
