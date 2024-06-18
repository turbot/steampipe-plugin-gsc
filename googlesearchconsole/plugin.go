package googlesearchconsole

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-googlesearchconsole",
		DefaultTransform: transform.FromCamel().NullIfZero(),
		ConnectionKeyColumns: []plugin.ConnectionKeyColumn{
			{
				Name:    "project",
				Hydrate: getProject,
			},
		},
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"googlesearchconsole_indexing_status":               tableGoogleSearchConsoleIndexingStatus(ctx),
			"googlesearchconsole_pagespeed_analysis":            tableGoogleSearchConsolePagespeedAnalysis(ctx),
			"googlesearchconsole_pagespeed_analysis_aggregated": tableGoogleSearchConsolePagespeedAnalysisAggregated(ctx),
			"googlesearchconsole_site":                          tableGoogleSearchConsoleSite(ctx),
			"googlesearchconsole_sitemap":                       tableGoogleSearchConsoleSitemap(ctx),
		},
	}
	return p
}
