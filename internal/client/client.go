package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/resnickio/unifi-go-sdk/pkg/unifi"
)

type Client struct {
	Network      *unifi.NetworkClient
	Site         string
	BaseURL      string
	APIKey       string
	IsStandalone bool
	HTTPClient   *retryablehttp.Client
}

func NewClient(host, username, password, apiKey, site string, insecure, isStandalone bool) (*Client, error) {
	if site == "" {
		site = "default"
	}

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	retryClient.RetryWaitMin = 1 * time.Second
	retryClient.RetryWaitMax = 30 * time.Second
	retryClient.Logger = nil

	tr := &http.Transport{}
	if insecure {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	retryClient.HTTPClient.Transport = tr

	cfg := unifi.NetworkClientConfig{
		BaseURL:            host,
		Site:               site,
		Username:           username,
		Password:           password,
		APIKey:             apiKey,
		InsecureSkipVerify: insecure,
	}

	networkClient, err := unifi.NewNetworkClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize unifi network client: %w", err)
	}

	c := &Client{
		Network:      networkClient,
		Site:         site,
		BaseURL:      host,
		APIKey:       apiKey,
		IsStandalone: isStandalone,
		HTTPClient:   retryClient,
	}

	if apiKey == "" {
		if err := c.Login(context.Background()); err != nil {
			return nil, fmt.Errorf("failed to login to unifi: %w", err)
		}
	}

	return c, nil
}

func (c *Client) Login(ctx context.Context) error {
	return c.Network.Login(ctx)
}

func (c *Client) doRequest(ctx context.Context, method, path string, body any, result any) error {
	reqURL := c.BaseURL + path

	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	req, err := retryablehttp.NewRequestWithContext(ctx, method, reqURL, bodyBytes)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.APIKey != "" {
		req.Header.Set("X-API-KEY", c.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unifi api error (status %d): %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		var apiResp struct {
			Meta struct {
				RC string `json:"rc"`
			} `json:"meta"`
			Data json.RawMessage `json:"data"`
		}

		bodyContent, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(bodyContent, &apiResp); err == nil && apiResp.Meta.RC != "" {
			return json.Unmarshal(apiResp.Data, result)
		}

		return json.Unmarshal(bodyContent, result)
	}

	return nil
}

func (c *Client) doREST(ctx context.Context, method, endpoint string, body, result any) error {
	path := "/api/s/" + url.PathEscape(c.Site) + "/rest/" + endpoint
	if !c.IsStandalone {
		path = "/proxy/network" + path
	}
	return c.doRequest(ctx, method, path, body, result)
}

func (c *Client) doV2(ctx context.Context, method, endpoint string, body, result any) error {
	path := "/v2/api/site/" + url.PathEscape(c.Site) + "/" + endpoint
	if !c.IsStandalone {
		path = "/proxy/network" + path
	}
	return c.doRequest(ctx, method, path, body, result)
}

// Helpers

func createResource[T any](ctx context.Context, c *Client, endpoint string, item *T) (*T, error) {
	var items []T
	if err := c.doREST(ctx, "POST", endpoint, item, &items); err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("empty response from unifi")
	}
	return &items[0], nil
}

func getResource[T any](ctx context.Context, c *Client, endpoint, id string) (*T, error) {
	var items []T
	if err := c.doREST(ctx, "GET", endpoint+"/"+id, nil, &items); err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("resource not found")
	}
	return &items[0], nil
}

func listResources[T any](ctx context.Context, c *Client, endpoint string) ([]T, error) {
	var items []T
	if err := c.doREST(ctx, "GET", endpoint, nil, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func updateResource[T any](ctx context.Context, c *Client, endpoint, id string, item *T) (*T, error) {
	var items []T
	if err := c.doREST(ctx, "PUT", endpoint+"/"+id, item, &items); err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return getResource[T](ctx, c, endpoint, id)
	}
	return &items[0], nil
}

func deleteResource(ctx context.Context, c *Client, endpoint, id string) error {
	return c.doREST(ctx, "DELETE", endpoint+"/"+id, nil, nil)
}

// Network CRUD

func (c *Client) CreateNetwork(ctx context.Context, network *unifi.Network) (*unifi.Network, error) {
	return createResource(ctx, c, "networkconf", network)
}

func (c *Client) GetNetwork(ctx context.Context, id string) (*unifi.Network, error) {
	return getResource[unifi.Network](ctx, c, "networkconf", id)
}

func (c *Client) ListNetworks(ctx context.Context) ([]unifi.Network, error) {
	return listResources[unifi.Network](ctx, c, "networkconf")
}

func (c *Client) UpdateNetwork(ctx context.Context, id string, network *unifi.Network) (*unifi.Network, error) {
	return updateResource(ctx, c, "networkconf", id, network)
}

func (c *Client) DeleteNetwork(ctx context.Context, id string) error {
	return deleteResource(ctx, c, "networkconf", id)
}

// FirewallRule CRUD

func (c *Client) CreateFirewallRule(ctx context.Context, rule *unifi.FirewallRule) (*unifi.FirewallRule, error) {
	return createResource(ctx, c, "firewallrule", rule)
}

func (c *Client) GetFirewallRule(ctx context.Context, id string) (*unifi.FirewallRule, error) {
	return getResource[unifi.FirewallRule](ctx, c, "firewallrule", id)
}

func (c *Client) UpdateFirewallRule(ctx context.Context, id string, rule *unifi.FirewallRule) (*unifi.FirewallRule, error) {
	return updateResource(ctx, c, "firewallrule", id, rule)
}

func (c *Client) DeleteFirewallRule(ctx context.Context, id string) error {
	return deleteResource(ctx, c, "firewallrule", id)
}

// PortProfile CRUD

func (c *Client) CreatePortProfile(ctx context.Context, profile *unifi.PortConf) (*unifi.PortConf, error) {
	return createResource(ctx, c, "portconf", profile)
}

func (c *Client) GetPortProfile(ctx context.Context, id string) (*unifi.PortConf, error) {
	return getResource[unifi.PortConf](ctx, c, "portconf", id)
}

func (c *Client) ListPortProfiles(ctx context.Context) ([]unifi.PortConf, error) {
	return listResources[unifi.PortConf](ctx, c, "portconf")
}

func (c *Client) UpdatePortProfile(ctx context.Context, id string, profile *unifi.PortConf) (*unifi.PortConf, error) {
	return updateResource(ctx, c, "portconf", id, profile)
}

func (c *Client) DeletePortProfile(ctx context.Context, id string) error {
	return deleteResource(ctx, c, "portconf", id)
}

// UserGroup CRUD

func (c *Client) CreateUserGroup(ctx context.Context, group *unifi.UserGroup) (*unifi.UserGroup, error) {
	return createResource(ctx, c, "usergroup", group)
}

func (c *Client) GetUserGroup(ctx context.Context, id string) (*unifi.UserGroup, error) {
	return getResource[unifi.UserGroup](ctx, c, "usergroup", id)
}

func (c *Client) ListUserGroups(ctx context.Context) ([]unifi.UserGroup, error) {
	return listResources[unifi.UserGroup](ctx, c, "usergroup")
}

func (c *Client) UpdateUserGroup(ctx context.Context, id string, group *unifi.UserGroup) (*unifi.UserGroup, error) {
	return updateResource(ctx, c, "usergroup", id, group)
}

func (c *Client) DeleteUserGroup(ctx context.Context, id string) error {
	return deleteResource(ctx, c, "usergroup", id)
}

// APGroup CRUD

type apGroupCreateRequest struct {
	Name        string   `json:"name"`
	DeviceMACs  []string `json:"device_macs"`
	ForWLANConf bool     `json:"for_wlanconf"`
}

func (c *Client) CreateAPGroup(ctx context.Context, group *unifi.APGroup) (*unifi.APGroup, error) {
	req := apGroupCreateRequest{
		Name:       group.Name,
		DeviceMACs: group.DeviceMACs,
	}
	if group.ForWLANConf != nil {
		req.ForWLANConf = *group.ForWLANConf
	}
	if req.DeviceMACs == nil {
		req.DeviceMACs = []string{}
	}

	var created unifi.APGroup
	err := c.doV2(ctx, "POST", "apgroups", req, &created)
	return &created, err
}

func (c *Client) GetAPGroup(ctx context.Context, id string) (*unifi.APGroup, error) {
	var group unifi.APGroup
	err := c.doV2(ctx, "GET", "apgroups/"+id, nil, &group)
	if err != nil {
		groups, _ := c.ListAPGroups(ctx)
		for _, g := range groups {
			if g.ID == id {
				return &g, nil
			}
		}
		return nil, err
	}
	return &group, nil
}

func (c *Client) ListAPGroups(ctx context.Context) ([]unifi.APGroup, error) {
	var groups []unifi.APGroup
	err := c.doV2(ctx, "GET", "apgroups", nil, &groups)
	return groups, err
}

func (c *Client) UpdateAPGroup(ctx context.Context, id string, group *unifi.APGroup) (*unifi.APGroup, error) {
	req := apGroupCreateRequest{
		Name:       group.Name,
		DeviceMACs: group.DeviceMACs,
	}
	if group.ForWLANConf != nil {
		req.ForWLANConf = *group.ForWLANConf
	}
	if req.DeviceMACs == nil {
		req.DeviceMACs = []string{}
	}

	var updated unifi.APGroup
	err := c.doV2(ctx, "PUT", "apgroups/"+id, req, &updated)
	return &updated, err
}

func (c *Client) DeleteAPGroup(ctx context.Context, id string) error {
	return c.doV2(ctx, "DELETE", "apgroups/"+id, nil, nil)
}

// WLAN CRUD

func (c *Client) CreateWLAN(ctx context.Context, wlan *unifi.WLANConf) (*unifi.WLANConf, error) {
	return createResource(ctx, c, "wlanconf", wlan)
}

func (c *Client) GetWLAN(ctx context.Context, id string) (*unifi.WLANConf, error) {
	return getResource[unifi.WLANConf](ctx, c, "wlanconf", id)
}

func (c *Client) ListWLANs(ctx context.Context) ([]unifi.WLANConf, error) {
	return listResources[unifi.WLANConf](ctx, c, "wlanconf")
}

func (c *Client) UpdateWLAN(ctx context.Context, id string, wlan *unifi.WLANConf) (*unifi.WLANConf, error) {
	return updateResource(ctx, c, "wlanconf", id, wlan)
}

func (c *Client) DeleteWLAN(ctx context.Context, id string) error {
	return deleteResource(ctx, c, "wlanconf", id)
}

// FirewallGroup CRUD

func (c *Client) CreateFirewallGroup(ctx context.Context, group *unifi.FirewallGroup) (*unifi.FirewallGroup, error) {
	return createResource(ctx, c, "firewallgroup", group)
}

func (c *Client) GetFirewallGroup(ctx context.Context, id string) (*unifi.FirewallGroup, error) {
	return getResource[unifi.FirewallGroup](ctx, c, "firewallgroup", id)
}

func (c *Client) ListFirewallGroups(ctx context.Context) ([]unifi.FirewallGroup, error) {
	return listResources[unifi.FirewallGroup](ctx, c, "firewallgroup")
}

func (c *Client) UpdateFirewallGroup(ctx context.Context, id string, group *unifi.FirewallGroup) (*unifi.FirewallGroup, error) {
	return updateResource(ctx, c, "firewallgroup", id, group)
}

func (c *Client) DeleteFirewallGroup(ctx context.Context, id string) error {
	return deleteResource(ctx, c, "firewallgroup", id)
}

type User struct {
	ID          string `json:"_id,omitempty"`
	MAC         string `json:"mac"`
	Name        string `json:"name,omitempty"`
	Note        string `json:"note,omitempty"`
	UseFixedIP  *bool  `json:"use_fixedip,omitempty"`
	FixedIP     string `json:"fixed_ip,omitempty"`
	NetworkID   string `json:"network_id,omitempty"`
	UsergroupID string `json:"usergroup_id,omitempty"`
	Blocked     *bool  `json:"blocked,omitempty"`
	IsWired     *bool  `json:"is_wired,omitempty"`
	IsGuest     *bool  `json:"is_guest,omitempty"`
	OUI         string `json:"oui,omitempty"`
	Noted       *bool  `json:"noted,omitempty"`
	SiteID      string `json:"site_id,omitempty"`
}

func (c *Client) CreateUser(ctx context.Context, user *User) (*User, error) {
	var users []User
	err := c.doREST(ctx, "POST", "user", user, &users)
	if err != nil {
		return nil, err
	}
	return &users[0], nil
}

func (c *Client) GetUser(ctx context.Context, id string) (*User, error) {
	var users []User
	err := c.doREST(ctx, "GET", "user/"+id, nil, &users)
	if err != nil {
		return nil, err
	}
	return &users[0], nil
}

func (c *Client) UpdateUser(ctx context.Context, id string, user *User) (*User, error) {
	var users []User
	err := c.doREST(ctx, "PUT", "user/"+id, user, &users)
	if err != nil {
		return nil, err
	}
	return &users[0], nil
}

func (c *Client) DeleteUser(ctx context.Context, mac string) error {
	payload := map[string]any{
		"cmd":  "forget-sta",
		"macs": []string{mac},
	}
	path := "/api/s/" + url.PathEscape(c.Site) + "/cmd/stamgr"
	if !c.IsStandalone {
		path = "/proxy/network" + path
	}
	return c.doRequest(ctx, "POST", path, payload, nil)
}

// RADIUSProfile CRUD

func (c *Client) CreateRADIUSProfile(ctx context.Context, profile *unifi.RADIUSProfile) (*unifi.RADIUSProfile, error) {
	return createResource(ctx, c, "radiusprofile", profile)
}

func (c *Client) GetRADIUSProfile(ctx context.Context, id string) (*unifi.RADIUSProfile, error) {
	return getResource[unifi.RADIUSProfile](ctx, c, "radiusprofile", id)
}

func (c *Client) ListRADIUSProfiles(ctx context.Context) ([]unifi.RADIUSProfile, error) {
	return listResources[unifi.RADIUSProfile](ctx, c, "radiusprofile")
}

func (c *Client) UpdateRADIUSProfile(ctx context.Context, id string, profile *unifi.RADIUSProfile) (*unifi.RADIUSProfile, error) {
	return updateResource(ctx, c, "radiusprofile", id, profile)
}

func (c *Client) DeleteRADIUSProfile(ctx context.Context, id string) error {
	return deleteResource(ctx, c, "radiusprofile", id)
}

// PortForward CRUD

func (c *Client) CreatePortForward(ctx context.Context, forward *unifi.PortForward) (*unifi.PortForward, error) {
	return createResource(ctx, c, "portforward", forward)
}

func (c *Client) GetPortForward(ctx context.Context, id string) (*unifi.PortForward, error) {
	return getResource[unifi.PortForward](ctx, c, "portforward", id)
}

func (c *Client) UpdatePortForward(ctx context.Context, id string, forward *unifi.PortForward) (*unifi.PortForward, error) {
	return updateResource(ctx, c, "portforward", id, forward)
}

func (c *Client) DeletePortForward(ctx context.Context, id string) error {
	return deleteResource(ctx, c, "portforward", id)
}

// Routing CRUD

func (c *Client) CreateStaticRoute(ctx context.Context, route *unifi.Routing) (*unifi.Routing, error) {
	req := map[string]any{
		"name":                  route.Name,
		"type":                  "static-route",
		"enabled":               true,
		"static-route_network":  route.StaticRouteNetwork,
		"static-route_nexthop":  route.StaticRouteNexthop,
		"static-route_type":     "nexthop-route",
		"static-route_distance": 1,
	}
	if route.Enabled != nil {
		req["enabled"] = *route.Enabled
	}
	if route.StaticRouteDistance != nil {
		req["static-route_distance"] = *route.StaticRouteDistance
	}

	var routes []unifi.Routing
	err := c.doREST(ctx, "POST", "routing", req, &routes)
	if err != nil {
		return nil, err
	}
	return &routes[0], nil
}

func (c *Client) GetStaticRoute(ctx context.Context, id string) (*unifi.Routing, error) {
	return getResource[unifi.Routing](ctx, c, "routing", id)
}

func (c *Client) UpdateStaticRoute(ctx context.Context, id string, route *unifi.Routing) (*unifi.Routing, error) {
	req := map[string]any{
		"_id":                   id,
		"name":                  route.Name,
		"type":                  "static-route",
		"enabled":               true,
		"static-route_network":  route.StaticRouteNetwork,
		"static-route_nexthop":  route.StaticRouteNexthop,
		"static-route_type":     "nexthop-route",
		"static-route_distance": 1,
	}
	if route.Enabled != nil {
		req["enabled"] = *route.Enabled
	}
	if route.StaticRouteDistance != nil {
		req["static-route_distance"] = *route.StaticRouteDistance
	}

	var routes []unifi.Routing
	err := c.doREST(ctx, "PUT", "routing/"+id, req, &routes)
	if err != nil {
		return nil, err
	}
	if len(routes) == 0 {
		return c.GetStaticRoute(ctx, id)
	}
	return &routes[0], nil
}

func (c *Client) DeleteStaticRoute(ctx context.Context, id string) error {
	return deleteResource(ctx, c, "routing", id)
}

// StaticDNS CRUD

func (c *Client) CreateStaticDNS(ctx context.Context, record *unifi.StaticDNS) (*unifi.StaticDNS, error) {
	req := map[string]any{
		"key":         record.Key,
		"value":       record.Value,
		"record_type": record.RecordType,
		"enabled":     true,
		"ttl":         0,
		"port":        0,
		"priority":    0,
		"weight":      0,
	}
	if req["record_type"] == "" {
		req["record_type"] = "A"
	}
	if record.Enabled != nil {
		req["enabled"] = *record.Enabled
	}
	if record.TTL != nil && *record.TTL > 0 {
		req["ttl"] = *record.TTL
	}
	if record.RecordType == "SRV" || record.RecordType == "MX" {
		if record.Port != nil {
			req["port"] = *record.Port
		}
		if record.Priority != nil {
			req["priority"] = *record.Priority
		}
		if record.Weight != nil {
			req["weight"] = *record.Weight
		}
	}

	var created unifi.StaticDNS
	err := c.doV2(ctx, "POST", "static-dns", req, &created)
	return &created, err
}

func (c *Client) GetStaticDNS(ctx context.Context, id string) (*unifi.StaticDNS, error) {
	var record unifi.StaticDNS
	err := c.doV2(ctx, "GET", "static-dns/"+id, nil, &record)
	if err != nil {
		records, _ := c.ListStaticDNS(ctx)
		for _, r := range records {
			if r.ID == id {
				return &r, nil
			}
		}
		return nil, err
	}
	return &record, nil
}

func (c *Client) ListStaticDNS(ctx context.Context) ([]unifi.StaticDNS, error) {
	var records []unifi.StaticDNS
	err := c.doV2(ctx, "GET", "static-dns", nil, &records)
	return records, err
}

func (c *Client) UpdateStaticDNS(ctx context.Context, id string, record *unifi.StaticDNS) (*unifi.StaticDNS, error) {
	req := map[string]any{
		"_id":         id,
		"key":         record.Key,
		"value":       record.Value,
		"record_type": record.RecordType,
		"enabled":     true,
		"ttl":         0,
		"port":        0,
		"priority":    0,
		"weight":      0,
	}
	if req["record_type"] == "" {
		req["record_type"] = "A"
	}
	if record.Enabled != nil {
		req["enabled"] = *record.Enabled
	}
	if record.TTL != nil && *record.TTL > 0 {
		req["ttl"] = *record.TTL
	}
	if record.RecordType == "SRV" || record.RecordType == "MX" {
		if record.Port != nil {
			req["port"] = *record.Port
		}
		if record.Priority != nil {
			req["priority"] = *record.Priority
		}
		if record.Weight != nil {
			req["weight"] = *record.Weight
		}
	}

	var updated unifi.StaticDNS
	err := c.doV2(ctx, "PUT", "static-dns/"+id, req, &updated)
	return &updated, err
}

func (c *Client) DeleteStaticDNS(ctx context.Context, id string) error {
	return c.doV2(ctx, "DELETE", "static-dns/"+id, nil, nil)
}

// TrafficRule CRUD

func (c *Client) CreateTrafficRule(ctx context.Context, rule *unifi.TrafficRule) (*unifi.TrafficRule, error) {
	req := map[string]any{
		"name":            rule.Name,
		"action":          rule.Action,
		"matching_target": rule.MatchingTarget,
		"description":     rule.Description,
		"enabled":         true,
		"target_devices":  rule.TargetDevices,
	}
	if rule.Enabled != nil {
		req["enabled"] = *rule.Enabled
	}
	if len(rule.TargetDevices) == 0 {
		req["target_devices"] = []unifi.TrafficRuleTarget{{Type: "ALL_CLIENTS"}}
	}

	var created unifi.TrafficRule
	err := c.doV2(ctx, "POST", "trafficrules", req, &created)
	return &created, err
}

func (c *Client) GetTrafficRule(ctx context.Context, id string) (*unifi.TrafficRule, error) {
	var rule unifi.TrafficRule
	err := c.doV2(ctx, "GET", "trafficrules/"+id, nil, &rule)
	if err != nil {
		rules, _ := c.ListTrafficRules(ctx)
		for _, r := range rules {
			if r.ID == id {
				return &r, nil
			}
		}
		return nil, err
	}
	return &rule, nil
}

func (c *Client) ListTrafficRules(ctx context.Context) ([]unifi.TrafficRule, error) {
	var rules []unifi.TrafficRule
	err := c.doV2(ctx, "GET", "trafficrules", nil, &rules)
	return rules, err
}

func (c *Client) UpdateTrafficRule(ctx context.Context, id string, rule *unifi.TrafficRule) (*unifi.TrafficRule, error) {
	req := map[string]any{
		"name":            rule.Name,
		"action":          rule.Action,
		"matching_target": rule.MatchingTarget,
		"description":     rule.Description,
		"enabled":         true,
		"target_devices":  rule.TargetDevices,
	}
	if rule.Enabled != nil {
		req["enabled"] = *rule.Enabled
	}
	if len(rule.TargetDevices) == 0 {
		req["target_devices"] = []unifi.TrafficRuleTarget{{Type: "ALL_CLIENTS"}}
	}

	var updated unifi.TrafficRule
	err := c.doV2(ctx, "PUT", "trafficrules/"+id, req, &updated)
	return &updated, err
}

func (c *Client) DeleteTrafficRule(ctx context.Context, id string) error {
	return c.doV2(ctx, "DELETE", "trafficrules/"+id, nil, nil)
}
