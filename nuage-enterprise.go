package nuagewim

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/nuagenetworks/go-bambou/bambou"
	"github.com/nuagenetworks/vspk-go/vspk"
)

func nuageEnterprise(enterpriseCfg map[string]interface{}, parent *vspk.Me) *vspk.Enterprise {
	enterprise := &vspk.Enterprise{}

	enterprises, err := parent.Enterprises(&bambou.FetchingInfo{
		Filter: enterpriseCfg["Name"].(string)})
	handleError(err, "Enterprise", "READ")

	// init the enterprise struct that will hold either the received object
	// or will be created from the enterpriseCfg
	if enterprises != nil {
		fmt.Println("Enterpise already exists")

		enterprise = enterprises[0]
		mergo.Map(enterprise, enterpriseCfg, mergo.WithOverride)
		enterprise.Save()

	} else {
		mergo.Map(enterprise, enterpriseCfg, mergo.WithOverride)
		err := parent.CreateEnterprise(enterprise)
		handleError(err, "Enterprise", "CREATE")

		fmt.Println("Enterprise created")
	}
	return enterprise
}
