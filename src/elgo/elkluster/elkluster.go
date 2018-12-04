// Copyright 2018-present Ovi Chis www.ovios.org All rights reserved.
// Use of this source code is governed by a MIT-license.

package elkluster

import (
    "context"
    "elgo/logger"
    "gopkg.in/olivere/elastic.v6"
    "fmt"
    "os"
    "encoding/json"
)

func check(e error) {
    if e != nil {
        logger.LogError(e)
        panic(e)
    }
}

func InnitiateClient(ctx context.Context, url string) *elastic.Client {
    client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(url))
    check(err)
    fmt.Println("Using Elasticsearch URL: " , url)
    return client
}

func ClusterInfo(ctx context.Context, client *elastic.Client)  {
    res, err := client.ClusterHealth().Do(ctx)
    check(err)
    fmt.Printf(`
    ClusterName: %s
    Status: %s
    Number Of Nodes: %d
    Number Of Data Nodes: %d
    Active Primary Shards: %d
    Active Shards: %d
    Relocating Shards: %d
    Initializing Shards: %d
    Unassigned Shards: %d
    Delayed Unassigned Shards: %d
    Number Of Pending Tasks: %d 
    Number Of InFlight Fetch: %d
    Task MaxWait Time In Queue In Millis: %d
    Active Shards Percent As Number: %.1f
    `, res.ClusterName, res.Status, res.NumberOfNodes, res.NumberOfDataNodes, 
    res.ActivePrimaryShards, res.ActiveShards, res.RelocatingShards, 
    res.InitializingShards, res.UnassignedShards, res.DelayedUnassignedShards, 
    res.NumberOfPendingTasks, res.NumberOfInFlightFetch, res.TaskMaxWaitTimeInQueueInMillis, 
    res.ActiveShardsPercentAsNumber)
    fmt.Printf("\n")
}

func IndexExists(ctx context.Context, client *elastic.Client, index string) bool {
    exists, err := client.IndexExists(index).Do(ctx)
    check(err)
    return exists
}

func CreateIndex(ctx context.Context, client *elastic.Client, index string, indexbody string) ( bool, error) {
    exists, _ := client.IndexExists(index).Do(ctx)
    if exists {
        fmt.Printf("Index %s already exists. \n" , index)
        os.Exit(1)
    }
    createIndex, err := client.CreateIndex(index).BodyJson(indexbody).Do(ctx)
    check(err)
	return createIndex.Acknowledged, err
}

func ListIndexes(ctx context.Context, client *elastic.Client) ([]string, error) {
    indexes, err := client.IndexNames()
    check(err)
    return indexes, err
}

func RemoveIndex(ctx context.Context, client *elastic.Client, index string) ( bool, error) {
    exists, _ := client.IndexExists(index).Do(ctx)
    if !exists {
        fmt.Printf("Index %s doesn't exist. \n", index)
        os.Exit(1)
    }
    deleteIndex, err := client.DeleteIndex(index).Do(ctx)
    check(err)
    return deleteIndex.Acknowledged, err
}

func IndexDoc(ctx context.Context, client *elastic.Client, index string, 
        doctype string, indexbody string) (*elastic.IndexResponse, error) {
    put, err := client.Index().Index(index).Type(doctype).BodyJson(indexbody).Do(ctx)
    check(err)
    return put, err
}

func CreateRepo(ctx context.Context, client *elastic.Client, 
        reponame string, repotype string, repolocation string) bool {
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
        logger.LogError(serr)
		fmt.Println(serr)
	}
    
    src, berr := service.Do(ctx)
	if berr != nil {
        logger.LogError(berr)
		fmt.Println(berr)
	}
	_, jerr := json.Marshal(src)
	if jerr != nil {
        logger.LogInfo(fmt.Sprintf(`Marshaling to JSON failed: %v`, jerr))
		fmt.Printf("Marshaling to JSON failed: %v\n", jerr)
	}
    return src.Acknowledged
}

func RemoveRepo(ctx context.Context, client *elastic.Client, reponame string) bool {
    res, err := client.SnapshotDeleteRepository(reponame).Do(ctx)
    check(err)
    return res.Acknowledged
}

func SnapCreate(ctx context.Context, client *elastic.Client, 
        reponame string, snapname string, snapindex string) *bool {
    snapbody := fmt.Sprintf(`
    {
        "indices": "%s",
        "ignore_unavailable": "true",
        "include_global_state": "false",
        "wait_for_completion": "true"
    }`, snapindex)
    res , err := client.SnapshotCreate(reponame, snapname).BodyString(snapbody).Do(ctx)
    check(err)
    return res.Accepted
}

func SnapDelete(ctx context.Context, client *elastic.Client, 
        reponame string, snapname string) bool {
    res, err := client.SnapshotDelete(reponame, snapname).Do(ctx)
    check(err)
    return res.Acknowledged
}

func SnapRestore(ctx context.Context, client *elastic.Client, 
        reponame string, snapname string, snapindex string) bool {
    snapbody := fmt.Sprintf(`
    {
        "indices": "%s",
        "ignore_unavailable": "true",
        "include_global_state": "false"
    }`, snapindex)
    res, err := client.SnapshotRestore(reponame, snapname).BodyString(snapbody).Do(ctx)
    check(err)
    return res.Accepted
}
