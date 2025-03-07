package gen

type Component string

const (
	Component_Clients Component = "clients"
	Component_Csproj  Component = "csproj"
	Component_Readme  Component = "readme"
	Component_Queries Component = "queries"
	Component_Client  Component = "client"
)

func Component_values() []Component {
	return []Component{Component_Clients, Component_Csproj, Component_Readme, Component_Queries, Component_Client}
}

func Component_stringValues() (stringValues []string) {
	for _, val := range Component_values() {
		stringValues = append(stringValues, string(val))
	}
	return
}
