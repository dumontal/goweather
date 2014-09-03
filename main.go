package main

import (
    "encoding/json"
    "net/http"
    "strings"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello! I am a weather server.\n"))
    })

    http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
        city := strings.SplitN(r.URL.Path, "/", 3)[2]

        data, err := query(city)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        json.NewEncoder(w).Encode(data)
    })

    http.ListenAndServe(":8080", nil)
}

func query(city string) (WeatherData, error) {
    resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city)
    if err != nil {
        return WeatherData{}, err
    }

    defer resp.Body.Close()

    var data WeatherData
    err = json.NewDecoder(resp.Body).Decode(&data)

    if err != nil {
        return WeatherData{}, err
    }

    return data, nil
}

type WeatherData struct {
    Name string `json:"name"`
    Main struct {
        Kelvin float64 `json:"temp"`
    } `json:"main"`
}
