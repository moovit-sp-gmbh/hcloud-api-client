package views

type View struct {
	Name              string
	HumanReadableName string
}

const (
	HCLOUD = iota
	DEBUG

	SERVICE
	SERVICE_IDP
	SERVICE_IDP_ACCOUNT
	SERVICE_IDP_ORGANIZATIONS
	SERVICE_HIGH5
	SERVICE_HIGH5_APPS

	CONFIG
	CONFIG_CONTEXT

	HELP
)

var Views = make(map[int]View)

func init() {
	Views[HCLOUD] = View{Name: "hcloud", HumanReadableName: "hcloud"}

	Views[SERVICE] = View{Name: "service", HumanReadableName: "Service"}
	Views[SERVICE_IDP] = View{Name: "idp", HumanReadableName: "IDP"}
	Views[SERVICE_IDP_ACCOUNT] = View{Name: "account", HumanReadableName: "Account"}
	Views[SERVICE_IDP_ORGANIZATIONS] = View{Name: "organizations", HumanReadableName: "Organizations"}
	Views[SERVICE_HIGH5] = View{Name: "high5", HumanReadableName: "High5"}
	Views[SERVICE_HIGH5_APPS] = View{Name: "apps", HumanReadableName: "Apps"}

	Views[CONFIG] = View{Name: "config", HumanReadableName: "Config"}
	Views[CONFIG_CONTEXT] = View{Name: "context", HumanReadableName: "Context"}

	Views[HELP] = View{Name: "help", HumanReadableName: "Help"}
	Views[DEBUG] = View{Name: "debug", HumanReadableName: "Debug"}
}
