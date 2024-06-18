package googlesearchconsole

import (
	"context"

	"google.golang.org/api/searchconsole/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableGoogleSearchConsoleSitemap(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "googlesearchconsole_sitemap",
		Description: "Lists the sitemaps-entries submitted for sites, or included in the sitemap index file.",
		List: &plugin.ListConfig{
			KeyColumns:    plugin.OptionalColumns([]string{"site_url"}),
			Hydrate:       listSitemaps,
			ParentHydrate: listSites,
		},
		Columns: []*plugin.Column{
			{
				Name:        "site_url",
				Description: "The URL of the site.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "path",
				Description: "The url of the sitemap.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WmxSitemap.Path"),
			},
			{
				Name:        "errors",
				Description: "Number of errors in the sitemap.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("WmxSitemap.Errors"),
			},
			{
				Name:        "is_pending",
				Description: "If true, the sitemap has not been processed.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("WmxSitemap.IsPending"),
			},
			{
				Name:        "is_sitemaps_index",
				Description: "If true, the sitemap is a collection of sitemaps.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("WmxSitemap.IsSitemapsIndex"),
			},
			{
				Name:        "last_downloaded",
				Description: "Date & time in which this sitemap was last downloaded.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("WmxSitemap.LastDownloaded"),
			},
			{
				Name:        "last_submitted",
				Description: "Date & time in which this sitemap was last submitted.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("WmxSitemap.LastSubmitted"),
			},
			{
				Name:        "type",
				Description: "The type of the sitemap.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WmxSitemap.Type"),
			},
			{
				Name:        "warnings",
				Description: "Number of warnings in the sitemap.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("WmxSitemap.Warnings"),
			},
			{
				Name:        "contents",
				Description: "The various content types in the sitemap.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WmxSitemap.Contents"),
			},
			{
				Name:        "project",
				Description: "The GCP Project associated with the credentials in use.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getProject,
				Transform:   transform.FromValue(),
			},
		},
	}
}

type SitemapInfo struct {
	SiteUrl string
	*searchconsole.WmxSitemap
}

//// LIST FUNCTION

func listSitemaps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	site := h.Item.(*searchconsole.WmxSite)
	var siteUrl string
	if d.Quals["site_url"] != nil {
		siteUrl = d.EqualsQualString("site_url")

		if siteUrl != site.SiteUrl {
			return nil, nil
		}
	} else {
		siteUrl = site.SiteUrl
	}

	// Create client
	opts, err := getSearchConsoleSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("googlesearchconsole_sitemap.listSitemaps", "connection_error", err)
		return nil, err
	}

	// Create service
	svc, err := searchconsole.NewService(ctx, opts...)
	if err != nil {
		plugin.Logger(ctx).Error("googlesearchconsole_sitemap.listSitemaps", "service_creation_error", err)
		return nil, err
	}

	req := svc.Sitemaps.List(siteUrl)

	resp, err := req.Context(ctx).Do()
	if err != nil {
		plugin.Logger(ctx).Error("googlesearchconsole_sitemap.listSitemaps", "api_error", err)
		return nil, err
	}

	if resp.Sitemap != nil {
		for _, sitemap := range resp.Sitemap {
			info := &SitemapInfo{SiteUrl: siteUrl, WmxSitemap: sitemap}
			d.StreamListItem(ctx, info)
		}
	}

	return nil, nil
}
