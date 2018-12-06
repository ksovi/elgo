### elgo is an Elasticsearch management tool written in Go.<br/>
elgo connects by default to http://localhost:9200. <br/>
Specify -host and -port to create a different URL.<br/>
Specify -host WITHOUT "http://".<br/>
<br/>
So far elgo supports index actions like create, remove, list, exists, <br/>
repository actions like create and remove, snapshot actions like create, <br/>
delete and restore, and cluster health information.

elgo also can index documents to indexes, using the index-doc action. <br/>
The document and metadata must be provided as a json file. this repository contains an example file <br/>
that can be used as a demo.<br/> 
The following example indexes the doc defined in o.json o index "ovi" with type "doc".<br/>


```$ ./elgo -action index-doc -i ovindex -type doc -f o.json```
Bulk request has been added and an example input file can be found here (b.json). Notice the single quotes ' at the beginning and ' end 
of 
the 
file. <br/> All bulk request must be provided including index, type and ID. <br/>
All bulk requests must be provided in a json file, inside ticks or single quotes, just like the example file b.json.<br/>

```
$./elgo.go -action bulk-request -f b.json 
Using Elasticsearch URL:  http://localhost:9200
Created documents:  4
Indexed documents:  6
Updated documents:  2
Deleted documents:  2

```

elgo writes a logfile called elgo.out in the current working directory. <br/>

When creating an index an input file can also be specified (examples: i.json or simple.index.json) to define <br/>
index settings at creation, or mappings etc.<br/>

Creating a repository requirs a name, a type and location. location can be a directory or directory tree <br/>
that will be created inside ```path.repo``` as specified in ```elasticsearch.yml```.

```
$ ./elgo -action create-repo -r repo1 -type fs -l tmp/repo1
```
In this example we create a repository of type fs called repo1 in tmp/repo1. <br/>
If the path.repo path is /var/backups, the full directory for this repository will be /var/backups/tmp/repo1/<br/>
Make sure repo.path is defined and the path exists.

snap-create and snap-restore support all indices, or just specific indices.
For ex: 

```$ ./elgo -action snap-create -s snap2 -r repo0 -i ovi*``` 

creates a snapshot for all indices starting with "ovi", <br/>
but if "-i ovi" is used , only index "ovi" will be snapshoted. <br/>
If no index is specified with "-i", all indexes will be implied. Restore action works the same way. <br/>



```
$ ./elgo

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

```

