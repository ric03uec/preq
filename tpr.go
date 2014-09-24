package main
import (
	"fmt"
	"strings"
	"strconv"
	"bufio"
	"os"
	"os/exec"
	"github.com/codegangsta/cli"
)

var CONFIG = make(map[string]string)
var REMOTE_REFS = make([]string, 0)

func listPr(c *cli.Context) {
	if err := validateRepo(c); err != nil {
		fmt.Println("Could not list Pull Reqests")
		os.Exit(1)
	}

	outputString , _ := exec.Command("git", "branch", "-r").Output()
	branches := fmt.Sprintf("%s", string(outputString[:]))
	refs := strings.Split(branches, "\n")

	for i := range refs {
		remoteBranch := refs[i]
		refSplits := strings.Split(remoteBranch, "/")
		if length := len(refSplits); length == 3  {
			if strings.TrimSpace(refSplits[0]) == CONFIG["DEFAULT_REMOTE_REF"] {
				fmt.Printf("%s\n", remoteBranch)
			}
		}
	}
}

func applyPr(c *cli.Context) {
	if err := validateRepo(c); err != nil {
		fmt.Println("Could not apply the Pull Request")
		os.Exit(1)
	}

	fmt.Print(fmt.Sprintf("Enter the PR number to apply for ref %s (e.g. 42): ", CONFIG["DEFAULT_REMOTE_REF"]))
	reader := bufio.NewReader(os.Stdin)
	pr, _ := reader.ReadString('\n')
	pr = strings.TrimSpace(pr)
	prNumber, _ := strconv.Atoi(pr)

	outputString , _ := exec.Command("git", "branch", "-r").Output()
	branches := fmt.Sprintf("%s", string(outputString[:]))
	refs := strings.Split(branches, "\n")

	prExists := false
	for i := range refs {
		remoteBranch := refs[i]
		refSplits := strings.Split(remoteBranch, "/")
		if length := len(refSplits); length == 3  {
			if strings.TrimSpace(refSplits[0]) == CONFIG["DEFAULT_REMOTE_REF"] {
				remotePRNumber, _ := strconv.Atoi(refSplits[2])
				if remotePRNumber == prNumber {
					prExists = true
				}
			}
		}
	}
	if prExists {
		remoteRefPath := fmt.Sprintf("%s/pr/%s", CONFIG["DEFAULT_REMOTE_REF"], strconv.Itoa(prNumber))
		_, err := exec.Command("git", "checkout", remoteRefPath).Output()
		if err != nil {
			fmt.Println("Error occured while patching ")
			os.Exit(1)
		}
		fmt.Println("Successfully patched branch with ", remoteRefPath)
		os.Exit(0)
	} else {
		fmt.Println("No PR exists for ref : \n Try refreshing using 'tpr fetch'", CONFIG["DEFAULT_REMOTE_REF"])
		os.Exit(0)
	}
}

func revertMaster(c *cli.Context) {
	outputString, err := exec.Command("git", "checkout", "master").Output()
	if err != nil {
		fmt.Println("Error occured while reverting to master")
		os.Exit(1)
	}
	fmt.Println(fmt.Sprintf("%s", outputString))
	os.Exit(0)
}

func switchRef(c *cli.Context) {
	// switch ref

}

func fetch(c *cli.Context) {
	if err := validateRepo(c); err != nil {
		fmt.Println("Could not list Pull Reqests")
		os.Exit(1)
	}

	fmt.Println("Fetching pull requests for remote ref: ", CONFIG["DEFAULT_REMOTE_REF"])
	refSpec := fmt.Sprintf("refs/pull/*/head:refs/remotes/%s/pr/*", CONFIG["DEFAULT_REMOTE_REF"])

	_, err := exec.Command("git", "fetch", CONFIG["DEFAULT_REMOTE_REF"], refSpec).Output()
	if err != nil {
		fmt.Println("Could not fetch remote Pull Requests")
		os.Exit(1)
	}
	fmt.Println("Successfully fetched PR's for remote ref: ", CONFIG["DEFAULT_REMOTE_REF"])
}

func validateRepo(c *cli.Context) (err error){
	_, gitErr := exec.Command("git", "rev-parse").Output()
	if gitErr != nil {
		fmt.Println("Current directory not under git version control")
		return gitErr
	} else {
		initializeConfig()
		return nil
	}
}

func initializeConfig() {
	output, err := exec.Command("git", "remote", "show").Output()

	if err != nil {
		fmt.Println("Error running 'git remote show'")
		os.Exit(1)
	}
	outputString := fmt.Sprintf("%s", string(output[:]))
	refs := strings.Split(outputString, "\n")
	for index := range refs {
		if len(refs[index]) != 0 {
			REMOTE_REFS = append(REMOTE_REFS, refs[index])
		}
	}
	if len(REMOTE_REFS) == 0 {
		fmt.Println("No remote refs defined")
		refName, refUrl := getRef()
		_, err := exec.Command("git", "remote", "add", refName, refUrl).Output()
		if err != nil {
			fmt.Println("Error while inserting new git ref")
			os.Exit(1)
		}
		CONFIG["DEFAULT_REMOTE_REF"] = refName
		REMOTE_REFS[0] = refName
	} else if len(REMOTE_REFS) == 1 {
		CONFIG["DEFAULT_REMOTE_REF"] = REMOTE_REFS[0]
	} else {
		CONFIG["DEFAULT_REMOTE_REF"] = REMOTE_REFS[0]
		CONFIG["DEFAULT_REMOTE_REF"] = getDefaultRef()
	}

	CONFIG["DEFAULT_BRANCH"] = "master"
}

func getDefaultRef() (string) {
	fmt.Println("Choose ref to set as remote")
	for index := range REMOTE_REFS {
		fmt.Println("\t","(", (index+1), ") ", REMOTE_REFS[index])
	}
	reader := bufio.NewReader(os.Stdin)
	selected, _ := reader.ReadString('\n')
	selected = strings.TrimSpace(selected)
	index, _ := strconv.Atoi(selected)
	return REMOTE_REFS[index - 1]
}

func getRef() (string, string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter ref name (e.g. parent): ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter the url (e.g. git@github.com:ric03uec/tpr.git: ")
	refUrl, _ := reader.ReadString('\n')
	refUrl = strings.TrimSpace(refUrl)

	return name, refUrl
}

func main() {
	app := cli.NewApp()
	app.Name = "tpr"
	app.Usage = "Test github pull requests locally"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name: "list",
			ShortName: "l",
			Usage: "List of all the Pull Requests",
			Action: listPr,

		},
		{
			Name: "apply",
			ShortName: "a",
			Usage: "Apply the specified Pull Request",
			Action: applyPr,
		},
		{
			Name: "revert",
			ShortName: "r",
			Usage: "Revert back to master branch",
			Action: revertMaster,

		},
		{
			Name: "fetch",
			ShortName: "f",
			Usage: "Fetch latest upstream Pull Requests",
			Action: fetch,
		},
	}

	app.Run(os.Args)
}
