package main

import ("fmt"
        "net/http"
        "encoding/json")

type Response struct {
    CelleryTool struct {
        Version          string `json:"version"`
        APIVersion       string `json:"apiVersion"`
        BallerinaVersion string `json:"ballerinaVersion"`
        GitCommit        string `json:"gitCommit"`
        Built            string `json:"built"`
        OsArch           string `json:"osArch"`
        Experimental     bool   `json:"experimental"`
    } `json:"celleryTool"`
    CelleryRepository struct {
        Server        string `json:"server"`
        APIVersion    string `json:"apiVersion"`
        Authenticated bool   `json:"authenticated"`
    } `json:"celleryRepository"`
    Kubernetes struct {
        Version string `json:"version"`
        Crd     string `json:"crd"`
    } `json:"kubernetes"`
    Docker struct {
        Registry string `json:"registry"`
    } `json:"docker"`
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "%s", json_msg)
}

func main() {
    response := Response{}
    celleryTool := CelleryTool{}
    celleryRepository := CelleryRepository{}
    kubernetes := Kubernetes{}
    docker := Docker{}
	
    celleryTool.version = "0.50-beta1"
    celleryTool.apiVersion = "0.95"
    celleryTool.ballerinaVersion = "0.987.0"
    celleryTool.gitCommit = "78a6b6b"
    celleryTool.built = "Tue Oct 23 22:41:53 2018"
    celleryTool.

    http.HandleFunc("/version", versionHandler)
    http.ListenAndServe(":8080", nil)
}

