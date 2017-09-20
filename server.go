package main

import ("github.com/go-martini/martini"
		"github.com/kellydunn/golang-geo"
        "strings"
		"strconv")

// Lon stands for Longitude
// Lat stands for Latitude
func findGreatCircleDistance(fromLon float64, 
                             fromLat float64,
                             toLon float64,
                             toLat float64) string {

     start := geo.NewPoint(fromLon, fromLat)
     finish := geo.NewPoint(toLon, toLat)

     // find the great circle distance between them
     return strconv.FormatFloat(start.GreatCircleDistance(finish), 'f', 6, 64)
}

func main() {

    m := martini.Classic()
    m.Get("/", func() string {
    	return "This program helps you calculate Great Circle Distance"
    	})

    m.Get("/find/:from/:to", func(params martini.Params) string {
        from := strings.Split(params["from"], ",")
        to := strings.Split(params["to"], ",")
        
        fromLon, err := strconv.ParseFloat(from[0], 64)
        fromLat, err := strconv.ParseFloat(from[1], 64)
        toLon, err := strconv.ParseFloat(to[0], 64)
        toLat, err := strconv.ParseFloat(to[1], 64)
        
        if err != nil {
            return "False parameters"
        }

        return "Great Circle Distance is " + findGreatCircleDistance(fromLon,
                                                                     fromLat,
                                                                     toLon,
                                                                     toLat) + " km"
        })

    m.Get("/test", func() string {
        return "Great Circle Distance is " + findGreatCircleDistance(42.25, 
                                                                     120.2,
                                                                     30.25,
                                                                     112.2) + " km"
    	})
    m.Run()
}
