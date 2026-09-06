// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package bedrockagentcore

import (
	"context"
	"iter"

	"github.com/YakDriver/smarterr"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentcorecontrol"
	awstypes "github.com/aws/aws-sdk-go-v2/service/bedrockagentcorecontrol/types"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/fwdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	"github.com/hashicorp/terraform-provider-aws/internal/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/smerr"
)

// @FrameworkListResource("aws_bedrockagentcore_capacity_provider")
func newCapacityProviderResourceAsListResource() list.ListResourceWithConfigure {
	return &capacityProviderListResource{}
}

type capacityProviderListResource struct {
	capacityProviderResource
	framework.WithList
}

func (l *capacityProviderListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream) {
	conn := l.Meta().BedrockAgentCoreClient(ctx)
	stream.Results = func(yield func(list.ListResult) bool) {
		for item, err := range listCapacityProviders(ctx, conn, &bedrockagentcorecontrol.ListCapacityProvidersInput{}) {
			if err != nil {
				yield(fwdiag.NewListResultErrorDiagnostic(err))
				return
			}
			result := request.NewListResult(ctx)
			var data capacityProviderResourceModel
			data.ID = types.StringPointerValue(item.CapacityProviderId)
			data.ARN = types.StringPointerValue(item.CapacityProviderArn)
			var output *bedrockagentcorecontrol.GetCapacityProviderOutput
			if request.IncludeResource {
				output, err = findCapacityProviderByID(ctx, conn, aws.ToString(item.CapacityProviderId))
				if retry.NotFound(err) {
					continue
				}
				if err != nil {
					yield(fwdiag.NewListResultErrorDiagnostic(err))
					return
				}
			}
			l.SetResult(ctx, l.Meta(), request.IncludeResource, &data, &result, func() {
				if request.IncludeResource {
					smerr.AddEnrich(ctx, &result.Diagnostics, l.flatten(ctx, output, &data))
				}
				result.DisplayName = aws.ToString(item.Name)
			})
			if !yield(result) {
				return
			}
		}
	}
}

func listCapacityProviders(ctx context.Context, conn *bedrockagentcorecontrol.Client, input *bedrockagentcorecontrol.ListCapacityProvidersInput) iter.Seq2[awstypes.CapacityProviderSummary, error] {
	return func(yield func(awstypes.CapacityProviderSummary, error) bool) {
		pages := bedrockagentcorecontrol.NewListCapacityProvidersPaginator(conn, input)
		for pages.HasMorePages() {
			page, err := pages.NextPage(ctx)
			if err != nil {
				yield(awstypes.CapacityProviderSummary{}, smarterr.NewError(err))
				return
			}
			for _, item := range page.CapacityProviders {
				if !yield(item, nil) {
					return
				}
			}
		}
	}
}
