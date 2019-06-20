package nuagewim

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/nuagenetworks/go-bambou/bambou"
	"github.com/nuagenetworks/vspk-go/vspk"
)

func NuageInfraVSCProf(vscProfCfg map[string]interface{}, parent *vspk.Me) *vspk.InfrastructureVscProfile {
	fmt.Println("########################################")
	fmt.Println("#####     Infra VSC Profil    ##########")
	fmt.Println("########################################")

	vscProfs, err := parent.InfrastructureVscProfiles(&bambou.FetchingInfo{
		Filter: vscProfCfg["Name"].(string)})
	handleError(err, "READ", "VSC profile")

	vscProf := &vspk.InfrastructureVscProfile{}

	if vscProfs != nil {
		fmt.Println("VSC already exists")

		vscProf = vscProfs[0]
		mergo.Map(vscProf, vscProfCfg, mergo.WithOverride)

		vscProf.Save()
	} else {
		mergo.Map(vscProf, vscProfCfg, mergo.WithOverride)
		err := parent.CreateInfrastructureVscProfile(vscProf)
		handleError(err, "CREATE", "VSC Profile ")

		fmt.Println("VSC profile created")
	}

	fmt.Printf("%#v \n", vscProf)
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	return vscProf
}

func NuageUnderlay(underlayCfg map[string]interface{}, parent *vspk.Me) *vspk.Underlay {
	fmt.Println("########################################")
	fmt.Println("#####     Underlay            ##########")
	fmt.Println("########################################")

	underlays, err := parent.Underlays(&bambou.FetchingInfo{
		Filter: underlayCfg["Name"].(string)})
	handleError(err, "READ", "Underlay")

	underlay := &vspk.Underlay{}

	if underlays != nil {
		fmt.Println("Underlay already exists")

		underlay = underlays[0]
		mergo.Map(underlay, underlayCfg, mergo.WithOverride)

		underlay.Save()
	} else {
		mergo.Map(underlay, underlayCfg, mergo.WithOverride)
		err := parent.CreateUnderlay(underlay)
		handleError(err, "CREATE", "Underlay ")

		fmt.Println("Underlay created")
	}

	fmt.Printf("%#v \n", underlay)
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	return underlay
}
