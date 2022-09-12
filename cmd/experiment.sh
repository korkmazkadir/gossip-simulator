#!/bin/bash

type="classic"


fanout=8

batch_size=100

# echo "${type}"

# for i in {1..9}
# do
#     fault=$(jq -n ${i}/100)
#     echo -n "Fanout ${fanout} Fault ${fault} => " 
#     ./cmd -type=${type} -d=${fanout} -f=${fault}
# done

# echo "------------------"

type="ida"

echo "------${type}------"

fault_percents=(5 10 15 20 25 30 35 40 45 50 55 60 65 70 75 80 85 90)
#fault_percents=(45)

#for i in "${fault_percents[@]}"
for i in {0..100}
do
    fault=$(jq -n ${i}/100)
    echo -n "Fanout ${fanout} Fault ${fault} => " 
    ./cmd -type=${type} -d=${fanout} -f=${fault} -e=${batch_size} >> simulation.out
done


type="classic"

echo "------${type}------"

#for i in "${fault_percents[@]}"
for i in {0..100}
do
    fault=$(jq -n ${i}/100)
    echo -n "Fanout ${fanout} Fault ${fault} => " 
    ./cmd -type=${type} -d=${fanout} -f=${fault} -e=${batch_size} >> simulation.out
done
