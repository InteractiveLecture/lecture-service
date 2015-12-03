#!/bin/bash
cat ddl/* functions/* | psql -U postgres -h $PGHOST
