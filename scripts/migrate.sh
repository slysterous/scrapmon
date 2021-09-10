#!/bin/bash
set -e 

migrate_db() {

    echo "Number of arguments: $#";
    echo "First argument: $1";
    echo "Second argument: $2";

    # Attempt to run migrations, retry if the database container is not yet ready
    echo "Running migrations on $1"
    i=0
    
    #docker run -it -v $(pwd)/internal/migrations/:/migrations --network host migrate/migrate -path=/migrations/ -database "postgres://$4:$5@$1:$2/$3?sslmode=disable" up  

    migrate -database postgres://$4:$5@$1:$2/$3?sslmode=disable -verbose -source file:../migrations up
    #./migrate -database postgres://postgres:password@scrapmon-db:5432/scrapmon?sslmode=disable -verbose -source file://./migrations up
    while [ $? -ne 0 -a $i -lt 60 ]; do
        echo "Database not ready (attempt #$i), retrying.."
        echo "migrate -database postgres://$4:$5@$1:$2/$3?sslmode=disable -verbose -source file://./internal/migrations up"
        sleep 2
        i=`expr $i + 1`
         migrate -database postgres://$4:$5@$1:$2/$3?sslmode=disable -verbose -source file:../migrations up
         #docker run -it -v $(pwd)/internal/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgres://$4:$5@$1:$2/$3?sslmode=disable" up  
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
migrate_db ${SCRAPMON_DB_HOST} ${SCRAPMON_DB_PORT} ${SCRAPMON_DB_NAME} ${SCRAPMON_DB_USER} ${SCRAPMON_DB_PASSWORD}
check_success

# In case test database config exists, run the migrations on it too
if [ ! -z ${SCRAPMON_TEST_DB_HOST+x} ] && [ ! -z ${SCRAPMON_TEST_DB_PORT+x} ] && [ ! -z ${SCRAPMON_TEST_DB_USER+x} ]; then
    # Run test migration
    echo "Migrating test database..."
    migrate_db ${SCRAPMON_TEST_DB_HOST} ${SCRAPMON_TEST_DB_PORT} ${SCRAPMON_DB_NAME} ${SCRAPMON_TEST_DB_USER} ${SCRAPMON_DB_PASSWORD}
    check_success
fi
echo "Migration succeeded, starting scrapmon..."