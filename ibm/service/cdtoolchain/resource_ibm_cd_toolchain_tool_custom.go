// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cdtoolchain

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/continuous-delivery-go-sdk/cdtoolchainv2"
)

func ResourceIBMCdToolchainToolCustom() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCdToolchainToolCustomCreate,
		ReadContext:   resourceIBMCdToolchainToolCustomRead,
		UpdateContext: resourceIBMCdToolchainToolCustomUpdate,
		DeleteContext: resourceIBMCdToolchainToolCustomDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"toolchain_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_toolchain_tool_custom", "toolchain_id"),
				Description:  "ID of the toolchain to bind the tool to.",
			},
			"parameters": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "Unique key-value pairs representing parameters to be used to create the tool.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type the name of the tool that you are integrating; for example: Delivery Pipeline.",
						},
						"lifecycle_phase": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Select the lifecycle phase of the IBM Cloud Garage Method that is the most closely associated with this tool.",
						},
						"image_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type the URL of the icon to show on your tool integration's card.",
						},
						"documentation_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type the URL for your tool's documentation.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type a name for this specific tool integration; for example: My Build and Deploy Pipeline.",
						},
						"dashboard_url": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type the URL that you want to navigate to when you click the tool integration card.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type a description for the tool instance.",
						},
						"additional_properties": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "(Advanced) Type any information that is needed to integrate with other tools in your toolchain.",
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_toolchain_tool_custom", "name"),
				Description:  "Name of tool.",
			},
			"resource_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource group where tool can be found.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tool CRN.",
			},
			"toolchain_crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of toolchain which the tool is bound to.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URI representing the tool.",
			},
			"referent": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information on URIs to access this resource through the UI or API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ui_href": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "URI representing the this resource through the UI.",
						},
						"api_href": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "URI representing the this resource through an API.",
						},
					},
				},
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Latest tool update timestamp.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current configuration state of the tool.",
			},
			"tool_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tool ID.",
			},
		},
	}
}

func ResourceIBMCdToolchainToolCustomValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "toolchain_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([^\x00-\x7F]|[a-zA-Z0-9-._ ])+$`,
			MinValueLength:             0,
			MaxValueLength:             128,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_cd_toolchain_tool_custom", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMCdToolchainToolCustomCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createToolOptions := &cdtoolchainv2.CreateToolOptions{}

	createToolOptions.SetToolchainID(d.Get("toolchain_id").(string))
	createToolOptions.SetToolTypeID("customtool")
	remapFields := map[string]string{
		"lifecycle_phase":       "lifecyclePhase",
		"image_url":             "imageUrl",
		"documentation_url":     "documentationUrl",
		"additional_properties": "additional-properties",
	}
	parametersModel := GetParametersForCreate(d, ResourceIBMCdToolchainToolCustom(), remapFields)
	createToolOptions.SetParameters(parametersModel)
	if _, ok := d.GetOk("name"); ok {
		createToolOptions.SetName(d.Get("name").(string))
	}

	toolchainToolPost, response, err := cdToolchainClient.CreateToolWithContext(context, createToolOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateToolWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateToolWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createToolOptions.ToolchainID, *toolchainToolPost.ID))

	return resourceIBMCdToolchainToolCustomRead(context, d, meta)
}

func resourceIBMCdToolchainToolCustomRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getToolByIDOptions := &cdtoolchainv2.GetToolByIDOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getToolByIDOptions.SetToolchainID(parts[0])
	getToolByIDOptions.SetToolID(parts[1])

	toolchainTool, response, err := cdToolchainClient.GetToolByIDWithContext(context, getToolByIDOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetToolByIDWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetToolByIDWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("toolchain_id", toolchainTool.ToolchainID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting toolchain_id: %s", err))
	}
	remapFields := map[string]string{
		"lifecycle_phase":       "lifecyclePhase",
		"image_url":             "imageUrl",
		"documentation_url":     "documentationUrl",
		"additional_properties": "additional-properties",
	}
	parametersMap := GetParametersFromRead(toolchainTool.Parameters, ResourceIBMCdToolchainToolCustom(), remapFields)
	if err = d.Set("parameters", []map[string]interface{}{parametersMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting parameters: %s", err))
	}
	if err = d.Set("name", toolchainTool.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("resource_group_id", toolchainTool.ResourceGroupID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group_id: %s", err))
	}
	if err = d.Set("crn", toolchainTool.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("toolchain_crn", toolchainTool.ToolchainCRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting toolchain_crn: %s", err))
	}
	if err = d.Set("href", toolchainTool.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	referentMap, err := resourceIBMCdToolchainToolCustomToolModelReferentToMap(toolchainTool.Referent)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("referent", []map[string]interface{}{referentMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting referent: %s", err))
	}
	if err = d.Set("updated_at", flex.DateTimeToString(toolchainTool.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}
	if err = d.Set("state", toolchainTool.State); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting state: %s", err))
	}
	if err = d.Set("tool_id", toolchainTool.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting tool_id: %s", err))
	}

	return nil
}

func resourceIBMCdToolchainToolCustomUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	updateToolOptions := &cdtoolchainv2.UpdateToolOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateToolOptions.SetToolchainID(parts[0])
	updateToolOptions.SetToolID(parts[1])

	hasChange := false

	patchVals := &cdtoolchainv2.ToolchainToolPrototypePatch{}
	if d.HasChange("toolchain_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "toolchain_id"))
	}
	if d.HasChange("parameters") {
		remapFields := map[string]string{
			"lifecycle_phase":       "lifecyclePhase",
			"image_url":             "imageUrl",
			"documentation_url":     "documentationUrl",
			"additional_properties": "additional-properties",
		}
		parameters := GetParametersForUpdate(d, ResourceIBMCdToolchainToolCustom(), remapFields)
		patchVals.Parameters = parameters
		hasChange = true
	}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}

	if hasChange {
		updateToolOptions.ToolchainToolPrototypePatch, _ = patchVals.AsPatch()
		_, response, err := cdToolchainClient.UpdateToolWithContext(context, updateToolOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateToolWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateToolWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMCdToolchainToolCustomRead(context, d, meta)
}

func resourceIBMCdToolchainToolCustomDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteToolOptions := &cdtoolchainv2.DeleteToolOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteToolOptions.SetToolchainID(parts[0])
	deleteToolOptions.SetToolID(parts[1])

	response, err := cdToolchainClient.DeleteToolWithContext(context, deleteToolOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteToolWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteToolWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMCdToolchainToolCustomToolModelReferentToMap(model *cdtoolchainv2.ToolModelReferent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UIHref != nil {
		modelMap["ui_href"] = model.UIHref
	}
	if model.APIHref != nil {
		modelMap["api_href"] = model.APIHref
	}
	return modelMap, nil
}