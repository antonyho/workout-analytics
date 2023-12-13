Workout Analytics
====

## Assumptions

### Project hierarchy
This project is not taking the project structure with `pkg/` and `cmd/`.
Because this is a simple project. A flat hierarchy is easier to read.

### Fail fast approach
This application fails fast and replies with the corresponding HTTP status code.
No pretty JSON response body will be given. Even though it is not hard to implement description error response. 
I am not overdoing on non requirement feature.

### When today is not Monday
The analysed period will be dated back from last Monday to the (N)th week(s) before.

### Running this application on Sunday at 23:59:59
This application will not paranoid about the execution time which might cause big difference to the outcome due to that.

### Go version
I am not using latest, at the time of writing, Go v1.21.5.
Because the linter `golangci-lint` latest version is not supporting latest Go version.

### No custom error
I did not use custom error because this project is really simple.
I cannot find a good place to use my custom error. 


### Usage

#### Start server locally
You may use the fast way to start the server without building it into binary.

Suppose you are using a POSIX shell environment with `Make` utilities installed.
```shell
make run
```

#### Start server on Docker image
Of course, you need to have Docker installed on your system.

##### Build Docker image
```shell
make docker-image
```

##### Start Docker image
```shell
docker run --rm -p 8080:8080 twaiv/workout-analytics
```

##### Docker Compose
Apart from the above `docker build` and `docker run` commands. `docker compose` manifest file is also provided.

To start as a detached suite
```shell
docker compose up -d
```

To tear down the suite
```shell
docker compose down
```


#### Testing
There is a JSON file in `testdata` under the project directory for you to test. But the test is related to current date time.
The test data is during 2023-11-04 to 2023-11-26. You might not get any analysis, if you test very late and set the number of weeks too small.

Suppose you are using POSIX shell with `cURL` installed. Following command can test the localhost server with the project test data file. 
```shell
curl -v -H "Content-Type: application/json" --data-binary "@testdata/request.json" "http://localhost:8080/analyse?nweeks=48"
```