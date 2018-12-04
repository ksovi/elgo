// Copyright 2018-present Ovi Chis www.ovios.org All rights reserved.
// Use of this source code is governed by a MIT-license.

package actions

import (
    "fmt"
    "elk2/elkluster"
    "elk2/logger"
    "context"
    "os"
    "io/ioutil"
)

/* Add actions:
 * create-repo +
 * delete-repo +
 * snap-create +
 * snap-delete +
 * snap-restore +
 * cluster-info 
 */


func check(e error) {
    if e != nil {
        logger.LogError(e)
        panic(e)
    }
}

func PassAction(action, url, input_file, indexname,actiontype, reponame, repolocation, snapname, ElkUsage string) {
    ctx := context.Background()
    client := elkluster.InnitiateClient(ctx, url)    
    switch action {
        case "create-index":
            if indexname == "" {
                fmt.Println("An index name is required for create-index. [-i <indexname>]")
                os.Exit(1)
            }
            indexbody := ""
            if input_file != "" {
                data, err := ioutil.ReadFile(input_file)
                check(err)
                indexbody = string(data)
            }
            res, err := elkluster.CreateIndex(ctx, client, indexname, indexbody)
            check(err)
            if res {
                logger.LogInfo(fmt.Sprintf(`Successfully created index %s`, indexname))
                fmt.Println("Successfully created index ", indexname)
            }
        case "remove-index":
            if indexname == "" {
                fmt.Println("An indexname is required for remove-index. [-i <indexname>]")
                os.Exit(1)
            }
            res, err := elkluster.RemoveIndex(ctx, client, indexname)
            check(err)
            if res {
                logger.LogInfo(fmt.Sprintf(`Successfully removed index %s`, indexname))
                fmt.Println("Successfully removed index ", indexname)
            }
        case "list-indexes":
            indexes, err := elkluster.ListIndexes(ctx, client)
            check(err)
            i := 0
            for _, name := range indexes {
                i ++
                fmt.Println(i, ":", name)
            }
        case "index-exists":
            if indexname == "" {
                fmt.Println("An indexname is required for index-exists. [-i <indexname>]")
                os.Exit(1)
            }
            exists := elkluster.IndexExists(ctx, client, indexname)
            if exists {
                fmt.Printf("==> Index %s exists.\n", indexname)
            } else {
                fmt.Printf("==> Index %s doesn't exist.\n" , indexname)
            }
        case "index-doc":
            if indexname == "" {
                fmt.Println("An indexname is required for index-doc. [-i <indexname>]")
                os.Exit(1)
            }
            if actiontype == "" {
                fmt.Println("Doc type is required for index-doc. [-type <type>]")
                os.Exit(1)
            }
            indexbody := ""
            if input_file == "" {
                fmt.Println("An input file is required for index-doc. [-f <path to file>]")
                os.Exit(1)
            } else {
                data, err := ioutil.ReadFile(input_file)
                check(err)
                indexbody = string(data)
            }
            res, err := elkluster.IndexDoc(ctx, client, indexname, actiontype, indexbody)
            check(err)
            logger.LogInfo(fmt.Sprintf(`Indexed document with ID: %s to index: %s, type: %s`, res.Id, res.Index, res.Type))
            fmt.Printf("Indexed document with ID: %s to index: %s, type: %s\n", res.Id, res.Index, res.Type)
        case "create-repo":
            if actiontype ==  "" {
                fmt.Println("Type, name and location are all required for create-repo. [-type <type> -r <repo name> -l <repo path>]")
                os.Exit(1)
            }
            if reponame == "" {
                fmt.Println("Type, name and location are all required for create-repo. [-type <type> -r <repo name> -l <repo path>]")
                os.Exit(1)
            }
            if repolocation == "" {
                fmt.Println("Type, name and location are all required for create-repo. [-type <type> -r <repo name> -l <repo path>]")
                os.Exit(1)
            }
            result := elkluster.CreateRepo(ctx, client, reponame, actiontype, repolocation)
            if result == true {
                fmt.Printf("Successfully created repo: %s of type: %s at location: %s\n", reponame, actiontype, repolocation)
                logger.LogInfo(fmt.Sprintf(`Successfully created repo: %s of type: %s at location: %s`, reponame, actiontype, repolocation))
            }
        case "remove-repo":
            if reponame == "" {
                fmt.Println("Repo name is required for remove-repo. [-r <repo name>]")
                os.Exit(1)
            }
            result := elkluster.RemoveRepo(ctx, client, reponame)
            if result == true {
                fmt.Printf("Successfully removed repo: %s.\n", reponame)
                logger.LogInfo(fmt.Sprintf(`Successfully removed repo %s`, reponame))
            }
        case "snap-create":
            if reponame == "" {
                fmt.Println("Repo name is required for snap-create. [-r <repo name>]")
                os.Exit(1)
            }
            if snapname == "" {
                fmt.Println("Snap name is required for snap-create. [-s <snap name>]")
                os.Exit(1)
            }
            snapindex := indexname
            result := elkluster.SnapCreate(ctx, client, reponame, snapname, snapindex)
            if *result {
                fmt.Printf("Successfully created snap %s in repo %s.\n", snapname, reponame)
                logger.LogInfo(fmt.Sprintf(`Successfully created snap %s in repo %s.`, snapname, reponame))
            }
        case "snap-delete":
            if reponame == "" {
                fmt.Println("Repo name is required for snap-delete. [-r <repo name>]")
                os.Exit(1)
            }
            if snapname == "" {
                fmt.Println("Snap name is required for snap-delete. [-s <snap name>]")
                os.Exit(1)
            }
            result := elkluster.SnapDelete(ctx, client, reponame, snapname)
            if result == true {
                fmt.Printf("Successfully removed snapshot %s from repository %s.\n", snapname, reponame)
                logger.LogInfo(fmt.Sprintf(`Successfully removed snapshot %s from repository %s`, snapname, reponame))
            }
        case "snap-restore": 
            if reponame == "" {
                fmt.Println("Repo name is required for snap-restore. [-r <repo name>]")
                os.Exit(1)
            }
            if snapname == "" {
                fmt.Println("Snap name is required for snap-restore. [-s <snap name>]")
                os.Exit(1)
            }
            result := elkluster.SnapRestore(ctx, client, reponame, snapname, indexname)
            if result == true {
                fmt.Printf("Successfully restored snapshot %s from repository %s.\n", snapname, reponame)
                logger.LogInfo(fmt.Sprintf(`Successfully restored snapshot %s from repository %s`, snapname, reponame))
            }
        case "cluster-info":
            elkluster.ClusterInfo(ctx, client)
        default:
            fmt.Printf("Action %s is not valid. \n" , action)
            fmt.Println(ElkUsage)
    }
}

