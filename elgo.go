// Copyright 2018-present Ovi Chis www.ovios.org All rights reserved.
// Use of this source code is governed by a MIT-license.

package main

import (
    "flag"
    "os"
    "elgo/actions"
    ///"strconv"
    "fmt"
)

func main() {
    ElkUsage := `
    At least one action is required. 
    Use -action with one of the following supported actions:
    ==> create-index - requires at least -i <index name>. Optional -f input body in json format to pass index settings.
    ==> remove-index - requires -i <index name>
    ==> list-indexes - returns a list of all indexes.
    ==> index-exists - requires -i <index name>
    ==> index-doc - requires -i <index name> -type <type> -f input json file to be indexed.
    ==> create-repo - requires -r <repo name> -type <type> -l <location>.
    ==> remove-repo - required -r <repo name>.
    ==> snap-create - required -r <repo name> -s <snap name>. Optional -i <index name>. * or multiple indexes accepted.
    ==> snap-delete - required -r <repo name> -s <snap name>
    ==> snap-restore - requires -r <repo name> -s <snap name>
    ==> cluster-info - returns cluster information
    ==> bulk-request - create/index/delete or update using the bulk interface. -f <json file> required containing the bulk request to be sent to elasticsearch.
    ==> search - searches through indexed documents for a specific field, like Id, Type, etc. in one or more or all indexes.`
    
    urlPtr := flag.String("url", "http://localhost:9200", "Elastic host or IP and port.")
    actionPtr := flag.String("action", "", "Action to execute")
    indexPtr := flag.String("i", "", "Index name")
    inputfilePtr := flag.String("f", "", "Input json file.")
    repoNamePtr := flag.String("r", "", "Repo name")
    repoLocPtr := flag.String("l", "", "Repo location.")
    snapNamePtr := flag.String("s", "", "Snap name.")
    typePtr := flag.String("type", "", "Doc type for indexing")
    maxReturnPtr := flag.Int("mr", 10, "Search max number of returns.")
    sfieldPtr := flag.String("sf", "", "Search field.")
    svaluePtr := flag.String("sv", "", "Search value.")
    
    flag.Parse()
    
    if flag.NFlag() < 1 {
        fmt.Println(ElkUsage)
        os.Exit(1)
    }
    
    action := *actionPtr
    url := *urlPtr
    input_file := *inputfilePtr
    indexname := *indexPtr
    actiontype := *typePtr
    reponame := *repoNamePtr
    repolocation := *repoLocPtr
    snapname := *snapNamePtr
    maxreturns := *maxReturnPtr
    sfield := *sfieldPtr
    svalue := *svaluePtr
    
   
    actions.PassAction(action, url, input_file, indexname, actiontype, 
                       reponame, repolocation, snapname, sfield, svalue,  ElkUsage, maxreturns)
}
