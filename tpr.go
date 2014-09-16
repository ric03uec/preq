package main
import (
	"fmt"
	"os"
	"os/exec"
	"github.com/codegangsta/cli"
)


func list_pr(c *cli.Context) {
	if err := validate_repo(c); err != nil {
		fmt.Println("Could not list Pull Reqests")
		os.Exit(1)
	}
}

func apply_pr(c *cli.Context) {
	if err := validate_repo(c); err != nil {
		fmt.Println("Could not apply the Pull Request")
		os.Exit(1)
	}
}

func revert_master(c *cli.Context) {
	if err := validate_repo(c); err != nil {
		fmt.Println("Could not revert to master branch")
		os.Exit(1)
	}
}

func parse_args(args []string) {
	// parse the command line paramters to call correct method

}

func validate_repo(c *cli.Context) (err error){
	// check if the current directory is a valid git repo
	_, gitErr := exec.Command("git", "rev-parse").Output()
	if gitErr != nil {
		fmt.Println("Current directory not under git version control")
		return gitErr
	} else {
		return nil
	}
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

		}
	}

	app.Run(os.Args)
}
