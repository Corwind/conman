package utils

import (
	"bytes"
	"encoding/gob"

	"helm.sh/helm/v3/pkg/repo"

	"github.com/Corwind/utils/dbutils"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
	helmclient "github.com/mittwald/go-helm-client"
)

type RepoParams struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func InitRepos(client helmclient.Client, db fdb.Database, namespace string) error {
	repoParamsList, err := FetchRepoParamsList(db, namespace)
	if err != nil {
		return err
	}
	for _, repoParams := range repoParamsList {
		client.AddOrUpdateChartRepo(
			repo.Entry{
				Name: repoParams.Name,
				URL:  repoParams.URL,
			},
		)
	}
	client.UpdateChartRepos()
	return nil
}

func NewRepoParams(name string, url string) *RepoParams {
	return &RepoParams{
		Name: name,
		URL:  url,
	}
}

func SaveRepoParams(db fdb.Database, namespace string, params RepoParams) (interface{}, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	t := tuple.Tuple{"repo", namespace, "name", params.Name}
	enc.Encode(&params)
	return dbutils.DbSaveEntity(db, t, &buffer)
}

func FetchRepoParams(db fdb.Database, namespace string, name string) (interface{}, error) {
	ret, err := dbutils.DbFetchEntity(db, tuple.Tuple{"repo", namespace, "name", name})
	if err != nil {
		return nil, err
	}
	return DecodeRepoParams(ret)
}

func FetchRepoParamsList(db fdb.Database, namespace string) ([]RepoParams, error) {
	ret, err := dbutils.DbFetchRange(db, tuple.Tuple{"repo", namespace, "name"})
	if err != nil {
		return nil, err
	}
	values := make([]RepoParams, 0, len(ret))

	for _, value := range ret {
		repoParams, err := DecodeRepoParams(value)
		if err != nil {
			return nil, err
		}
		values = append(values, repoParams)
	}

	if namespace != "public" {
		public_repos, err := FetchRepoParamsList(db, "public")
		if err != nil {
			return nil, err
		}
		values = append(values, public_repos...)
	}

	return values, nil
}

func DecodeRepoParams(value interface{}) (RepoParams, error) {
	var repoParams RepoParams
	var b bytes.Buffer
	b.Write(value.([]byte))
	decoder := gob.NewDecoder(&b)
	err := decoder.Decode(&repoParams)
	return repoParams, err
}
