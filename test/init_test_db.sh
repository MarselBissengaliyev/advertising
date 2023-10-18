#!/bin/bash

# Set the database connection parameters
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=
DB_NAME=advertising

# Execute the SQL script
psql -h $DB_HOST -p $DB_PORT -U $DB_USERNAME -d $DB_NAME -a -f init_test_db.sql
