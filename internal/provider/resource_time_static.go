package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-provider-time/internal/validators/timevalidator"
)

var (
	_ resource.Resource                = (*timeStaticResource)(nil)
	_ resource.ResourceWithImportState = (*timeStaticResource)(nil)
)

func NewTimeStaticResource() resource.Resource {
	return &timeStaticResource{}
}

type timeStaticResource struct{}

func (t timeStaticResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_static"
}

func (t timeStaticResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Manages a static time resource, which keeps a locally sourced UTC timestamp stored in the Terraform state. " +
			"This prevents perpetual differences caused by using " +
			"the [`timestamp()` function](https://www.terraform.io/docs/configuration/functions/timestamp.html).",
		Attributes: map[string]tfsdk.Attribute{
			"day": {
				Description: "Number day of timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"hour": {
				Description: "Number hour of timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"triggers": {
				Description: "Arbitrary map of values that, when changed, will trigger a new base timestamp value to be saved. " +
					"See [the main provider documentation](../index.md) for more information.",
				Type: types.MapType{
					ElemType: types.StringType,
				},
				Optional: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"minute": {
				Description: "Number minute of timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"month": {
				Description: "Number month of timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"rfc3339": {
				Description: "Base timestamp in " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format " +
					"(see [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) e.g., " +
					"`YYYY-MM-DDTHH:MM:SSZ`). Defaults to the current time.",
				Type:     types.StringType,
				Optional: true,
				Computed: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
				Validators: []tfsdk.AttributeValidator{
					timevalidator.IsRFC3339Time(),
				},
			},
			"second": {
				Description: "Number second of timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"unix": {
				Description: "Number of seconds since epoch time, e.g. `1581489373`.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"year": {
				Description: "Number year of timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"id": {
				Description: "RFC3339 format of the offset timestamp, e.g. `2020-02-12T06:36:13Z`.",
				Type:        types.StringType,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.UseStateForUnknown(),
				},
			},
		},
	}, nil
}

func (t timeStaticResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	timestamp, err := time.Parse(time.RFC3339, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Import time static error",
			"The id that was supplied could not be parsed as RFC3339.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	formattedTimestamp := timestamp.Format(time.RFC3339)

	state := timeStaticModelV0{
		Year:    types.Int64Value(int64(timestamp.Year())),
		Month:   types.Int64Value(int64(timestamp.Month())),
		Day:     types.Int64Value(int64(timestamp.Day())),
		Hour:    types.Int64Value(int64(timestamp.Hour())),
		Minute:  types.Int64Value(int64(timestamp.Minute())),
		Second:  types.Int64Value(int64(timestamp.Second())),
		RFC3339: types.StringValue(formattedTimestamp),
		Unix:    types.Int64Value(timestamp.Unix()),
		ID:      types.StringValue(formattedTimestamp),
	}
	state.Triggers = types.MapValueMust(types.StringType, map[string]attr.Value{})

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (t timeStaticResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan timeStaticModelV0

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	timestamp := time.Now().UTC()

	if plan.RFC3339.ValueString() != "" {
		var err error

		if timestamp, err = time.Parse(time.RFC3339, plan.RFC3339.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				"Create time static error",
				"The rfc3339 timestamp that was supplied could not be parsed as RFC3339.\n\n+"+
					fmt.Sprintf("Original Error: %s", err),
			)
			return
		}
	}

	formattedTimestamp := timestamp.Format(time.RFC3339)

	state := timeStaticModelV0{
		Triggers: plan.Triggers,
		Year:     types.Int64Value(int64(timestamp.Year())),
		Month:    types.Int64Value(int64(timestamp.Month())),
		Day:      types.Int64Value(int64(timestamp.Day())),
		Hour:     types.Int64Value(int64(timestamp.Hour())),
		Minute:   types.Int64Value(int64(timestamp.Minute())),
		Second:   types.Int64Value(int64(timestamp.Second())),
		RFC3339:  types.StringValue(formattedTimestamp),
		Unix:     types.Int64Value(timestamp.Unix()),
		ID:       types.StringValue(formattedTimestamp),
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (t timeStaticResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

func (t timeStaticResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data timeStaticModelV0

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (t timeStaticResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

type timeStaticModelV0 struct {
	Day      types.Int64  `tfsdk:"day"`
	Hour     types.Int64  `tfsdk:"hour"`
	Triggers types.Map    `tfsdk:"triggers"`
	Minute   types.Int64  `tfsdk:"minute"`
	Month    types.Int64  `tfsdk:"month"`
	RFC3339  types.String `tfsdk:"rfc3339"`
	Second   types.Int64  `tfsdk:"second"`
	Unix     types.Int64  `tfsdk:"unix"`
	Year     types.Int64  `tfsdk:"year"`
	ID       types.String `tfsdk:"id"`
}
