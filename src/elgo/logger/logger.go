// Copyright 2018-present Ovi Chis www.ovios.org All rights reserved.
// Use of this source code is governed by a MIT-license.

package logger

import (
    "log"
    "os"
)

var logfile string = "elgo.out"
//var logprefix string = "elkgo"

func LogError(logerror error) {
    f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    log.SetOutput(f)
    //log.SetPrefix(logprefix)
    log.Println(logerror)
}

func LogInfo(logstring string) {
    f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    log.SetOutput(f)
    log.Println(logstring)
}

