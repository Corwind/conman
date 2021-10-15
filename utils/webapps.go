package utils

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/Corwind/utils/dbutils"

	"github.com/google/uuid"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
	"github.com/dchest/uniuri"
	helmclient "github.com/mittwald/go-helm-client"
)

type WebappParams struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	ReleaseNameFmt string `json:"-"`
	ChartName      string `json:"chartname"`
	ValuesYamlFmt  string `json:"-"`
}

func InstallWebapp(
	client helmclient.Client,
	db fdb.Database,
	params WebappParams,
	namespace string,
	hostname string,
	additionnalYaml []string,
) (interface{}, error) {
	name := fmt.Sprintf(params.ReleaseNameFmt, strings.ToLower(uniuri.New()))
	fqdn := fmt.Sprintf("%s.%s", name, hostname)
	params.ValuesYamlFmt = strings.ReplaceAll(params.ValuesYamlFmt, "{hostname}", fqdn)
	var valuesYaml = []string{
		params.ValuesYamlFmt,
		strings.Join(additionnalYaml, "\n"),
	}
	chart := helmclient.ChartSpec{
		ReleaseName: name,
		ChartName:   params.ChartName,
		Namespace:   namespace,
		ValuesYaml:  strings.Join(valuesYaml, "\n"),
	}

	ret, err := client.InstallOrUpgradeChart(context.Background(), &chart)

	if err != nil {
		return nil, err
	}
	rel := Release{
		Name:      ret.Name,
		Id:        uuid.New().String(),
		FQDN:      fqdn,
		Namespace: ret.Namespace,
	}
	return SaveRelease(db, rel)
}

func NewWebappParams(name string, chartname string, values string) *WebappParams {
	return &WebappParams{
		Id:             uuid.New().String(),
		Name:           name,
		ReleaseNameFmt: strings.ToLower(name) + "-%s",
		ChartName:      chartname,
		ValuesYamlFmt:  values,
	}
}

func SaveWebappParams(db fdb.Database, namespace string, params WebappParams) (interface{}, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	var tp tuple.Tuple = tuple.Tuple{"webapp", namespace, "id", params.Id}
	enc.Encode(&tp)
	dbutils.DbSaveEntity(db, tuple.Tuple{"webapp", namespace, "name", params.Name}, &buffer)
	var buffer2 bytes.Buffer
	enc = gob.NewEncoder(&buffer2)
	enc.Encode(params)
	return dbutils.DbSaveEntity(db, tp, &buffer2)
}

func fetchWebappParams(db fdb.Database, t tuple.Tuple) (interface{}, error) {
	ret, err := dbutils.DbFetchEntity(db, t)
	if err != nil {
		return nil, err
	}
	return DecodeWebappParams(ret)
}

func FetchWebappParamsById(db fdb.Database, namespace string, id string) (interface{}, error) {
	return fetchWebappParams(db, tuple.Tuple{"webapp", namespace, "id", id})
}

func FetchWebappParamsByName(db fdb.Database, namespace string, name string) (interface{}, error) {
	ret, err := dbutils.DbFetchEntity(db, tuple.Tuple{"webapp", namespace, "name", name})
	if err != nil {
		ret, err = dbutils.DbFetchEntity(db, tuple.Tuple{"webapp", "public", "name", name})
		if err != nil {
			return nil, err
		}
	}

	var t tuple.Tuple
	var buffer bytes.Buffer
	buffer.Write(ret.([]byte))
	dec := gob.NewDecoder(&buffer)
	err = dec.Decode(&t)
	if err != nil {
		return nil, err
	}
	return fetchWebappParams(db, t)
}

func FetchWebappParamsList(db fdb.Database, namespace string) ([]WebappParams, error) {
	ret, err := dbutils.DbFetchRange(db, tuple.Tuple{"webapp", namespace, "id"})
	if err != nil {
		return nil, err
	}
	values := make([]WebappParams, 0, len(ret))

	for _, value := range ret {
		webappParam, err := DecodeWebappParams(value)
		if err != nil {
			return nil, err
		}
		values = append(values, webappParam)
	}

	if namespace != "public" {
		public_webapps, err := FetchWebappParamsList(db, "public")
		if err != nil {
			return nil, err
		}
		values = append(values, public_webapps...)
	}

	return values, nil
}

func DecodeWebappParams(ret interface{}) (WebappParams, error) {
	var webappParams WebappParams
	var b bytes.Buffer
	b.Write(ret.([]byte))
	decoder := gob.NewDecoder(&b)
	err := decoder.Decode(&webappParams)
	return webappParams, err
}
