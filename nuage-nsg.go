package nuagewim

import (
	"fmt"
	"strings"

	"github.com/imdario/mergo"
	"github.com/nuagenetworks/go-bambou/bambou"
	"github.com/nuagenetworks/vspk-go/vspk"
)

type NuageNSGCfg struct {
	Name            string `json:"Name"`
	NSGTemplateName string `json:"NSGTemplateName"`
	NSGTemplateID   string `json:"NSGTemplateID"`
	NetworkPorts    []struct {
		Name             string `json:"Name"`
		VlanID           int    `json:"vlanId"`
		Role             string `json:"Role"`
		AddressFamily    string `json:"addressFamily"`
		Mode             string `json:"mode"`
		UnderlayName     string `json:"underlayName,omitempty"`
		UnderlayID       string `json:"underlayId,omitempty"`
		VscName          string `json:"vscName,omitempty"`
		VscID            string `json:"vscId,omitempty"`
		Address          string `json:"address,omitempty"`
		Netmask          string `json:"netmask,omitempty"`
		DNS              string `json:"dns,omitempty"`
		Gateway          string `json:"gateway,omitempty"`
		LteConfiguration struct {
			Apn     string `json:"apn"`
			PdpType string `json:"pdp-type"`
			PinCode string `json:"pin-code"`
		} `json:"lteConfiguration,omitempty"`
	} `json:"NetworkPorts"`
	ShuntPorts []struct {
		Name          string `json:"Name"`
		VlanID        int    `json:"vlanId"`
		Role          string `json:"Role"`
		AddressFamily string `json:"addressFamily"`
		Mode          string `json:"mode"`
		Address       string `json:"address"`
		Netmask       string `json:"netmask"`
		DNS           string `json:"dns"`
		Gateway       string `json:"gateway"`
	} `json:"ShuntPorts"`
	AccessPorts []struct {
		Name   string `json:"Name"`
		VlanID int    `json:"VlanID"`
	} `json:"AccessPorts"`
	WifiPorts []struct {
		Name string `json:"Name"`
		Ssid string `json:"ssid"`
	} `json:"WifiPorts"`
}

func NuageCreateEntireNSG(nsgCfg NuageNSGCfg, parent *vspk.Enterprise) *vspk.NSGateway {
	fmt.Println("########################################")
	fmt.Println("#####  Create Entire NSG GW   ##########")
	fmt.Println("########################################")

	networkAcceleration := "NONE"
	if nsgCfg.ShuntPorts != nil {
		networkAcceleration = "PERFORMANCE"
	}

	nsGatewayCfg := map[string]interface{}{
		"Name":                  nsgCfg.Name,
		"TCPMSSEnabled":         true,
		"TCPMaximumSegmentSize": 1330,
		"NetworkAcceleration":   networkAcceleration,
		"TemplateID":            nsgCfg.NSGTemplateID,
	}

	nsGateway := NuageNSG(nsGatewayCfg, parent)

	for i, port := range nsgCfg.NetworkPorts {
		fmt.Printf("NSG Network Port %d Name: %s \n", i, port.Name)

		nsPortCfg := map[string]interface{}{
			"Name":            port.Name,
			"PhysicalName":    port.Name,
			"PortType":        "NETWORK",
			"VLANRange":       "0-4094",
			"EnableNATProbes": true,
			"NATTraversal":    "FULL_NAT",
		}
		nsPort := NuageNSGPort(nsPortCfg, nsGateway)

		fmt.Printf("Port: %#v \n", nsPort)

		nsVlanCfg := map[string]interface{}{
			"Value":                  port.VlanID,
			"AssociatedVSCProfileID": port.VscID,
		}
		fmt.Printf("VLANCfg: %#v \n", nsVlanCfg)
		nsVlan := NuageVlan(nsVlanCfg, nsPort)
		fmt.Printf("Port: %#v \n", nsVlan)

		var patEnabled = true
		var underlayEnabled = true
		var dnsV4 = ""
		var addressV4 = ""
		var netmaskV4 = ""
		var gatewayV4 = ""
		var dnsV6 = ""
		var addressV6 = ""
		var gatewayV6 = ""

		if port.AddressFamily == "IPV4" && port.Mode == "static" {
			dnsV4 = port.DNS
			addressV4 = port.Address
			netmaskV4 = port.Netmask
			gatewayV4 = port.Gateway
			fmt.Printf("BRANCH IPV4 AND STATIC ADDRESSING MODE \n")
		} else if port.AddressFamily == "IPV6" && port.Mode == "static" {
			dnsV6 = port.DNS
			addressV6 = port.Address
			gatewayV6 = port.Gateway
			patEnabled = false
			underlayEnabled = false
			fmt.Printf("BRANCH IPV6 AND STATIC ADDRESSING MODE \n")
		} else if port.AddressFamily == "IPV6" {
			patEnabled = false
			underlayEnabled = false
			fmt.Printf("BRANCH IPV6 AND DYNAMIC ADDRESSING MODE \n")
		}

		uplinkConnectionCfg := map[string]interface{}{
			"PATEnabled":      patEnabled,
			"UnderlayEnabled": underlayEnabled,
			"Role":            port.Role,
			"Mode":            port.Mode,
			"AddressFamily":   port.AddressFamily,
			"DNSAddress":      dnsV4,
			"Gateway":         gatewayV4,
			"Address":         addressV4,
			"Netmask":         netmaskV4,
			"DNSAddressV6":    dnsV6,
			"GatewayV6":       gatewayV6,
			"AddressV6":       addressV6,
			"UnderlayID":      port.UnderlayID,
		}
		fmt.Println(uplinkConnectionCfg)

		uplinkConn := NuageUplinkConnection(uplinkConnectionCfg, nsVlan)

		if strings.Contains(port.Name, "lte") {
			fmt.Println("LTE")

			customePropCfg := map[string]interface{}{
				"AttributeName":  "apn",
				"AttributeValue": port.LteConfiguration.Apn,
			}

			NuageCustomProperty(customePropCfg, uplinkConn)

			customePropCfg = map[string]interface{}{
				"AttributeName":  "pdp-type",
				"AttributeValue": port.LteConfiguration.PdpType,
			}

			NuageCustomProperty(customePropCfg, uplinkConn)

			customePropCfg = map[string]interface{}{
				"AttributeName":  "pin-code",
				"AttributeValue": port.LteConfiguration.PinCode,
			}

			NuageCustomProperty(customePropCfg, uplinkConn)

		} else {
			fmt.Println("ETHERNET")
		}
	}
	for i, port := range nsgCfg.ShuntPorts {
		fmt.Printf("NSG Shunt Port %d Name: %s \n", i, port.Name)

		nsPortCfg := map[string]interface{}{
			"Name":            port.Name,
			"PhysicalName":    port.Name,
			"PortType":        "NETWORK",
			"VLANRange":       "0-4094",
			"EnableNATProbes": true,
			"NATTraversal":    "FULL_NAT",
			"Mtu":             2000,
		}
		nsPort := NuageNSGPort(nsPortCfg, nsGateway)

		nsVlanCfg := map[string]interface{}{
			"Value":       port.VlanID,
			"Name":        "shunt",
			"Description": "shunt",
			"ShuntVLAN":   true,
		}
		fmt.Printf("VLANCfg: %s \n", nsVlanCfg)
		nsVlan := NuageVlan(nsVlanCfg, nsPort)
		fmt.Println(nsVlan)

		uplinkConnCfg := map[string]interface{}{
			"PATEnabled":    true,
			"Role":          port.Role,
			"Mode":          port.Mode,
			"AddressFamily": port.AddressFamily,
			"DNSAddress":    port.DNS,
			"Gateway":       port.Gateway,
			"Address":       port.Address,
			"Netmask":       port.Netmask,
			//	"DNSAddressV6":    dnsV6,
			//	"GatewayV6":       gatewayV6,
			//	"AddressV6":       addressV6,
			"UnderlayEnabled": true,
			//"UnderlayID":      port.UnderlayID,
		}

		uplinkConn := NuageUplinkConnection(uplinkConnCfg, nsVlan)
		fmt.Println(uplinkConn)

	}
	for i, port := range nsgCfg.AccessPorts {
		fmt.Printf("NSG Access Port %d Name: %s \n", i, port.Name)

		nsPortCfg := map[string]interface{}{
			"Name":         port.Name,
			"PhysicalName": port.Name,
			"PortType":     "ACCESS",
			"VLANRange":    "0-4094",
		}
		nsPort := NuageNSGPort(nsPortCfg, nsGateway)

		nsVlanCfg := map[string]interface{}{
			"Value": port.VlanID,
		}
		fmt.Printf("Access VLANCfg: %#v \n", nsVlanCfg)
		nsVlan := NuageVlan(nsVlanCfg, nsPort)
		fmt.Println(nsVlan)
	}
	for i, port := range nsgCfg.WifiPorts {
		fmt.Printf("NSG Wifi Port %d Name: %s \n", i, port.Name)

		nsPortCfg := map[string]interface{}{
			"Name":              port.Name,
			"WifiFrequencyBand": "FREQ_2_4_GHZ",
			"WifiMode":          "WIFI_B_G_N",
			"CountryCode":       "BE",
		}
		nsPort := NuageNSGWirelessPort(nsPortCfg, nsGateway)

		ssidConnCfg := map[string]interface{}{
			"Name":               port.Ssid,
			"Passphrase":         "4no*heydQ",
			"AuthenticationMode": "WPA2",
			"BroadcastSSID":      true,
		}
		ssidConn := NuageSSIDConnection(ssidConnCfg, nsPort)
		fmt.Println(ssidConn)
	}

	fmt.Printf("%#v \n", nsGateway)
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	return nsGateway
}

func NuageNSGatewayTemplate(nsGatewayTemplCfg map[string]interface{}, parent *vspk.Me) *vspk.NSGatewayTemplate {
	fmt.Println("########################################")
	fmt.Println("#####  NSG Gateway Profile    ##########")
	fmt.Println("########################################")

	nsGatewayTempls, err := parent.NSGatewayTemplates(&bambou.FetchingInfo{
		Filter: nsGatewayTemplCfg["Name"].(string)})
	handleError(err, "READ", "NS Gateway Template")

	// init the nsGateway struct that will hold either the received object
	// or will be created from the nsGatewayCfg
	nsGatewayTempl := &vspk.NSGatewayTemplate{}

	if nsGatewayTempls != nil {
		fmt.Println("NS Gateway Template already exists")

		nsGatewayTempl = nsGatewayTempls[0]
		mergo.Map(nsGatewayTempl, nsGatewayTemplCfg, mergo.WithOverride)

		nsGatewayTempl.Save()
	} else {
		mergo.Map(nsGatewayTempl, nsGatewayTemplCfg, mergo.WithOverride)
		err := parent.CreateNSGatewayTemplate(nsGatewayTempl)
		handleError(err, "CREATE", "NS Gateway Template ")

		fmt.Println("NS Gateway Template created")
	}

	fmt.Printf("%#v \n", nsGatewayTempl)
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	return nsGatewayTempl
}

func NuageNSG(nsGatewayCfg map[string]interface{}, parent *vspk.Enterprise) *vspk.NSGateway {
	fmt.Println("########################################")
	fmt.Println("#####        NSG Gateway      ##########")
	fmt.Println("########################################")

	nsGateways, err := parent.NSGateways(&bambou.FetchingInfo{
		Filter: nsGatewayCfg["Name"].(string)})
	handleError(err, "READ", "NS Gateway")

	// init the nsGateway struct that will hold either the received object
	// or will be created from the nsGatewayCfg
	nsGateway := &vspk.NSGateway{}

	if nsGateways != nil {
		fmt.Println("NS Gateway already exists")

		nsGateway = nsGateways[0]
		mergo.Map(nsGateway, nsGatewayCfg, mergo.WithOverride)

		nsGateway.Save()
	} else {
		mergo.Map(nsGateway, nsGatewayCfg, mergo.WithOverride)
		err := parent.CreateNSGateway(nsGateway)
		handleError(err, "CREATE", "NS Gateway ")

		fmt.Println("NS Gateway created")
	}

	fmt.Printf("%#v \n", nsGateway)
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	return nsGateway
}

func NuageNSGPort(nsPortCfg map[string]interface{}, parent *vspk.NSGateway) *vspk.NSPort {
	fmt.Println("########################################")
	fmt.Println("#####        NSG Port         ##########")
	fmt.Println("########################################")

	nsPorts, err := parent.NSPorts(&bambou.FetchingInfo{
		Filter: nsPortCfg["Name"].(string)})
	handleError(err, "READ", "NSG Port")

	// init the nsPort struct that will hold either the received object
	// or will be created from the nsPortCfg
	nsPort := &vspk.NSPort{}

	if nsPorts != nil {
		fmt.Println("NS Port already exists")

		nsPort = nsPorts[0]
		mergo.Map(nsPort, nsPortCfg, mergo.WithOverride)

		nsPort.Save()
	} else {
		mergo.Map(nsPort, nsPortCfg, mergo.WithOverride)
		fmt.Println(nsPortCfg)
		err := parent.CreateNSPort(nsPort)
		handleError(err, "CREATE", "NS Port ")

		fmt.Println("NS Port created")
	}

	fmt.Printf("%#v \n", nsPort)
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	return nsPort
}

func NuageNSGWirelessPort(nsPortCfg map[string]interface{}, parent *vspk.NSGateway) *vspk.WirelessPort {
	fmt.Println("########################################")
	fmt.Println("#####  NSG Wireless Port      ##########")
	fmt.Println("########################################")

	nsPorts, err := parent.WirelessPorts(&bambou.FetchingInfo{
		Filter: nsPortCfg["Name"].(string)})
	handleError(err, "READ", "Wireless Port")

	// init the nsPort struct that will hold either the received object
	// or will be created from the nsPortCfg
	nsPort := &vspk.WirelessPort{}

	if nsPorts != nil {
		fmt.Println("Wireless Port already exists")

		nsPort = nsPorts[0]
		mergo.Map(nsPort, nsPortCfg, mergo.WithOverride)

		nsPort.Save()
	} else {
		mergo.Map(nsPort, nsPortCfg, mergo.WithOverride)
		err := parent.CreateWirelessPort(nsPort)
		handleError(err, "CREATE", "Wireless Port ")

		fmt.Println("Wireless Port created")
	}

	fmt.Printf("%#v \n", nsPort)
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	return nsPort
}

func NuageVlan(nsVlanCfg map[string]interface{}, parent *vspk.NSPort) *vspk.VLAN {
	fmt.Println("########################################")
	fmt.Println("#####        NSG Vlan         ##########")
	fmt.Println("########################################")

	fmt.Printf("VLAN Cfg: %#v \n", nsVlanCfg)
	nsVlans, err := parent.VLANs(&bambou.FetchingInfo{
		Filter: fmt.Sprintf("value == %d", nsVlanCfg["Value"])})

	handleError(err, "READ", "NSG VLAN")

	// init the nsVlan struct that will hold either the received object
	// or will be created from the nsVlanCfg
	nsVlan := &vspk.VLAN{}

	fmt.Printf("VLANs %s \n", nsVlans)

	if nsVlans != nil {
		fmt.Println("NS VLAN already exists")

		nsVlan = nsVlans[0]
		mergo.Map(nsVlan, nsVlanCfg, mergo.WithOverride)

		nsVlan.Save()
	} else {
		mergo.Map(nsVlan, nsVlanCfg, mergo.WithOverride)
		//nsVlan.Value, _ = strconv.Atoi("0")
		//nsVlan.Value = 0
		fmt.Printf("VLAN: %#v \n", nsVlan)
		fmt.Printf("Port: %#v \n", parent)
		err := parent.CreateVLAN(nsVlan)
		handleError(err, "CREATE", "NS VLAN ")

		fmt.Println("NS VLAN created")
	}

	fmt.Printf("%#v \n", nsVlan)
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	return nsVlan
}

func NuageSSIDConnection(ssidConnCfg map[string]interface{}, parent *vspk.WirelessPort) *vspk.SSIDConnection {
	fmt.Println("########################################")
	fmt.Println("#####   SSID Connection       ##########")
	fmt.Println("########################################")

	ssidConns, err := parent.SSIDConnections(&bambou.FetchingInfo{
		Filter: ssidConnCfg["Name"].(string)})
	handleError(err, "READ", "SSiD Connection")

	ssidConn := &vspk.SSIDConnection{}

	if ssidConns != nil {
		fmt.Println("SSiD Connection already exists")

		ssidConn = ssidConns[0]
		mergo.Map(ssidConn, ssidConnCfg, mergo.WithOverride)

		ssidConn.Save()
	} else {
		mergo.Map(ssidConn, ssidConnCfg, mergo.WithOverride)

		err := parent.CreateSSIDConnection(ssidConn)
		handleError(err, "CREATE", "SSiD Connection")

		fmt.Println("SSiD Connection created")
	}

	fmt.Printf("%#v \n", ssidConn)
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	return ssidConn
}

func NuageUplinkConnection(uplinkConnCfg map[string]interface{}, parent *vspk.VLAN) *vspk.UplinkConnection {
	fmt.Println("########################################")
	fmt.Println("#####   NSG Uplink Connection ##########")
	fmt.Println("########################################")

	uplinkConns, err := parent.UplinkConnections(&bambou.FetchingInfo{})

	handleError(err, "READ", "Uplink Connection")

	uplinkConn := &vspk.UplinkConnection{}

	if uplinkConns != nil {
		fmt.Println("Uplink Connection already exists")

		uplinkConn = uplinkConns[0]
		mergo.Map(uplinkConn, uplinkConnCfg, mergo.WithOverride)

		uplinkConn.Save()
	} else {
		fmt.Printf("Uplink Connection Cfg: %#v \n", uplinkConnCfg)
		mergo.Map(uplinkConn, uplinkConnCfg, mergo.WithOverride)
		fmt.Printf("Uplink Connection: %#v \n", uplinkConn)
		fmt.Printf("Uplink Connection Cfg: %#v \n", uplinkConnCfg)

		err := parent.CreateUplinkConnection(uplinkConn)
		handleError(err, "CREATE", "Uplink Connection")

		fmt.Println("Uplink Connection created")
	}

	fmt.Printf("%#v \n", uplinkConn)
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	return uplinkConn
}

func NuageCustomProperty(customePropCfg map[string]interface{}, parent *vspk.UplinkConnection) *vspk.CustomProperty {
	fmt.Println("########################################")
	fmt.Println("#####   Custom property       ##########")
	fmt.Println("########################################")

	customeProps, err := parent.CustomProperties(&bambou.FetchingInfo{
		Filter: customePropCfg["AttributeName"].(string)})
	handleError(err, "READ", "Custome Property")

	// init the nsVlan struct that will hold either the received object
	// or will be created from the nsVlanCfg
	customeProp := &vspk.CustomProperty{}

	if customeProps != nil {
		fmt.Println("Custom property already exists")

		customeProp = customeProps[0]
		mergo.Map(customeProp, customePropCfg, mergo.WithOverride)

		customeProp.Save()
	} else {
		mergo.Map(customeProp, customePropCfg, mergo.WithOverride)

		err := parent.CreateCustomProperty(customeProp)
		handleError(err, "CREATE", "Custome Property")

		fmt.Println("Custome Property created")
	}

	fmt.Printf("%#v \n", customeProp)
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	return customeProp
}
