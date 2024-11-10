package provider

import (
	"context"

	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/generated"
	providerdata "github.com/chrismarget-j/jtaf/terraform-provider-jtaf/provider_data"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.ProviderWithValidateConfig = (*jtafProvider)(nil)

type jtafProvider struct{}

func (j *jtafProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = generated.ProviderName
}

func (j *jtafProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: (*config)(nil).attributes(),
	}
}

func (j *jtafProvider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	var cfg config
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cfg.loadEnv(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	cfg.validate(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (j *jtafProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var cfg config
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cfg.loadEnv(&resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.ResourceData = providerdata.NewResourceData(cfg.session(ctx, &resp.Diagnostics))
	if resp.Diagnostics.HasError() {
		return
	}
}

func (j *jtafProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (j *jtafProvider) Resources(_ context.Context) []func() resource.Resource {
	return generated.Resources
}

func NewProvider() provider.Provider {
	return &jtafProvider{}
}
