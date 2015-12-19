[![Build Status](https://ci.matthewbrown.io/api/badges/mnbbrown/csv2api/status.svg)](https://ci.matthewbrown.io/mnbbrown/csv2api)

**csv2api** is a no frills daemon that turns csv files into an http API.
A couple of things to note:

- The API returns JSON by default. 
- It assumes the first line of the CSV file is the headers. 
- csv2api will replace all spaces with underscores and make the headers lowercase.

#### Install from docker

```sh
docker pull mnbbrown/csv2api:latest
docker run -d -v data:/tmp/data -p 8080:8080 mnbbrown/csv2api:latest
```

#### Install from source

Clone the repository to your Go workspace:

```
git clone git://github.com/mnbbrown/csv2api.git $GOPATH/src/github.com/mnbbrown/csv2api
cd $GOPATH/src/github.com/mnbbrown/csv2api
```

Add data:
```sh
mkdir -p /tmp/data
cat <<EOT >> /tmp/data/people.csv
Name,DOB
Matthew Brown,18/09/1992
EOT
```

Build and run:

```sh
go build
APP_ENV="production" PORT=8080 SERVE_FROM="/tmp/data" API_KEY="xyz" ./csv2api
```

#### API

Request data using [httpie](http://httpie.org):

```sh
# as JSON
http -v GET http://localhost:8080/api/v1/people Authorization:"Bearer xyz"
# [{"name":"Matthew Brown","dob":"18/09/1992"},...]

# as CSV
http -v GET http://localhost:8080/api/v1/people Authorization:"Bearer xyz" Accept:text/csv
# Name,DOB
# Matthew Brown,18/09/1992
```

##### Filtering

You can filter the fields the API returns with the `fields` query parameter as a comma separated list of fields. i.e. `?fields=name,dob` will return name and dob.

```sh
# as JSON
http -v GET http://localhost:8080/api/v1/people?fields=dob Authorization:"Bearer xyz"
# [{"dob":"18/09/1992"},...]
```