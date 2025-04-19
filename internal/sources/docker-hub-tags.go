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

const DockerHubBaseUrl = "https://hub.docker.com"
const DockerHubBaseApiUrl = "https://hub.docker.com/v2"

type FetchDockerHubTagsInput struct {
	Owner      string
	Repository string
	Limit      int
	TagFilter  *regexp.Regexp
}

type FetchDockerHubTagsOutput struct {
	Input FetchDockerHubTagsInput
	Data  FetchDockerHubTagsOutputData
}

type FetchDockerHubTagsOutputData struct {
	Results []FetchDockerHubTagsOutputDataResult `json:"results"`
}

type FetchDockerHubTagsOutputDataResult struct {
	Name        string    `json:"name"`
	Digest      string    `json:"digest"`
	LastUpdated time.Time `json:"last_updated"`
}

func (output *FetchDockerHubTagsOutput) Rss() *rss.RssXml {
	items := []rss.RssXmlChannelItem{}
	for _, x := range output.Data.Results {
		if !output.Input.TagFilter.MatchString(x.Name) {
			continue
		}
		item := rss.RssXmlChannelItem{
			Title:       x.Name,
			Description: fmt.Sprintf("%s/%s:%s", output.Input.Owner, output.Input.Repository, x.Name),
			Link:        fmt.Sprintf("%s/layers/%s/%s/%s/images/%s", DockerHubBaseUrl, output.Input.Owner, output.Input.Repository, x.Name, x.Digest),
			PubDate:     rss.MakeRssXmlChannelItemPubDate(x.LastUpdated),
		}
		items = append(items, item)
	}
	channel := rss.RssXmlChannel{
		Title:       fmt.Sprintf("%s/%s tags", output.Input.Owner, output.Input.Repository),
		Description: fmt.Sprintf("%s/%s tags", output.Input.Owner, output.Input.Repository),
		Link:        fmt.Sprintf("%s/r/%s/%s/tags", DockerHubBaseUrl, output.Input.Owner, output.Input.Repository),
		Items:       items,
	}
	rss := rss.RssXml{
		Version: 2,
		Channel: channel,
	}
	return &rss
}

func FetchDockerHubTags(input FetchDockerHubTagsInput) (*FetchDockerHubTagsOutput, error) {
	apiUrl := fmt.Sprintf("%s/namespaces/%s/repositories/%s/tags?page_size=%d", DockerHubBaseApiUrl, input.Owner, input.Repository, input.Limit)
	resp, err := common.GlobalHttpClient.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	var data FetchDockerHubTagsOutputData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	output := FetchDockerHubTagsOutput{
		Input: input,
		Data:  data,
	}
	return &output, nil
}

func RouteDockerHubTags(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	query := req.URL.Query()
	ownerValid, owner := common.StringAsContentString(query.Get("owner"))
	repositoryValid, repository := common.StringAsContentString(query.Get("repository"))
	limitValid, limit := common.StringAsNumberOrNil(query.Get("limit"))
	tagFilterValid, tagFilter := common.StringAsRegExpOrNil(query.Get("tagFilter"))
	if !ownerValid || !repositoryValid || !limitValid || !tagFilterValid {
		fmt.Printf("%v %v %v %v", ownerValid, repositoryValid, limitValid, tagFilterValid)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	input := FetchDockerHubTagsInput{
		Owner:      owner,
		Repository: repository,
		Limit:      limit,
		TagFilter:  tagFilter,
	}
	output, err := FetchDockerHubTags(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	output.Rss().WriteToHttpResponseWriter(w)
}
