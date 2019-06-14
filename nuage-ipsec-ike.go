package nuagewim

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/nuagenetworks/go-bambou/bambou"
	"github.com/nuagenetworks/vspk-go/vspk"
)

func nuageCreateIKEPSK(ikePSKCfg map[string]interface{}, parent *vspk.Enterprise) *vspk.IKEPSK {
	//create PSK
	fmt.Println("########################################")
	fmt.Println("#####     IKE PSK          #############")
	fmt.Println("########################################")

	ikePSK := &vspk.IKEPSK{}

	ikePSKs, err := parent.IKEPSKs(&bambou.FetchingInfo{
		Filter: ikePSKCfg["Name"].(string)})
	handleError(err, "IKE PSK", "READ")

	// init the ikePSK struct that will hold either the received object
	// or will be created from the ikePSKCfg
	if ikePSKs != nil {
		fmt.Println("IKE PSK already exists")

		ikePSK = ikePSKs[0]
		mergo.Map(ikePSK, ikePSKCfg, mergo.WithOverride)
		ikePSK.Save()

	} else {

		//ikePSK = &vspk.IKEPSK{}
		mergo.Map(ikePSK, ikePSKCfg, mergo.WithOverride)
		ikePSKErr := parent.CreateIKEPSK(ikePSK)
		handleError(ikePSKErr, "IKE PSK", "CREATE")

		fmt.Println("IKE PSK created")
	}
	fmt.Printf("%#v \n", ikePSK)
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	return ikePSK
}

func nuageCreateIKEGateway(ikeGatewayCfg map[string]interface{}, parent *vspk.Enterprise) *vspk.IKEGateway {
	fmt.Println("########################################")
	fmt.Println("#####     IKE Gateway      #############")
	fmt.Println("########################################")

	ikeGateways, err := parent.IKEGateways(&bambou.FetchingInfo{
		Filter: ikeGatewayCfg["Name"].(string)})
	handleError(err, "READ", "IKE Gateway")

	// init the ikeGateway struct that will hold either the received object
	// or will be created from the ikeGatewayCfg
	ikeGateway := &vspk.IKEGateway{}

	if ikeGateways != nil {
		fmt.Println("IKE Gateway already exists")

		ikeGateway = ikeGateways[0]
		mergo.Map(ikeGateway, ikeGatewayCfg, mergo.WithOverride)

		ikeGateway.Save()
	} else {
		//ikeGateway1 = &vspk.IKEGateway{}
		mergo.Map(ikeGateway, ikeGatewayCfg, mergo.WithOverride)
		ikeGatewayErr := parent.CreateIKEGateway(ikeGateway)
		handleError(ikeGatewayErr, "CREATE", "IKE Gateway")

		fmt.Println("IKE Gateway created")

		ikeSubnet := &vspk.IKESubnet{}
		ikeSubnet.Prefix = "0.0.0.0/0"
		ikeSubnetErr := ikeGateway.CreateIKESubnet(ikeSubnet)
		handleError(ikeSubnetErr, "CREATE", "IKE Subnet")
		fmt.Printf("IKE Subnet created: %s\n", ikeSubnet)
	}
	fmt.Printf("%#v \n", ikeGateway)
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	return ikeGateway
}

func nuageCreateIKEEncryptionProfile(ikeEncryptionProfileCfg map[string]interface{}, parent *vspk.Enterprise) *vspk.IKEEncryptionprofile {
	fmt.Println("########################################")
	fmt.Println("#####IKE Encryption Profile#############")
	fmt.Println("########################################")

	ikeEncryptionProfiles, err := parent.IKEEncryptionprofiles(&bambou.FetchingInfo{
		Filter: ikeEncryptionProfileCfg["Name"].(string)})
	handleError(err, "READ", "IKE Encryption Profile")

	// init the IKEEncryptionprofile struct that will hold either the received object
	// or will be created from the IKEEncryptionprofileCfg
	ikeEncryptionProfile := &vspk.IKEEncryptionprofile{}

	if ikeEncryptionProfiles != nil {
		fmt.Println("IKE Encryption Profile already exists")

		ikeEncryptionProfile = ikeEncryptionProfiles[0]
		mergo.Map(ikeEncryptionProfile, ikeEncryptionProfileCfg, mergo.WithOverride)

		ikeEncryptionProfile.Save()
	} else {
		//ikeEncryptionProfile = &vspk.IKEEncryptionprofile{}
		mergo.Map(ikeEncryptionProfile, ikeEncryptionProfileCfg, mergo.WithOverride)
		ikeEncryptionProfileErr := parent.CreateIKEEncryptionprofile(ikeEncryptionProfile)
		handleError(ikeEncryptionProfileErr, "CREATE", "IKE Encryption Profile")

		fmt.Println("IKE IKE Encryption Profile created")
	}
	fmt.Printf("%#v \n", ikeEncryptionProfile)
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	return ikeEncryptionProfile
}

func nuageCreateIKEGatewayProfile(ikeGatewayProfileCfg map[string]interface{}, parent *vspk.Enterprise) *vspk.IKEGatewayProfile {
	fmt.Println("########################################")
	fmt.Println("#####   IKE Gateway Profile   ##########")
	fmt.Println("########################################")

	ikeGatewayProfiles, err := parent.IKEGatewayProfiles(&bambou.FetchingInfo{
		Filter: ikeGatewayProfileCfg["Name"].(string)})
	handleError(err, "READ", "IKE Gateway Profile")

	// init the ikeGatewayProfile struct that will hold either the received object
	// or will be created from the ikeGatewayProfileCfg
	ikeGatewayProfile := &vspk.IKEGatewayProfile{}

	if ikeGatewayProfiles != nil {
		fmt.Println("IKE Gateway Profile already exists")

		ikeGatewayProfile = ikeGatewayProfiles[0]
		mergo.Map(ikeGatewayProfile, ikeGatewayProfileCfg, mergo.WithOverride)

		ikeGatewayProfile.Save()
	} else {
		//ikeGatewayProfile1 = &vspk.IKEGatewayProfile{}
		mergo.Map(ikeGatewayProfile, ikeGatewayProfileCfg, mergo.WithOverride)
		ikeGatewayProfileErr := parent.CreateIKEGatewayProfile(ikeGatewayProfile)
		handleError(ikeGatewayProfileErr, "CREATE", "IKE Gateway Profile1")

		fmt.Println("IKE Gateway Profile1 created")
	}

	fmt.Printf("%#v \n", ikeGatewayProfile)
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	return ikeGatewayProfile
}

func nuageIKEGatewayConnection(ikeGatewayConnCfg map[string]interface{}, parent *vspk.VLAN) *vspk.IKEGatewayConnection {
	fmt.Println("########################################")
	fmt.Println("#####  IKE GW Connection      ##########")
	fmt.Println("########################################")

	ikeGatewayConns, err := parent.IKEGatewayConnections(&bambou.FetchingInfo{
		Filter: ikeGatewayConnCfg["Name"].(string)})
	handleError(err, "READ", "IKE GW Connection")

	// init the nsPort struct that will hold either the received object
	// or will be created from the nsPortCfg
	ikeGatewayConn := &vspk.IKEGatewayConnection{}

	if ikeGatewayConns != nil {
		fmt.Println("IKE GW Connection already exists")

		ikeGatewayConn = ikeGatewayConns[0]
		mergo.Map(ikeGatewayConn, ikeGatewayConnCfg, mergo.WithOverride)

		ikeGatewayConn.Save()
	} else {
		mergo.Map(ikeGatewayConn, ikeGatewayConnCfg, mergo.WithOverride)
		ikeGatewayConnErr := parent.CreateIKEGatewayConnection(ikeGatewayConn)
		handleError(ikeGatewayConnErr, "CREATE", "IKE GW Connection ")

		fmt.Println("IKE GW Connection created")
	}

	fmt.Printf("%#v \n", ikeGatewayConn)
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	return ikeGatewayConn
}
