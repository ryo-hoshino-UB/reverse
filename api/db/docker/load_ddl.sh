#!/bin/bash
cat data/001_init.sql | docker-compose exec -T mysql mysql --user=root --password=rootpass 