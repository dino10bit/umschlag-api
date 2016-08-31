package model

// TeamOrgs is simply a collection of team org structs.
type TeamOrgs []*TeamOrg

// TeamOrg represents a team org model definition.
type TeamOrg struct {
	TeamID int    `json:"team_id" sql:"index"`
	Team   *Team  `json:"team,omitempty"`
	OrgID  int    `json:"org_id" sql:"index"`
	Org    *Org   `json:"org,omitempty"`
	Perm   string `json:"perm,omitempty"`
}
