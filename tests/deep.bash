#!/bin/bash

PROG=./clon

T_deepString() {
    result="$($PROG \
    name=John \
    wife[mother][name]=Agatha \
    wife[name]=Jane \
  )"
  expected='{"name":"John","wife":{"mother":{"name":"Agatha"},"name":"Jane"}}'
  [[ "${result}" == "${expected}" ]]
}

T_deepComplex() {
    result="$($PROG \
    name=John \
    wife[age]:=30 \
    wife[mother][name]=Agatha \
    wife[mother][age]:=87 \
    wife[name]=Jane \
    wife[mother][female]:=true \
  )"
  expected='{"name":"John","wife":{"age":30,"mother":{"name":"Agatha","age":87,"female":true},"name":"Jane"}}'
  [[ "${result}" == $expected ]]
}
