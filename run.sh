until $(curl --output /dev/null --silent --location --request GET --fail 'http://accountapi:8080/v1/health' --header 'Content-Type: application>
go test -v ./... -coverprofile cover.out