# Instagram Sniffer

<<<<<<< HEAD
### Test branch Create User info struct and rewrite recursive code
=======
##### Test branch Create User info struct and rewrite recursive code
>>>>>>> c18ca1b52cdab4e5644fe07ade9a0dffe7bedbf6

Build:
#### docker build -t insta .
Start:
<<<<<<< HEAD
#### docker run --rm --name insta -p 8000:8000 -e THR=2 insta
=======
#### docker run --rm --name insta -p 8080:8080 -e THR=2 -e BS=10 insta

Check docker space: 
#### docker exec -it insta sh

DB Migration:
#### migrate create -ext sql -dir db/migrations -seq create_items_table

Check docker-compose:
#### docker-compose up -d
>>>>>>> c18ca1b52cdab4e5644fe07ade9a0dffe7bedbf6
