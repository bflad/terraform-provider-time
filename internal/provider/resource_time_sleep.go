package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.ResourceType = (*timeSleepResourceType)(nil)

type timeSleepResourceType struct{}

func (t timeSleepResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	//TODO implement me
	panic("implement me")
}

func (t timeSleepResourceType) NewResource(ctx context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	//TODO implement me
	panic("implement me")
}

var (
	_ tfsdk.Resource                = (*timeSleepResource)(nil)
	_ tfsdk.ResourceWithImportState = (*timeSleepResource)(nil)
)

type timeSleepResource struct {
}

func (t timeSleepResource) ImportState(ctx context.Context, request tfsdk.ImportResourceStateRequest, response *tfsdk.ImportResourceStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (t timeSleepResource) Create(ctx context.Context, request tfsdk.CreateResourceRequest, response *tfsdk.CreateResourceResponse) {
	//TODO implement me
	panic("implement me")
}

func (t timeSleepResource) Read(ctx context.Context, request tfsdk.ReadResourceRequest, response *tfsdk.ReadResourceResponse) {
	//TODO implement me
	panic("implement me")
}

func (t timeSleepResource) Update(ctx context.Context, request tfsdk.UpdateResourceRequest, response *tfsdk.UpdateResourceResponse) {
	//TODO implement me
	panic("implement me")
}

func (t timeSleepResource) Delete(ctx context.Context, request tfsdk.DeleteResourceRequest, response *tfsdk.DeleteResourceResponse) {
	//TODO implement me
	panic("implement me")
}
