package nuagewim

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/nuagenetworks/go-bambou/bambou"
	"github.com/nuagenetworks/vspk-go/vspk"
)

func NuageDomain(domainCfg map[string]interface{}, parent *vspk.Enterprise) *vspk.Domain {
	fmt.Println("########################################")
	fmt.Println("#####            Domain       ##########")
	fmt.Println("########################################")

	domains, err := parent.Domains(&bambou.FetchingInfo{
		Filter: domainCfg["Name"].(string)})
	handleError(err, "READ", "Domain")

	domain := &vspk.Domain{}

	if domains != nil {
		fmt.Println("DOmain already exists")

		domain = domains[0]
		mergo.Map(domain, domainCfg, mergo.WithOverride)

		domain.Save()
	} else {
		mergo.Map(domain, domainCfg, mergo.WithOverride)

		err := parent.CreateDomain(domain)
		handleError(err, "CREATE", "Domain")

		fmt.Println("Domain created")
	}

	fmt.Printf("%#v \n", domain)
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	fmt.Println("****************************************")
	return domain
}
