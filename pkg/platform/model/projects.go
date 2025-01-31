package model

import (
	"fmt"
	"strings"

	"github.com/go-openapi/strfmt"

	"github.com/ActiveState/cli/internal/errs"
	"github.com/ActiveState/cli/pkg/platform/api/graphql"
	"github.com/ActiveState/cli/pkg/platform/api/graphql/model"
	"github.com/ActiveState/cli/pkg/platform/api/graphql/request"

	"github.com/ActiveState/cli/internal/constants"
	"github.com/ActiveState/cli/internal/locale"
	"github.com/ActiveState/cli/internal/logging"
	"github.com/ActiveState/cli/pkg/platform/api"
	"github.com/ActiveState/cli/pkg/platform/api/mono/mono_client/projects"
	clientProjects "github.com/ActiveState/cli/pkg/platform/api/mono/mono_client/projects"
	"github.com/ActiveState/cli/pkg/platform/api/mono/mono_models"
	"github.com/ActiveState/cli/pkg/platform/authentication"
)

type ErrProjectNameConflict struct{ *locale.LocalizedError }

type ErrProjectNotFound struct{ *locale.LocalizedError }

// FetchProjectByName fetches a project for an organization.
func FetchProjectByName(orgName string, projectName string) (*mono_models.Project, error) {
	logging.Debug("fetching project (%s) in organization (%s)", projectName, orgName)

	request := request.ProjectByOrgAndName(orgName, projectName)

	gql := graphql.New()
	response := model.Projects{}
	err := gql.Run(request, &response)
	if err != nil {
		return nil, errs.Wrap(err, "GraphQL request failed")
	}

	if len(response.Projects) == 0 {
		if !authentication.LegacyGet().Authenticated() {
			return nil, locale.NewInputError("err_api_project_not_found_unauthenticated", "", orgName, projectName)
		}
		return nil, &ErrProjectNotFound{locale.NewInputError("err_api_project_not_found", "", projectName, orgName)}
	}

	return response.Projects[0].ToMonoProject()
}

// FetchOrganizationProjects fetches the projects for an organization
func FetchOrganizationProjects(orgName string) ([]*mono_models.Project, error) {
	projParams := clientProjects.NewListProjectsParams()
	projParams.SetOrganizationName(orgName)
	orgProjects, err := authentication.Client().Projects.ListProjects(projParams, authentication.ClientAuth())
	if err != nil {
		return nil, processProjectErrorResponse(err)
	}
	return orgProjects.Payload, nil
}

func LanguageByCommit(commitID strfmt.UUID) (Language, error) {
	languages, err := FetchLanguagesForCommit(commitID)
	if err != nil {
		return Language{}, err
	}

	if len(languages) == 0 {
		return Language{}, locale.NewInputError("err_no_languages")
	}

	return languages[0], nil
}

// DefaultBranchForProjectName retrieves the default branch for the given project owner/name.
func DefaultBranchForProjectName(owner, name string) (*mono_models.Branch, error) {
	proj, err := FetchProjectByName(owner, name)
	if err != nil {
		return nil, err
	}

	return DefaultBranchForProject(proj)
}

func BranchesForProject(owner, name string) ([]*mono_models.Branch, error) {
	proj, err := FetchProjectByName(owner, name)
	if err != nil {
		return nil, err
	}
	return proj.Branches, nil
}

func BranchNamesForProjectFiltered(owner, name string, excludes ...string) ([]string, error) {
	proj, err := FetchProjectByName(owner, name)
	if err != nil {
		return nil, err
	}
	branches := make([]string, 0)
	for _, branch := range proj.Branches {
		for _, exclude := range excludes {
			if branch.Label != exclude {
				branches = append(branches, branch.Label)
			}
		}
	}
	return branches, nil
}

// DefaultBranchForProject retrieves the default branch for the given project
func DefaultBranchForProject(pj *mono_models.Project) (*mono_models.Branch, error) {
	for _, branch := range pj.Branches {
		if branch.Default {
			return branch, nil
		}
	}
	return nil, locale.NewError("err_no_default_branch")
}

// BranchForProjectNameByName retrieves the named branch for the given project
// org/name
func BranchForProjectNameByName(owner, name, branch string) (*mono_models.Branch, error) {
	proj, err := FetchProjectByName(owner, name)
	if err != nil {
		return nil, err
	}

	return BranchForProjectByName(proj, branch)
}

// BranchForProjectByName retrieves the named branch for the given project
func BranchForProjectByName(pj *mono_models.Project, name string) (*mono_models.Branch, error) {
	if name == "" {
		return nil, locale.NewInputError("err_empty_branch", "Empty branch name provided.")
	}

	for _, branch := range pj.Branches {
		if branch.Label != "" && branch.Label == name {
			return branch, nil
		}
	}

	return nil, locale.NewInputError(
		"err_no_matching_branch_label",
		"This project has no branch with label matching [NOTICE]{{.V0}}[/RESET].",
		name,
	)
}

// CreateEmptyProject will create the project on the platform
func CreateEmptyProject(owner, name string, private bool) (*mono_models.Project, error) {
	addParams := projects.NewAddProjectParams()
	addParams.SetOrganizationName(owner)
	addParams.SetProject(&mono_models.Project{Name: name, Private: private})
	pj, err := authentication.Client().Projects.AddProject(addParams, authentication.ClientAuth())
	if err != nil {
		msg := api.ErrorMessageFromPayload(err)
		if _, ok := err.(*projects.AddProjectConflict); ok {
			return nil, &ErrProjectNameConflict{locale.WrapInputError(err, msg)}
		}
		return nil, locale.WrapError(err, msg)
	}

	return pj.Payload, nil
}

func CreateCopy(sourceOwner, sourceName, targetOwner, targetName string, makePrivate bool) (*mono_models.Project, error) {
	// Retrieve the source project that we'll be forking
	sourceProject, err := FetchProjectByName(sourceOwner, sourceName)
	if err != nil {
		return nil, locale.WrapInputError(err, "err_fork_fetchProject", "Could not find the source project: {{.V0}}/{{.V1}}", sourceOwner, sourceName)
	}

	// Create the target project
	targetProject, err := CreateEmptyProject(targetOwner, targetName, false)
	if err != nil {
		return nil, locale.WrapError(err, "err_fork_createProject", "Could not create project: {{.V0}}/{{.V1}}", targetOwner, targetName)
	}

	sourceBranch, err := DefaultBranchForProject(sourceProject)
	if err != nil {
		return nil, locale.WrapError(err, "err_branch_nodefault", "Project has no default branch.")
	}
	if sourceBranch.CommitID != nil {
		targetBranch, err := DefaultBranchForProject(targetProject)
		if err != nil {
			return nil, locale.WrapError(err, "err_branch_nodefault", "Project has no default branch.")
		}
		if err := UpdateBranchCommit(targetBranch.BranchID, *sourceBranch.CommitID); err != nil {
			return nil, locale.WrapError(err, "err_fork_branchupdate", "Failed to update branch.")
		}
	}

	// Turn the target project private if this was requested (unfortunately this can't be done int the Creation step)
	if makePrivate {
		if err := MakeProjectPrivate(targetOwner, targetName); err != nil {
			logging.Debug("Cannot make forked project private; deleting public fork.")
			deleteParams := projects.NewDeleteProjectParams()
			deleteParams.SetOrganizationName(targetOwner)
			deleteParams.SetProjectName(targetName)
			if _, err := authentication.Client().Projects.DeleteProject(deleteParams, authentication.ClientAuth()); err != nil {
				return nil, locale.WrapError(
					err, "err_fork_private_but_project_created",
					"Your project was created but could not be made private, please head over to https://{{.V0}}/{{.V1}}/{{.V2}} to manually update your privacy settings.",
					constants.PlatformURL, targetOwner, targetName)
			}
			return nil, locale.WrapError(err, "err_fork_private", "Your fork could not be made private.")
		}
	}

	return targetProject, nil
}

// MakeProjectPrivate turns the given project private
func MakeProjectPrivate(owner, name string) error {
	editParams := projects.NewEditProjectParams()
	yes := true
	editParams.SetProject(&mono_models.ProjectEditable{
		Private: &yes,
	})
	editParams.SetOrganizationName(owner)
	editParams.SetProjectName(name)

	_, err := authentication.Client().Projects.EditProject(editParams, authentication.ClientAuth())
	if err != nil {
		msg := api.ErrorMessageFromPayload(err)
		return locale.WrapError(err, msg)
	}

	return nil
}

// ProjectURL creates a valid platform URL for the given project parameters
func ProjectURL(owner, name, commitID string) string {
	url := fmt.Sprintf("https://%s/%s/%s", constants.PlatformURL, owner, name)
	if commitID != "" {
		url = url + "?commitID=" + commitID
	}
	return url
}

// CommitURL creates a valid platform commit URL for the given commit
func CommitURL(commitID string) string {
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(constants.DashboardCommitURL, "/"), commitID)
}

func processProjectErrorResponse(err error, params ...string) error {
	switch statusCode := api.ErrorCode(err); statusCode {
	case 401:
		return locale.WrapInputError(err, "err_api_not_authenticated")
	case 404:
		p := append([]string{""}, params...)
		return &ErrProjectNotFound{locale.WrapInputError(err, "err_api_project_not_found", p...)}
	default:
		return locale.WrapError(err, "err_api_unknown", "Unexpected API error")
	}
}

func AddBranch(projectID strfmt.UUID, label string) (strfmt.UUID, error) {
	var branchID strfmt.UUID
	addParams := projects.NewAddBranchParams()
	addParams.SetProjectID(projectID)
	addParams.Body.Label = label

	res, err := authentication.Client().Projects.AddBranch(addParams, authentication.ClientAuth())
	if err != nil {
		msg := api.ErrorMessageFromPayload(err)
		return branchID, locale.WrapError(err, msg)
	}

	return res.Payload.BranchID, nil
}
