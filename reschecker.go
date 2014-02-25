package main

import (
  "github.com/jamesduncombe/iron_payload"
  "time"
  "bytes"
  "sync"
  "fmt"
  "image/jpeg"
  "net/http"
  "encoding/json"
)

type Config struct {
  Validation struct {
    Width, Height int
  }
  Token string
  CallbackUrl string
  Urls []string
}

type Image struct {
  Image string  `json:"image"`
  Result bool   `json:"result"`
}

type Rules struct {
  Width  int  `json:"width"`
  Height int  `json:"height"`
}

type Result struct {
  CallbackUrl string  `json:"callbackUrl"`
  Token       string  `json:"token"`
  Rules  Rules    `json:"rules"`
  Images []Image  `json:"images"`
}

// Main functions

func checkRes(url string, i int, wg *sync.WaitGroup, iron *Config, result *Result) {

  defer wg.Done()

  start := time.Now()

  fmt.Println("Getting image...", i)
  resp, err := http.Get(url)
  defer resp.Body.Close()

  if err != nil {
    panic("Failed to get image from URL given")
  }

  elapsed := time.Since(start)
  fmt.Printf("Took %s GET for image %d\n", elapsed, i)

  fmt.Println("Decoding...", i)
  config, err := jpeg.DecodeConfig(resp.Body)

  if err != nil {
    fmt.Println("Invalid JPEG...", i)
  }

  // Collect results
  if config.Width >= iron.Validation.Width && config.Height >= iron.Validation.Height {
    result.Images = append(result.Images, Image{ url, true })
  } else {
    result.Images = append(result.Images, Image{ url, false })
  }

}

func main() {

  // Set initial start time (for elapsed time taken) and start new sync group
  // for go routines
  start := time.Now()
  wg := new(sync.WaitGroup)

  // Iron hold the config from params passed in
  iron := new(Config)

  // Result holds the result to send back to the API
  result := new(Result)

  // Get the payload from params passed in and parse
  m := iron_payload.GetPayload().(map[string]interface{})

  config := m["config"].(map[string]interface{})["validation"]

  iron.Token = m["config"].(map[string]interface{})["token"].(string)

  iron.CallbackUrl = m["config"].(map[string]interface{})["callbackUrl"].(string)

  iron.Validation.Width = int(config.(map[string]interface{})["width"].(float64))
  iron.Validation.Height = int(config.(map[string]interface{})["height"].(float64))

  images := m["config"].(map[string]interface{})["images"]
  for _, v := range images.([]interface{}) {
    iron.Urls = append(iron.Urls, v.(string))
  }

  // Setup initial return results
  result.CallbackUrl  = iron.CallbackUrl
  result.Token        = iron.Token
  result.Rules.Width  = iron.Validation.Width
  result.Rules.Height = iron.Validation.Height

  // Perform the actual work
  fmt.Println("Crunching... ", len(iron.Urls))
  for i, url := range(iron.Urls) {
    wg.Add(1)
    go checkRes(url, i, wg, iron, result)
  }

  // Wait for all go routines to stop before passing here
  wg.Wait()

  // Give our final elapsed times
  elapsed := time.Since(start)
  fmt.Printf("Took %s in all\n", elapsed)

  // Marshal the JSON
  json, err := json.Marshal(result)
  if err != nil {
    panic("Failed to Marshal JSON, check the structure")
  }

  res, err := http.Post(result.CallbackUrl, "application/json", bytes.NewReader(json))
  defer res.Body.Close()
  if err != nil {
    panic("Could not post back to server, check callback url")
  }

  fmt.Println("Posted back to server")

}
