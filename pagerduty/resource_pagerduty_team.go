package pagerduty

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyTeamCreate,
		Read:   resourcePagerDutyTeamRead,
		Update: resourcePagerDutyTeamUpdate,
		Delete: resourcePagerDutyTeamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
		},
	}
}

func buildTeamStruct(d *schema.ResourceData) *pagerduty.Team {
	team := &pagerduty.Team{
		Name: d.Get("name").(string),
	}

	if attr, ok := d.GetOk("description"); ok {
		team.Description = attr.(string)
	}

	return team
}

func resourcePagerDutyTeamCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	team := buildTeamStruct(d)

	log.Printf("[INFO] Creating PagerDuty team %s", team.Name)

	team, _, err := client.Teams.Create(team)
	if err != nil {
		return err
	}

	d.SetId(team.ID)

	return nil

}

func resourcePagerDutyTeamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty team %s", d.Id())

	team, _, err := client.Teams.Get(d.Id())
	if err != nil {
		return handleNotFoundError(err, d)
	}

	d.Set("name", team.Name)
	d.Set("description", team.Description)

	return nil
}

func resourcePagerDutyTeamUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	team := buildTeamStruct(d)

	log.Printf("[INFO] Updating PagerDuty team %s", d.Id())

	if _, _, err := client.Teams.Update(d.Id(), team); err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyTeamDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Deleting PagerDuty team %s", d.Id())

	if _, err := client.Teams.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}
