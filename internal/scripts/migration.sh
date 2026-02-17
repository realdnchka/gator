#!/bin/bash
cd sql/schema/ && goose postgres "postgres://rd:@localhost:5432/gator" down
