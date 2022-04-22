package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kayteh/podio-go"
)

func resourceSpace() *schema.Resource {
	return &schema.Resource{
		Description: `Manage a Space/Workspace within Podio. See https://developers.podio.com/doc/spaces/ for more information.
		
			Caveats:
			- Deleting may be expected to fail due to API trust limits. Set *ignore_delete_errors* to true to ignore errors.
			
		`,

		CreateContext: resourceSpaceCreate,
		ReadContext:   resourceSpaceRead,
		UpdateContext: resourceSpaceUpdate,
		DeleteContext: resourceSpaceDelete,

		Schema: map[string]*schema.Schema{
			"org_id": {
				Description: "The ID of the Organization the Space belongs to. Changing this forces a new resource to be created.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "Name of the space. Changing this does not update the space URL.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"url": {
				Description: "URL of the space",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"url_label": {
				Description: "URL label/slug of the space",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"privacy": {
				Description: "Privacy of the space, one of `open`, `closed`. Defaults to `closed`",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "closed",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "open" && value != "closed" {
						errors = append(errors, fmt.Errorf("%q must be one of `open`, `closed`", k))
					}
					return
				},
			},
			"auto_join": {
				Description: "If true, new employees automatically join this space. Defaults to `false`",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"post_on_new_app": {
				Description: "If true, new apps are posted as a status update to this space. Defaults to `false`",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"post_on_new_member": {
				Description: "If true, new members are posted as a status update to this space. Defaults to `false`",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"ignore_delete_errors": {
				Description: "If true, errors are ignored when deleting a space. Defaults to `false`",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceSpaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*podio.Client)

	space, err := client.CreateSpace(podio.CreateSpaceParams{
		Name:            d.Get("name").(string),
		OrgID:           d.Get("org_id").(int),
		Privacy:         d.Get("privacy").(string),
		AutoJoin:        d.Get("auto_join").(bool),
		PostOnNewApp:    d.Get("post_on_new_app").(bool),
		PostOnNewMember: d.Get("post_on_new_member").(bool),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", space.ID))
	d.Set("url", space.URL)
	d.Set("url_label", space.URLLabel)
	d.Set("name", space.Name)
	d.Set("privacy", space.Privacy)
	d.Set("auto_join", space.AutoJoin)
	d.Set("post_on_new_app", space.PostOnNewApp)
	d.Set("post_on_new_member", space.PostOnNewMember)

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a space")

	return nil
}

func resourceSpaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*podio.Client)

	space, err := client.GetSpace(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("url", space.URL)
	d.Set("url_label", space.URLLabel)
	d.Set("name", space.Name)
	d.Set("privacy", space.Privacy)
	d.Set("auto_join", space.AutoJoin)
	d.Set("post_on_new_app", space.PostOnNewApp)
	d.Set("post_on_new_member", space.PostOnNewMember)

	tflog.Trace(ctx, "read a space")

	return nil
}

func resourceSpaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*podio.Client)

	space, err := client.UpdateSpace(d.Id(), podio.CreateSpaceParams{
		Name:            d.Get("name").(string),
		Privacy:         d.Get("privacy").(string),
		AutoJoin:        d.Get("auto_join").(bool),
		PostOnNewApp:    d.Get("post_on_new_app").(bool),
		PostOnNewMember: d.Get("post_on_new_member").(bool),
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("url", space.URL)
	d.Set("url_label", space.URLLabel)
	d.Set("name", space.Name)
	d.Set("privacy", space.Privacy)
	d.Set("auto_join", space.AutoJoin)
	d.Set("post_on_new_app", space.PostOnNewApp)
	d.Set("post_on_new_member", space.PostOnNewMember)

	return diag.Errorf("not implemented")
}

func resourceSpaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*podio.Client)

	err := client.DeleteSpace(d.Id())
	if err != nil {
		if d.Get("ignore_delete_errors").(bool) {
			tflog.Info(ctx, "ignoring error when deleting space")
			return nil
		}
		return diag.FromErr(err)
	}

	tflog.Trace(ctx, "deleted space")

	return nil
}
