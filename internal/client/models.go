package client

import (
	"encoding/json"
)

// WANProviderCapabilities describes ISP bandwidth capabilities for a WAN network.
type WANProviderCapabilities struct {
	DownloadKilobitsPerSecond *int `json:"download_kilobits_per_second,omitempty"`
	UploadKilobitsPerSecond   *int `json:"upload_kilobits_per_second,omitempty"`
}

// QoSProfile describes Quality of Service settings for a port profile.
type QoSProfile struct {
	QoSPolicies    []json.RawMessage `json:"qos_policies,omitempty"`
	QoSProfileMode string            `json:"qos_profile_mode,omitempty"`
}

// NetworkVLAN contains VLAN configuration for a network.
type NetworkVLAN struct {
	VLAN        *int   `json:"vlan,omitempty"`
	VLANEnabled *bool  `json:"vlan_enabled,omitempty"`
	IPSubnet    string `json:"ip_subnet,omitempty"`
}

// NetworkDHCPGateway contains DHCP gateway override settings.
type NetworkDHCPGateway struct {
	DHCPDGatewayEnabled *bool  `json:"dhcpd_gateway_enabled,omitempty"`
	DHCPDGateway        string `json:"dhcpd_gateway,omitempty"`
}

// NetworkDHCPDNS contains DHCP DNS server settings.
type NetworkDHCPDNS struct {
	DHCPDDNSEnabled *bool  `json:"dhcpd_dns_enabled,omitempty"`
	DHCPDDns1       string `json:"dhcpd_dns_1,omitempty"`
	DHCPDDns2       string `json:"dhcpd_dns_2,omitempty"`
	DHCPDDns3       string `json:"dhcpd_dns_3,omitempty"`
	DHCPDDns4       string `json:"dhcpd_dns_4,omitempty"`
}

// NetworkDHCPBoot contains DHCP boot/PXE settings.
type NetworkDHCPBoot struct {
	DHCPDBootEnabled  *bool  `json:"dhcpd_boot_enabled,omitempty"`
	DHCPDBootServer   string `json:"dhcpd_boot_server,omitempty"`
	DHCPDBootFilename string `json:"dhcpd_boot_filename,omitempty"`
	DHCPDTFTPServer   string `json:"dhcpd_tftp_server,omitempty"`
}

// NetworkDHCPNTP contains DHCP NTP server settings.
type NetworkDHCPNTP struct {
	DHCPDNTPEnabled *bool  `json:"dhcpd_ntp_enabled,omitempty"`
	DHCPDNtp1       string `json:"dhcpd_ntp_1,omitempty"`
	DHCPDNtp2       string `json:"dhcpd_ntp_2,omitempty"`
}

// NetworkDHCP contains all DHCP-related configuration for a network.
type NetworkDHCP struct {
	DHCPDEnabled           *bool  `json:"dhcpd_enabled,omitempty"`
	DHCPDStart             string `json:"dhcpd_start,omitempty"`
	DHCPDStop              string `json:"dhcpd_stop,omitempty"`
	DHCPDLeasetime         *int   `json:"dhcpd_leasetime,omitempty"`
	DHCPRelayEnabled       *bool  `json:"dhcp_relay_enabled,omitempty"`
	DHCPDTimeOffsetEnabled *bool  `json:"dhcpd_time_offset_enabled,omitempty"`
	DHCPDUnifiController   string `json:"dhcpd_unifi_controller,omitempty"`
	DHCPDWPADUrl           string `json:"dhcpd_wpad_url,omitempty"`
	DHCPGuardingEnabled    *bool  `json:"dhcpguard_enabled,omitempty"`
	NetworkDHCPGateway
	NetworkDHCPDNS
	NetworkDHCPBoot
	NetworkDHCPNTP
}

// NetworkWANIPv6 contains WAN IPv6-specific settings.
type NetworkWANIPv6 struct {
	WANTypeV6            string `json:"wan_type_v6,omitempty"`
	WANIPv6DNS1          string `json:"wan_ipv6_dns1,omitempty"`
	WANIPv6DNS2          string `json:"wan_ipv6_dns2,omitempty"`
	WANIPv6DNSPreference string `json:"wan_ipv6_dns_preference,omitempty"`
	WANDHCPv6Cos         *int   `json:"wan_dhcpv6_cos,omitempty"`
	WANDHCPv6PDSizeAuto  *bool  `json:"wan_dhcpv6_pd_size_auto,omitempty"`
}

// NetworkWANQoS contains WAN Quality of Service settings.
type NetworkWANQoS struct {
	WANSmartQEnabled *bool  `json:"wan_smartq_enabled,omitempty"`
	WANEgressQOS     string `json:"wan_egress_qos,omitempty"`
	WANDHCPCos       *int   `json:"wan_dhcp_cos,omitempty"`
}

// NetworkWANLoadBalance contains WAN load balancing and failover settings.
type NetworkWANLoadBalance struct {
	WANFailoverPriority  *int   `json:"wan_failover_priority,omitempty"`
	WANLoadBalanceType   string `json:"wan_load_balance_type,omitempty"`
	WANLoadBalanceWeight *int   `json:"wan_load_balance_weight,omitempty"`
}

// NetworkWANVLAN contains WAN VLAN tagging settings.
type NetworkWANVLAN struct {
	WANVLANEnabled *bool `json:"wan_vlan_enabled,omitempty"`
	WANVLAN        *int  `json:"wan_vlan,omitempty"`
}

// NetworkWAN contains all WAN-specific configuration for a network.
type NetworkWAN struct {
	WAN                     string                   `json:"wan,omitempty"`
	WANType                 string                   `json:"wan_type,omitempty"`
	WANIP                   string                   `json:"wan_ip,omitempty"`
	WANNetmask              string                   `json:"wan_netmask,omitempty"`
	WANGateway              string                   `json:"wan_gateway,omitempty"`
	WANNetworkGroup         string                   `json:"wan_networkgroup,omitempty"`
	WANIPAliases            []string                 `json:"wan_ip_aliases,omitempty"`
	WANDNSPreference        string                   `json:"wan_dns_preference,omitempty"`
	WANDHCPOptions          []json.RawMessage        `json:"wan_dhcp_options,omitempty"`
	WANDsliteRemoteHost     string                   `json:"wan_dslite_remote_host,omitempty"`
	WANDsliteRemoteHostAuto *bool                    `json:"wan_dslite_remote_host_auto,omitempty"`
	WANProviderCapabilities *WANProviderCapabilities `json:"wan_provider_capabilities,omitempty"`
	ReportWANEvent          *bool                    `json:"report_wan_event,omitempty"`
	NetworkWANIPv6
	NetworkWANQoS
	NetworkWANLoadBalance
	NetworkWANVLAN
}

// NetworkIPv6 contains IPv6 configuration settings.
type NetworkIPv6 struct {
	IPv6SettingPreference string `json:"ipv6_setting_preference,omitempty"`
	IPv6WANDelegationType string `json:"ipv6_wan_delegation_type,omitempty"`
}

// NetworkMulticast contains multicast and IGMP settings.
type NetworkMulticast struct {
	IGMPSnooping      *bool  `json:"igmp_snooping,omitempty"`
	IGMPProxyUpstream *bool  `json:"igmp_proxy_upstream,omitempty"`
	IGMPProxyFor      string `json:"igmp_proxy_for,omitempty"`
	DomainName        string `json:"domain_name,omitempty"`
}

// NetworkAccess contains network access and NAT settings.
type NetworkAccess struct {
	InternetAccessEnabled     *bool    `json:"internet_access_enabled,omitempty"`
	IntraNetworkAccessEnabled *bool    `json:"intra_network_access_enabled,omitempty"`
	IsNAT                     *bool    `json:"is_nat,omitempty"`
	NATOutboundIPAddresses    []string `json:"nat_outbound_ip_addresses,omitempty"`
	MACOverrideEnabled        *bool    `json:"mac_override_enabled,omitempty"`
	MDNSEnabled               *bool    `json:"mdns_enabled,omitempty"`
	LteLANEnabled             *bool    `json:"lte_lan_enabled,omitempty"`
	UpnpLANEnabled            *bool    `json:"upnp_lan_enabled,omitempty"`
	PptpcServerEnabled        *bool    `json:"pptpc_server_enabled,omitempty"`
}

// NetworkRouting contains routing and firewall zone configuration.
type NetworkRouting struct {
	NetworkGroup     string `json:"networkgroup,omitempty"`
	RoutingTableID   *int   `json:"routing_table_id,omitempty"`
	SingleNetworkLAN string `json:"single_network_lan,omitempty"`
	FirewallZoneID   string `json:"firewall_zone_id,omitempty"`
}

// Network represents a UniFi network/VLAN configuration.
type Network struct {
	ID                string `json:"_id,omitempty"`
	SiteID            string `json:"site_id,omitempty"`
	Name              string `json:"name"`
	Purpose           string `json:"purpose,omitempty"`
	Enabled           *bool  `json:"enabled,omitempty"`
	SettingPreference string `json:"setting_preference,omitempty"`
	GatewayType       string `json:"gateway_type,omitempty"`
	GatewayDevice     string `json:"gateway_device,omitempty"`
	AutoScaleEnabled  *bool  `json:"auto_scale_enabled,omitempty"`
	AttrHiddenID      string `json:"attr_hidden_id,omitempty"`
	AttrNoDelete      *bool  `json:"attr_no_delete,omitempty"`
	NetworkVLAN
	NetworkDHCP
	NetworkWAN
	NetworkIPv6
	NetworkMulticast
	NetworkAccess
	NetworkRouting
}

// FirewallRule represents a UniFi firewall rule.
type FirewallRule struct {
	ID                    string   `json:"_id,omitempty"`
	SiteID                string   `json:"site_id,omitempty"`
	Name                  string   `json:"name"`
	Enabled               *bool    `json:"enabled,omitempty"`
	RuleIndex             *int     `json:"rule_index,omitempty"`
	Ruleset               string   `json:"ruleset,omitempty"`
	Action                string   `json:"action,omitempty"`
	Protocol              string   `json:"protocol,omitempty"`
	ProtocolMatchExcepted *bool    `json:"protocol_match_excepted,omitempty"`
	ProtocolV6            string   `json:"protocol_v6,omitempty"`
	ICMPTypename          string   `json:"icmp_typename,omitempty"`
	ICMPv6Typename        string   `json:"icmp_v6_typename,omitempty"`
	Logging               *bool    `json:"logging,omitempty"`
	StateEstablished      *bool    `json:"state_established,omitempty"`
	StateInvalid          *bool    `json:"state_invalid,omitempty"`
	StateNew              *bool    `json:"state_new,omitempty"`
	StateRelated          *bool    `json:"state_related,omitempty"`
	IPSec                 string   `json:"ipsec,omitempty"`
	SrcFirewallGroupIDs   []string `json:"src_firewallgroup_ids,omitempty"`
	SrcMACAddress         string   `json:"src_mac_address,omitempty"`
	SrcAddress            string   `json:"src_address,omitempty"`
	SrcNetworkConfID      string   `json:"src_networkconf_id,omitempty"`
	SrcNetworkConfType    string   `json:"src_networkconf_type,omitempty"`
	SrcPort               string   `json:"src_port,omitempty"`
	DstFirewallGroupIDs   []string `json:"dst_firewallgroup_ids,omitempty"`
	DstAddress            string   `json:"dst_address,omitempty"`
	DstNetworkConfID      string   `json:"dst_networkconf_id,omitempty"`
	DstNetworkConfType    string   `json:"dst_networkconf_type,omitempty"`
	DstPort               string   `json:"dst_port,omitempty"`
}

// FirewallGroup represents a UniFi firewall group.
type FirewallGroup struct {
	ID           string   `json:"_id,omitempty"`
	SiteID       string   `json:"site_id,omitempty"`
	Name         string   `json:"name"`
	GroupType    string   `json:"group_type,omitempty"`
	GroupMembers []string `json:"group_members,omitempty"`
}

// PortForward represents a UniFi port forwarding rule.
type PortForward struct {
	ID                 string   `json:"_id,omitempty"`
	SiteID             string   `json:"site_id,omitempty"`
	Name               string   `json:"name"`
	Enabled            *bool    `json:"enabled,omitempty"`
	PfwdInterface      string   `json:"pfwd_interface,omitempty"`
	Proto              string   `json:"proto,omitempty"`
	Src                string   `json:"src,omitempty"`
	DstPort            string   `json:"dst_port,omitempty"`
	Fwd                string   `json:"fwd,omitempty"`
	FwdPort            string   `json:"fwd_port,omitempty"`
	Log                *bool    `json:"log,omitempty"`
	DestinationIP      string   `json:"destination_ip,omitempty"`
	DestinationIPs     []string `json:"destination_ips,omitempty"`
	SrcLimitingEnabled *bool    `json:"src_limiting_enabled,omitempty"`
}

// APGroup represents an access point group.
type APGroup struct {
	ID           string   `json:"_id,omitempty"`
	Name         string   `json:"name"`
	AttrHiddenID string   `json:"attr_hidden_id,omitempty"`
	AttrNoDelete *bool    `json:"attr_no_delete,omitempty"`
	DeviceMACs   []string `json:"device_macs,omitempty"`
	ForWLANConf  *bool    `json:"for_wlanconf,omitempty"`
}

// WLANConf represents a UniFi wireless network (SSID) configuration.
type WLANConf struct {
	ID                          string            `json:"_id,omitempty"`
	SiteID                      string            `json:"site_id,omitempty"`
	Name                        string            `json:"name"`
	Enabled                     *bool             `json:"enabled,omitempty"`
	Security                    string            `json:"security,omitempty"`
	WPAMode                     string            `json:"wpa_mode,omitempty"`
	WPAEnc                      string            `json:"wpa_enc,omitempty"`
	WPA3Support                 *bool             `json:"wpa3_support,omitempty"`
	WPA3Transition              *bool             `json:"wpa3_transition,omitempty"`
	WPA3Enhanced192             *bool             `json:"wpa3_enhanced_192,omitempty"`
	WPA3FastRoaming             *bool             `json:"wpa3_fast_roaming,omitempty"`
	XPassphrase                 string            `json:"x_passphrase,omitempty"`
	XIappKey                    string            `json:"x_iapp_key,omitempty"`
	PassphraseAutogenerated     *bool             `json:"passphrase_autogenerated,omitempty"`
	PrivatePresharedKeys        []json.RawMessage `json:"private_preshared_keys,omitempty"`
	PrivatePresharedKeysEnabled *bool             `json:"private_preshared_keys_enabled,omitempty"`
	NetworkConfID               string            `json:"networkconf_id,omitempty"`
	Usergroup                   string            `json:"usergroup_id,omitempty"`
	IsGuest                     *bool             `json:"is_guest,omitempty"`
	HideSsid                    *bool             `json:"hide_ssid,omitempty"`
	WLANBand                    string            `json:"wlan_band,omitempty"`
	WLANBands                   []string          `json:"wlan_bands,omitempty"`
	APGroupIDs                  []string          `json:"ap_group_ids,omitempty"`
	APGroupMode                 string            `json:"ap_group_mode,omitempty"`
	Vlan                        *int              `json:"vlan,omitempty"`
	VlanEnabled                 *bool             `json:"vlan_enabled,omitempty"`
	MacFilterEnabled            *bool             `json:"mac_filter_enabled,omitempty"`
	MacFilterList               []string          `json:"mac_filter_list,omitempty"`
	MacFilterPolicy             string            `json:"mac_filter_policy,omitempty"`
	RadiusProfileID             string            `json:"radiusprofile_id,omitempty"`
	RadiusDasEnabled            *bool             `json:"radius_das_enabled,omitempty"`
	RadiusMacAuthEnabled        *bool             `json:"radius_mac_auth_enabled,omitempty"`
	RadiusMacaclFormat          string            `json:"radius_macacl_format,omitempty"`
	ScheduleEnabled             *bool             `json:"schedule_enabled,omitempty"`
	Schedule                    []string          `json:"schedule,omitempty"`
	ScheduleWithDuration        []json.RawMessage `json:"schedule_with_duration,omitempty"`
	SettingPreference           string            `json:"setting_preference,omitempty"`
	MinrateNgEnabled            *bool             `json:"minrate_ng_enabled,omitempty"`
	MinrateNgDataRateKbps       *int              `json:"minrate_ng_data_rate_kbps,omitempty"`
	MinrateNgAdvertisingRates   *bool             `json:"minrate_ng_advertising_rates,omitempty"`
	MinrateNaEnabled            *bool             `json:"minrate_na_enabled,omitempty"`
	MinrateNaDataRateKbps       *int              `json:"minrate_na_data_rate_kbps,omitempty"`
	MinrateNaAdvertisingRates   *bool             `json:"minrate_na_advertising_rates,omitempty"`
	MinrateSettingPreference    string            `json:"minrate_setting_preference,omitempty"`
	No2GhzOui                   *bool             `json:"no2ghz_oui,omitempty"`
	NoIPv6Ndp                   *bool             `json:"no_ipv6_ndp,omitempty"`
	OptimizeIotWifiConn         *bool             `json:"optimize_iot_wifi_connectivity,omitempty"`
	PmfMode                     string            `json:"pmf_mode,omitempty"`
	BcastEnhanceEnabled         *bool             `json:"bcastenhance_enabled,omitempty"`
	McastEnhanceEnabled         *bool             `json:"mcastenhance_enabled,omitempty"`
	GroupRekey                  *int              `json:"group_rekey,omitempty"`
	DtimMode                    string            `json:"dtim_mode,omitempty"`
	DtimNa                      *int              `json:"dtim_na,omitempty"`
	DtimNg                      *int              `json:"dtim_ng,omitempty"`
	Dtim6e                      *int              `json:"dtim_6e,omitempty"`
	Uapsd                       *bool             `json:"uapsd_enabled,omitempty"`
	FastRoamingEnabled          *bool             `json:"fast_roaming_enabled,omitempty"`
	ProxyArp                    *bool             `json:"proxy_arp,omitempty"`
	BssTransition               *bool             `json:"bss_transition,omitempty"`
	L2Isolation                 *bool             `json:"l2_isolation,omitempty"`
	IappEnabled                 *bool             `json:"iapp_enabled,omitempty"`
}

// PortConf represents a UniFi switch port profile.
type PortConf struct {
	ID                            string      `json:"_id,omitempty"`
	SiteID                        string      `json:"site_id,omitempty"`
	Name                          string      `json:"name"`
	Forward                       string      `json:"forward,omitempty"`
	NativeNetworkconfID           string      `json:"native_networkconf_id,omitempty"`
	TaggedNetworkconfIDs          []string    `json:"tagged_networkconf_ids,omitempty"`
	ExcludedNetworkconfIDs        []string    `json:"excluded_networkconf_ids,omitempty"`
	VoiceNetworkconfID            string      `json:"voice_networkconf_id,omitempty"`
	Autoneg                       *bool       `json:"autoneg,omitempty"`
	Dot1xCtrl                     string      `json:"dot1x_ctrl,omitempty"`
	Dot1xIDleTimeout              *int        `json:"dot1x_idle_timeout,omitempty"`
	EgressRateLimitKbps           *int        `json:"egress_rate_limit_kbps,omitempty"`
	EgressRateLimitEnabled        *bool       `json:"egress_rate_limit_kbps_enabled,omitempty"`
	FullDuplex                    *bool       `json:"full_duplex,omitempty"`
	Isolation                     *bool       `json:"isolation,omitempty"`
	LldpmedEnabled                *bool       `json:"lldpmed_enabled,omitempty"`
	LldpmedNotifyEnabled          *bool       `json:"lldpmed_notify_enabled,omitempty"`
	MulticastRouterNetworkconfIDs []string    `json:"multicast_router_networkconf_ids,omitempty"`
	OpMode                        string      `json:"op_mode,omitempty"`
	PoeMode                       string      `json:"poe_mode,omitempty"`
	PortKeepaliveEnabled          *bool       `json:"port_keepalive_enabled,omitempty"`
	PortSecurityEnabled           *bool       `json:"port_security_enabled,omitempty"`
	PortSecurityMacAddress        []string    `json:"port_security_mac_address,omitempty"`
	QosProfile                    *QoSProfile `json:"qos_profile,omitempty"`
	SettingPreference             string      `json:"setting_preference,omitempty"`
	Speed                         *int        `json:"speed,omitempty"`
	StormctrlBcastEnabled         *bool       `json:"stormctrl_bcast_enabled,omitempty"`
	StormctrlBcastRate            *int        `json:"stormctrl_bcast_rate,omitempty"`
	StormctrlMcastEnabled         *bool       `json:"stormctrl_mcast_enabled,omitempty"`
	StormctrlMcastRate            *int        `json:"stormctrl_mcast_rate,omitempty"`
	StormctrlUcastEnabled         *bool       `json:"stormctrl_ucast_enabled,omitempty"`
	StormctrlUcastRate            *int        `json:"stormctrl_ucast_rate,omitempty"`
	StpPortMode                   *bool       `json:"stp_port_mode,omitempty"`
	TaggedVlanMgmt                string      `json:"tagged_vlan_mgmt,omitempty"`
}

// Routing represents a UniFi static route.
type Routing struct {
	ID                   string `json:"_id,omitempty"`
	SiteID               string `json:"site_id,omitempty"`
	Name                 string `json:"name"`
	Enabled              *bool  `json:"enabled,omitempty"`
	Type                 string `json:"type,omitempty"`
	GatewayType          string `json:"gateway_type,omitempty"`
	GatewayDevice        string `json:"gateway_device,omitempty"`
	StaticRouteNetwork   string `json:"static-route_network,omitempty"`
	StaticRouteNexthop   string `json:"static-route_nexthop,omitempty"`
	StaticRouteDistance  *int   `json:"static-route_distance,omitempty"`
	StaticRouteInterface string `json:"static-route_interface,omitempty"`
	StaticRouteType      string `json:"static-route_type,omitempty"`
}

// UserGroup represents a UniFi user group.
type UserGroup struct {
	ID             string `json:"_id,omitempty"`
	SiteID         string `json:"site_id,omitempty"`
	Name           string `json:"name"`
	QosRateMaxDown *int   `json:"qos_rate_max_down,omitempty"`
	QosRateMaxUp   *int   `json:"qos_rate_max_up,omitempty"`
	AttrHiddenID   string `json:"attr_hidden_id,omitempty"`
	AttrNoDelete   *bool  `json:"attr_no_delete,omitempty"`
}

// RADIUSProfile represents a UniFi RADIUS profile.
type RADIUSProfile struct {
	ID                    string         `json:"_id,omitempty"`
	SiteID                string         `json:"site_id,omitempty"`
	Name                  string         `json:"name"`
	UseUsgAcctServer      *bool          `json:"use_usg_acct_server,omitempty"`
	UseUsgAuthServer      *bool          `json:"use_usg_auth_server,omitempty"`
	VlanEnabled           *bool          `json:"vlan_enabled,omitempty"`
	VlanWlanMode          string         `json:"vlan_wlan_mode,omitempty"`
	AcctServers           []RADIUSServer `json:"acct_servers,omitempty"`
	AuthServers           []RADIUSServer `json:"auth_servers,omitempty"`
	InterimUpdateEnabled  *bool          `json:"interim_update_enabled,omitempty"`
	InterimUpdateInterval *int           `json:"interim_update_interval,omitempty"`
	AttrHiddenID          string         `json:"attr_hidden_id,omitempty"`
	AttrNoDelete          *bool          `json:"attr_no_delete,omitempty"`
	AttrNoEdit            *bool          `json:"attr_no_edit,omitempty"`
}

// RADIUSServer represents a RADIUS server configuration.
type RADIUSServer struct {
	IP      string `json:"ip,omitempty"`
	Port    *int   `json:"port,omitempty"`
	XSecret string `json:"x_secret,omitempty"`
}

// PolicySchedule defines when a firewall policy is active.
type PolicySchedule struct {
	Mode           string   `json:"mode,omitempty"`
	TimeRangeStart string   `json:"time_range_start,omitempty"`
	TimeRangeEnd   string   `json:"time_range_end,omitempty"`
	DaysOfWeek     []string `json:"days_of_week,omitempty"`
}

// StaticDNS represents a static DNS record (v2 API).
type StaticDNS struct {
	ID         string `json:"_id,omitempty"`
	Key        string `json:"key,omitempty"`
	Value      string `json:"value,omitempty"`
	RecordType string `json:"record_type,omitempty"`
	Enabled    *bool  `json:"enabled,omitempty"`
	TTL        *int   `json:"ttl,omitempty"`
	Port       *int   `json:"port,omitempty"`
	Priority   *int   `json:"priority,omitempty"`
	Weight     *int   `json:"weight,omitempty"`
}

// TrafficRule represents a traffic management rule (v2 API).
type TrafficRule struct {
	ID             string              `json:"_id,omitempty"`
	Name           string              `json:"name"`
	Enabled        *bool               `json:"enabled,omitempty"`
	Action         string              `json:"action,omitempty"`
	MatchingTarget string              `json:"matching_target,omitempty"`
	TargetDevices  []TrafficRuleTarget `json:"target_devices,omitempty"`
	Schedule       *PolicySchedule     `json:"schedule,omitempty"`
	Description    string              `json:"description,omitempty"`
	AppCategoryIDs []string            `json:"app_category_ids,omitempty"`
	AppIDs         []int               `json:"app_ids,omitempty"`
	Domains        []TrafficDomain     `json:"domains,omitempty"`
	IPAddresses    []string            `json:"ip_addresses,omitempty"`
	IPRanges       []string            `json:"ip_ranges,omitempty"`
	Regions        []string            `json:"regions,omitempty"`
	NetworkID      string              `json:"network_id,omitempty"`
	BandwidthLimit *TrafficBandwidth   `json:"bandwidth_limit,omitempty"`
}

// TrafficRuleTarget specifies a device target for a traffic rule.
type TrafficRuleTarget struct {
	ClientMAC string `json:"client_mac,omitempty"`
	Type      string `json:"type,omitempty"`
	NetworkID string `json:"network_id,omitempty"`
}

// TrafficBandwidth specifies bandwidth limits for a traffic rule.
type TrafficBandwidth struct {
	DownloadLimitKbps *int  `json:"download_limit_kbps,omitempty"`
	UploadLimitKbps   *int  `json:"upload_limit_kbps,omitempty"`
	Enabled           *bool `json:"enabled,omitempty"`
}

// TrafficDomain represents a domain entry for traffic rules and routes.
type TrafficDomain struct {
	Domain      string `json:"domain"`
	Description string `json:"description,omitempty"`
	Ports       []int  `json:"ports,omitempty"`
}

// User represents a UniFi user/client device record (legacy REST API).
type User struct {
	ID          string `json:"_id,omitempty"`
	SiteID      string `json:"site_id,omitempty"`
	MAC         string `json:"mac"`
	Name        string `json:"name,omitempty"`
	Note        string `json:"note,omitempty"`
	Noted       *bool  `json:"noted,omitempty"`
	UseFixedIP  *bool  `json:"use_fixedip,omitempty"`
	FixedIP     string `json:"fixed_ip,omitempty"`
	NetworkID   string `json:"network_id,omitempty"`
	UsergroupID string `json:"usergroup_id,omitempty"`
	Blocked     *bool  `json:"blocked,omitempty"`
	IsWired     *bool  `json:"is_wired,omitempty"`
	IsGuest     *bool  `json:"is_guest,omitempty"`
	OUI         string `json:"oui,omitempty"`
	FirstSeen   *int64 `json:"first_seen,omitempty"`
	LastSeen    *int64 `json:"last_seen,omitempty"`
}
