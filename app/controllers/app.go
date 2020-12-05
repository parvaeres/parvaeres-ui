package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"github.com/yosssi/gohtml"
	"go8s/app"
	"go8s/app/clients"
	"go8s/app/models"
	"strings"
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
		c.ViewArgs["giterrstatus"] = true
		c.ViewArgs["giterrmessage"] = err.Error()
		log.Warn("Problem with git url detected: " + err.Error())
	} else {
		files, err := gitHandler.GetDirList(fs)
		if err != nil {
			log.Warn("Problem with git folder detected: " + err.Error())
		} else {
			fileContentMap := make(map[string]string)
			fileValidationMap := make(map[string]string)

			c.ViewArgs["files"] = files
			for _, file := range files {
				log.Debug("File within folder " + path + " :- " + file.Name())
				fileContent, err := gitHandler.GetFileContent(fs, file.Name())
				if err != nil {
					log.Warn(err.Error())
				} else {
					fileContentMap[file.Name()] = fileContent
					log.Debug("File content: \n" + fileContent)
					jsonObj, err := clients.Yaml2Json(fileContent)
					if err != nil {
						log.Warn(err.Error())
					} else {
						log.Debug("Json version of file is:\n" + jsonObj)
						status, msg := app.MAOHandler.ValidateApp(jsonObj)
						if status {
							log.Debug("MAO check response: " + msg)

							var response models.MAOResponse
							json.Unmarshal([]byte(msg), &response)
							log.Debug("MAO response received: ", response)

							if response.Reports != nil {
								fileValidationMap[file.Name()] = generateHTML(*response.Reports)
							}

							//var objmap map[string]interface{}
							//if err := json.Unmarshal([]byte(msg), &objmap); err != nil {
							//	log.Warn(err.Error())
							//}
							//if err == nil {
							//	jsonIndent, err := json.MarshalIndent(objmap,"", "  ")
							//	if err != nil {
							//
							//	} else {
							//		fileValidationMap[file.Name()] = string(jsonIndent)
							//	}
							//}
						} else {
							log.Debug("MAO check call resulted in error state: msg=" + msg)
						}
					}
				}
			}
			c.ViewArgs["filescontent"] = fileContentMap
			c.ViewArgs["filesvalidation"] = fileValidationMap
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

	appData := models.ParvaeresApplicationData{
		email,
		path,
		repository,
	}
	status, msg := app.ParvaeresHandler.RegisterApplication(appData)

	//status, msg := true, "{\"Message\":\"FOUND\",\"Items\":[{\"UUID\":\"some-uuid-string\",\"RepoURL\":\"https://github.com/riccardomc/parvaeres-examples.git\",\"Path\":\"guestbook-lb\",\"Email\":\"test@test.ch\",\"Status\":\"DEPLOYED\",\"LiveURLs\":[\"http://a.b.c.d:8081/\"],\"LogsURL\":\"https://www.poc.parvaeres.io/deployment/some-uuid-string/logs\"}]}"

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

func (c App) DeleteDeployment() revel.Result {
	log := c.Log.New("route", "delete application")
	id := c.Params.Route.Get("id")
	log.Debug("Application id: " + id)
	status, msg := app.ParvaeresHandler.DeleteApplication(id)
	c.ViewArgs["apiStatus"] = status

	if status {
		var response models.ParvaeresAPIResponse
		json.Unmarshal([]byte(msg), &response)
		log.Debug("received application delete raw-response: " + msg)
		c.ViewArgs["deploymentErrorFlag"] = response.Error
		if response.Error {
			log.Error("error in application delete, reason: " + response.Message)
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

func generateHTML(reports models.MAOReports) (response string) {
	if reports.OpsConformance != nil {
		response += "<div style=\"border-width:0px;text-align:left;font-family:Helvetica;font-size:14px;padding-bottom:10px;border-color:black;border-style:solid;\">"
		response += "ResourceName=" + reports.WorkloadSupplyChain.ResourceName + "<br>"
		response += "ResourceKind=" + reports.OpsConformance.ResourceKind + "<br>"
		response += "ResourceNamespace=" + reports.OpsConformance.ResourceNamespace + "<br>"
		for _, result := range reports.OpsConformance.Results {
			if strings.Contains(result.Severity, "Medium") {
				response += "<div style=\"background-color: #FFDFD3;padding: 5px 5px 5px 5px;border-width:1px;text-align:left;font-family:Helvetica;font-size:14px;padding-bottom:5px;border-color:orange;border-style:solid;\">"
			} else if strings.Contains(result.Severity, "High") {
				response += "<div style=\"background-color: #D291BC;padding: 5px 5px 5px 5px;border-width:1px;text-align:left;font-family:Helvetica;font-size:14px;padding-bottom:5px;border-color:red;border-style:solid;\">"
			} else {
				response += "<div>"
			}
			response += "<b>For resource:</b> " + result.Resource.Name + "<br>"
			response += "<b>Message:</b> " + strings.TrimSpace(result.Message) + "<br>"
			response += "<b>Recommendation:</b> " + strings.TrimSpace(result.Recommendation) + "<br>"
			response += "</div>"
		}
		response += "</div>"
	}

	if reports.PodSecurity != nil {
		response += "<div style=\"border-width:0px;text-align:left;font-family:Helvetica;font-size:14px;padding-bottom:10px;border-color:black;border-style:solid;\">"
		response += "ResourceName=" + reports.WorkloadSupplyChain.ResourceName + "<br>"
		response += "ResourceKind=" + reports.PodSecurity.ResourceKind + "<br>"
		response += "ResourceNamespace=" + reports.PodSecurity.ResourceNamespace + "<br>"
		for _, result := range reports.PodSecurity.Results {
			if strings.Contains(result.Severity, "Medium") {
				response += "<div style=\"background-color: #FFDFD3;padding: 5px 5px 5px 5px;border-width:1px;text-align:left;font-family:Helvetica;font-size:14px;padding-bottom:5px;border-color:orange;border-style:solid;\">"
			} else if strings.Contains(result.Severity, "High") {
				response += "<div style=\"background-color: #D291BC;padding: 5px 5px 5px 5px;border-width:1px;text-align:left;font-family:Helvetica;font-size:14px;padding-bottom:5px;border-color:red;border-style:solid;\">"
			} else {
				response += "<div>"
			}
			response += "<b>For resource:</b> " + result.Resource.Name + "<br>"
			response += "<b>Message:</b> " + strings.TrimSpace(result.Message) + "<br>"
			response += "<b>Recommendation:</b> " + strings.TrimSpace(result.Recommendation) + "<br>"
			response += "</div>"
		}
		response += "</div>"
	}

	if reports.WorkloadSupplyChain != nil {
		response += "<div style=\"border-width:0px;text-align:left;font-family:Helvetica;font-size:14px;padding-bottom:10px;border-color:black;border-style:solid;\">"
		response += "ResourceName=" + reports.WorkloadSupplyChain.ResourceName + "<br>"
		response += "ResourceKind=" + reports.WorkloadSupplyChain.ResourceKind + "<br>"
		response += "ResourceNamespace=" + reports.WorkloadSupplyChain.ResourceNamespace + "<br>"
		for _, result := range reports.WorkloadSupplyChain.Results {
			if strings.Contains(result.Severity, "Medium") {
				response += "<div style=\"background-color: #FFDFD3;padding: 5px 5px 5px 5px;border-width:1px;text-align:left;font-family:Helvetica;font-size:14px;padding-bottom:5px;border-color:orange;border-style:solid;\">"
			} else if strings.Contains(result.Severity, "High") {
				response += "<div style=\"background-color: #D291BC;padding: 5px 5px 5px 5px;border-width:1px;text-align:left;font-family:Helvetica;font-size:14px;padding-bottom:5px;border-color:red;border-style:solid;\">"
			} else {
				response += "<div>"
			}
			response += "<b>For resource:</b> " + result.Resource.Name + "<br>"
			response += "<b>Message:</b> " + strings.TrimSpace(result.Message) + "<br>"
			response += "<b>Recommendation:</b> " + strings.TrimSpace(result.Recommendation) + "<br>"
			response += "</div>"
		}
		response += "</div>"
	}
	return
}
