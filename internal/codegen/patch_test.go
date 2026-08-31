package codegen

import (
	"go/format"
	"os"
	"strings"
	"testing"
)

// mustParse fails the test if an edit produced source that does not compile as
// Go — the whole point of splicing rather than reprinting is that the result
// stays valid.
func mustParse(t *testing.T, label string, src []byte) {
	t.Helper()
	if _, err := format.Source(src); err != nil {
		t.Fatalf("%s: result does not parse: %v\n%s", label, err, numbered(string(src)))
	}
}

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestAddImports(t *testing.T) {
	src := read(t, "../cli/mothership/stats.go")

	out, err := AddImports(src, []string{"net/http", "github.com/google/uuid", "encoding/json"})
	if err != nil {
		t.Fatal(err)
	}
	mustParse(t, "AddImports", out)

	got := string(out)
	for _, want := range []string{`"net/http"`, `"github.com/google/uuid"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing import %s", want)
		}
	}
	// encoding/json was already imported and must not be duplicated.
	if n := strings.Count(got, `"encoding/json"`); n != 1 {
		t.Errorf("encoding/json appears %d times, want 1", n)
	}
	// The standard library and the rest stay in separate groups.
	imports := got[strings.Index(got, "import ("):strings.Index(got, ")\n")]
	if !strings.Contains(imports, "\n\n") {
		t.Errorf("import groups collapsed:\n%s", imports)
	}

	// Idempotent.
	again, err := AddImports(out, []string{"net/http"})
	if err != nil {
		t.Fatal(err)
	}
	if string(again) != got {
		t.Error("re-adding an existing import changed the file")
	}
}

func TestAddFlagsToInit(t *testing.T) {
	// stats.go has no init() at all — the bug that leaves --orgName unregistered.
	out, err := AddFlagsToInit(read(t, "../cli/mothership/stats.go"), []string{
		`statsCommand.Flags().String("orgName", "", "")`,
	})
	if err != nil {
		t.Fatal(err)
	}
	mustParse(t, "AddFlagsToInit/new", out)
	if !strings.Contains(string(out), "func init() {") {
		t.Error("no init() created")
	}

	// connect.go already has one, so the new line must join it rather than
	// declaring a second init.
	out, err = AddFlagsToInit(read(t, "../cli/mothership/connect.go"), []string{
		`connectCommand.Flags().String("extra", "", "")`,
	})
	if err != nil {
		t.Fatal(err)
	}
	mustParse(t, "AddFlagsToInit/existing", out)
	if n := strings.Count(string(out), "func init() {"); n != 1 {
		t.Errorf("%d init functions, want 1", n)
	}
	if !strings.Contains(string(out), `connectCommand.Flags().String("extra"`) {
		t.Error("flag not added")
	}
	// The existing registrations survive.
	if !strings.Contains(string(out), `searchConnectedCommand.Flags().String("token", "", "JWT of the user")`) {
		t.Error("existing flag registrations lost")
	}
}

func TestAddSubcommandToFactory(t *testing.T) {
	// mothership.go collects subcommands in its Command() factory.
	out, err := AddSubcommand(read(t, "../cli/mothership/mothership.go"), "mothershipCommand", "poolCommand")
	if err != nil {
		t.Fatal(err)
	}
	mustParse(t, "AddSubcommand/factory", out)

	got := string(out)
	if !strings.Contains(got, "mothershipCommand.AddCommand(poolCommand)") {
		t.Errorf("subcommand not registered:\n%s", got)
	}
	if !strings.Contains(got, "mothershipCommand.AddCommand(statsCommand)") {
		t.Error("existing registrations lost")
	}
	if !strings.Contains(got, "return mothershipCommand") {
		t.Error("factory return lost")
	}

	if again, _ := AddSubcommand(out, "mothershipCommand", "poolCommand"); string(again) != got {
		t.Error("registering twice changed the file")
	}
}

func TestAddSubcommandToInit(t *testing.T) {
	// agent.go collects its subcommand in init().
	out, err := AddSubcommand(read(t, "../cli/mothership/agent.go"), "agentCommand", "agentListCommand")
	if err != nil {
		t.Fatal(err)
	}
	mustParse(t, "AddSubcommand/init", out)
	got := string(out)
	if !strings.Contains(got, "agentCommand.AddCommand(agentListCommand)") {
		t.Errorf("subcommand not registered:\n%s", got)
	}
	if !strings.Contains(got, "agentCommand.AddCommand(agentDeleteCommand)") {
		t.Error("existing registration lost")
	}
}

// TestAddSubcommandToParentWithoutChildren covers a leaf command becoming a
// group: statsCommand has no AddCommand call anywhere.
func TestAddSubcommandToParentWithoutChildren(t *testing.T) {
	out, err := AddSubcommand(read(t, "../cli/mothership/stats.go"), "statsCommand", "statsHistoryCommand")
	if err != nil {
		t.Fatal(err)
	}
	mustParse(t, "AddSubcommand/leaf", out)
	if !strings.Contains(string(out), "statsCommand.AddCommand(statsHistoryCommand)") {
		t.Errorf("subcommand not registered:\n%s", out)
	}
}

func TestAddPushField(t *testing.T) {
	src := read(t, "../service.go")

	out, err := AddPushField(src, "AgentPoolId")
	if err != nil {
		t.Fatal(err)
	}
	mustParse(t, "AddPushField", out)

	got := string(out)
	if !strings.Contains(got, "AgentPoolId: arg.AgentPoolId,") {
		t.Errorf("field not added:\n%s", got)
	}
	// Every pre-existing field must survive, or a push starts dropping state.
	for _, want := range []string{"Jwt:", "OrgName:", "AgentHelloToken:", "Timestamp:"} {
		if !strings.Contains(got, want) {
			t.Errorf("existing field %s lost", want)
		}
	}

	again, err := AddPushField(out, "AgentPoolId")
	if err != nil {
		t.Fatal(err)
	}
	if string(again) != got {
		t.Error("adding the same field twice changed the file")
	}
}

func TestAddClientConstructor(t *testing.T) {
	src := read(t, "../clients/clients.go")

	out, err := AddClientConstructor(src, "cosmoClient", "NewCosmoClient")
	if err != nil {
		t.Fatal(err)
	}
	mustParse(t, "AddClientConstructor", out)

	got := string(out)
	if !strings.Contains(got, "func NewCosmoClient(cmd *cobra.Command) (*cosmoClient.ClientWithResponses, error)") {
		t.Errorf("constructor not added:\n%s", got)
	}
	if !strings.Contains(got, "command.ServerCosmo(cmd)") {
		t.Error("wrong base-URL helper")
	}
	if !strings.Contains(got, `"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/clients/cosmoClient"`) {
		t.Error("client package not imported")
	}
	if !strings.Contains(got, "func NewMothershipClient") {
		t.Error("existing constructors lost")
	}

	// idp's existing constructor is spelled NewIDPClient and uses ServerIDP;
	// re-adding it must be a no-op rather than a second spelling.
	again, err := AddClientConstructor(out, "idpClient", "NewIDPClient")
	if err != nil {
		t.Fatal(err)
	}
	if string(again) != got {
		t.Error("re-adding an existing constructor changed the file")
	}
}

func TestServiceGroupFileAndRegistration(t *testing.T) {
	group, err := ServiceGroupFile("dali", "", "")
	if err != nil {
		t.Fatal(err)
	}
	mustParse(t, "ServiceGroupFile", group)
	got := string(group)
	for _, want := range []string{
		"package dali",
		"func Command() *cobra.Command {",
		`daliCommand.PersistentFlags().String("server"`,
		"return daliCommand",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}

	out, err := RegisterService(read(t, "../../cmd/root.go"), "dali")
	if err != nil {
		t.Fatal(err)
	}
	mustParse(t, "RegisterService", out)
	rooted := string(out)
	if !strings.Contains(rooted, "rootCmd.AddCommand(dali.Command())") {
		t.Errorf("service not registered:\n%s", rooted)
	}
	if !strings.Contains(rooted, `"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/cli/dali"`) {
		t.Error("service package not imported")
	}
	if !strings.Contains(rooted, "rootCmd.AddCommand(mothership.Command())") {
		t.Error("existing registrations lost")
	}
}
