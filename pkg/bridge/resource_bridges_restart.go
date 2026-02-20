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
	RestartEndpoint = "bridge-client/api/v1/bridges/{bridge_ID}/restart"
	RestartEndpoint  = "bridge-client/api/v1/bridges/{bridge_ID}/restart"
)

func NewBridgesRestartResource() resource.Resource {
	return &BridgesRestartResource{
		TypeName: "bridge_restart",
	}
}

type BridgesRestartResource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

type BridgesRestartResourceModel struct {
	BridgeID types.String `tfsdk:"bridge_ID"`
}

type RestartRequestAPIModel struct {
	BridgeID string `json:"bridge_ID"`
}

type RestartAPIModel struct {
	BridgeID string `json:"bridge_ID"`
}

func (r *BridgesRestartResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName
}

func (r *BridgesRestartResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bridge_ID": schema.StringAttribute{
				Required: true,
				Description: "The bridge_ID of the resource.",
			},
		},
		MarkdownDescription: "Manages restart in JFrog Bridge.",
	}
}

func (r *BridgesRestartResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.ProviderData = req.ProviderData.(util.ProviderMetadata)
}


func (r *BridgesRestartResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	go util.SendUsageResourceCreate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan BridgesRestartResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody := BridgesRestartRequestAPIModel{
		BridgeID: plan.BridgeID.ValueString(),
	}

	var result BridgesRestartAPIModel

	response, err := r.ProviderData.Client.R().
		SetBody(requestBody).
		SetResult(&result).
		Post(RestartEndpoint)
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


func (r *BridgesRestartResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	go util.SendUsageResourceRead(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state BridgesRestartResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result BridgesRestartAPIModel

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
			"bridge_ID": state.BridgeID.ValueString(),
		}).
		SetResult(&result).
		Get(RestartEndpoint)
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




