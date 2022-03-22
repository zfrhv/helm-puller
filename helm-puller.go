package main

import (
  "flag"
  "fmt"
  "os"
  "archive/tar"
  "compress/gzip"
  "io"
  "io/ioutil"
  "log"
  "strings"
)

type CV struct {
  chart string
  values string
}

func ScanChart(gzipStream io.Reader) map[string]*CV {
  files := make(map[string]*CV)

  uncompressedStream, err := gzip.NewReader(gzipStream)
  if err != nil {
    log.Fatal("ExtractTarGz: file not found")
  }

  tarReader := tar.NewReader(uncompressedStream)

  for true {
    header, err := tarReader.Next()

    if err == io.EOF {
      break
    }

    if err != nil {
      log.Fatalf("ExtractTarGz: Next() failed: %s", err.Error())
    }

    path := strings.Split(header.Name, "/")
    fileName   := path[len(path)-1]
    fileFolder := path[len(path)-2]
    if fileName == "Chart.yaml" || fileName == "values.yaml" {
      if files[fileFolder] == nil {
        files[fileFolder] = &CV{}
      }
      content, _ := ioutil.ReadAll(tarReader)
      if fileName == "Chart.yaml" {
        files[fileFolder].chart = string(content)
      }
      if fileName == "values.yaml" {
        files[fileFolder].values = string(content)
      }
    }
    if false {
      fmt.Println(files)
    }
  }
  return files
}

func main() {
  // Subcommands
  pullCommand := flag.NewFlagSet("pull", flag.ExitOnError)
  pushCommand := flag.NewFlagSet("push", flag.ExitOnError)

  chartName := pullCommand.String("chart", "", "chart to collect the images from, for example --chart kube-prometheus-stack-32.2.1.tgz")

  // Verify that a subcommand has been provided
  // os.Arg[0] is the main command
  // os.Arg[1] will be the subcommand
  if len(os.Args) < 2 {
    fmt.Println("Please use one of the commands bellow:\n\nhelm-puller pull\nhelm-puller push\n")
    os.Exit(1)
  }

  switch os.Args[1] {
  case "pull":
    pullCommand.Parse(os.Args[2:])
  case "push":
    pushCommand.Parse(os.Args[2:])
  default:
    fmt.Println("Please use one of the commands bellow:\n\nhelm-puller pull\nhelm-puller push\n")
    os.Exit(1)
  }

  if pullCommand.Parsed() {
    // Required Flags
    if *chartName == "" {
      fmt.Println("error: flag --chart cannot be empty\n")
      pullCommand.PrintDefaults()
      os.Exit(1)
    }


    chart, err := os.Open(*chartName)
    if err != nil {
      fmt.Println("error opening %s: %s", *chartName, err.Error())
    }
    files := ScanChart(chart)
    for pack := range files {
      fmt.Println(files[pack].values)

    }
  }

  if pushCommand.Parsed() {
      fmt.Printf("in parsing section\n")
  }
}
