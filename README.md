tpr
===

Testing github pull requests locally.  
``tpr`` allows you to:
  - list all the pull requests for a pre-defined ```remote``` ref
  - patch your current branch with any of the pull requests
  - revert back to ```master``` branch once you're done testing the pr changes

Assumptions:  
  - read access to the PR repository, in case of private repositories
  - works out of box for public github repositories

Usage:  

```
$ tpr

NAME:
   tpr - Test github pull requests locally

USAGE:
   tpr [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   list, l      List of all the Pull Requests
   apply, a     Apply the specified Pull Request
   revert, r    Revert back to master branch
   fetch, f     Fetch latest upstream Pull Requests
   help, h      Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --version, -v        print the version
   --help, -h           show help
```

Installation:  
```
wget -qO- https://raw.githubusercontent.com/ric03uec/tpr/v0.1.0/install.sh | bash 
```
