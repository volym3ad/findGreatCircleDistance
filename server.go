package main

import ("github.com/go-martini/martini"
        "github.com/kellydunn/golang-geo" // main lib for great distance search
        "github.com/go-redis/redis" // go client for redis
        
        "strings"
        "strconv"
        "regexp"
        )

// refactoring is needed!!!!

var (
    match = "^.*[0-9],.*[0-9]$"
    iter = 1
)

// add exception to server initialization
func redisDBInitialization(address string,
                           password string,
                           database int) *redis.Client {

    return redis.NewClient(&redis.Options{
        Addr:     address,
        Password: password,
        DB:       database,
    })
}

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

    client := redisDBInitialization("127.0.0.1:6379", "", 0) // 0 means to use default Redis db
    martini.Env = martini.Prod // cause development is too easy =p
    m := martini.Classic()

    m.Get("/", func() string {
    	return "This program helps you calculate Great Circle Distance"
    	})

    m.Get("/find/:from/:to", func(params martini.Params) (int, string) {

        key := iter
        value := params["from"]+" ; "+params["to"]

        // regexp matching    
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

        err := client.Set(strconv.Itoa(key), value, 0).Err()
        if err != nil {
            panic(err)
        }
        
        iter += 1
        return 200, "Great Circle Distance is " + findGreatCircleDistance(x[0],x[1],y[0],y[1]) + " km"    
        })

    m.Get("/history", func() string {
        
        // latest history requests
        // not sure we need to declare array here -_-

        values := make([]string, iter)
        var err_value error
        for i:=0; i < iter-1; i++ {
            values[i], err_value = client.Get(strconv.Itoa(i+1)).Result()
            if err_value == redis.Nil {
                return "key does not exists"
            } else if err_value != nil {
                panic(err_value)
            }
        }

        return strings.Join(values[:], "\n")
        })

    // delete history
    m.Delete("/history", func() string {
        
        for i:=0; i < iter-1; i++ {
            client.Del(strconv.Itoa(i+1))
        }
        // add exceptions when fetching key id
        // and database connection exception....maybe..

        iter = 1
        return "Delete successful"
        })

    m.Get("/test", func() string {
        // just some testing values to validate function's work
        return "Great Circle Distance is " + findGreatCircleDistance(42.25,120.2,30.25,112.2) + " km"
    	})

    // processing 404 errors
    m.NotFound(func() string { 
        return "Seems like you doing something nasty.. 404 error"
        })
    m.Run()
}
