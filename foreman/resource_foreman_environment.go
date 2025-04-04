package foreman

import (
	"context"
	"fmt"
	"strconv"

	"github.com/HanseMerkur/terraform-provider-utils/autodoc"
	"github.com/HanseMerkur/terraform-provider-utils/log"
	"github.com/masikrus/terraform-provider-foreman/foreman/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceForemanEnvironment() *schema.Resource {
	return &schema.Resource{

		CreateContext: resourceForemanEnvironmentCreate,
		ReadContext:   resourceForemanEnvironmentRead,
		UpdateContext: resourceForemanEnvironmentUpdate,
		DeleteContext: resourceForemanEnvironmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{

			autodoc.MetaAttribute: {
				Type:     schema.TypeBool,
				Computed: true,
				Description: fmt.Sprintf(
					"%s A puppet environment, branch.",
					autodoc.MetaSummary,
				),
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: fmt.Sprintf(
					"Name of the environment. Usually maps to the name of "+
						"a puppet branch. "+
						"%s \"production\"",
					autodoc.MetaExample,
				),
			},
		},
	}
}

// -----------------------------------------------------------------------------
// Conversion Helpers
// -----------------------------------------------------------------------------

// buildForemanEnvironment constructs a ForemanEnvironment reference from a
// resource data reference.  The struct's  members are populated from the data
// populated in the resource data.  Missing members will be left to the zero
// value for that member's type.
func buildForemanEnvironment(d *schema.ResourceData) *api.ForemanEnvironment {
	log.Tracef("resource_foreman_environment.go#buildForemanEnvironment")

	environment := api.ForemanEnvironment{}

	obj := buildForemanObject(d)
	environment.ForemanObject = *obj

	var attr interface{}
	var ok bool

	if attr, ok = d.GetOk("name"); ok {
		environment.Name = attr.(string)
	}

	return &environment
}

// setResourceDataFromForemanEnvironment sets a ResourceData's attributes from
// the attributes of the supplied ForemanEnvironment reference
func setResourceDataFromForemanEnvironment(d *schema.ResourceData, fe *api.ForemanEnvironment) {
	log.Tracef("resource_foreman_environment.go#setResourceDataFromForemanEnvironment")

	d.SetId(strconv.Itoa(fe.Id))
	d.Set("name", fe.Name)
}

// -----------------------------------------------------------------------------
// Resource CRUD Operations
// -----------------------------------------------------------------------------

func resourceForemanEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("resource_foreman_environment.go#Create")

	client := meta.(*api.Client)
	e := buildForemanEnvironment(d)

	log.Debugf("ForemanEnvironment: [%+v]", e)

	createdEnv, createErr := client.CreateEnvironment(ctx, e)
	if createErr != nil {
		return diag.FromErr(createErr)
	}

	log.Debugf("Created ForemanEnvironment: [%+v]", createdEnv)

	setResourceDataFromForemanEnvironment(d, createdEnv)

	return nil
}

func resourceForemanEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("resource_foreman_environment.go#Read")

	client := meta.(*api.Client)
	e := buildForemanEnvironment(d)

	log.Debugf("ForemanEnvironment: [%+v]", e)

	readEnvironment, readErr := client.ReadEnvironment(ctx, e.Id)
	if readErr != nil {
		return diag.FromErr(api.CheckDeleted(d, readErr))
	}

	log.Debugf("Read ForemanEnvironment: [%+v]", readEnvironment)

	setResourceDataFromForemanEnvironment(d, readEnvironment)

	return nil
}

func resourceForemanEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("resource_foreman_environment.go#Update")

	client := meta.(*api.Client)
	e := buildForemanEnvironment(d)

	log.Debugf("ForemanEnvironment: [%+v]", e)

	updatedEnv, updateErr := client.UpdateEnvironment(ctx, e)
	if updateErr != nil {
		return diag.FromErr(updateErr)
	}

	log.Debugf("Updated ForemanEnvironment: [%+v]", updatedEnv)

	setResourceDataFromForemanEnvironment(d, updatedEnv)

	return nil
}

func resourceForemanEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Tracef("resource_foreman_environment.go#Delete")

	client := meta.(*api.Client)
	e := buildForemanEnvironment(d)

	// NOTE(ALL): d.SetId("") is automatically called by terraform assuming delete
	//   returns no errors

	return diag.FromErr(api.CheckDeleted(d, client.DeleteEnvironment(ctx, e.Id)))
}
