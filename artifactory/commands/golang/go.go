package golang

import (
	"github.com/jfrog/gocmd"
	"github.com/jfrog/jfrog-cli-go/artifactory/utils"
	"github.com/jfrog/jfrog-cli-go/artifactory/utils/golang"
	"github.com/jfrog/jfrog-cli-go/artifactory/utils/golang/project"
	"github.com/jfrog/jfrog-cli-go/utils/config"
)

func ExecuteGo(noRegistry, publishDeps bool, goArg []string, targetRepo, buildName, buildNumber string, details *config.ArtifactoryDetails) error {
	err := golang.LogGoVersion()
	if err != nil {
		return err
	}
	isCollectBuildInfo := len(buildName) > 0 && len(buildNumber) > 0
	if isCollectBuildInfo {
		err := utils.SaveBuildGeneralDetails(buildName, buildNumber)
		if err != nil {
			return err
		}
	}
	// The version is not necessary because we are collecting the dependencies only.
	goProject, err := project.Load("-")
	if err != nil {
		return err
	}

	serviceManager, err := utils.CreateServiceManager(details, false)
	if err != nil {
		return err
	}

	err = gocmd.RunWithFallbacksAndPublish(goArg, targetRepo, noRegistry, publishDeps, serviceManager)
	if err != nil {
		return err
	}
	if isCollectBuildInfo {
		err = goProject.LoadDependencies()
		if err != nil {
			return err
		}
		err = utils.SaveBuildInfo(buildName, buildNumber, goProject.BuildInfo(false))
	}

	return err
}