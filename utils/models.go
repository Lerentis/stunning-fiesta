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

type TemplatesListResponse struct {
	Templates []string `json:"templates"`
}

type Stage struct {
	Name    string `json:"name"`
	Account string `json:"account"`
	Domain  string `json:"domain"`
}
type Cluster struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Stages      []Stage `json:"stages"`
}

type ClusterList struct {
	Clusters []Cluster `json:"clusters"`
}

type MountPath struct {
	ClusterName string `json:"clusterName"`
	NonProd     string `json:"nonprod"`
	Prod        string `json:"prod"`
}

type DBMS struct {
	Name       string      `json:"name"`
	MountPaths []MountPath `json:"mountPaths"`
}

type DBMSList struct {
	DBMS []DBMS `json:"dbms"`
}

type Project struct {
	Name               string   `json:"name"`
	Topic              string   `json:"topic"`
	ProjectSSH         string   `json:"project_ssh"`
	ProjectURL         string   `json:"project_url"`
	ProjectPath        string   `json:"project_path"`
	KubernetesRepoPath string   `json:"git_kubernetes_path"`
	DBMS               string   `json:"dbms"`
	Clusters           []string `json:"clusters"`
	Team               string   `json:"ownerTeam"`
}

type ProjectList struct {
	Projects []Project `json:"projects"`
}

type Team struct {
	Name string `json:"name"`
}

type TeamList struct {
	Teams []Team `json:"teams"`
}

type HelmData struct {
	Repository string `json:"repository"`
	Version    string `json:"version"`
}

type ApplicationTypes struct {
	Name       string   `json:"name"`
	Helm       HelmData `json:"helm"`
	HasService bool     `json:"hasService"`
}

type ApplicationTypesList struct {
	ApplicationTypes []ApplicationTypes `json:"applicationTypes"`
}

type VaultClusterData struct {
	ClusterName string `json:"clusterName"`
	LoginPath   string `json:"loginPath"`
	AuthRole    string `json:"authRole"`
}

type Vault struct {
	Account   string             `json:"account"`
	Addr      string             `json:"addr"`
	LoginPath string             `json:"loginPath"`
	AuthRole  string             `json:"authRole"`
	Clusters  []VaultClusterData `json:"clusters"`
}
type VaultList struct {
	Vaults []Vault `json:"vaults"`
}
