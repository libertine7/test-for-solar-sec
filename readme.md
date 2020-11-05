# Test for SolarSecurity

to run use:
go run main.go postgres://postgres:mysecretpassword@172.17.0.2/postgres?sslmode=disable

to run docker if you need use:
docker run --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres

for test run client in client folder