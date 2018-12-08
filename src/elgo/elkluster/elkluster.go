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

func BulkAction(ctx context.Context, client *elastic.Client, bulkbody string) {
    elgoReq := elastic.NewBulkIndexRequest().Doc(bulkbody).Index("_").Type("_")
    bulkRequest := client.Bulk()
    bulkRequest = bulkRequest.Add(elgoReq)
    bulkResponse, err := bulkRequest.Do(ctx)
    check(err)
    
    cr := bulkResponse.Created()
    de := bulkResponse.Deleted()
    ix := bulkResponse.Indexed()
    up := bulkResponse.Updated() 
    
    logger.LogInfo(fmt.Sprintf(`Total number of documents created: %d`, len(cr)))
    logger.LogInfo(fmt.Sprintf(`Total number of documents indexed: %d`, len(ix)))
    logger.LogInfo(fmt.Sprintf(`Total number of documents updated: %d`, len(up)))
    logger.LogInfo(fmt.Sprintf(`Total number of documents deleted: %d`, len(de)))
    
    if len(cr) > 0 {
        for i, _ := range cr {
            prstring := fmt.Sprintf(`Created document ID: %s Index: %s Type: %s Version: %d Elasticsearch response: %d Error: %v ` , cr[i].Id, cr[i].Index, cr[i].Type, cr[i].Version, cr[i].Status, cr[i].Error)
            logger.LogInfo(prstring)
        } 
    } else {
        prstring := "No documents created."
        fmt.Println(prstring)
        logger.LogInfo(prstring)
    }
    if len(de) > 0 {
        for i, _ := range de {
            prstring := fmt.Sprintf(`Deleted document ID: %s Index: %s Type: %s Version: %d Elasticsearch response: %d Error: %v` , de[i].Id, de[i].Index, de[i].Type, de[i].Version, de[i].Status, de[i].Error)
            logger.LogInfo(prstring)
        } 
    } else {
        prstring := "No documents deleted."
        fmt.Println(prstring)
        logger.LogInfo(prstring)
    }
    if len(ix) > 0 {
        for i, _ := range ix {
            prstring := fmt.Sprintf(`Indexed document ID: %s Index: %s Type: %s Version: %d Elasticsearch response: %d Error: %v` , ix[i].Id, ix[i].Index, ix[i].Type, ix[i].Version, ix[i].Status, ix[i].Error)
            logger.LogInfo(prstring)
        } 
    } else {
        prstring := "No documents indexed."
        fmt.Println(prstring)
        logger.LogInfo(prstring)
    }
    if len(up) > 0 {
        for i, _ := range up {
            prstring := fmt.Sprintf(`Updated document ID: %s Index: %s Type: %s Version: %d Elasticsearch response: %d Error: %v` , up[i].Id, up[i].Index, up[i].Type, up[i].Version, up[i].Status, up[i].Error)
            logger.LogInfo(prstring)
        } 
    } else {
        prstring := "No documents updated."
        fmt.Println(prstring)
        logger.LogInfo(prstring)
    }
    
    fmt.Println("Created documents: ", len(cr))
    fmt.Println("Indexed documents: ", len(ix)) 
    fmt.Println("Updated documents: ", len(up))
    fmt.Println("Deleted documents: ", len(de))
}

func ElgoSearch(ctx context.Context, client *elastic.Client, indexname, sfield, svalue string, maxreturns int) {
    termQuery := elastic.NewTermQuery(sfield, svalue)
    logger.LogInfo(fmt.Sprintf("Running search on index: %s , search field: %s, search value: %s", indexname, sfield, svalue))     
    SR, err := client.Search().Index(indexname).Query(termQuery).From(0).Size(maxreturns).Pretty(true).Do(ctx)
    check(err)
    fmt.Printf("Search took %d milliseconds.\n", SR.TookInMillis)
    logger.LogInfo(fmt.Sprintf("Search took %d milliseconds.", SR.TookInMillis))
    fmt.Printf("Found %d results. \n", SR.Hits.TotalHits)
    logger.LogInfo(fmt.Sprintf("Found %d results.", SR.Hits.TotalHits))
    
    for _, v := range SR.Hits.Hits {
        fmt.Println("Document ID: ", v.Id)
        fmt.Println("Document Type: ", v.Type)
        items := make(map[string]interface{})
        err := json.Unmarshal(*v.Source, &items)
        check(err)
        for m , item := range items {
            fmt.Println(m, ": ", item)
        }
    }
}
