package nuagewim

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/nuagenetworks/go-bambou/bambou"
	"github.com/nuagenetworks/vspk-go/vspk"
)

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
		nsGatewayErr := parent.CreateNSGateway(nsGateway)
		handleError(nsGatewayErr, "CREATE", "NS Gateway ")

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
		nsPortErr := parent.CreateNSPort(nsPort)
		handleError(nsPortErr, "CREATE", "NS Port ")

		fmt.Println("NS Gateway created")
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

	nsVlans, err := parent.VLANs(&bambou.FetchingInfo{
		Filter: fmt.Sprintf("value == %s", nsVlanCfg["Value"])})

	handleError(err, "READ", "NSG VLAN")

	// init the nsVlan struct that will hold either the received object
	// or will be created from the nsVlanCfg
	nsVlan := &vspk.VLAN{}

	if nsVlans != nil {
		fmt.Println("NS VLAN already exists")

		nsVlan = nsVlans[0]
		mergo.Map(nsVlan, nsVlanCfg, mergo.WithOverride)

		nsVlan.Save()
	} else {
		mergo.Map(nsVlan, nsVlanCfg, mergo.WithOverride)
		nsVlanErr := parent.CreateVLAN(nsVlan)
		handleError(nsVlanErr, "CREATE", "NS VLAN ")

		fmt.Println("NS VLAN created")
	}

	fmt.Printf("%#v \n", nsVlan)
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	return nsVlan
}
