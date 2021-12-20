package grants

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type RoleResponse struct {
	Account struct {
		Address string `json:"address"`
		Name    string `json:"name"`
		Type    string `json:"type"`
	} `json:"account"`
	Grants []struct {
		Org struct {
			Address string `json:"address"`
			Name    string `json:"name"`
			Type    string `json:"type"`
		} `json:"org"`
		Roles       []string      `json:"roles"`
		AuthzGrants []string      `json:"authzGrants"`
		Apps        []interface{} `json:"apps"`
	} `json:"grants"`
}

type Grants struct {
	Orgs []Org `json:"orgs"`
}

type Org struct {
	Name        string   `json:"name"`
	Roles       []string `json:"roles"`
	AuthzGrants []string `json:"authzGrants"`
}

var (
	Client HTTPClient
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func init() {
	Client = &http.Client{}
}

func GetGrants(grantsURL string, address string) (*Grants, error) {
	uri := strings.ReplaceAll(grantsURL, "{addr}", address)
	roleReq, _ := http.NewRequest("GET", uri, nil)
	roleReq.Header.Add("x-sender", address)

	resp, err := Client.Do(roleReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var roleResponse RoleResponse
	if err := json.Unmarshal(body, &roleResponse); err != nil {
		fmt.Println("Can not unmarshal JSON")
		return nil, err
	}

	if len(roleResponse.Grants) == 0 {
		return nil, fmt.Errorf("account doesn't exist in role service")
	}

	var grants Grants
	for _, grant := range roleResponse.Grants {
		org := Org{
			Name:        grant.Org.Name,
			Roles:       grant.Roles,
			AuthzGrants: grant.AuthzGrants,
		}

		grants.Orgs = append(grants.Orgs, org)
	}
	return &grants, nil
}
