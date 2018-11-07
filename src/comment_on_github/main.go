package main

import (
	"os"
    "fmt"
    "log"
    "strings"
    "strconv"
    "context"
    "io/ioutil"
    "encoding/json"
    "github.com/google/go-github/github"
    "github.com/securego/gosec"
    "golang.org/x/oauth2"
)

func usage() {
	fmt.Printf("usage: %s FILE1 FILE2 ...\n", os.Args[0])
	fmt.Printf("\tFILE is output of gosec (https://github.com/securego/gosec) in JSON format\n")
	os.Exit(0)
}

func getFileNamesFromTree(tree *github.Tree) (file_names []string) {
	// log.Printf("%s", tree.Entries)
	for _, x := range tree.Entries {
		// log.Printf("%s", x.Path)
		file_names = append(file_names, *x.Path)
	}
	return
}

// orig. gosec
// Issue is returned by a gosec rule if it discovers an issue with the scanned code.
type Issue struct {
	Severity   string `json:"severity"`   // issue severity (how problematic it is)
	Confidence string `json:"confidence"` // issue confidence (how sure we are we found it)
	RuleID     string `json:"rule_id"`    // Human readable explanation
	What       string `json:"details"`    // Human readable explanation
	File       string `json:"file"`       // File name we found it in
	Code       string `json:"code"`       // Impacted code line
	Line       string `json:"line"`       // Line number in file
}

type GosecResult struct {
	Issues []*Issue
	Stats  *gosec.Metrics
}

func loadGosecJsonFile(filename string) (result *GosecResult) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Fatal(err)
	}
	return
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	// Load parameters
	ACCESS_TOKEN := os.Getenv("GITHUB_ACCESS_TOKEN")
	if ACCESS_TOKEN == "" {
		panic("GITHUB_ACCESS_TOKEN is not set")
	}
	BASE_DIR, _ := os.Getwd()
	GOSEC_RESULT_FILES := os.Args[1:]
	log.Printf("BASE_DIR = %s\n", BASE_DIR)
	log.Printf("GOSEC_RESULT_FILES = %s\n", GOSEC_RESULT_FILES)

	// Connect ot GitHub API with token
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ACCESS_TOKEN},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// FIXME: Hard coded
	owner := "K-atc"
	repo_name := "play-with-gosec"
	sha := "ac94ad47762cc5735c24854d331d4cf9b2a52ad5"
	commit, res, err := client.Git.GetCommit(ctx, owner, repo_name, sha)
	if err != nil {
		panic(res.Status)
	}
	log.Printf("Commit: %s\n", commit.GetMessage())

	// Fetch commited files
	tree, res, err := client.Git.GetTree(ctx, owner, repo_name, commit.Tree.GetSHA(), true) // recursive: True
	if err != nil {
		panic(res.Status)
	}
	commit_files := make(map[string]bool)
	for _, x := range getFileNamesFromTree(tree) {
		commit_files[x] = true
		// or just keys, without values: elementMap[s] = ""
	}
	log.Printf("Commit files: %v\n", commit_files)

	// Remove comments from this commit
	comments, res, err := client.Repositories.ListCommitComments(ctx, owner, repo_name, sha, nil)
	if err != nil {
	    panic(res.Status)
	}
    log.Printf("comments = %v", comments)
	for _, x := range comments {
	    log.Printf("deleting: %s", x)
	    res, err := client.Repositories.DeleteComment(ctx, owner, repo_name, *x.ID)
	    if err != nil {
	        panic(res.Status)
	    }
	}

	// Add comments
	for _, result_file := range GOSEC_RESULT_FILES {
		result := *loadGosecJsonFile(result_file)
		for _, issue := range result.Issues {
            path := strings.Trim(strings.Replace(issue.File, BASE_DIR, "", -1), "/")
            if _, has_key := commit_files[path]; !has_key {
                // log.Printf("[REJECTED] path = %q", path)
                continue
            }
            log.Printf("path = %q", path)
            log.Printf("issue: What = %s", issue.What)

            // Convert absolute path to relative path
            var position int
            if strings.Contains(issue.Line, "-") {
                position, _ = strconv.Atoi(strings.Split(issue.Line, "-")[0])
            } else {
                position, _ = strconv.Atoi(issue.Line)
            }
            log.Printf("position = %d", position)

            // Build comment message
			var body string
			body += "### issue reported by gosec\n"
			body += "**" + issue.Severity + ":" + issue.Confidence + "** "
			body += issue.What + " (" + issue.RuleID + ")\n"
			body += path + ":" + issue.Line + " `" + issue.Code + "`"
			log.Printf(body)

            // Create comment
            comment := &github.RepositoryComment{
                Body: &body,
                Path: &path,
                Position: &position,
            }
            _, res, err := client.Repositories.CreateComment(ctx, owner, repo_name, sha, comment)
            if err != nil {
                panic(res.Status)
            }
		}
	}
	return
}
