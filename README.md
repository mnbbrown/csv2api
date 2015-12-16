
**csv2api** is a no frills daemon that turns csv files into an http API.

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

The API returns JSON by default. It assumes the first line of the CSV file is the headers. It will replace all spaces with underscores (" " => "_") and make the headers lowercase.
