# Instagram Sniffer

##### Test branch Create User info struct and rewrite recursive code

Build:
#### docker build -t insta .
Start:
#### docker run --rm --name insta -p 8080:8080 -e THR=2 -e BS=10 insta

Check docker space: 
#### docker exec -it insta sh

DB Migration:
#### migrate create -ext sql -dir db/migrations -seq create_items_table

Check docker-compose:
#### docker-compose up -d
<<<<<<< HEAD

We are trying new methods
=======
>>>>>>> 26c7092724eeb1f900ecaa8d4fd94ec0dafe9b45
