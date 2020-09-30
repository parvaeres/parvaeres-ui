package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"github.com/yosssi/gohtml"
	"go8s/app"
	"go8s/app/models"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Contact() revel.Result {
	return c.Render()
}

func (c App) DoSubmit() revel.Result {
	log := c.Log.New("route", "create application")
	email := c.Params.Form.Get("email")
	path := c.Params.Form.Get("dir")
	repository := c.Params.Form.Get("repo")
	log.Debug("Received data: email=" + email + ", folder=" + path + ", repository=" + repository)
	appData := models.ParvaeresApplicationData{
		email,
		path,
		repository,
	}
	status, msg := app.ParvaeresHandler.RegisterApplication(appData)
	c.ViewArgs["apiStatus"] = status

	if status {
		var response models.ParvaeresAPIResponse
		json.Unmarshal([]byte(msg), &response)
		log.Debug("received application creation raw-response: " + msg)
		c.ViewArgs["deploymentErrorFlag"] = response.Error
		if response.Error {
			log.Error("error in application creation, reason: " + response.Message)
			c.ViewArgs["deploymentErrorMessage"] = response.Message
		} else {
			log.Debug("application created successfully, UUID: " + response.Items[0].UUID + ", status: " + response.Items[0].Status)
			c.ViewArgs["appUUID"] = response.Items[0].UUID
			c.ViewArgs["appStatus"] = response.Items[0].Status
		}
	}
	return c.Render()
}

func (c App) GetDeployment() revel.Result {
	log := c.Log.New("route", "fetch application")
	id := c.Params.Route.Get("id")
	log.Debug("Application id: " + id)
	status, msg := app.ParvaeresHandler.GetApplication(id)
	c.ViewArgs["apiStatus"] = status

	if status {
		var response models.ParvaeresAPIResponse
		json.Unmarshal([]byte(msg), &response)
		log.Debug("received application fetch raw-response: " + msg)
		c.ViewArgs["deploymentErrorFlag"] = response.Error
		if response.Error {
			log.Error("error in application fetch, reason: " + response.Message)
			c.ViewArgs["deploymentErrorMessage"] = response.Message
		} else {
			c.ViewArgs["appUUID"] = response.Items[0].UUID
			c.ViewArgs["appStatus"] = response.Items[0].Status
			c.ViewArgs["logsUrl"] = response.Items[0].LogsURL
		}
	}
	return c.Render()
}

func (c App) GetDeploymentLogs() revel.Result {
	log := c.Log.New("route", "fetch application logs")
	id := c.Params.Route.Get("id")
	log.Debug("Application id: " + id)
	status, msg := app.ParvaeresHandler.GetApplicationLogs(id)
	c.ViewArgs["apiStatus"] = status

	if status {
		var response models.ParvaeresAPIResponse
		json.Unmarshal([]byte(msg), &response)
		log.Debug("received application logs fetch raw-response: " + msg)
		c.ViewArgs["deploymentErrorFlag"] = response.Error
		if response.Error {
			log.Error("error in application fetch, reason: " + response.Message)
			c.ViewArgs["deploymentErrorMessage"] = response.Message
		} else {
			c.ViewArgs["container"] = response.Items[0].Container
			c.ViewArgs["pod"] = response.Items[0].Pod
			c.ViewArgs["logs"] = gohtml.Format(response.Items[0].Logs)
		}
	}
	return c.Render()
}
