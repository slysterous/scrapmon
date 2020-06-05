#!/bin/sh
set -e 

migrate_db() {
    # Attempt to run migrations, retry if the database container is not yet ready
    echo "Running migrations on $1"
    i=0
    
    ./migrate -database postgres://$4:$5@$1:$2/$3?sslmode=disable -verbose -source file://./migrations up
    #./migrate -database postgres://postgres:password@print-scrape-db:5432/print-scrape?sslmode=disable -verbose -source file://./migrations up
    while [ $? -ne 0 -a $i -lt 10 ]; do
        echo "Database not ready (attempt #$i), retrying.."
        sleep 2
        i=`expr $i + 1`
         ./migrate -database postgres://$4:$5@$1:$2/$3?sslmode=disable -verbose -source file://./migrations up
    done
}

check_success() {
    # Exit if the last command failed
    if [ $? -ne 0 ]; then
        echo "Last command failed, exiting.."
        exit 1
    fi
}

set +e

# Run default migration
echo "Migrating database..."
migrate_db ${PRINT_SCRAPE_DB_HOST} ${PRINT_SCRAPE_DB_PORT} ${PRINT_SCRAPE_DB_NAME} ${PRINT_SCRAPE_DB_USER} ${PRINT_SCRAPE_DB_PASSWORD}
check_success

# In case test database config exists, run the migrations on it too
if [ ! -z ${PRINT_SCRAPE_TEST_DB_HOST+x} ] && [ ! -z ${PRINT_SCRAPE_TEST_DB_PORT+x} ] && [ ! -z ${PRINT_SCRAPE_TEST_DB_USER+x} ]; then
    # Run test migration
    echo "Migrating test database..."
    migrate_db ${PRINT_SCRAPE_TEST_DB_HOST} ${PRINT_SCRAPE_TEST_DB_PORT} ${PRINT_SCRAPE_DB_NAME} ${PRINT_SCRAPE_TEST_DB_USER} ${PRINT_SCRAPE_DB_PASSWORD}
    check_success
fi
echo "Migration succeeded, starting print-scrape..."