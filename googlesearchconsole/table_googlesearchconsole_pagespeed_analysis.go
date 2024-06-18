package googlesearchconsole

import (
	"context"
	"strings"
	"sync"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	sitemapper "github.com/yterajima/go-sitemap"
	"google.golang.org/api/pagespeedonline/v5"
)

//// TABLE DEFINITION

func tableGoogleSearchConsolePagespeedAnalysis(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "googlesearchconsole_pagespeed_analysis",
		Description: "Lists the pagespeed analysis for the URLs in the sitemap.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "sitemap_url",
					Require: plugin.Required,
				},
				{
					Name:       "strategy",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
			},
			Hydrate: listPagespeedAnalyses,
		},
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "loc",
					Require: plugin.Required,
				},
				{
					Name:    "strategy",
					Require: plugin.Optional,
				},
			},
			Hydrate: getPagespeedAnalysis,
		},
		Columns: getPagespeedAnalysisColumns(),
	}
}

func getPagespeedAnalysisColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "sitemap_url",
			Description: "The URL of the sitemap.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromQual("sitemap_url").NullIfZero(),
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
			Name:        "loc",
			Description: "The URL of the page.",
			Type:        proto.ColumnType_STRING,
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
			Transform:   transform.FromField("UrlInspectionResult.LoadingExperience.OverallCategory"),
		},
		{
			Name:        "cls",
			Description: "The Cumulative Layout Shift (CLS) of the page.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("UrlInspectionResult.LoadingExperience.Metrics.CUMULATIVE_LAYOUT_SHIFT_SCORE.Category"),
		},
		{
			Name:        "cls_percentile",
			Description: "The percentile of the Cumulative Layout Shift (CLS) of the page.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("UrlInspectionResult.LoadingExperience.Metrics.CUMULATIVE_LAYOUT_SHIFT_SCORE.Percentile"),
		},
		{
			Name:        "ttfb",
			Description: "The Time to First Byte (TTFB) of the page.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("UrlInspectionResult.LoadingExperience.Metrics.EXPERIMENTAL_TIME_TO_FIRST_BYTE.Category"),
		},
		{
			Name:        "ttfb_percentile",
			Description: "The percentile of the Time to First Byte (TTFB) of the page.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("UrlInspectionResult.LoadingExperience.Metrics.EXPERIMENTAL_TIME_TO_FIRST_BYTE.Percentile"),
		},
		{
			Name:        "fcp",
			Description: "The First Contentful Paint (FCP) of the page.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("UrlInspectionResult.LoadingExperience.Metrics.FIRST_CONTENTFUL_PAINT_MS.Category"),
		},
		{
			Name:        "fcp_percentile",
			Description: "The percentile of the First Contentful Paint (FCP) of the page.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("UrlInspectionResult.LoadingExperience.Metrics.FIRST_CONTENTFUL_PAINT_MS.Percentile"),
		},
		{
			Name:        "fid",
			Description: "The First Input Delay (FID) of the page.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("UrlInspectionResult.LoadingExperience.Metrics.FIRST_INPUT_DELAY_MS.Category"),
		},
		{
			Name:        "fid_percentile",
			Description: "The percentile of the First Input Delay (FID) of the page.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("UrlInspectionResult.LoadingExperience.Metrics.FIRST_INPUT_DELAY_MS.Percentile"),
		},
		{
			Name:        "inp",
			Description: "The Interaction to Next Paint (INP) of the page.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("UrlInspectionResult.LoadingExperience.Metrics.INTERACTION_TO_NEXT_PAINT.Category"),
		},
		{
			Name:        "inp_percentile",
			Description: "The percentile of the Interaction to Next Paint (INP) of the page.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("UrlInspectionResult.LoadingExperience.Metrics.INTERACTION_TO_NEXT_PAINT.Percentile"),
		},
		{
			Name:        "lcp",
			Description: "The Largest Contentful Paint (LCP) of the page.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("UrlInspectionResult.LoadingExperience.Metrics.LARGEST_CONTENTFUL_PAINT_MS.Category"),
		},
		{
			Name:        "lcp_percentile",
			Description: "The percentile of the Largest Contentful Paint (LCP) of the page.",
			Type:        proto.ColumnType_INT,
			Transform:   transform.FromField("UrlInspectionResult.LoadingExperience.Metrics.LARGEST_CONTENTFUL_PAINT_MS.Percentile"),
		},
		{
			Name:        "cls_bucket_range",
			Description: "The bucket range of the Cumulative Layout Shift (CLS) of the page.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("UrlInspectionResult.LoadingExperience.Metrics.CUMULATIVE_LAYOUT_SHIFT_SCORE.Distributions"),
		},
		{
			Name:        "ttfb_bucket_range",
			Description: "The bucket range of the Time to First Byte (TTFB) of the page.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("UrlInspectionResult.LoadingExperience.Metrics.EXPERIMENTAL_TIME_TO_FIRST_BYTE.Distributions"),
		},
		{
			Name:        "fcp_bucket_range",
			Description: "The bucket range of the First Contentful Paint (FCP) of the page.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("UrlInspectionResult.LoadingExperience.Metrics.FIRST_CONTENTFUL_PAINT_MS.Distributions"),
		},
		{
			Name:        "fid_bucket_range",
			Description: "The bucket range of the First Input Delay (FID) of the page.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("UrlInspectionResult.LoadingExperience.Metrics.FIRST_INPUT_DELAY_MS.Distributions"),
		},
		{
			Name:        "inp_bucket_range",
			Description: "The bucket range of the Interaction to Next Paint (INP) of the page.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("UrlInspectionResult.LoadingExperience.Metrics.INTERACTION_TO_NEXT_PAINT.Distributions"),
		},
		{
			Name:        "lcp_bucket_range",
			Description: "The bucket range of the Largest Contentful Paint (LCP) of the page.",
			Type:        proto.ColumnType_JSON,
			Transform:   transform.FromField("UrlInspectionResult.LoadingExperience.Metrics.LARGEST_CONTENTFUL_PAINT_MS.Distributions"),
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

type AnalysisPerURL struct {
	Loc                 string
	Strategy            string
	UrlInspectionResult *pagespeedonline.PagespeedApiPagespeedResponseV5
}

//// LIST FUNCTION

func listPagespeedAnalyses(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	smUrl := d.EqualsQualString("sitemap_url")
	strategy := d.EqualsQualString("strategy")

	if smUrl == "" {
		plugin.Logger(ctx).Error("googlesearchconsole_pagespeed_analysis.listPagespeedAnalyses", "validation_error", "The sitemap_url must be specified.")
		return nil, nil
	}
	if strategy != "" {
		if strings.ToLower(strategy) != "mobile" && strings.ToLower(strategy) != "desktop" {
			plugin.Logger(ctx).Error("lgooglesearchconsole_pagespeed_analysis.istPagespeedAnalysis", "validation_error", "Invalid strategy. The strategy should be either 'mobile' or 'desktop'.")
			return nil, nil
		}
	}
	if strategy == "" {
		strategy = "desktop"
	}

	sitemapURLs, err := sitemapper.Get(smUrl, nil)
	if err != nil {
		plugin.Logger(ctx).Error("googlesearchconsole_pagespeed_analysis.listPagespeedAnalyses", "sitemap_error", err)
		return nil, err
	}

	batches := createBatches(sitemapURLs.URL, 50) // Assuming a batchSize of 50

	var wg sync.WaitGroup
	wg.Add(len(batches))

	for i, batch := range batches {
		go processPagespeedAnalysisBatch(ctx, d, strategy, batch, i, &wg)
	}
	wg.Wait() // Wait for all batches to complete

	for _, sitemapURL := range sitemapURLs.URL {
		status := AnalysisPerURL{
			Loc:                 sitemapURL.Loc,
			Strategy:            strategy,
			UrlInspectionResult: pagespeedAnalysisPerUrl[sitemapURL.Loc],
		}
		d.StreamListItem(ctx, status)
	}

	for k := range pagespeedAnalysisPerUrl {
		delete(pagespeedAnalysisPerUrl, k)
	}

	return nil, nil
}

func getPagespeedAnalysis(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pageUrl := d.EqualsQualString("loc")
	strategy := d.EqualsQualString("strategy")

	if pageUrl == "" {
		plugin.Logger(ctx).Error("googlesearchconsole_pagespeed_analysis.getPagespeedAnalysis", "validation_error", "The loc must be specified.")
		return nil, nil
	}
	if strategy != "" {
		if strings.ToLower(strategy) != "mobile" && strings.ToLower(strategy) != "desktop" {
			plugin.Logger(ctx).Error("googlesearchconsole_pagespeed_analysis.getPagespeedAnalysis", "validation_error", "Invalid strategy. The strategy should be either 'mobile' or 'desktop'.")
			return nil, nil
		}
	}
	if strategy == "" {
		strategy = "desktop"
	}

	resp, err := getPagespeedAnalysisService(ctx, d, pageUrl, strategy)
	if err != nil {
		plugin.Logger(ctx).Error("googlesearchconsole_pagespeed_analysis.getPagespeedAnalysis", "api_error", err)
		return nil, err
	}

	status := AnalysisPerURL{
		Loc:                 pageUrl,
		Strategy:            strategy,
		UrlInspectionResult: resp,
	}

	return status, nil
}
