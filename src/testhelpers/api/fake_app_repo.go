package api

import (
	"cf"
	"cf/net"
	"generic"
)

type FakeApplicationRepository struct {
	FindAllApps []cf.Application

	ReadName      string
	ReadApp       cf.Application
	ReadErr       bool
	ReadAuthErr   bool
	ReadNotFound  bool

	CreateAppParams    cf.AppParams

	UpdateParams    cf.AppParams
	UpdateAppGuid   string
	UpdateAppResult cf.Application
	UpdateErr       bool

	DeletedAppGuid string
}

func (repo *FakeApplicationRepository) Read(name string) (app cf.Application, apiResponse net.ApiResponse) {
	repo.ReadName = name
	app = repo.ReadApp

	if repo.ReadErr {
		apiResponse = net.NewApiResponseWithMessage("Error finding app by name.")
	}
	if repo.ReadAuthErr {
		apiResponse = net.NewApiResponse("Authentication failed.", "1000", 401)
	}
	if repo.ReadNotFound {
		apiResponse = net.NewNotFoundApiResponse("%s %s not found", "App", name)
	}

	return
}

func (repo *FakeApplicationRepository) Create(params cf.AppParams) (resultApp cf.Application, apiResponse net.ApiResponse) {
	repo.CreateAppParams = params

	resultApp.Guid = params.Get("name").(string) + "-guid"
	resultApp.Name = params.Get("name").(string)
	resultApp.State = "stopped"

	if params.Has("space_guid") {
		resultApp.SpaceGuid = params.Get("space_guid").(string)
	}
	if params.Has("buildpack") {
		resultApp.BuildpackUrl = params.Get("buildpack").(string)
	}
	if params.Has("command") {
		resultApp.Command = params.Get("command").(string)
	}
	if params.Has("disk_quota") {
		resultApp.DiskQuota = params.Get("disk_quota").(uint64)
	}
	if params.Has("instances") {
		resultApp.InstanceCount = params.Get("instances").(int)
	}
	if params.Has("memory") {
		resultApp.Memory = params.Get("memory").(uint64)
	}

	if params.Has("env_vars") {
		for key, val := range params.Get("env_vars").(generic.Map) {
			resultApp.EnvironmentVars[key.(string)] = val.(string)
		}
	}
	return
}

func (repo *FakeApplicationRepository) Update(appGuid string, params cf.AppParams) (updatedApp cf.Application, apiResponse net.ApiResponse) {
	repo.UpdateAppGuid = appGuid
	repo.UpdateParams = params
	updatedApp = repo.UpdateAppResult
	if repo.UpdateErr {
		apiResponse = net.NewApiResponseWithMessage("Error updating app.")
	}
	return
}

func (repo *FakeApplicationRepository) Delete(appGuid string) (apiResponse net.ApiResponse) {
	repo.DeletedAppGuid = appGuid
	return
}
