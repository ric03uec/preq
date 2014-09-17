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

func list_pr(c *cli.Context) {
	if err := validate_repo(c); err != nil {
		fmt.Println("Could not list Pull Reqests")
		os.Exit(1)
	}

	fmt.Printf("%s\n", REMOTE_REFS)
	fmt.Printf("%s", CONFIG)
}

func apply_pr(c *cli.Context) {
	if err := validate_repo(c); err != nil {
		fmt.Println("Could not apply the Pull Request")
		os.Exit(1)
	}
	args := c.Args()
	fmt.Printf("%s\n", args)
}

func revert_master(c *cli.Context) {
	if err := validate_repo(c); err != nil {
		fmt.Println("Could not revert to master branch")
		os.Exit(1)
	}
}

func switchRef(c *cli.Context) {
	// switch ref

}

func validate_repo(c *cli.Context) (err error){
	_, gitErr := exec.Command("git", "rev-parse").Output()
	if gitErr != nil {
		fmt.Println("Current directory not under git version control")
		return gitErr
	} else {
		initialize_config()
		return nil
	}
}

func initialize_config() {
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
		refName, refUrl := get_ref()
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
		CONFIG["DEFAULT_REMOTE_REF"] = get_default_ref()
	}

	CONFIG["DEFAULT_BRANCH"] = "master"
}

func get_default_ref() (string) {
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

func get_ref() (string, string) {
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
			Action: list_pr,

		},
		{
			Name: "apply",
			ShortName: "a",
			Usage: "Apply the specified Pull Request",
			Action: apply_pr,
		},
		{
			Name: "revert",
			ShortName: "r",
			Usage: "Revert back to master branch",
			Action: revert_master,

		},
		{
			Name: "switch",
			ShortName: "s",
			Usage: "Switch default remote ref",
			Action: switchRef,
		},
	}

	app.Run(os.Args)
}
