package goterramodel

import "go.mongodb.org/mongo-driver/bson/primitive"

// Model defines a set of VM which can be used to generate some terraform templates for openstack, ...
type Model struct {
	Name             string `json:"name"`
	Count            int64  `json:"count"`
	PublicIP         string `json:"public_ip"`
	EphemeralStorage string `json:"ephemeral_disk"`
	SharedStorage    string `json:"shared_storage"`
}

// Recipe describe a recipe for an app
type Recipe struct {
	ID            primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	Script        string              `json:"script"`
	Public        bool                `json:"public"`
	Namespace     string              `json:"namespace"`
	BaseImages    []string            `json:"base"`
	ParentRecipe  string              `json:"parent"`
	Timestamp     int64               `json:"ts"`
	Previous      string              `json:"prev"`   // Previous recipe id, for versioning
	Inputs        map[string]string   `json:"inputs"` // List of input variables needed when executing at app for this recipe, those variables should be sent as env_XX if XX is in requires: varname,label
	Tags          []string            `json:"tags"`
	Remote        string              `json:"remote"` // path in git repo
	RemoteVersion string              `json:"rversion"`
	Version       string              `json:"version"`
	Frozen        bool                `json:"frozen"`
	Defaults      map[string][]string `json:"defaults"`
}

// Template represents a terraform template
type Template struct {
	ID            primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Namespace     string              `json:"namespace"`
	Timestamp     int64               `json:"ts"`
	Public        bool                `json:"public"`
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	Data          map[string]string   `json:"data"` // map of cloud kind / terraform template
	Model         []Model             `json:"model"`
	Inputs        map[string]string   `json:"inputs"` // expected inputs varname, label
	Previous      string              `json:"prev"`   // Previous recipe id, for versioning
	Tags          []string            `json:"tags"`
	Remote        string              `json:"remote"`   // name of template in repo (dir)
	RemoteVersion string              `json:"rversion"` // version of template in repo (subdir)
	Version       string              `json:"version"`
	VarRecipes    []string            `json:"varrecipes"` // Name of variables containing recipe list
	Frozen        bool                `json:"frozen"`
	Defaults      map[string][]string `json:"defaults"`
}

// Application descripe an app to deploy
type Application struct {
	ID              primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	Public          bool                `json:"public"`
	Recipes         []string            `json:"recipes"` // recipe ids
	Namespace       string              `json:"namespace"`
	Template        string              `json:"template"` // template id
	Image           []string            `json:"image"`
	Timestamp       int64               `json:"ts"`
	Previous        string              `json:"prev"`            // Previous app id, for versioning
	TemplateRecipes map[string][]string `json:"templaterecipes"` // Recipes per template recipe variable
	Remote          string              `json:"remote"`          // name of template in repo (dir)
	RemoteVersion   string              `json:"rversion"`        // version of template in repo (subdir)
	Version         string              `json:"version"`
	Frozen          bool                `json:"frozen"`
	Defaults        map[string][]string `json:"defaults"`
}

// Event represent an action (deploy, destroy, etc.) on a run (historical data)
type Event struct {
	TS      int64  `json:"ts"`
	Action  string `json:"action"`
	Success bool   `json:"success"`
}

// EndPointSecret contains credentials (password crypted) for a user/endpoint combination
type EndPointSecret struct {
	UID       string `json:"uid"`
	UserName  string `json:"name"`
	Password  string `json:"password"`
	EndPoint  string `json:"endpoint"`
	Namespace string `json:"namespace"`
}

// Run represents a deployment info for an app
type Run struct {
	Name            string             `json:"name"`
	Description     string             `json:"description"`
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AppID           string             `json:"appID"` // Application id
	Inputs          map[string]string  `json:"inputs"`
	SensitiveInputs map[string]string  `json:"secretinputs"` // Secret variables (password, etc.) will be given to terraform as env variables
	Status          string             `json:"status"`
	Endpoint        string             `json:"endpoint"`
	Namespace       string             `json:"namespace"`
	UID             string
	Start           int64   `json:"start"`
	End             int64   `json:"end"`
	Duration        float64 `json:"duration"`
	Outputs         string  `json:"outputs"`
	Error           string  `json:"error"`
	Deployment      string  `json:"deployment"`
	Events          []Event `json:"events"`
}

// Openstack maps to openstack provider in openstack
type Openstack struct {
	UserName          string `json:"user_name"`
	Password          string `json:"password"`
	Flavor            string `json:"flavor_name"`
	KeyPair           string `json:"key_pair"`
	TenantName        string `json:"tenant_name"`
	TenantID          string `json:"tenant_id"`
	AuthURL           string `json:"auth_url"`
	Region            string `json:"region"`
	DomainName        string `json:"domain_name,omitempty"`
	DomainID          string `json:"domain_id,omitempty"`
	ProjectDomainID   string `json:"project_domain_id,omitempty"`
	ProjectDomainName string `json:"project_domain_name,omitempty"`
	UserDomainID      string `json:"user_domain_id,omitempty"`
	UserDomainName    string `json:"user_domain_name,omitempty"`

	Inputs map[string]string `json:"inputs"` // expected inputs (credentials, ...), varname, label
}

// EndPoint specifies a cloud endpoint data
type EndPoint struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name"`
	Kind      string             `json:"kind"` // openstack, etc.
	Namespace string             `json:"namespace"`
	// Openstack Openstack          `json:"openstack"` // for Kind=openstack
	Features map[string]string `json:"features"`
	Inputs   map[string]string `json:"inputs"` // expected inputs varname, label
	Config   map[string]string `json:"config"` // Preset some inputs like endpoints url, ... to be set in terraform variables
	Images   map[string]string `json:"images"` // map recipe image id to endpoint image id
	Public   bool              `json:"public"`

	Timestamp int64               `json:"ts"`
	Remote    string              `json:"remote"` // name of template in repo (dir)
	Defaults  map[string][]string `json:"defaults"`
}

// EndpointDefaults defines user data for endpoint id
type EndpointDefaults struct {
	UID       string              `json:"uid"`
	Endpoint  string              `json:"endpoint"`
	Namespace string              `json:"ns"`
	Defaults  map[string][]string `json:"defaults"`
}

// RunAction is message struct to be sent to the run component
// action: apply or destroy
// id: identifier of the run
type RunAction struct {
	Action  string            `json:"action"`
	ID      string            `json:"id"`
	Secrets map[string]string `json:"secrets"`
}

// UserAction is message sent to track events to gotevent exchange..
type UserAction struct {
	Action string
	UID    string
	Data   string
}

// NSData represent a namespace
type NSData struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name    string             `json:"name"`
	Owners  []string           `json:"owners"`
	Members []string           `json:"members"`
}
