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
	BlacklistsEndpoint = "bridge-server/api/v1/blacklist"
	BlacklistEndpoint  = "bridge-server/api/v1/blacklist/{bridge_ID}"
)

func NewBlacklistResource() resource.Resource {
	return &BlacklistResource{
		TypeName: "bridge_blacklist",
	}
}

type BlacklistResource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

type BlacklistResourceModel struct {
	BridgeID types.String `tfsdk:"bridge_ID"`
}

type BlacklistRequestAPIModel struct {
	BridgeID string `json:"bridge_ID"`
}

type BlacklistAPIModel struct {
	BridgeID string `json:"bridge_ID"`
}

func (r *BlacklistResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName
}

func (r *BlacklistResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bridge_ID": schema.StringAttribute{
				Required: true,
				Description: "The bridge_ID of the resource.",
			},
		},
		MarkdownDescription: "Manages blacklist in JFrog Bridge.",
	}
}

func (r *BlacklistResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.ProviderData = req.ProviderData.(util.ProviderMetadata)
}



func (r *BlacklistResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	go util.SendUsageResourceRead(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state BlacklistResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result BlacklistAPIModel

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
			"bridge_ID": state.BridgeID.ValueString(),
		}).
		SetResult(&result).
		Get(BlacklistEndpoint)
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


func (r *BlacklistResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	go util.SendUsageResourceUpdate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan BlacklistResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody := BlacklistRequestAPIModel{
		BridgeID: plan.BridgeID.ValueString(),
	}

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
			"bridge_ID": plan.BridgeID.ValueString(),
		}).
		SetBody(requestBody).
		Put(BlacklistEndpoint)
	if err != nil {
		utilfw.UnableToUpdateResourceError(resp, err.Error())
		return
	}

	if response.IsError() {
		utilfw.UnableToUpdateResourceError(resp, response.String())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}



func (r *BlacklistResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	go util.SendUsageResourceDelete(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state BlacklistResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
			"bridge_ID": state.BridgeID.ValueString(),
		}).
		Delete(BlacklistEndpoint)
	if err != nil {
		utilfw.UnableToDeleteResourceError(resp, err.Error())
		return
	}

	if response.StatusCode() == http.StatusNotFound {
		return
	}

	if response.IsError() {
		utilfw.UnableToDeleteResourceError(resp, response.String())
		return
	}
}

