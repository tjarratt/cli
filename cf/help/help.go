package help

import (
	"os"
	"strings"
	"text/tabwriter"
	"text/template"

	cli "github.com/tjarratt/cg_cli"
	"github.com/tjarratt/cli/cf/terminal"
)

type groupedCommands struct {
	Name             string
	CommandSubGroups [][]cmdPresenter
}

func (c groupedCommands) SubTitle(name string) string {
	return terminal.HeaderColor(name + ":")
}

type cmdPresenter struct {
	Name        string
	Description string
}

func NewCmdPresenter(app *cli.App, maxNameLen int, cmdName string) (presenter cmdPresenter) {
	cmd := app.Command(cmdName)

	presenter.Name = presentCmdName(*cmd)
	padding := strings.Repeat(" ", maxNameLen-len(presenter.Name))
	presenter.Name = presenter.Name + padding
	presenter.Description = cmd.Description

	return
}

func presentCmdName(cmd cli.Command) (name string) {
	name = cmd.Name
	if cmd.ShortName != "" {
		name = name + ", " + cmd.ShortName
	}
	return
}

type appPresenter struct {
	cli.App
	Commands []groupedCommands
}

func (p appPresenter) Title(name string) string {
	return terminal.HeaderColor(name)
}

func getMaxCmdNameLength(app *cli.App) (length int) {
	for _, cmd := range app.Commands {
		name := presentCmdName(cmd)
		if len(name) > length {
			length = len(name)
		}
	}
	return
}

func NewAppPresenter(app *cli.App) (presenter appPresenter) {
	maxNameLen := getMaxCmdNameLength(app)

	presenter.Name = app.Name
	presenter.Usage = app.Usage
	presenter.Version = app.Version
	presenter.Name = app.Name
	presenter.Flags = app.Flags
	presenter.Compiled = app.Compiled

	presenter.Commands = []groupedCommands{
		{
			Name: "GETTING STARTED",
			CommandSubGroups: [][]cmdPresenter{
				{
					NewCmdPresenter(app, maxNameLen, "login"),
					NewCmdPresenter(app, maxNameLen, "logout"),
					NewCmdPresenter(app, maxNameLen, "passwd"),
					NewCmdPresenter(app, maxNameLen, "target"),
				}, {
					NewCmdPresenter(app, maxNameLen, "api"),
					NewCmdPresenter(app, maxNameLen, "auth"),
				},
			},
		}, {
			Name: "APPS",
			CommandSubGroups: [][]cmdPresenter{
				{
					NewCmdPresenter(app, maxNameLen, "apps"),
					NewCmdPresenter(app, maxNameLen, "app"),
				}, {
					NewCmdPresenter(app, maxNameLen, "push"),
					NewCmdPresenter(app, maxNameLen, "scale"),
					NewCmdPresenter(app, maxNameLen, "delete"),
					NewCmdPresenter(app, maxNameLen, "rename"),
				}, {
					NewCmdPresenter(app, maxNameLen, "start"),
					NewCmdPresenter(app, maxNameLen, "stop"),
					NewCmdPresenter(app, maxNameLen, "restart"),
				}, {
					NewCmdPresenter(app, maxNameLen, "events"),
					NewCmdPresenter(app, maxNameLen, "files"),
					NewCmdPresenter(app, maxNameLen, "logs"),
				}, {
					NewCmdPresenter(app, maxNameLen, "env"),
					NewCmdPresenter(app, maxNameLen, "set-env"),
					NewCmdPresenter(app, maxNameLen, "unset-env"),
				}, {
					NewCmdPresenter(app, maxNameLen, "stacks"),
				},
			},
		}, {
			Name: "SERVICES",
			CommandSubGroups: [][]cmdPresenter{
				{
					NewCmdPresenter(app, maxNameLen, "marketplace"),
					NewCmdPresenter(app, maxNameLen, "services"),
					NewCmdPresenter(app, maxNameLen, "service"),
				}, {
					NewCmdPresenter(app, maxNameLen, "create-service"),
					NewCmdPresenter(app, maxNameLen, "delete-service"),
					NewCmdPresenter(app, maxNameLen, "rename-service"),
				}, {
					NewCmdPresenter(app, maxNameLen, "bind-service"),
					NewCmdPresenter(app, maxNameLen, "unbind-service"),
				}, {
					NewCmdPresenter(app, maxNameLen, "create-user-provided-service"),
					NewCmdPresenter(app, maxNameLen, "update-user-provided-service"),
				},
			},
		}, {
			Name: "ORGS",
			CommandSubGroups: [][]cmdPresenter{
				{
					NewCmdPresenter(app, maxNameLen, "orgs"),
					NewCmdPresenter(app, maxNameLen, "org"),
				}, {
					NewCmdPresenter(app, maxNameLen, "create-org"),
					NewCmdPresenter(app, maxNameLen, "delete-org"),
					NewCmdPresenter(app, maxNameLen, "rename-org"),
				},
			},
		}, {
			Name: "SPACES",
			CommandSubGroups: [][]cmdPresenter{
				{
					NewCmdPresenter(app, maxNameLen, "spaces"),
					NewCmdPresenter(app, maxNameLen, "space"),
				}, {
					NewCmdPresenter(app, maxNameLen, "create-space"),
					NewCmdPresenter(app, maxNameLen, "delete-space"),
					NewCmdPresenter(app, maxNameLen, "rename-space"),
				},
			},
		}, {
			Name: "DOMAINS",
			CommandSubGroups: [][]cmdPresenter{
				{
					NewCmdPresenter(app, maxNameLen, "domains"),
					NewCmdPresenter(app, maxNameLen, "create-domain"),
					NewCmdPresenter(app, maxNameLen, "delete-domain"),
					NewCmdPresenter(app, maxNameLen, "create-shared-domain"),
					NewCmdPresenter(app, maxNameLen, "delete-shared-domain"),
				},
			},
		}, {
			Name: "ROUTES",
			CommandSubGroups: [][]cmdPresenter{
				{
					NewCmdPresenter(app, maxNameLen, "routes"),
					NewCmdPresenter(app, maxNameLen, "create-route"),
					NewCmdPresenter(app, maxNameLen, "map-route"),
					NewCmdPresenter(app, maxNameLen, "unmap-route"),
					NewCmdPresenter(app, maxNameLen, "delete-route"),
					NewCmdPresenter(app, maxNameLen, "delete-orphaned-routes"),
				},
			},
		}, {
			Name: "BUILDPACKS",
			CommandSubGroups: [][]cmdPresenter{
				{
					NewCmdPresenter(app, maxNameLen, "buildpacks"),
					NewCmdPresenter(app, maxNameLen, "create-buildpack"),
					NewCmdPresenter(app, maxNameLen, "update-buildpack"),
					NewCmdPresenter(app, maxNameLen, "rename-buildpack"),
					NewCmdPresenter(app, maxNameLen, "delete-buildpack"),
				},
			},
		}, {
			Name: "USER ADMIN",
			CommandSubGroups: [][]cmdPresenter{
				{
					NewCmdPresenter(app, maxNameLen, "create-user"),
					NewCmdPresenter(app, maxNameLen, "delete-user"),
				}, {
					NewCmdPresenter(app, maxNameLen, "org-users"),
					NewCmdPresenter(app, maxNameLen, "set-org-role"),
					NewCmdPresenter(app, maxNameLen, "unset-org-role"),
				}, {
					NewCmdPresenter(app, maxNameLen, "space-users"),
					NewCmdPresenter(app, maxNameLen, "set-space-role"),
					NewCmdPresenter(app, maxNameLen, "unset-space-role"),
				},
			},
		}, {
			Name: "ORG ADMIN",
			CommandSubGroups: [][]cmdPresenter{
				{
					NewCmdPresenter(app, maxNameLen, "quotas"),
					NewCmdPresenter(app, maxNameLen, "quota"),
					NewCmdPresenter(app, maxNameLen, "set-quota"),
				}, {
					NewCmdPresenter(app, maxNameLen, "create-quota"),
					NewCmdPresenter(app, maxNameLen, "delete-quota"),
					NewCmdPresenter(app, maxNameLen, "update-quota"),
				},
			},
		}, {
			Name: "SERVICE ADMIN",
			CommandSubGroups: [][]cmdPresenter{
				{
					NewCmdPresenter(app, maxNameLen, "service-auth-tokens"),
					NewCmdPresenter(app, maxNameLen, "create-service-auth-token"),
					NewCmdPresenter(app, maxNameLen, "update-service-auth-token"),
					NewCmdPresenter(app, maxNameLen, "delete-service-auth-token"),
				}, {
					NewCmdPresenter(app, maxNameLen, "service-brokers"),
					NewCmdPresenter(app, maxNameLen, "create-service-broker"),
					NewCmdPresenter(app, maxNameLen, "update-service-broker"),
					NewCmdPresenter(app, maxNameLen, "delete-service-broker"),
					NewCmdPresenter(app, maxNameLen, "rename-service-broker"),
				}, {
					NewCmdPresenter(app, maxNameLen, "migrate-service-instances"),
					NewCmdPresenter(app, maxNameLen, "purge-service-offering"),
				},
			},
		}, {
			Name: "ADVANCED",
			CommandSubGroups: [][]cmdPresenter{
				{
					NewCmdPresenter(app, maxNameLen, "curl"),
					NewCmdPresenter(app, maxNameLen, "config"),
				},
			},
		},
	}
	return
}

func ShowAppHelp(helpTemplate string, appToPrint interface{}) {
	showAppHelp(helpTemplate, appToPrint)
}

func showAppHelp(helpTemplate string, appToPrint interface{}) {
	app := appToPrint.(*cli.App)
	presenter := NewAppPresenter(app)

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	t := template.Must(template.New("help").Parse(helpTemplate))
	t.Execute(w, presenter)
	w.Flush()
}

var appHelpTemplate = `{{.Title "NAME:"}}
   {{.Name}} - {{.Usage}}

{{.Title "USAGE:"}}
   [environment variables] {{.Name}} [global options] command [arguments...] [command options]

{{.Title "VERSION:"}}
   {{.Version}}

{{.Title "BUILD TIME:"}}
   {{.Compiled}}
   {{range .Commands}}
{{.SubTitle .Name}}{{range .CommandSubGroups}}
{{range .}}   {{.Name}} {{.Description}}
{{end}}{{end}}{{end}}
{{.Title "ENVIRONMENT VARIABLES"}}
   CF_COLOR=false                     Do not colorize output
   CF_HOME=path/to/dir/               Override path to default config directory
   CF_STAGING_TIMEOUT=15              Max wait time for buildpack staging, in minutes
   CF_STARTUP_TIMEOUT=5               Max wait time for app instance startup, in minutes
   CF_TRACE=true                      Print API request diagnostics to stdout
   CF_TRACE=path/to/trace.log         Append API request diagnostics to a log file
   HTTP_PROXY=proxy.example.com:8080  Enable HTTP proxying for API requests

{{.Title "GLOBAL OPTIONS"}}
   --version, -v                      Print the version
   --help, -h                         Show help
`

func AppHelpTemplate() string {
	return appHelpTemplate
}
