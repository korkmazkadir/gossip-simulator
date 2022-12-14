#!/bin/bash

type="classic"


fanout=8

batch_size=500

data_chunk_count=48

parity_chunk_count=80


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

#fault_percents=(5 10 15 20 25 30 35 40 45 50 55 60 65 70 75 80 85 90)
#fault_percents=(45)
#for i in "${fault_percents[@]}"


output_file="sim_fanout${fanout}_datachunk${data_chunk_count}_experiment${batch_size}.out"

echo "Output file: ${output_file}"

for i in {0..99}
do

    type="ida"

    start=`date +%s`

        fault=$(jq -n ${i}/100) 
        echo -n "Fanout ${fanout} Fault ${fault} => "
        ./cmd \
        -type=${type} \
        -d=${fanout} \
        -f=${fault} \
        -e=${batch_size} \
        -dc=${data_chunk_count} \
        -pc=${parity_chunk_count} \
        >> ${output_file}

    end=`date +%s`
    elapsed_time=$((end-start))

    echo  "ElapsedTime: ${elapsed_time}"

    type="classic"

    start=`date +%s`

        fault=$(jq -n ${i}/100) 
        echo -n "Fanout ${fanout} Fault ${fault} => "
        ./cmd \
        -type=${type} \
        -d=${fanout} \
        -f=${fault} \
        -e=${batch_size} \
        -dc=${data_chunk_count} \
        >> ${output_file}

    end=`date +%s`
    elapsed_time=$((end-start))

    echo  "ElapsedTime: ${elapsed_time}"

done


exit

type="classic"

echo "------${type}------"

#for i in "${fault_percents[@]}"
for i in {0..100}
do
    start=`date +%s`

        fault=$(jq -n ${i}/100) 
        echo -n "Fanout ${fanout} Fault ${fault} => "
        ./cmd -type=${type} -d=${fanout} -f=${fault} -e=${batch_size} >> simulation.out

    end=`date +%s`
    elapsed_time=$((end-start))

    echo  "ElapsedTime: ${elapsed_time}"
done
