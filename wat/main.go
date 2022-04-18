package main

import (
	"fmt"
	"wat/server"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	fmt.Println("Hello Wat Server!")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		fmt.Println(err)
	}

	svc := sqs.New(sess, aws.NewConfig().WithEndpoint("http://sqs:9324"))

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", server.Status())
	e.GET("/status", server.Status())

	e.POST("/job", server.SubmitJob(svc))
	e.POST("/event/:plugin", server.PushEvent(svc))

	e.Logger.Fatal(e.Start(":5000"))
}
