#!/bin/bash
cat data/001_init.sql data/002_ddl.sql | docker-compose exec -T mysql mysql --user=root --password=rootpass 