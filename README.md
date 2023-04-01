# Race condition detection
Race condition detection with a real world reverse proxy written in GO.

[![blog post](/docs/high-level-diagram.png)](https://medium.com/@kristicevic.antonio/detecting-race-conditions-in-go-265c9c049155)

This repo has two main branches
* `before` - with race condition defect
* `after` - without race condition

Check out the blog post [here](https://medium.com/@kristicevic.antonio/detecting-race-conditions-in-go-265c9c049155).

## 👟 How to run
Prerequisites:
* [Docker] - Your favourite container platform
* [go] - An open-source programming language supported by Google
* [k6] - The best developer experience for load testing

Run the following command to spin up the system using docker compose
```sh
$ docker-compose up --build -d 
```

The reverse proxy is reachable on port `8080`
```sh
$ curl -X POST localhost:8080/a/foo/bar
```

### 🧪 Testing 
To run the tests execute the following command
```sh
$ go test
```

If you want to enable race detector tool add the `-race` flag

```sh
$ go test -race
```

### 🪨 Load testing
While the system is running, execute the following command to start the load test
```sh
$ k6 run --vus 10 --iterations 1000 load-test.js
```

This command will execute the `load-test.js` script with 1000 iterations and 10 virual users.

## 🎞️ Rendering terminal GIFs
Make sure you have [terminalizer] on your system.

Run the following command to render the terminal recordings:
```sh
$ terminalizer render docs/<file_name>.yml
```

[Docker]: <https://www.docker.com>
[go]: <https://go.dev>
[k6]: <https://k6.io>
[terminalizer]: <https://www.terminalizer.com>