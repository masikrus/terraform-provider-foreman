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

func dataSourceForemanComputeProfile() *schema.Resource {
	r := resourceForemanComputeProfile()
	ds := helper.DataSourceSchemaFromResourceSchema(r.Schema)

	ds["name"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		Description: fmt.Sprintf(
			"Compute profile name."+
				"%s \"2-Medium\"",
			autodoc.MetaExample,
		),
	}

	return &schema.Resource{
		ReadContext: dataSourceForemanComputeProfileRead,
		Schema:      ds,
	}
}

// -----------------------------------------------------------------------------
// Resource CRUD Operations
// -----------------------------------------------------------------------------

func dataSourceForemanComputeProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("data_source_foreman_computeprofile.go#Read")

	client := meta.(*api.Client)
	t := buildForemanComputeProfile(d)

	log.Debugf("ForemanComputeProfile: [%+v]", t)

	queryResponse, queryErr := client.QueryComputeProfile(ctx, t)
	if queryErr != nil {
		return diag.FromErr(queryErr)
	}

	if queryResponse.Subtotal == 0 {
		return diag.Errorf("Data source template kind returned no results")
	} else if queryResponse.Subtotal > 1 {
		return diag.Errorf("Data source template kind returned more than 1 result")
	}

	var queryComputeProfile api.ForemanComputeProfile
	var ok bool
	if queryComputeProfile, ok = queryResponse.Results[0].(api.ForemanComputeProfile); !ok {
		return diag.Errorf(
			"Data source results contain unexpected type. Expected "+
				"[api.ForemanComputeProfile], got [%T]",
			queryResponse.Results[0],
		)
	}
	t = &queryComputeProfile

	log.Debugf("ForemanComputeProfile: [%+v]", t)

	setResourceDataFromForemanComputeProfile(d, t)

	return nil
}
