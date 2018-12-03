// Copyright 2018-present Ovi Chis www.ovios.org All rights reserved.
// Use of this source code is governed by a MIT-license.

package elkluster

import (
    "context"
    "gopkg.in/olivere/elastic.v6"
    "fmt"
    "os"
    "encoding/json"
)
    

func InnitiateClient(ctx context.Context, url string) *elastic.Client {
    //fmt.Println("trying to connect to ", url)
    client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(url))
    if err != nil {
        panic(err)
    } else {
        //fmt.Println("Innitiated client " , client)
        fmt.Println("Using Elasticsearch " , url)
    }
    return client
}

func IndexExists(ctx context.Context, client *elastic.Client, index string) ( bool, error) {
    exists, err := client.IndexExists(index).Do(ctx)
    if err != nil {
        fmt.Println(err)
        return exists, err
    }
    return exists, err
}

func CreateIndex(ctx context.Context, client *elastic.Client, index string, indexbody string) ( bool, error) {
    exists, _ := client.IndexExists(index).Do(ctx)
    if exists {
        fmt.Printf("Index %s already exists. \n" , index)
        os.Exit(1)
    }
    createIndex, err := client.CreateIndex(index).BodyJson(indexbody).Do(ctx)
    if err != nil {
        fmt.Println(err)
    }
	return createIndex.Acknowledged, err
}

func ListIndexes(ctx context.Context, client *elastic.Client) ([]string, error) {
    indexes, err := client.IndexNames()
    return indexes, err
}

func RemoveIndex(ctx context.Context, client *elastic.Client, index string) ( bool, error) {
    exists, _ := client.IndexExists(index).Do(ctx)
    if !exists {
        fmt.Printf("Index %s doesn't exist. \n", index)
        os.Exit(1)
    }
    deleteIndex, err := client.DeleteIndex(index).Do(ctx)
     if err != nil {
        fmt.Println(err)
    }
    return deleteIndex.Acknowledged, err
}

func IndexDoc(ctx context.Context, client *elastic.Client, index string, doctype string, indexbody string) (*elastic.IndexResponse, error) {
    put, err := client.Index().Index(index).Type(doctype).BodyJson(indexbody).Do(ctx)
    if err != nil {
        panic(err)
    }
    return put, err
}

func CreateRepo(ctx context.Context, client *elastic.Client, reponame string, repotype string, repolocation string) bool {
    repoBody := fmt.Sprintf( `
    {
        "type": "%s",
        "settings": {
            "location": "%s"
        }
    }`, repotype, repolocation)
    
    service := client.SnapshotCreateRepository(reponame)
    service = service.Type(repotype).
    BodyString(repoBody)
    
    if serr := service.Validate(); serr != nil {
		fmt.Println(serr)
	}
    
    src, berr := service.Do(ctx)
	if berr != nil {
		fmt.Println(berr)
	}
	_, jerr := json.Marshal(src)
	if jerr != nil {
		fmt.Printf("Marshaling to JSON failed: %v\n", jerr)
	}
    return src.Acknowledged
}

func RemoveRepo(ctx context.Context, client *elastic.Client, reponame string) bool {
    res, err := client.SnapshotDeleteRepository(reponame).Do(ctx)
    if err != nil {
        panic(err)
    }
    return res.Acknowledged
}

func SnapCreate(ctx context.Context, client *elastic.Client, reponame string, snapname string, snapindex string) *bool {
    snapbody := fmt.Sprintf(`
    {
        "indices": "%s",
        "ignore_unavailable": "true",
        "include_global_state": "false",
        "wait_for_completion": "true"
    }`, snapindex)
    res , err := client.SnapshotCreate(reponame, snapname).BodyString(snapbody).Do(ctx)
    if err != nil {
        panic(err)
    }
    return res.Accepted
}

func SnapDelete(ctx context.Context, client *elastic.Client, reponame string, snapname string) bool {
    res, err := client.SnapshotDelete(reponame, snapname).Do(ctx)
    if err != nil {
        panic(err)
    }
    return res.Acknowledged
}

func SnapRestore(ctx context.Context, client *elastic.Client, reponame string, snapname string, snapindex string) bool {
    snapbody := fmt.Sprintf(`
    {
        "indices": "%s",
        "ignore_unavailable": "true",
        "include_global_state": "false"
    }`, snapindex)
    res, err := client.SnapshotRestore(reponame, snapname).BodyString(snapbody).Do(ctx)
    if err != nil {
        panic(err)
    }
    return res.Accepted
}
