#!/bin/bash
cat migration.sql | docker exec -i scopex-mysql mysql -uhomestead -p!Secret1234 scopex-assignment
echo "Migration applied successfully."
