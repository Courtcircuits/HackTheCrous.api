#!/bin/bash

sql_script=$( cat create_tables.sql)
pg_database=$(cat ../.env | grep PG_DATABASE | cut -d '=' -f2)
pg_user=$(cat ../.env | grep PG_USER | cut -d '=' -f2)
pg_host=$(cat ../.env | grep PG_HOST | cut -d '=' -f2)

echo "$sql_script" | psql -U "$pg_user" -d "$pg_database" -h "$pg_host"
