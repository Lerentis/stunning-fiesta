package utils

import "encoding/xml"

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
	HelmTemplates           []string `json:"helm-templates"`
	ApplicationTemplates    []string `json:"application-templates"`
	InfrastructureTemplates []string `json:"infrastructure-templates"`
	NamespaceTemplates      []string `json:"namespace-templates"`
	Dependencies            string   `json:"dependencies"`
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

type Pom struct {
	XMLName      xml.Name        `xml:"project"`
	Xmlns        string          `xml:"xmlns,attr"`
	XmlnsXsi     string          `xml:"xmlns:xsi,attr"`
	XsiSchema    string          `xml:"xsi:schemaLocation,attr"`
	ModelVersion string          `xml:"modelVersion"`
	GroupID      string          `xml:"groupId"`
	ArtifactID   string          `xml:"artifactId"`
	Version      string          `xml:"version"`
	Dependencies []PomDependency `xml:"dependencies>dependency"`
	Build        *PomBuild       `xml:"build,omitempty"`
	Plugins      []PomPlugin     `xml:"-"`
}

type PomDependency struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
}

type PomBuild struct {
	Plugins []PomPlugin `xml:"plugins>plugin"`
}

type PomPlugin struct {
	GroupID       string                  `xml:"groupId"`
	ArtifactID    string                  `xml:"artifactId"`
	Version       string                  `xml:"version"`
	Executions    []PomPluginExecution    `xml:"executions>execution,omitempty"`
	Configuration *PomPluginConfiguration `xml:"configuration,omitempty"`
}

type PomPluginExecution struct {
	Phase string   `xml:"phase"`
	Goals []string `xml:"goals>goal"`
}

type PomPluginConfiguration map[string]interface{}
