package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type iproute2ProviderModel struct{
	Host			types.String `tfsdk:"host"`
	User			types.String `tfsdk:"user"`
	PrivateKey		types.String `tfsdk:"private_key"`
}

type iproute2Provider struct{
}
func New() provider.Provider {
    return &iproute2Provider{}
}

func (p *iproute2Provider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
    resp.TypeName = "iproute2"
}

func (p *iproute2Provider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host":			schema.StringAttribute{Required: true},
			"user":			schema.StringAttribute{Required: true},
			"private_key":	schema.StringAttribute{Required: true, Sensitive: true},
		},
	}
}

func (p *iproute2Provider) Resources(_ context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        func() resource.Resource { return &bridgeResource{} },
    }
}

func (p *iproute2Provider) DataSources(_ context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource{}
}


func (p *iproute2Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
    var config iproute2ProviderModel
    diags := req.Config.Get(ctx, &config)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    resp.ResourceData = config
}
