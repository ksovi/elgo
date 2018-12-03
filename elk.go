package main

import (
    "fmt"
    "elk2/elkluster"
    "elk2/logger"
    "strconv"
    "context"
    "flag"
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
        panic(e)
    }
}

// This is to be able to handle multiple actions at once.
// will probably be removed as multiple actions don't really make sense

type arrayFlags []string

func (i *arrayFlags) String() string {
    return "No default action."
}

func (i *arrayFlags) Set(action string) error {
    *i = append(*i, action)
    return nil
}

var elkActions arrayFlags
// 

// start main

func main() {
    // Setting up a context 
    ctx := context.Background()
    
    // Handling command loine arguments 
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
    ==> snap-restore - requires -r <repo name> -s <snap name>`
    
    hostPtr := flag.String("host", "localhost", "Elastic host or IP.")
    portPtr := flag.Int("port", 9200, "Elastic port number.")
    flag.Var(&elkActions, "action", "Action to be executed against the Elasticsearch cluster.")
    indexPtr := flag.String("i", "", "Index name")
    inputfilePtr := flag.String("f", "", "Input json file.")
    repoNamePtr := flag.String("r", "", "Repo name")
    repoLocPtr := flag.String("l", "", "Repo location.")
    snapNamePtr := flag.String("s", "", "Snap name.")
    typePtr := flag.String("type", "", "Doc type for indexing")
    
    flag.Parse()
    
    if flag.NFlag() < 1 {
        fmt.Println(ElkUsage)
        os.Exit(1)
    }
    
    host := *hostPtr
    port := *portPtr
    input_file := *inputfilePtr
    indexname := *indexPtr
    actiontype := *typePtr
    reponame := *repoNamePtr
    repolocation := *repoLocPtr
    snapname := *snapNamePtr
    
    p := strconv.Itoa(port)
    url := "http://"+host+":"+p
    
    client := elkluster.InnitiateClient(ctx, url)
    
    for _, action := range elkActions {
        switch action {
            case "create-index":
                if indexname == "" {
                    fmt.Println("An index name is required for create-index. [-i <indexname>]")
                    os.Exit(1)
                }
                indexbody := ""
                if input_file != "" {
                data, err := ioutil.ReadFile(input_file)
                logger.LogError(err)
                check(err)
                indexbody = string(data)
                }
                res, err := elkluster.CreateIndex(ctx, client, indexname, indexbody)
                check(err)
                if res {
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
                exists, err := elkluster.IndexExists(ctx, client, indexname)
                check(err)
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
                logstring := fmt.Sprintf(`Indexed document with ID: %s to index: %s, type: %s`, res.Id, res.Index, res.Type)
                logger.LogInfo(logstring)
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
                    logstring := fmt.Sprintf(`Successfully created repo: %s of type: %s at location: %s`, reponame, actiontype, repolocation)
                    logger.LogInfo(logstring)
                }
            case "remove-repo":
                if reponame == "" {
                    fmt.Println("Repo name is required for remove-repo. [-r <repo name>]")
                    os.Exit(1)
                }
                result := elkluster.RemoveRepo(ctx, client, reponame)
                if result == true {
                    fmt.Printf("Successfully removed repo: %s.\n", reponame)
                    logstring := fmt.Sprintf(`Successfully removed repo %s`, reponame)
                    logger.LogInfo(logstring)
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
                    logstring := fmt.Sprintf(`Successfully created snap %s in repo %s.`, snapname, reponame)
                    logger.LogInfo(logstring)
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
                    logstring := fmt.Sprintf(`Successfully removed snapshot %s from repository %s`, snapname, reponame)
                    logger.LogInfo(logstring)
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
                    logstring := fmt.Sprintf(`Successfully restored snapshot %s from repository %s`, snapname, reponame)
                    logger.LogInfo(logstring)
                }
            default:
                fmt.Printf("Action %s is not valid. \n" , action)
                fmt.Println(ElkUsage)
        }
    }
}
