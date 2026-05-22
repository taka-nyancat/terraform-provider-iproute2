package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type bridgeResource struct {
	host		string
	user		string
	privateKey	string
}

type bridgeResourceModel struct {
	Name	types.String `tfsdk:"name"`
}

func (r *bridgeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bridge"
}

func (r *bridgeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{Required: true},
		},
	}
}

func (r *bridgeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan bridgeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError(){
		return
	}

	err := runSSH(r.host, r.user, r.privateKey,
		"sudo ip link add "+plan.Name.ValueString()+" type bridge && ip link set "+plan.Name.ValueString()+" up")
	if err != nil {
		resp.Diagnostics.AddError("SSH error", err.Error())
		return
	}
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *bridgeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state bridgeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := runSSH(r.host, r.user, r.privateKey, 
		"sudo ip link show "+state.Name.ValueString())
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *bridgeResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *bridgeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state bridgeResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    err := runSSH(r.host, r.user, r.privateKey,
        "sudo ip link del "+state.Name.ValueString())
    if err != nil {
        resp.Diagnostics.AddError("SSH error", err.Error())
    }
}

func (r *bridgeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }
    config, ok := req.ProviderData.(iproute2ProviderModel)
    if !ok {
        resp.Diagnostics.AddError("Unexpected provider data", "Expected iproute2ProviderModel")
        return
    }
    r.host       = config.Host.ValueString()
    r.user       = config.User.ValueString()
    r.privateKey = config.PrivateKey.ValueString()
}
