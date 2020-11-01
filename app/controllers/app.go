package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"github.com/yosssi/gohtml"
	"go8s/app"
	"go8s/app/clients"
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

func (c App) DoValidate() revel.Result {
	log := c.Log.New("route", "validate application")
	email := c.Params.Form.Get("email")
	path := c.Params.Form.Get("dir")
	repository := c.Params.Form.Get("repo")
	c.ViewArgs["email"] = email
	c.ViewArgs["dir"] = path
	c.ViewArgs["repo"] = repository
	log.Debug("Received data: email=" + email + ", folder=" + path + ", repository=" + repository)
	//quickly checking if repo and folders are valid
	gitHandler := clients.GitParam{}
	gitHandler.URL = repository
	gitHandler.Folder = path
	fs, err := gitHandler.GetMemFS()
	if err != nil {
		log.Warn("Problem with git url detected: " + err.Error())
	} else {
		files, err := gitHandler.GetDirList(fs)
		if err != nil {
			log.Warn("Problem with git folder detected: " + err.Error())
		} else {
			c.ViewArgs["files"] = files
			for _, file := range files {
				log.Debug("File within folder " + path + " :- " + file.Name())
			}
		}
	}
	//TODO make this a 2 step process: validate and then deploy

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
			if len(response.Items) > 0 {
				log.Debug("application created successfully, UUID: " + response.Items[0].UUID + ", status: " + response.Items[0].Status)
				c.ViewArgs["appUUID"] = response.Items[0].UUID
				c.ViewArgs["appStatus"] = response.Items[0].Status
			}
		}
	}
	return c.Render()
}

func (c App) DoSubmit() revel.Result {
	log := c.Log.New("route", "create application")
	email := c.Params.Form.Get("email")
	path := c.Params.Form.Get("dir")
	repository := c.Params.Form.Get("repo")
	log.Debug("Received data: email=" + email + ", folder=" + path + ", repository=" + repository)
	//quickly checking if repo and folders are valid
	gitHandler := clients.GitParam{}
	gitHandler.URL = repository
	gitHandler.Folder = path
	fs, err := gitHandler.GetMemFS()
	if err != nil {
		log.Warn("Problem with git url detected: " + err.Error())
	} else {
		files, err := gitHandler.GetDirList(fs)
		if err != nil {
			log.Warn("Problem with git folder detected: " + err.Error())
		} else {
			for _, file := range files {
				log.Debug("File within folder " + path + " :- " + file.Name())
			}
		}
	}
	//TODO make this a 2 step process: validate and then deploy

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
			if len(response.Items) > 0 {
				log.Debug("application created successfully, UUID: " + response.Items[0].UUID + ", status: " + response.Items[0].Status)
				c.ViewArgs["appUUID"] = response.Items[0].UUID
				c.ViewArgs["appStatus"] = response.Items[0].Status
			}
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
			if len(response.Items) > 0 {
				c.ViewArgs["appUUID"] = response.Items[0].UUID
				c.ViewArgs["appStatus"] = response.Items[0].Status
				c.ViewArgs["logsUrl"] = response.Items[0].LogsURL
				c.ViewArgs["liveUrls"] = response.Items[0].LiveURLs
			}
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
			if len(response.Items) > 0 {
				c.ViewArgs["container"] = response.Items[0].Container
				c.ViewArgs["pod"] = response.Items[0].Pod
				c.ViewArgs["logs"] = gohtml.Format(response.Items[0].Logs)
			}

		}
	}
	return c.Render()
}
