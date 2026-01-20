// Copyright (c) JFrog Ltd. (2025)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bridge

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jfrog/terraform-provider-shared/util"
	utilfw "github.com/jfrog/terraform-provider-shared/util/fw"
)

const bridgeBasePath = "/bridge-client/api/v1/bridges"

var _ resource.Resource = &BridgeResource{}
var _ resource.ResourceWithImportState = &BridgeResource{}

func NewBridgeResource() resource.Resource {
	return &BridgeResource{}
}

type BridgeResource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

type bridgeProxyModel struct {
	Enabled             types.Bool   `tfsdk:"enabled"`
	CacheExpirationSecs types.Int64  `tfsdk:"cache_expiration_secs"`
	Key                 types.String `tfsdk:"key"`
	SchemeOverride      types.String `tfsdk:"scheme_override"`
}

type bridgeRemoteModel struct {
	Url      types.String      `tfsdk:"url"`
	Insecure types.Bool        `tfsdk:"insecure"`
	Proxy    *bridgeProxyModel `tfsdk:"proxy"`
}

type bridgeLocalModel struct {
	Url                types.String `tfsdk:"url"`
	AnonymousEndpoints types.List   `tfsdk:"anonymous_endpoints"`
}

type bridgeTargetUsageModel struct {
	Low  types.Int64 `tfsdk:"low"`
	High types.Int64 `tfsdk:"high"`
}

type bridgeTunnelCreationJobModel struct {
	IntervalMinutes types.Int64 `tfsdk:"interval_minutes"`
}

type bridgeTunnelClosingJobModel struct {
	CronExpr              types.String `tfsdk:"cron_expr"`
	AllowCloseUsedTunnels types.Bool   `tfsdk:"allow_close_used_tunnels"`
}

type bridgeJobsModel struct {
	TunnelCreation *bridgeTunnelCreationJobModel `tfsdk:"tunnel_creation"`
	TunnelClosing  *bridgeTunnelClosingJobModel  `tfsdk:"tunnel_closing"`
}

type BridgeResourceModel struct {
	ID           types.String            `tfsdk:"id"`
	BridgeID     types.String            `tfsdk:"bridge_id"`
	Remote       *bridgeRemoteModel      `tfsdk:"remote"`
	Local        *bridgeLocalModel       `tfsdk:"local"`
	PairingToken types.String            `tfsdk:"pairing_token"`
	MinTunnels   types.Int64             `tfsdk:"min_tunnels"`
	MaxTunnels   types.Int64             `tfsdk:"max_tunnels"`
	TargetUsage  *bridgeTargetUsageModel `tfsdk:"target_usage"`
	Jobs         *bridgeJobsModel        `tfsdk:"jobs"`
	CreatedAt    types.String            `tfsdk:"created_at"`
}

type bridgeProxyAPIModel struct {
	Enabled            *bool  `json:"enabled,omitempty"`
	CacheExpirationSec *int64 `json:"cache_expiration_secs,omitempty"`
	Key                string `json:"key,omitempty"`
	SchemeOverride     string `json:"scheme_override,omitempty"`
}

type bridgeRemoteAPIModel struct {
	Url      string               `json:"url,omitempty"`
	Token    string               `json:"token,omitempty"`
	Insecure *bool                `json:"insecure,omitempty"`
	Proxy    *bridgeProxyAPIModel `json:"proxy,omitempty"`
}

type bridgeLocalAPIModel struct {
	Url                string   `json:"url,omitempty"`
	AnonymousEndpoints []string `json:"anonymous_endpoints,omitempty"`
	DialTimeoutSecs    *int64   `json:"dial_timeout_secs,omitempty"`
}

type bridgeTargetUsageAPIModel struct {
	Low  *int64 `json:"low,omitempty"`
	High *int64 `json:"high,omitempty"`
}

type bridgeJobsAPIModel struct {
	TunnelCreation *struct {
		IntervalMinutes *int64 `json:"interval_minutes,omitempty"`
	} `json:"tunnel_creation,omitempty"`
	TunnelClosing *struct {
		CronExpr              string `json:"cron_expr,omitempty"`
		AllowCloseUsedTunnels *bool  `json:"allow_close_used_tunnels,omitempty"`
	} `json:"tunnel_closing,omitempty"`
}

type bridgeConfigAPIModel struct {
	BridgeID    string                     `json:"bridge_id,omitempty"`
	Type        string                     `json:"type,omitempty"`
	Remote      *bridgeRemoteAPIModel      `json:"remote,omitempty"`
	Local       *bridgeLocalAPIModel       `json:"local,omitempty"`
	MinTunnels  *int64                     `json:"min_tunnels,omitempty"`
	MaxTunnels  *int64                     `json:"max_tunnels,omitempty"`
	TargetUsage *bridgeTargetUsageAPIModel `json:"target_usage,omitempty"`
	Jobs        *bridgeJobsAPIModel        `json:"jobs,omitempty"`
}

// bridgeCreateRequestModel is for POST /bridges (create) - uses simple URL strings
type bridgeCreateRequestModel struct {
	BridgeID     string `json:"bridge_id"`
	Remote       string `json:"remote"`
	Local        string `json:"local"`
	PairingToken string `json:"pairing_token"`
}

// bridgeUpdateRequestModel is for PATCH /bridges/{id} (update) - uses object structures
type bridgeUpdateRequestModel struct {
	Remote      *bridgeRemoteAPIModel      `json:"remote,omitempty"`
	Local       *bridgeLocalAPIModel       `json:"local,omitempty"`
	MinTunnels  *int64                     `json:"min_tunnels,omitempty"`
	MaxTunnels  *int64                     `json:"max_tunnels,omitempty"`
	TargetUsage *bridgeTargetUsageAPIModel `json:"target_usage,omitempty"`
	Jobs        *bridgeJobsAPIModel        `json:"jobs,omitempty"`
}

type bridgeDebugResponse struct {
	Bridges []struct {
		ID      string               `json:"id"`
		Config  bridgeConfigAPIModel `json:"config"`
		Created string               `json:"created_at"`
	} `json:"bridges"`
}

func (r *BridgeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName // resource name is just "bridge"
	r.TypeName = resp.TypeName
}

func (r *BridgeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage JFrog Bridges via the bridge-client API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Internal Terraform resource ID (uses bridge_id).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bridge_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique identifier of the bridge. Changing this forces a new resource to be created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"remote": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Remote (bridge server) configuration.",
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "URL of the bridge server (remote JPD).",
					},
					"insecure": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Allow insecure TLS when connecting to the remote.",
					},
					"proxy": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Proxy configuration used by the remote connection.",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether proxy is enabled.",
							},
							"cache_expiration_secs": schema.Int64Attribute{
								Optional:            true,
								MarkdownDescription: "Proxy cache expiration in seconds.",
							},
							"key": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Proxy key.",
							},
							"scheme_override": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Override proxy scheme (e.g., http/https).",
							},
						},
					},
				},
			},
			"local": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Local (bridge client) configuration.",
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "URL of the bridge client (local JPD).",
					},
					"anonymous_endpoints": schema.ListAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "List of anonymous endpoints allowed through the bridge.",
					},
				},
			},
			"pairing_token": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Pairing token generated on the bridge server. Required on create; removed from state after creation.",
			},
			"min_tunnels": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Minimum tunnels.",
			},
			"max_tunnels": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Maximum tunnels.",
			},
			"target_usage": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Target usage thresholds.",
				Attributes: map[string]schema.Attribute{
					"low": schema.Int64Attribute{
						Optional: true,
					},
					"high": schema.Int64Attribute{
						Optional: true,
					},
				},
			},
			"jobs": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Job configuration for tunnel creation/closing.",
				Attributes: map[string]schema.Attribute{
					"tunnel_creation": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"interval_minutes": schema.Int64Attribute{
								Optional:            true,
								MarkdownDescription: "Interval in minutes for tunnel creation.",
							},
						},
					},
					"tunnel_closing": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"cron_expr": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Cron expression for tunnel closing.",
							},
							"allow_close_used_tunnels": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether to allow closing used tunnels.",
							},
						},
					},
				},
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp when the bridge was created (from debug snapshot).",
			},
		},
	}
}

func (r *BridgeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.ProviderData = req.ProviderData.(util.ProviderMetadata)
}

func (r *BridgeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	go util.SendUsageResourceCreate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan BridgeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.PairingToken.IsUnknown() || plan.PairingToken.IsNull() || plan.PairingToken.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing pairing_token",
			"pairing_token is required when defining a new bridge. Generate it on the bridge server and supply it on creation.",
		)
		return
	}

	// Create uses simple URL strings for remote/local
	payload := bridgeCreateRequestModel{
		BridgeID:     plan.BridgeID.ValueString(),
		Remote:       plan.Remote.Url.ValueString(),
		Local:        plan.Local.Url.ValueString(),
		PairingToken: plan.PairingToken.ValueString(),
	}
	response, err := r.ProviderData.Client.R().
		SetBody(payload).
		Post(bridgeBasePath)

	if err != nil {
		utilfw.UnableToCreateResourceError(resp, err.Error())
		return
	}
	if response.IsError() {
		utilfw.UnableToCreateResourceError(resp, response.String())
		return
	}

	plan.ID = plan.BridgeID
	plan.CreatedAt = types.StringNull()

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BridgeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	go util.SendUsageResourceRead(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state BridgeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Without the debug endpoint, keep existing state. User-driven updates will reconcile configuration.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BridgeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	go util.SendUsageResourceUpdate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan BridgeResourceModel
	var state BridgeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.PairingToken.IsNull() && plan.PairingToken.ValueString() != "" {
		resp.Diagnostics.AddWarning(
			"pairing_token ignored on update",
			"pairing_token is only used when defining a new bridge and will not be sent during updates.",
		)
	}

	plan.ID = state.ID
	plan.BridgeID = state.BridgeID
	plan.PairingToken = state.PairingToken // preserve from state
	plan.CreatedAt = state.CreatedAt

	// Update uses object structures for remote/local
	payload := buildUpdateRequest(plan)
	endpoint := fmt.Sprintf("%s/%s", bridgeBasePath, plan.BridgeID.ValueString())

	response, err := r.ProviderData.Client.R().
		SetBody(payload).
		Patch(endpoint)
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

func (r *BridgeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	go util.SendUsageResourceDelete(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state BridgeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := fmt.Sprintf("%s/%s", bridgeBasePath, state.BridgeID.ValueString())
	response, err := r.ProviderData.Client.R().
		Delete(endpoint)

	if err != nil {
		utilfw.UnableToDeleteResourceError(resp, err.Error())
		return
	}
	if response.IsError() {
		utilfw.UnableToDeleteResourceError(resp, response.String())
		return
	}
}

func (r *BridgeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("bridge_id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("created_at"), types.StringNull())...)
}

func buildUpdateRequest(model BridgeResourceModel) bridgeUpdateRequestModel {
	req := bridgeUpdateRequestModel{}

	if !model.MinTunnels.IsNull() && model.MinTunnels.ValueInt64() != 0 {
		val := model.MinTunnels.ValueInt64()
		req.MinTunnels = &val
	}
	if !model.MaxTunnels.IsNull() && model.MaxTunnels.ValueInt64() != 0 {
		val := model.MaxTunnels.ValueInt64()
		req.MaxTunnels = &val
	}
	if model.TargetUsage != nil {
		req.TargetUsage = &bridgeTargetUsageAPIModel{}
		if !model.TargetUsage.Low.IsNull() {
			val := model.TargetUsage.Low.ValueInt64()
			req.TargetUsage.Low = &val
		}
		if !model.TargetUsage.High.IsNull() {
			val := model.TargetUsage.High.ValueInt64()
			req.TargetUsage.High = &val
		}
	}
	if model.Jobs != nil {
		req.Jobs = &bridgeJobsAPIModel{}
		if model.Jobs.TunnelCreation != nil {
			req.Jobs.TunnelCreation = &struct {
				IntervalMinutes *int64 `json:"interval_minutes,omitempty"`
			}{}
			if !model.Jobs.TunnelCreation.IntervalMinutes.IsNull() {
				val := model.Jobs.TunnelCreation.IntervalMinutes.ValueInt64()
				req.Jobs.TunnelCreation.IntervalMinutes = &val
			}
		}
		if model.Jobs.TunnelClosing != nil {
			req.Jobs.TunnelClosing = &struct {
				CronExpr              string `json:"cron_expr,omitempty"`
				AllowCloseUsedTunnels *bool  `json:"allow_close_used_tunnels,omitempty"`
			}{}
			if model.Jobs.TunnelClosing.CronExpr.ValueString() != "" {
				req.Jobs.TunnelClosing.CronExpr = model.Jobs.TunnelClosing.CronExpr.ValueString()
			}
			if !model.Jobs.TunnelClosing.AllowCloseUsedTunnels.IsNull() {
				val := model.Jobs.TunnelClosing.AllowCloseUsedTunnels.ValueBool()
				req.Jobs.TunnelClosing.AllowCloseUsedTunnels = &val
			}
		}
	}

	if model.Remote != nil {
		req.Remote = &bridgeRemoteAPIModel{
			Url: model.Remote.Url.ValueString(),
		}
		if !model.Remote.Insecure.IsNull() {
			val := model.Remote.Insecure.ValueBool()
			req.Remote.Insecure = &val
		}
		if model.Remote.Proxy != nil {
			req.Remote.Proxy = &bridgeProxyAPIModel{}
			if !model.Remote.Proxy.Enabled.IsNull() {
				val := model.Remote.Proxy.Enabled.ValueBool()
				req.Remote.Proxy.Enabled = &val
			}
			if !model.Remote.Proxy.CacheExpirationSecs.IsNull() {
				val := model.Remote.Proxy.CacheExpirationSecs.ValueInt64()
				req.Remote.Proxy.CacheExpirationSec = &val
			}
			if model.Remote.Proxy.Key.ValueString() != "" {
				req.Remote.Proxy.Key = model.Remote.Proxy.Key.ValueString()
			}
			if model.Remote.Proxy.SchemeOverride.ValueString() != "" {
				req.Remote.Proxy.SchemeOverride = model.Remote.Proxy.SchemeOverride.ValueString()
			}
		}
	}

	if model.Local != nil {
		req.Local = &bridgeLocalAPIModel{
			Url: model.Local.Url.ValueString(),
		}
		if !model.Local.AnonymousEndpoints.IsNull() && model.Local.AnonymousEndpoints.Elements() != nil {
			var endpoints []string
			_ = model.Local.AnonymousEndpoints.ElementsAs(context.Background(), &endpoints, false)
			req.Local.AnonymousEndpoints = endpoints
		}
	}

	return req
}
