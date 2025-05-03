package sources

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"me.zyrouge.anything_to_rss/internal/common"
	"me.zyrouge.anything_to_rss/internal/rss"
)

const GitHubBaseUrl = "https://github.com"
const GitHubBaseApiUrl = "https://api.github.com"

type FetchGitHubContainerRegistryVersionsInput struct {
	AccountTypeRoute string
	Owner            string
	Repository       string
	Package          string
	Limit            int
	TagFilter        *regexp.Regexp
}

type FetchGitHubContainerRegistryVersionsOutput struct {
	Input FetchGitHubContainerRegistryVersionsInput
	Data  []FetchGitHubContainerRegistryVersionsOutputData
}

type FetchGitHubContainerRegistryVersionsOutputData struct {
	HtmlUrl   string                                                 `json:"html_url"`
	CreatedAt time.Time                                              `json:"created_at"`
	Metadata  FetchGitHubContainerRegistryVersionsOutputDataMetadata `json:"metadata"`
}

type FetchGitHubContainerRegistryVersionsOutputDataMetadata struct {
	Container FetchGitHubContainerRegistryVersionsOutputDataMetadataContainer `json:"container"`
}

type FetchGitHubContainerRegistryVersionsOutputDataMetadataContainer struct {
	Tags []string `json:"tags"`
}

func (output *FetchGitHubContainerRegistryVersionsOutput) Rss() *rss.RssXml {
	items := []rss.RssXmlChannelItem{}
	for _, x := range output.Data {
		var tag *string
		for _, t := range x.Metadata.Container.Tags {
			if output.Input.TagFilter.MatchString(t) {
				tag = &t
				break
			}
		}
		if tag == nil {
			continue
		}
		item := rss.RssXmlChannelItem{
			Title:       *tag,
			Description: fmt.Sprintf("%s/%s:%s", output.Input.Owner, output.Input.Repository, *tag),
			Author:      output.Input.Owner,
			Link:        x.HtmlUrl,
			PubDate:     rss.MakeRssXmlChannelItemPubDate(x.CreatedAt),
		}
		items = append(items, item)
	}
	channel := rss.RssXmlChannel{
		Title:       fmt.Sprintf("%s/%s versions", output.Input.Owner, output.Input.Repository),
		Description: fmt.Sprintf("%s/%s versions", output.Input.Owner, output.Input.Repository),
		Link:        fmt.Sprintf("%s/%s/%s/pkgs/container/%s", GitHubBaseUrl, output.Input.Owner, output.Input.Repository, output.Input.Package),
		Items:       items,
	}
	rss := rss.RssXml{
		Version: 2,
		Channel: channel,
	}
	return &rss
}

func FetchGitHubContainerRegistryVersions(input FetchGitHubContainerRegistryVersionsInput) (*FetchGitHubContainerRegistryVersionsOutput, error) {
	env, err := common.GetEnv()
	if err != nil {
		return nil, err
	}
	apiUrl := fmt.Sprintf("%s/%s/%s/packages/container/%s/versions?per_page=%d", GitHubBaseApiUrl, input.AccountTypeRoute, input.Owner, input.Package, input.Limit)
	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", env.GitHubContainerRegistryApiToken))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	resp, err := common.GlobalHttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	var data []FetchGitHubContainerRegistryVersionsOutputData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	output := FetchGitHubContainerRegistryVersionsOutput{
		Input: input,
		Data:  data,
	}
	return &output, nil
}

func RouteGitHubContainerRegistryVersions(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	query := req.URL.Query()
	accountTypeRouteValid, accountTypeRoute := common.StringAsContentString(query.Get("accountTypeRoute"))
	ownerValid, owner := common.StringAsContentString(query.Get("owner"))
	repositoryValid, repository := common.StringAsContentString(query.Get("repository"))
	packageValid, packageValue := common.StringAsContentString(query.Get("package"))
	limitValid, limit := common.StringAsNumberOrNil(query.Get("limit"))
	tagFilterValid, tagFilter := common.StringAsRegExpOrNil(query.Get("tagFilter"))
	if !accountTypeRouteValid || !ownerValid || !repositoryValid || !packageValid || !limitValid || !tagFilterValid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	input := FetchGitHubContainerRegistryVersionsInput{
		AccountTypeRoute: accountTypeRoute,
		Owner:            owner,
		Repository:       repository,
		Package:          packageValue,
		Limit:            limit,
		TagFilter:        tagFilter,
	}
	output, err := FetchGitHubContainerRegistryVersions(input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	output.Rss().WriteToHttpResponseWriter(w)
}
