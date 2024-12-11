package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/xanzy/go-gitlab"
)

type config struct {
	AccessToken string `env:"GITLAB_ACCESS_TOKEN,required"`
	GroupId     int    `env:"GITLAB_GROUP_ID,required"`
	ProjectsIds []int  `env:"GITLAB_PROJECTS_IDS,required"`
	Host        string `env:"GITLAB_HOST,required"`
}

func loadConfig() (c *config, err error) {
	c = &config{}
	err = env.Parse(c)

	return
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func prepareExcludedBranchesInProject(git *gitlab.Client, projectId int) (branchesToExcludeList map[string]bool) {
	branchesToExcludeList = map[string]bool{}

	for {
		optMr := &gitlab.ListProjectMergeRequestsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page:    1,
			},
			State: gitlab.String("opened"),
		}
		mrList, resp, err := git.MergeRequests.ListProjectMergeRequests(projectId, optMr)

		check(err)

		for _, mr := range mrList {
			log.Printf("Project id: %d. Found MR with name: %s, for: %s. Adding to excluded list...\n", projectId, mr.Title, mr.SourceBranch)
			branchesToExcludeList[mr.SourceBranch] = true
		}

		if resp.NextPage == 0 {
			return
		}

		optMr.Page = resp.NextPage
	}
}

func shouldSkipBranch(branch *gitlab.Branch, projectId int, branchesToExcludeList map[string]bool) bool {
	if branch.Protected {
		log.Printf("Project id: %d. Branch: %s is protected, skipping...\n", projectId, branch.Name)

		return true
	}

	if _, ok := branchesToExcludeList[branch.Name]; ok {
		log.Printf("Project id: %d. Branch: %s is excluded, skipping...\n", projectId, branch.Name)

		return true
	}

	if branch.Commit.CommittedDate.After(time.Now().AddDate(0, -2, 0)) {
		log.Printf("Project id: %d. Branch: %s is not stale, skipping...\n", projectId, branch.Name)

		return true
	}

	return false
}

func deleteBranchesInProject(git *gitlab.Client, projectId int, branchesToExcludeList map[string]bool, dryRun bool) (excludedBranches []string, deletedBranches []string) {
	optBranch := &gitlab.ListBranchesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	for {
		branchesList, resp, err := git.Branches.ListBranches(projectId, optBranch)
		check(err)

		for _, branch := range branchesList {
			if shouldSkipBranch(branch, projectId, branchesToExcludeList) {
				excludedBranches = append(excludedBranches, branch.Name)
				continue
			}

			log.Printf("Project id: %d. Branch: %s should be deleted, deleting...\n", projectId, branch.Name)
			//delete branch
			deletedBranches = append(deletedBranches, branch.Name)
			if dryRun {
				continue
			}

			_, err = git.Branches.DeleteBranch(projectId, branch.Name)
			check(err)
		}

		if resp.NextPage == 0 {
			return
		}

		optBranch.Page = resp.NextPage
	}
}

func handleBranchesToRemoveForProject(git *gitlab.Client, project *gitlab.Project, reportFileName string, dryRun bool) {
	log.Printf("\n====================\nFound project: %s with ID: %d\n====================\n", project.Name, project.ID)

	branchesToExcludeList := prepareExcludedBranchesInProject(git, project.ID)
	excludedBranches, deletedBranches := deleteBranchesInProject(git, project.ID, branchesToExcludeList, dryRun)
	updateReport(reportFileName, project, branchesToExcludeList, excludedBranches, deletedBranches)
}

func prepareReportFile() (reportName string, err error) {
	t := time.Now()
	reportName = "report_" + t.Format("2006-01-02_15-04-05") + ".txt"
	f, err := os.Create(reportName)

	check(err)

	defer f.Close()

	return
}

func handleBranchesToRemoveForGroup(git *gitlab.Client, groupId int, dryRun bool) {
	reportFileName, err := prepareReportFile()

	check(err)
	optProject := &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	for {
		projectsList, resp, err := git.Groups.ListGroupProjects(groupId, optProject)
		check(err)

		for _, project := range projectsList {
			handleBranchesToRemoveForProject(git, project, reportFileName, dryRun)
		}

		if resp.NextPage == 0 {
			return
		}

		optProject.Page = resp.NextPage
	}
}

func handleBranchesToRemoveByProjectsIds(git *gitlab.Client, projectsIds []int, dryRun bool) {
	reportFileName, err := prepareReportFile()

	check(err)

	for _, projectId := range projectsIds {
		project, _, err := git.Projects.GetProject(projectId, &gitlab.GetProjectOptions{})
		check(err)
		handleBranchesToRemoveForProject(git, project, reportFileName, dryRun)
	}
}

func updateReport(reportName string, project *gitlab.Project, withMergeRequest map[string]bool, excludedBranches []string, deletedBranches []string) {
	var reportText strings.Builder
	reportText.WriteString("========== Project name: " + project.Name + " with ID: " + strconv.Itoa(project.ID) + " ==========\n")
	reportText.WriteString("=== Branches with MR: ===\n")
	for key := range withMergeRequest {
		reportText.WriteString(key + "\n")
	}

	reportText.WriteString("=== Excluded branches: ===\n")
	for _, val := range excludedBranches {
		reportText.WriteString(val + "\n")
	}

	reportText.WriteString("=== Deleted branches: ===\n")
	for _, val := range deletedBranches {
		reportText.WriteString(val + "\n")
	}
	reportText.WriteString("\n")

	f, err := os.OpenFile(reportName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err)

	defer f.Close()

	_, err = f.WriteString(reportText.String())
	check(err)
}

func main() {
	config, err := loadConfig()
	check(err)

	runType := flag.String("scope", "projects", "project|group set scope for the run")
	dryRun := flag.Bool("dry-run", true, "Should show actions without deleting branches?")
	if *dryRun {
		log.Printf("Running with dry run mode\n")
	}

	flag.Parse()
	flag.Usage = func() {
		log.Println("Usage:")
		flag.PrintDefaults()
	}

	git, err := gitlab.NewClient(config.AccessToken, gitlab.WithBaseURL(config.Host))
	check(err)
	log.Printf("Connection to %s established\n", config.Host)

	switch *runType {
	case "projects":
		// We need to use IDs.
		// the "/projects?search=test" endpoint uses LIKE which returns multiple records
		handleBranchesToRemoveByProjectsIds(git, config.ProjectsIds, *dryRun)
	case "group":
		if !*dryRun {
			log.Fatal("This option is available only with dry-run=true option")
		}
		handleBranchesToRemoveForGroup(git, config.GroupId, *dryRun)
	default:
		log.Printf("The option %s is not supported\n", *runType)
	}
}
