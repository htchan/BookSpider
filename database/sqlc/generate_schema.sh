docker run --rm -d --name bookspider-sqlc-generator \
  -e POSTGRES_USER=book_spider -e POSTGRES_PASSWORD=password -e POSTGRES_DB=db \
  -v ${PWD}/../migrations:/migrations -v ./:/sqlc/ postgres:17.5-alpine

# check if the container is ready
while ! docker exec -i bookspider-sqlc-generator pg_isready -U book_spider; do sleep 1; done

# run migration and dump schema
docker exec bookspider-sqlc-generator bash -c 'for filename in /migrations/*.up.sql; do psql -U book_spider -d db -f $filename; done' && \
docker exec bookspider-sqlc-generator bash -c "pg_dump -U book_spider -d db -t books -t writers -t errors --schema-only > /sqlc/schema.sql"

# kill container
docker kill bookspider-sqlc-generator