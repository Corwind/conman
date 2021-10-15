package utils

import (
	"bytes"
	"encoding/gob"

	"github.com/Corwind/utils/dbutils"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
)

type Release struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	FQDN      string `json:"fqdn"`
}

func SaveRelease(db fdb.Database, release Release) (interface{}, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	enc.Encode(&release)
	t := tuple.Tuple{"releases", release.Namespace, "id", release.Id}
	_, err := dbutils.DbSaveEntity(db, t, &buffer)
	if err != nil {
		return nil, err
	}

	var buffer2 bytes.Buffer
	enc2 := gob.NewEncoder(&buffer2)
	enc2.Encode(&t)
	t4 := tuple.Tuple{"releases", release.Namespace, "name", release.Name}
	_, err = dbutils.DbSaveEntity(db, t4, &buffer2)
	if err != nil {
		return nil, err
	}

	return release, nil

}

func fetchRelease(db fdb.Database, t tuple.Tuple) (interface{}, error) {
	ret, err := dbutils.DbFetchEntity(db, t)
	if err != nil {
		return nil, err
	}
	return DecodeRelease(ret)
}

func FetchReleaseById(db fdb.Database, namespace string, id string) (interface{}, error) {
	return fetchRelease(db, tuple.Tuple{"releases", namespace, "id", id})
}

func FetchReleaseByName(db fdb.Database, namespace string, name string) (interface{}, error) {
	ret, err := dbutils.DbFetchEntity(db, tuple.Tuple{"releases", namespace, "name", name})
	if err != nil {
		return nil, err
	}

	var t tuple.Tuple
	var buffer bytes.Buffer
	buffer.Write(ret.([]byte))
	dec := gob.NewDecoder(&buffer)
	err = dec.Decode(&t)
	if err != nil {
		return nil, err
	}
	return fetchRelease(db, t)
}

func FetchReleases(db fdb.Database, namespace string) ([]Release, error) {
	ret, err := dbutils.DbFetchRange(db, tuple.Tuple{"releases", namespace, "name"})
	if err != nil {
		return nil, err
	}
	values := make([]Release, 0, len(ret))

	for _, value := range ret {
		t, err := decodeTuple(value)
		if err != nil {
			return nil, err
		}
		release, err := fetchRelease(db, t)
		if err != nil {
			return nil, err
		}
		values = append(values, release.(Release))
	}

	return values, nil
}

func decodeTuple(ret interface{}) (tuple.Tuple, error) {
	var tpl tuple.Tuple
	var b bytes.Buffer
	b.Write(ret.([]byte))
	decoder := gob.NewDecoder(&b)
	err := decoder.Decode(&tpl)
	return tpl, err
}

func DecodeRelease(ret interface{}) (Release, error) {
	var release Release
	var b bytes.Buffer
	b.Write(ret.([]byte))
	decoder := gob.NewDecoder(&b)
	err := decoder.Decode(&release)
	return release, err
}

func DeleteRelease(db fdb.Database, namespace string, name string) error {
	release, err := FetchReleaseByName(db, namespace, name)
	if err != nil {
		return err
	}

	dbutils.DbClearEntity(db, tuple.Tuple{"releases", namespace, "id", release.(Release).Id})
	dbutils.DbClearEntity(db, tuple.Tuple{"releases", namespace, "name", name})
	return nil
}
