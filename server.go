package main

import ("github.com/go-martini/martini"
		"github.com/kellydunn/golang-geo"
        "strings"
		"strconv"
        "regexp")

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

func regexpStringMatching(text string, reg string) bool {
    r, err := regexp.Compile(reg)

    if err != nil {
        return false
    }

    if r.MatchString(text) == true {
        return true
    } else {
        return false
    }
}

func main() {

    m := martini.Classic()
    m.Get("/", func() string {
    	return "This program helps you calculate Great Circle Distance"
    	})

    m.Get("/find/:from/:to", func(params martini.Params) (int, string) {

        // regexp matching
        match := "^.*[0-9],.*[0-9]$"        
        if !regexpStringMatching(params["from"], match) || !regexpStringMatching(params["to"], match) {
            return 400, "Bad Request => from and to params must be ^.*[0-9],.*[0-9]$"
        }

        from := strings.Split(params["from"], ",")
        to := strings.Split(params["to"], ",")

        var x, y [2] float64
        var err1, err2 error

        for i := 0; i <= 1; i++ {
            x[i], err1 = strconv.ParseFloat(from[i], 64)
            y[i], err2 = strconv.ParseFloat(to[i], 64)
            
            if err1 != nil || err2 != nil  {
                return 406, "Not Acceptable => False format of parameters, must be float64"
            }
        }

        return 200, "Great Circle Distance is " + findGreatCircleDistance(x[0],x[1],y[0],y[1]) + " km"
        })

    m.Get("/test", func() string {
        return "Great Circle Distance is " + findGreatCircleDistance(42.25,120.2,30.25,112.2) + " km"
    	})
    m.Run()
}
