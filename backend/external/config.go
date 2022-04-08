package external

import (
	"fmt"

	"github.com/dghubble/oauth1"
)

const (
	TASK_SERVICE_ID_ASANA     = "asana"
	TASK_SERVICE_ID_ATLASSIAN = "atlassian"
	TASK_SERVICE_ID_GT        = "gt"
	TASK_SERVICE_ID_GITHUB    = "github"
	TASK_SERVICE_ID_GOOGLE    = "google"
	TASK_SERVICE_ID_SLACK     = "slack"
	TASK_SERVICE_ID_TRELLO    = "trello"

	TASK_SOURCE_ID_ASANA     = "asana_task"
	TASK_SOURCE_ID_GCAL      = "gcal"
	TASK_SOURCE_ID_GITHUB_PR = "github_pr"
	TASK_SOURCE_ID_GT_TASK   = "gt_task"
	TASK_SOURCE_ID_GMAIL     = "gmail"
	TASK_SOURCE_ID_JIRA      = "jira"
)

type Config struct {
	Github                OauthConfigWrapper
	GoogleLoginConfig     OauthConfigWrapper
	GoogleAuthorizeConfig OauthConfigWrapper
	Slack                 OauthConfigWrapper
	Trello                *oauth1.Config
	Asana                 OauthConfigWrapper
	GoogleOverrideURLs    GoogleURLOverrides
	Atlassian             AtlassianConfig
}

func GetConfig() Config {
	return Config{
		GoogleLoginConfig:     getGoogleLoginConfig(),
		GoogleAuthorizeConfig: getGoogleLinkConfig(),
		Github:                getGithubConfig(),
		Slack:                 getSlackConfig(),
		Trello:                getTrelloConfig(),
		Asana:                 getAsanaConfig(),
		Atlassian:             AtlassianConfig{OauthConfig: getAtlassianOauthConfig()},
	}
}

type TaskServiceResult struct {
	Service TaskService
	Details TaskServiceDetails
	Sources []TaskSource
}

type TaskSourceResult struct {
	Source  TaskSource
	Details TaskSourceDetails
}

func (config Config) GetTaskServiceResult(serviceID string) (*TaskServiceResult, error) {
	nameToService := config.GetNameToService()
	result, ok := nameToService[serviceID]
	if !ok {
		return nil, fmt.Errorf("task service %s not found", serviceID)
	}
	return &result, nil
}

func (config Config) GetTaskSourceResult(sourceID string) (*TaskSourceResult, error) {
	nameToSource := config.getNameToSource()
	result, ok := nameToSource[sourceID]
	if !ok {
		return nil, fmt.Errorf("task source %s not found", sourceID)
	}
	return &result, nil
}

func (config Config) getNameToSource() map[string]TaskSourceResult {
	atlassianService := AtlassianService{Config: config.Atlassian}
	googleService := GoogleService{
		LoginConfig:  config.GoogleLoginConfig,
		LinkConfig:   config.GoogleAuthorizeConfig,
		OverrideURLs: config.GoogleOverrideURLs,
	}
	asanaService := AsanaService{Config: config.Asana}
	githubService := GithubService{Config: config.Github}
	return map[string]TaskSourceResult{
		TASK_SOURCE_ID_GCAL: {
			Details: TaskSourceGoogleCalendar,
			Source:  GoogleCalendarSource{Google: googleService},
		},
		TASK_SOURCE_ID_GT_TASK: {
			Details: TaskSourceGeneralTask,
			Source:  GeneralTaskTaskSource{},
		},
		TASK_SOURCE_ID_GMAIL: {
			Details: TaskSourceGmail,
			Source:  GmailSource{Google: googleService},
		},
		TASK_SOURCE_ID_JIRA: {
			Details: TaskSourceJIRA,
			Source:  JIRASource{Atlassian: atlassianService},
		},
		TASK_SOURCE_ID_ASANA: {
			Details: TaskSourceAsana,
			Source:  AsanaTaskSource{Asana: asanaService},
		},
		TASK_SOURCE_ID_GITHUB_PR: {
			Details: TaskSourceGithubPR,
			Source:  GithubPRSource{Github: githubService},
		},
	}
}

func (config Config) GetNameToService() map[string]TaskServiceResult {
	asanaService := AsanaService{Config: config.Asana}
	atlassianService := AtlassianService{Config: config.Atlassian}
	googleService := GoogleService{
		LoginConfig:  config.GoogleLoginConfig,
		LinkConfig:   config.GoogleAuthorizeConfig,
		OverrideURLs: config.GoogleOverrideURLs,
	}
	githubService := GithubService{Config: config.Github}
	return map[string]TaskServiceResult{
		TASK_SERVICE_ID_ATLASSIAN: {
			Service: atlassianService,
			Details: TaskServiceAtlassian,
			Sources: []TaskSource{JIRASource{Atlassian: atlassianService}},
		},
		TASK_SERVICE_ID_GT: {
			Service: GeneralTaskService{},
			Details: TaskServiceGeneralTask,
			Sources: []TaskSource{GeneralTaskTaskSource{}},
		},
		TASK_SERVICE_ID_GOOGLE: {
			Service: googleService,
			Details: TaskServiceGoogle,
			Sources: []TaskSource{
				GmailSource{Google: googleService},
				GoogleCalendarSource{Google: googleService},
			},
		},
		TASK_SERVICE_ID_SLACK: {
			Service: SlackService{Config: config.Slack},
			Details: TaskServiceSlack,
			Sources: []TaskSource{},
		},
		TASK_SERVICE_ID_GITHUB: {
			Service: githubService,
			Details: TaskServiceGithub,
			Sources: []TaskSource{GithubPRSource{Github: githubService}},
		},
		TASK_SERVICE_ID_TRELLO: {
			Service: TrelloService{Config: config.Trello},
			Details: TaskServiceTrello,
			Sources: []TaskSource{},
		},
		TASK_SERVICE_ID_ASANA: {
			Service: asanaService,
			Details: TaskServiceAsana,
			Sources: []TaskSource{AsanaTaskSource{Asana: asanaService}},
		},
	}
}

type AuthType string

var AuthTypeOauth2 AuthType = "oauth2"
var AuthTypeOauth1 AuthType = "oauth1"

type TaskServiceDetails struct {
	ID           string
	Name         string
	Logo         string
	LogoV2       string
	AuthType     AuthType
	IsLinkable   bool
	IsSignupable bool
}

var TaskServiceAtlassian = TaskServiceDetails{
	ID:           TASK_SERVICE_ID_ATLASSIAN,
	Name:         "Atlassian",
	Logo:         "/images/jira.svg",
	LogoV2:       "jira-v2",
	AuthType:     AuthTypeOauth2,
	IsLinkable:   false,
	IsSignupable: false,
}
var TaskServiceGeneralTask = TaskServiceDetails{
	ID:           TASK_SERVICE_ID_ATLASSIAN,
	Name:         "General Task",
	Logo:         "/images/generaltask.svg",
	LogoV2:       "generaltask",
	AuthType:     AuthTypeOauth2,
	IsLinkable:   false,
	IsSignupable: false,
}
var TaskServiceGithub = TaskServiceDetails{
	ID:           TASK_SERVICE_ID_GITHUB,
	Name:         "Github",
	Logo:         "/images/github.svg",
	LogoV2:       "github",
	AuthType:     AuthTypeOauth2,
	IsLinkable:   true,
	IsSignupable: false,
}
var TaskServiceGoogle = TaskServiceDetails{
	ID:           TASK_SERVICE_ID_GOOGLE,
	Name:         "Google",
	Logo:         "/images/gmail.svg",
	LogoV2:       "gmail",
	AuthType:     AuthTypeOauth2,
	IsLinkable:   true,
	IsSignupable: true,
}
var TaskServiceSlack = TaskServiceDetails{
	ID:           TASK_SERVICE_ID_SLACK,
	Name:         "Slack",
	Logo:         "/images/slack.svg",
	LogoV2:       "slack",
	AuthType:     AuthTypeOauth2,
	IsLinkable:   false,
	IsSignupable: false,
}
var TaskServiceTrello = TaskServiceDetails{
	ID:           TASK_SERVICE_ID_TRELLO,
	Name:         "Trello",
	Logo:         "/images/trello.svg",
	LogoV2:       "trello",
	AuthType:     AuthTypeOauth1,
	IsLinkable:   false,
	IsSignupable: false,
}
var TaskServiceAsana = TaskServiceDetails{
	ID:           TASK_SERVICE_ID_ASANA,
	Name:         "Asana",
	Logo:         "/images/asana.svg",
	LogoV2:       "asana",
	AuthType:     AuthTypeOauth2,
	IsLinkable:   false,
	IsSignupable: false,
}

type TaskSourceDetails struct {
	ID                     string
	Name                   string
	Logo                   string
	LogoV2                 string
	IsCompletable          bool
	CanCreateTask          bool
	IsReplyable            bool
	CanCreateCalendarEvent bool
}

var TaskSourceGeneralTask = TaskSourceDetails{
	ID:                     TASK_SOURCE_ID_GT_TASK,
	Name:                   "General Task",
	Logo:                   "/images/generaltask.svg",
	LogoV2:                 "generaltask",
	IsCompletable:          true,
	CanCreateTask:          true,
	IsReplyable:            false,
	CanCreateCalendarEvent: false,
}
var TaskSourceGoogleCalendar = TaskSourceDetails{
	ID:                     TASK_SOURCE_ID_GCAL,
	Name:                   "Google Calendar",
	Logo:                   "/images/gcal.svg",
	LogoV2:                 "gcal",
	IsCompletable:          false,
	CanCreateTask:          false,
	IsReplyable:            false,
	CanCreateCalendarEvent: true,
}
var TaskSourceGithubPR = TaskSourceDetails{
	ID:                     TASK_SOURCE_ID_GITHUB_PR,
	Name:                   "Git PR",
	Logo:                   "/images/github.svg",
	LogoV2:                 "github",
	IsCompletable:          true,
	CanCreateTask:          false,
	IsReplyable:            false,
	CanCreateCalendarEvent: false,
}
var TaskSourceGmail = TaskSourceDetails{
	ID:                     TASK_SOURCE_ID_GMAIL,
	Name:                   "Gmail",
	Logo:                   "/images/gmail.svg",
	LogoV2:                 "gmail",
	IsCompletable:          true,
	CanCreateTask:          false,
	IsReplyable:            true,
	CanCreateCalendarEvent: false,
}
var TaskSourceJIRA = TaskSourceDetails{
	ID:                     TASK_SOURCE_ID_JIRA,
	Name:                   "Jira",
	Logo:                   "/images/jira.svg",
	LogoV2:                 "jira",
	IsCompletable:          true,
	CanCreateTask:          false,
	IsReplyable:            false,
	CanCreateCalendarEvent: false,
}
var TaskSourceAsana = TaskSourceDetails{
	ID:                     TASK_SOURCE_ID_ASANA,
	Name:                   "Asana",
	Logo:                   "/images/asana.svg",
	LogoV2:                 "asana",
	IsCompletable:          true,
	CanCreateTask:          false,
	IsReplyable:            false,
	CanCreateCalendarEvent: false,
}
