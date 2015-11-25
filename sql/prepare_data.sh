#!/bin/bash
cat lecture_database.sql permissions.sql dummy_data.sql materialized_view.sql dml_functions.sql projection_functions.sql | psql -U postgres -h $PGHOST
