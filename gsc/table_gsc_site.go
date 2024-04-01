package gsc

import (
	"context"

	"google.golang.org/api/searchconsole/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableGSCSite(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gsc_site",
		Description: "Lists the user's Search Console sites.",
		List: &plugin.ListConfig{
			Hydrate: listSites,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("site_url"),
			Hydrate:    getSite,
		},
		Columns: []*plugin.Column{
			{
				Name:        "site_url",
				Description: "The URL of the site.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "permission_level",
				Description: "The user's permission level for the site.",
				Type:        proto.ColumnType_STRING,
			},
		},
	}
}

//// LIST FUNCTION

func listSites(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create client
	opts, err := getSearchConsoleSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gsc_site.listSites", "connection_error", err)
		return nil, err
	}

	// Create service
	svc, err := searchconsole.NewService(ctx, opts...)
	if err != nil {
		plugin.Logger(ctx).Error("gsc_site.listSites", "service_creation_error", err)
		return nil, err
	}

	req := svc.Sites.List()

	resp, err := req.Context(ctx).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gsc_site.listSites", "api_error", err)
		return nil, err
	}

	if resp.SiteEntry != nil {
		for _, site := range resp.SiteEntry {
			d.StreamListItem(ctx, site)
		}
	}

	return nil, nil
}

//// GET FUNCTION

func getSite(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	siteUrl := d.EqualsQualString("site_url")

	// Create client
	opts, err := getSearchConsoleSessionConfig(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("gsc_site.getSite", "connection_error", err)
		return nil, err
	}

	// Create service
	svc, err := searchconsole.NewService(ctx, opts...)
	if err != nil {
		plugin.Logger(ctx).Error("gsc_site.getSite", "service_creation_error", err)
		return nil, err
	}

	req := svc.Sites.Get(siteUrl)

	resp, err := req.Context(ctx).Do()
	if err != nil {
		plugin.Logger(ctx).Error("gsc_site.getSite", "api_error", err)
		return nil, err
	}

	return resp, nil
}
