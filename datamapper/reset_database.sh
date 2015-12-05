#!/bin/bash

cat ../sql/ddl* ../sql/functions/* ./dummy_data.sql | psql -h $PGHOST -U postgres -d lecture
