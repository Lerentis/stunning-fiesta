package utils

type App struct {
	Name               string `json:"applicationName"`
	Registry           string `json:"registry"`
	Project            string `json:"project"`
	IngressEnabeld     string `json:"ingressEnabled"`
	DatabaseEnabled    string `json:"databaseEnabled"`
	Cluster            string `json:"cluster"`
	ServiceNowOffering string `json:"serviceNowOffering"`
	ProjectSSH         string `json:"projectSSH"`
}
