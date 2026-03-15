package controllers

import s "github.com/verywelloo/3-go-echo-task-management/app/services"

var (
	projectCollection           = s.AppInstance.Collections.Projects
	projectPermissionCollection = s.AppInstance.Collections.ProjectPermission
	userCollection              = s.AppInstance.Collections.Users
)
