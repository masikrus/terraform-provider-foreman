package foreman

import (
	"context"
	"fmt"

	"github.com/HanseMerkur/terraform-provider-utils/autodoc"
	"github.com/HanseMerkur/terraform-provider-utils/helper"
	"github.com/HanseMerkur/terraform-provider-utils/log"
	"github.com/masikrus/terraform-provider-foreman/foreman/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceForemanProvisioningTemplate() *schema.Resource {
	// copy attributes from resource definition
	r := resourceForemanProvisioningTemplate()
	ds := helper.DataSourceSchemaFromResourceSchema(r.Schema)

	// define searchable attributes for the data source
	ds["name"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		Description: fmt.Sprintf(
			"The name of the provisioning template. "+
				"%s \"ESXi 6.0 Kickstart - BO1\"",
			autodoc.MetaExample,
		),
	}

	return &schema.Resource{

		ReadContext: dataSourceForemanProvisioningTemplateRead,

		// NOTE(ALL): See comments in the corresponding resource file
		Schema: ds,
	}
}

func dataSourceForemanProvisioningTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("data_source_foreman_provisioningtemplate.go#Read")

	client := meta.(*api.Client)
	t := buildForemanProvisioningTemplate(d)

	log.Debugf("ForemanProvisioningTemplate: [%+v]", t)

	queryResponse, queryErr := client.QueryProvisioningTemplate(ctx, t)
	if queryErr != nil {
		return diag.FromErr(queryErr)
	}

	if queryResponse.Subtotal == 0 {
		return diag.Errorf("Data source provisioning template returned no results")
	} else if queryResponse.Subtotal > 1 {
		return diag.Errorf("Data source provisioning template returned more than 1 result")
	}

	var queryTemplate api.ForemanProvisioningTemplate
	var ok bool
	if queryTemplate, ok = queryResponse.Results[0].(api.ForemanProvisioningTemplate); !ok {
		return diag.Errorf(
			"Data source results contain unexpected type. Expected "+
				"[api.ForemanProvisioningTemplate], got [%T]",
			queryResponse.Results[0],
		)
	}
	t = &queryTemplate

	log.Debugf("ForemanProvisioningTemplate: [%+v]", t)

	setResourceDataFromForemanProvisioningTemplate(d, t)

	return nil
}
