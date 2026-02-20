package bridge

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jfrog/terraform-provider-shared/util"
	utilfw "github.com/jfrog/terraform-provider-shared/util/fw"
)

const (
	StopEndpoint = "bridge-client/api/v1/bridges/{bridge_ID}/stop"
	StopEndpoint  = "bridge-client/api/v1/bridges/{bridge_ID}/stop"
)

func NewBridgesStopResource() resource.Resource {
	return &BridgesStopResource{
		TypeName: "bridge_stop",
	}
}

type BridgesStopResource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

type BridgesStopResourceModel struct {
	BridgeID types.String `tfsdk:"bridge_ID"`
}

type StopRequestAPIModel struct {
	BridgeID string `json:"bridge_ID"`
}

type StopAPIModel struct {
	BridgeID string `json:"bridge_ID"`
}

func (r *BridgesStopResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName
}

func (r *BridgesStopResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bridge_ID": schema.StringAttribute{
				Required: true,
				Description: "The bridge_ID of the resource.",
			},
		},
		MarkdownDescription: "Manages stop in JFrog Bridge.",
	}
}

func (r *BridgesStopResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.ProviderData = req.ProviderData.(util.ProviderMetadata)
}


func (r *BridgesStopResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	go util.SendUsageResourceCreate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan BridgesStopResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody := BridgesStopRequestAPIModel{
		BridgeID: plan.BridgeID.ValueString(),
	}

	var result BridgesStopAPIModel

	response, err := r.ProviderData.Client.R().
		SetBody(requestBody).
		SetResult(&result).
		Post(StopEndpoint)
	if err != nil {
		utilfw.UnableToCreateResourceError(resp, err.Error())
		return
	}

	if response.IsError() {
		utilfw.UnableToCreateResourceError(resp, response.String())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}


func (r *BridgesStopResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	go util.SendUsageResourceRead(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state BridgesStopResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result BridgesStopAPIModel

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
			"bridge_ID": state.BridgeID.ValueString(),
		}).
		SetResult(&result).
		Get(StopEndpoint)
	if err != nil {
		utilfw.UnableToRefreshResourceError(resp, err.Error())
		return
	}

	if response.StatusCode() == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}

	if response.IsError() {
		utilfw.UnableToRefreshResourceError(resp, response.String())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}




