#!/bin/bash

cat lecture_database.sql permissions.sql dummy_data.sql materialized_view.sql functions.sql | psql -U postgres -h localhost
