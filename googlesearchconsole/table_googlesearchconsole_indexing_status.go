package googlesearchconsole

import (
	"context"
	"sync"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	sitemapper "github.com/yterajima/go-sitemap"
	"google.golang.org/api/searchconsole/v1"
)

//// TABLE DEFINITION

func tableGoogleSearchConsoleIndexingStatus(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "googlesearchconsole_indexing_status",
		Description: "Lists the indexing status of the URLs in the sitemap.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"site_url", "sitemap_url"}),
			Hydrate:    listIndexingStatuses,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"site_url", "loc"}),
			Hydrate:    getIndexingStatus,
		},
		Columns: []*plugin.Column{
			{
				Name:        "loc",
				Description: "The URL of the page.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "site_url",
				Description: "The URL of the site.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("site_url"),
			},
			{
				Name:        "sitemap_url",
				Description: "The URL of the sitemap.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("sitemap_url"),
			},
			{
				Name:        "coverage_state",
				Description: "Could Google find and index the page.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UrlInspectionResult.IndexStatusResult.CoverageState"),
			},
			{
				Name:        "crawled_as",
				Description: "Primary crawler that was used by Google to crawl your site.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UrlInspectionResult.IndexStatusResult.CrawledAs"),
			},
			{
				Name:        "google_canonical",
				Description: "The URL of the page that Google selected as canonical. If the page was not indexed, this field is absent.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UrlInspectionResult.IndexStatusResult.GoogleCanonical"),
			},
			{
				Name:        "indexing_state",
				Description: "Whether or not the page blocks indexing through a noindex rule.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UrlInspectionResult.IndexStatusResult.IndexingState"),
			},
			{
				Name:        "last_crawl_time",
				Description: "Last time this URL was crawled by Google using the primary crawler.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UrlInspectionResult.IndexStatusResult.LastCrawlTime"),
			},
			{
				Name:        "page_fetch_state",
				Description: "Whether or not Google could retrieve the page from your server.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UrlInspectionResult.IndexStatusResult.PageFetchState"),
			},
			{
				Name:        "result_link",
				Description: "Link to Search Console URL inspection.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UrlInspectionResult.InspectionResultLink"),
			},
			{
				Name:        "robots_txt_state",
				Description: "Whether or not the page is blocked to Google by a robots.txt rule.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UrlInspectionResult.IndexStatusResult.RobotsTxtState"),
			},
			{
				Name:        "user_canonical",
				Description: "The URL that your page or site declares as canonical.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UrlInspectionResult.IndexStatusResult.UserCanonical"),
			},
			{
				Name:        "verdict",
				Description: "High level verdict about whether the URL is indexed (indexed status), or can be indexed (live inspection)",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UrlInspectionResult.IndexStatusResult.Verdict"),
			},
			{
				Name:        "referring_urls",
				Description: "URLs that link to the inspected URL, directly and indirectly.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("UrlInspectionResult.IndexStatusResult.ReferringUrls"),
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

type StatusPerURL struct {
	Loc                 string
	ChangeFreq          string
	LastMod             string
	Priority            float32
	UrlInspectionResult *searchconsole.UrlInspectionResult
}

//// LIST FUNCTION

func listIndexingStatuses(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	siteUrl := d.EqualsQualString("site_url")
	smUrl := d.EqualsQualString("sitemap_url")

	if smUrl == "" {
		plugin.Logger(ctx).Error("sitemap_url is required")
		return nil, nil
	}

	sitemapURLs, err := sitemapper.Get(smUrl, nil)
	if err != nil {
		plugin.Logger(ctx).Error("googlesearchconsole_indexing_status.listIndexingStatuses", "sitemap_error", err)
		return nil, err
	}

	batches := createBatches(sitemapURLs.URL, 50) // Assuming a batchSize of 50

	var wg sync.WaitGroup
	wg.Add(len(batches))

	for i, batch := range batches {
		go processPageIndexingStatusBatch(ctx, d, siteUrl, batch, i, &wg)
	}
	wg.Wait() // Wait for all batches to complete

	for _, sitemapURL := range sitemapURLs.URL {
		status := StatusPerURL{
			Loc:                 sitemapURL.Loc,
			ChangeFreq:          sitemapURL.ChangeFreq,
			LastMod:             sitemapURL.LastMod,
			Priority:            sitemapURL.Priority,
			UrlInspectionResult: statusPerUrl[sitemapURL.Loc],
		}
		d.StreamListItem(ctx, status)
	}

	for k := range statusPerUrl {
		delete(statusPerUrl, k)
	}

	return nil, nil
}

//// GET FUNCTION

func getIndexingStatus(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pageUrl := d.EqualsQualString("loc")
	siteUrl := d.EqualsQualString("site_url")

	if siteUrl == "" || pageUrl == "" {
		return nil, nil
	}

	resp, err := getPageIndexingStatusService(ctx, d, pageUrl, siteUrl)
	if err != nil {
		plugin.Logger(ctx).Error("googlesearchconsole_indexing_status.getIndexingStatus", "api_error", err)
		return nil, err
	}

	status := StatusPerURL{
		Loc:                 pageUrl,
		UrlInspectionResult: resp,
	}

	return status, nil
}
