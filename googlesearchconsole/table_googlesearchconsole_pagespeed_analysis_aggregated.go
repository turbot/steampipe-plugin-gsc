package googlesearchconsole

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/pagespeedonline/v5"
)

//// TABLE DEFINITION

func tableGoogleSearchConsolePagespeedAnalysisAggregated(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "googlesearchconsole_pagespeed_analysis_aggregated",
		Description: "Lists the aggregated pagespeed analysis for the URLs in the sitemap.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "site_url",
					Require: plugin.Required,
				},
				{
					Name:       "strategy",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
			},
			Hydrate: listPagespeedAnalysesAggregated,
		},
		Columns: getPagespeedAnalysisAggregatedColumns(),
	}
}

func getPagespeedAnalysisAggregatedColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "site_url",
			Description: "The URL of the site.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromQual("site_url"),
		},
		{
			Name:        "strategy",
			Description: "The analysis strategy (desktop or mobile) to use. Default is desktop.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "id",
			Description: "The ID of the page.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("UrlInspectionResult.Id"),
		},
		{
			Name:        "analysis_utc_timestamp",
			Description: "The timestamp of the analysis.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("UrlInspectionResult.AnalysisUTCTimestamp"),
		},
		{
			Name:        "overall_loading_experience",
			Description: "The loading experience of the page.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("UrlInspectionResult.OriginLoadingExperience.OverallCategory"),
		},
		{
			Name:        "cls",
			Description: "The Cumulative Layout Shift (CLS) of the page.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("UrlInspectionResult.OriginLoadingExperience.Metrics.CUMULATIVE_LAYOUT_SHIFT_SCORE.Category"),
		},
		{
			Name:        "cls_percentile",
			Description: "The percentile of the Cumulative Layout Shift (CLS) of the page.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("UrlInspectionResult.OriginLoadingExperience.Metrics.CUMULATIVE_LAYOUT_SHIFT_SCORE.Percentile"),
		},
		{
			Name:        "ttfb",
			Description: "The Time to First Byte (TTFB) of the page.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("UrlInspectionResult.OriginLoadingExperience.Metrics.EXPERIMENTAL_TIME_TO_FIRST_BYTE.Category"),
		},
		{
			Name:        "ttfb_percentile",
			Description: "The percentile of the Time to First Byte (TTFB) of the page.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("UrlInspectionResult.OriginLoadingExperience.Metrics.EXPERIMENTAL_TIME_TO_FIRST_BYTE.Percentile"),
		},
		{
			Name:        "fcp",
			Description: "The First Contentful Paint (FCP) of the page.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("UrlInspectionResult.OriginLoadingExperience.Metrics.FIRST_CONTENTFUL_PAINT_MS.Category"),
		},
		{
			Name:        "fcp_percentile",
			Description: "The percentile of the First Contentful Paint (FCP) of the page.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("UrlInspectionResult.OriginLoadingExperience.Metrics.FIRST_CONTENTFUL_PAINT_MS.Percentile"),
		},
		{
			Name:        "fid",
			Description: "The First Input Delay (FID) of the page.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("UrlInspectionResult.OriginLoadingExperience.Metrics.FIRST_INPUT_DELAY_MS.Category"),
		},
		{
			Name:        "fid_percentile",
			Description: "The percentile of the First Input Delay (FID) of the page.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("UrlInspectionResult.OriginLoadingExperience.Metrics.FIRST_INPUT_DELAY_MS.Percentile"),
		},
		{
			Name:        "inp",
			Description: "The Interaction to Next Paint (INP) of the page.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("UrlInspectionResult.OriginLoadingExperience.Metrics.INTERACTION_TO_NEXT_PAINT.Category"),
		},
		{
			Name:        "inp_percentile",
			Description: "The percentile of the Interaction to Next Paint (INP) of the page.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("UrlInspectionResult.OriginLoadingExperience.Metrics.INTERACTION_TO_NEXT_PAINT.Percentile"),
		},
		{
			Name:        "lcp",
			Description: "The Largest Contentful Paint (LCP) of the page.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("UrlInspectionResult.OriginLoadingExperience.Metrics.LARGEST_CONTENTFUL_PAINT_MS.Category"),
		},
		{
			Name:        "lcp_percentile",
			Description: "The percentile of the Largest Contentful Paint (LCP) of the page.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("UrlInspectionResult.OriginLoadingExperience.Metrics.LARGEST_CONTENTFUL_PAINT_MS.Percentile"),
		},
		{
			Name:        "cls_bucket_range",
			Description: "The bucket range of the Cumulative Layout Shift (CLS) of the page.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("UrlInspectionResult.OriginLoadingExperience.Metrics.CUMULATIVE_LAYOUT_SHIFT_SCORE.Distributions"),
		},
		{
			Name:        "ttfb_bucket_range",
			Description: "The bucket range of the Time to First Byte (TTFB) of the page.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("UrlInspectionResult.OriginLoadingExperience.Metrics.EXPERIMENTAL_TIME_TO_FIRST_BYTE.Distributions"),
		},
		{
			Name:        "fcp_bucket_range",
			Description: "The bucket range of the First Contentful Paint (FCP) of the page.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("UrlInspectionResult.OriginLoadingExperience.Metrics.FIRST_CONTENTFUL_PAINT_MS.Distributions"),
		},
		{
			Name:        "fid_bucket_range",
			Description: "The bucket range of the First Input Delay (FID) of the page.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("UrlInspectionResult.OriginLoadingExperience.Metrics.FIRST_INPUT_DELAY_MS.Distributions"),
		},
		{
			Name:        "inp_bucket_range",
			Description: "The bucket range of the Interaction to Next Paint (INP) of the page.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("UrlInspectionResult.OriginLoadingExperience.Metrics.INTERACTION_TO_NEXT_PAINT.Distributions"),
		},
		{
			Name:        "lcp_bucket_range",
			Description: "The bucket range of the Largest Contentful Paint (LCP) of the page.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("UrlInspectionResult.OriginLoadingExperience.Metrics.LARGEST_CONTENTFUL_PAINT_MS.Distributions"),
		},
		{
			Name:        "project",
			Description: "The GCP Project associated with the credentials in use.",
			Type:        proto.ColumnType_STRING,
			Hydrate:     getProject,
			Transform:   transform.FromValue(),
		},
	}
}

//// LIST FUNCTION

func listPagespeedAnalysesAggregated(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	siteUrl := d.EqualsQualString("site_url")
	strategy := d.EqualsQualString("strategy")

	if siteUrl == "" {
		plugin.Logger(ctx).Error("googlesearchconsole_pagespeed_analysis_aggregated.listPagespeedAnalysesAggregated", "validation_error", "site_url must be provided")
		return nil, nil
	}
	if strategy != "" {
		if strings.ToLower(strategy) != "mobile" && strings.ToLower(strategy) != "desktop" {
			plugin.Logger(ctx).Error("googlesearchconsole_pagespeed_analysis_aggregated.listPagespeedAnalysesAggregated", "validation_error", "Invalid strategy. The strategy should be either 'mobile' or 'desktop'.")
			return nil, nil
		}
	}
	if strategy == "" {
		strategy = "desktop"
	}

	// Create client
	opts, err := getPagespeedSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("googlesearchconsole_pagespeed_analysis_aggregated.listPagespeedAnalysesAggregated", "connection_error", err)
		return nil, err
	}

	// Create service
	svc, err := pagespeedonline.NewService(ctx, opts...)
	if err != nil {
		plugin.Logger(ctx).Error("googlesearchconsole_pagespeed_analysis_aggregated.listPagespeedAnalysesAggregated", "service_creation_error", err)
		return nil, err
	}

	req := svc.Pagespeedapi.Runpagespeed(siteUrl).Fields("id,originLoadingExperience,analysisUTCTimestamp")
	if strategy != "" {
		req.Strategy(strings.ToUpper(strategy))
	}

	resp, err := req.Context(ctx).Do()
	if err != nil {
		plugin.Logger(ctx).Error("googlesearchconsole_pagespeed_analysis_aggregated.listPagespeedAnalysesAggregated", "api_error", err)
		return nil, err
	}

	status := AnalysisPerURL{
		Loc:                 siteUrl,
		UrlInspectionResult: resp,
		Strategy:            strategy,
	}
	d.StreamListItem(ctx, status)

	return nil, nil
}
