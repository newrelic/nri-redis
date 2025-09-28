# Add data to Database 0
redis-cli -n 0 LPUSH list_db0 "item1" "item2" "item3"

# Add data to Database 1
redis-cli -n 1 LPUSH list_db1 "a" "b" "c"

# Add data to Database 2
redis-cli -n 2 LPUSH list_db2 "x" "y" "z"

echo "Finished initializing Redis data."
