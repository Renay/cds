package pipeline

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ovh/cds/sdk"
)

func init() {
	pipelineHookCmd.AddCommand(pipelineAddHookCmd())
	pipelineHookCmd.AddCommand(pipelineDeleteHookCmd())
	pipelineHookCmd.AddCommand(pipelineListHookCmd())
}

var pipelineHookCmd = &cobra.Command{
	Use:   "hook",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func pipelineAddHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "cds pipeline hook add <projectKey> <applicationName> <pipelineName> [<host>/<project>/<slug>]",
		Long:  ``,
		Run:   addPipelineHook,
	}

	return cmd
}

func pipelineDeleteHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "cds pipeline hook delete <projectKey> <applicationName> <pipelineName> [<host>/<project>/<slug>]",
		Long:  ``,
		Run:   deletePipelineHook,
	}

	return cmd
}

func pipelineListHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "cds pipeline hook list <projectKey> <applicationName> <pipelineName>",
		Long:  ``,
		Run:   listPipelineHook,
	}

	return cmd
}

func addPipelineHook(cmd *cobra.Command, args []string) {

	if len(args) < 3 {
		sdk.Exit("Wrong usage: See %s\n", cmd.Short)
	}

	pipelineProject := args[0]
	appName := args[1]
	pipelineName := args[2]

	p, err := sdk.GetPipeline(pipelineProject, pipelineName)
	if err != nil {
		sdk.Exit("✘ Error: Cannot retrieve pipeline %s-%s (%s)\n", pipelineProject, pipelineName, err)
	}

	a, err := sdk.GetApplication(pipelineProject, appName)
	if err != nil {
		sdk.Exit("✘ Error: Cannot retrieve application %s-%s (%s)\n", pipelineProject, appName, err)
	}

	//If the application is attached to a repositories manager, parameter <host>/<project>/<slug> aren't taken in account
	if a.RepositoriesManager != nil {
		err = sdk.AddHookOnRepositoriesManager(pipelineProject, appName, a.RepositoriesManager.Name, a.RepositoryFullname, pipelineName)
		if err != nil {
			sdk.Exit("✘ Error: Cannot add hook to pipeline %s-%s-%s (%s)\n", pipelineProject, appName, pipelineName, err)
		}
		fmt.Println("✔ Success")
	} else {
		if len(args) != 4 {
			sdk.Exit("✘ Error: Your application has to be attached to a repositories manager. Try : cds application reposmanager attach")
		}
		t := strings.Split(args[3], "/")
		if len(t) != 3 {
			sdk.Exit("✘ Error: Expected repository like <host>/<project>/<slug>. Got %d elements\n", len(t))
		}
		h, err := sdk.AddHook(a, p, t[0], t[1], t[2])
		if err != nil {
			sdk.Exit("✘ Error: Cannot add hook to pipeline %s-%s-%s (%s)\n", pipelineProject, appName, pipelineName, err)
		}
		if strings.Contains(t[0], "stash") {
			fmt.Printf(`Hook created on CDS.
	You now need to configure hook on stash. Use "Http Request Post Receive Hook" to create:
	POST https://<url-to-cds>/hook?&uid=%s&project=%s&name=%s&branch=${refChange.name}&hash=${refChange.toHash}&message=${refChange.type}&author=${user.name}

	`, h.UID, t[1], t[2])
		}
	}
}

func deletePipelineHook(cmd *cobra.Command, args []string) {

	if len(args) < 3 {
		sdk.Exit("Wrong usage: See %s\n", cmd.Short)
	}

	pipelineProject := args[0]
	appName := args[1]
	pipelineName := args[2]

	a, err := sdk.GetApplication(pipelineProject, appName)
	if err != nil {
		sdk.Exit("Cannot retrieve application %s-%s (%s)\n", pipelineProject, appName, err)
	}

	//If the application is attached to a repositories manager, parameter <host>/<project>/<slug> aren't taken in account
	if a.RepositoriesManager != nil {
		err = sdk.DeleteHookOnRepositoriesManager(pipelineProject, appName, a.RepositoriesManager.Name, a.RepositoryFullname, pipelineName)
		if err != nil {
			sdk.Exit("Cannot delete on pipeline %s-%s-%s (%s)\n", pipelineProject, appName, pipelineName, err)
		}
		fmt.Println("✔ Success")
	} else {

		t := strings.Split(args[3], "/")
		if len(t) != 3 {
			sdk.Exit("Expected repository like <host>/<project>/<slug>. Got %d elements\n", len(t))
		}

		hooks, err := sdk.GetHooks(pipelineProject, appName, pipelineName)
		if err != nil {
			sdk.Exit("Cannot retrieve hooks from %s/%s/%s (%s)\n", pipelineProject, appName, pipelineName, err)
		}

		for _, h := range hooks {
			if h.Host == t[0] && h.Project == t[1] && h.Repository == t[2] {
				err = sdk.DeleteHook(pipelineProject, appName, pipelineName, h.ID)
				if err != nil {
					sdk.Exit("Cannot delete hook from %s/%s/%s (%s)", pipelineProject, appName, pipelineName, err)
				}
				return
			}
		}
	}
}

func listPipelineHook(cmd *cobra.Command, args []string) {

	if len(args) != 3 {
		sdk.Exit("Wrong usage: See %s\n", cmd.Short)
	}

	pipelineProject := args[0]
	appName := args[1]
	pipelineName := args[2]

	hooks, err := sdk.GetHooks(pipelineProject, appName, pipelineName)
	if err != nil {
		sdk.Exit("Cannot retrieve hooks from %s/%s/%s (%s)\n", pipelineProject, appName, pipelineName, err)
	}

	for _, h := range hooks {
		fmt.Printf("- %s/%s/%s\n", h.Host, h.Project, h.Repository)
	}

}
